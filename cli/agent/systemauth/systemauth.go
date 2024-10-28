package systemauth

import (
	"fmt"
	"math"
	"time"

	"github.com/quexten/goldwarden/cli/agent/config"
	"github.com/quexten/goldwarden/cli/agent/sockets"
	"github.com/quexten/goldwarden/cli/agent/systemauth/biometrics"
	"github.com/quexten/goldwarden/cli/agent/systemauth/pinentry"
	"github.com/quexten/goldwarden/cli/logging"
)

var log = logging.GetLogger("Goldwarden", "Systemauth")

const tokenExpiry = 60 * time.Minute
const SSHTTL = 60 * time.Minute

type SessionType string

const (
	AccessVault SessionType = "com.quexten.goldwarden.accessvault"
	SSHKey      SessionType = "com.quexten.goldwarden.usesshkey"
	Pin         SessionType = "com.quexten.goldwarden.pin"
)

var sessionStore = SessionStore{
	Store: []Session{},
}

type Session struct {
	Pid            int
	ParentPid      int
	GrandParentPid int
	Expires        time.Time
	sessionType    SessionType
}

type SessionStore struct {
	Store []Session
}

func (s *SessionStore) CreateSession(pid int, parentpid int, grandparentpid int, sessionType SessionType, ttl time.Duration) Session {
	var session = Session{
		Pid:            pid,
		ParentPid:      parentpid,
		GrandParentPid: grandparentpid,
		Expires:        time.Now().Add(ttl),
		sessionType:    sessionType,
	}
	s.Store = append(s.Store, session)
	return session
}

func (s *SessionStore) verifySession(ctx sockets.CallingContext, sessionType SessionType) bool {
	for _, session := range s.Store {
		if session.sessionType == sessionType {
			if session.Expires.After(time.Now()) {
				return true
			}
		}
	}
	return false
}

// with session
func GetPermission(sessionType SessionType, ctx sockets.CallingContext, config *config.Config) (bool, error) {
	if ctx.Authenticated {
		return true, nil
	}

	log.Info("Checking permission for " + ctx.ProcessName + " with session type " + string(sessionType))
	var actionDescription = ""
	biometricsApprovalType := biometrics.AccessVault
	switch sessionType {
	case AccessVault:
		actionDescription = "access the vault"
		biometricsApprovalType = biometrics.AccessVault
	case SSHKey:
		actionDescription = "use an SSH key for signing"
		biometricsApprovalType = biometrics.SSHKey
	}
	var message = fmt.Sprintf("Do you want to authorize %s>%s>%s to %s? (This choice will be remembered for %d minutes)", ctx.GrandParentProcessName, ctx.ParentProcessName, ctx.ProcessName, actionDescription, int(math.Floor(tokenExpiry.Minutes())))

	if sessionStore.verifySession(ctx, sessionType) {
		log.Info("Permission granted from cached session")
	} else {
		if !sessionStore.verifySession(ctx, Pin) {
			if biometrics.BiometricsWorking() {
				biometricsApproval := biometrics.CheckBiometrics(biometricsApprovalType)
				if !biometricsApproval {
					return false, nil
				}
			} else {
				log.Warn("Biometrics is not available, asking for pin")
				pin, err := pinentry.GetPassword("Enter PIN", "Biometrics is not available. Enter your pin to authorize this action. "+message)
				if err != nil {
					return false, err
				}
				if !config.VerifyPin(pin) {
					return false, nil
				}
			}
		}

		// approval, err := pinentry.GetApproval("Goldwarden authorization", message)
		// if err != nil || !approval {
		// 	return false, err
		// }

		log.Info("Permission granted, creating session")
		sessionStore.CreateSession(ctx.ProcessPid, ctx.ParentProcessPid, ctx.GrandParentProcessPid, sessionType, tokenExpiry)
	}
	return true, nil
}

// no session
func CheckBiometrics(callingContext *sockets.CallingContext, approvalType biometrics.Approval) bool {
	var message = fmt.Sprintf("Do you want to grant %s>%s>%s one-time access your vault?", callingContext.GrandParentProcessName, callingContext.ParentProcessName, callingContext.ProcessName)
	var bioApproval = biometrics.CheckBiometrics(approvalType)
	if !bioApproval {
		return false
	}

	approval, err := pinentry.GetApproval("Goldwarden authorization", message)
	if err != nil {
		log.Error(err.Error())
	}

	return approval
}

func CreatePinSession(ctx sockets.CallingContext, ttl time.Duration) Session {
	return sessionStore.CreateSession(ctx.ProcessPid, ctx.ParentProcessPid, ctx.GrandParentProcessPid, Pin, ttl)
}

func VerifyPinSession(ctx sockets.CallingContext) bool {
	return sessionStore.verifySession(ctx, Pin)
}

func CreateSSHSession(ctx sockets.CallingContext) Session {
	return sessionStore.CreateSession(ctx.ProcessPid, ctx.ParentProcessPid, ctx.GrandParentProcessPid, SSHKey, SSHTTL)
}

func GetSSHSession(ctx sockets.CallingContext) bool {
	return sessionStore.verifySession(ctx, SSHKey)
}

func WipeSessions() {
	sessionStore.Store = []Session{}
}
