package systemauth

import (
	"time"

	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/systemauth/biometrics"
	"github.com/quexten/goldwarden/agent/systemauth/pinentry"
)

const tokenExpiry = 10 * time.Minute

var sessionStore = SessionStore{
	Store: []Session{},
}

type Session struct {
	Pid            int
	ParentPid      int
	GrandParentPid int
	Expires        time.Time
}

type SessionStore struct {
	Store []Session
}

func (s *SessionStore) CreateSession(pid int, parentpid int, grandparentpid int) Session {
	var session = Session{
		Pid:            pid,
		ParentPid:      parentpid,
		GrandParentPid: grandparentpid,
		Expires:        time.Now().Add(tokenExpiry),
	}
	s.Store = append(s.Store, session)
	return session
}

func (s *SessionStore) VerifySession(ctx sockets.CallingContext) bool {
	for _, session := range s.Store {
		if session.ParentPid == ctx.ParentProcessPid && session.GrandParentPid == ctx.GrandParentProcessPid {
			if session.Expires.After(time.Now()) {
				return true
			}
		}
	}
	return false
}

func GetApproval(title string, description string, requriesBiometrics bool) (bool, error) {
	approval, err := pinentry.GetApproval(title, description)
	if err != nil {
		return false, err
	}
	if requriesBiometrics {
		biometricsApproval := biometrics.CheckBiometrics(biometrics.AccessCredential)
		if !biometricsApproval {
			return false, nil
		}
	}
	return approval, nil
}

func CheckBiometrics(callingContext *sockets.CallingContext, approvalType biometrics.Approval) bool {
	if sessionStore.VerifySession(*callingContext) {
		return true
	}

	var approval = biometrics.CheckBiometrics(approvalType)
	if !approval {
		return false
	}

	sessionStore.CreateSession(callingContext.ProcessPid, callingContext.ParentProcessPid, callingContext.GrandParentProcessPid)
	return true
}

func CreateSession(ctx sockets.CallingContext) {
	sessionStore.CreateSession(ctx.ProcessPid, ctx.ParentProcessPid, ctx.GrandParentProcessPid)
}
