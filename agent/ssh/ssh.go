package ssh

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/sockets"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/systemauth/pinentry"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/logging"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var log = logging.GetLogger("Goldwarden", "SSH")

type vaultAgent struct {
	vault               *vault.Vault
	config              *config.Config
	unlockRequestAction func() bool
	context             sockets.CallingContext
}

func (vaultAgent) Add(key agent.AddedKey) error {
	return nil
}

func (vaultAgent vaultAgent) List() ([]*agent.Key, error) {
	if vaultAgent.vault.Keyring.IsLocked() {
		if !vaultAgent.unlockRequestAction() {
			return nil, errors.New("vault is locked")
		}

		systemauth.CreatePinSession(vaultAgent.context)
	}

	vaultSSHKeys := (*vaultAgent.vault).GetSSHKeys()
	var sshKeys []*agent.Key
	for _, vaultSSHKey := range vaultSSHKeys {
		signer, err := ssh.ParsePrivateKey([]byte(vaultSSHKey.Key))
		if err != nil {
			continue
		}
		pub := signer.PublicKey()
		sshKeys = append(sshKeys, &agent.Key{
			Format:  pub.Type(),
			Blob:    pub.Marshal(),
			Comment: vaultSSHKey.Name})
	}

	return sshKeys, nil
}

func (vaultAgent) Lock(passphrase []byte) error {
	return nil
}

func (vaultAgent) Remove(key ssh.PublicKey) error {
	return nil
}

func (vaultAgent) RemoveAll() error {
	return nil
}

func Eq(a, b ssh.PublicKey) bool {
	return 0 == bytes.Compare(a.Marshal(), b.Marshal())
}

func (vaultAgent vaultAgent) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	log.Info("Sign Request for key: %s", ssh.FingerprintSHA256(key))
	if vaultAgent.vault.Keyring.IsLocked() {
		if !vaultAgent.unlockRequestAction() {
			return nil, errors.New("vault is locked")
		}

		systemauth.CreatePinSession(vaultAgent.context)
	}

	var signer ssh.Signer
	var sshKey *vault.SSHKey

	vaultSSHKeys := (*vaultAgent.vault).GetSSHKeys()
	for _, vaultSSHKey := range vaultSSHKeys {
		sg, err := ssh.ParsePrivateKey([]byte(vaultSSHKey.Key))
		if err != nil {
			return nil, err
		}
		if Eq(sg.PublicKey(), key) {
			signer = sg
			sshKey = &vaultSSHKey
			break
		}
	}

	message := fmt.Sprintf("%s on %s>%s>%s is requesting signage with key %s", vaultAgent.context.UserName, vaultAgent.context.GrandParentProcessName, vaultAgent.context.ParentProcessName, vaultAgent.context.ProcessName, sshKey.Name)

	if approved, err := pinentry.GetApproval("SSH Key Signing Request", message); err != nil || !approved {
		log.Info("Sign Request for key: %s denied", sshKey.Name)
		return nil, errors.New("Approval not given")
	}

	if permission, err := systemauth.GetPermission(systemauth.SSHKey, vaultAgent.context, vaultAgent.config); err != nil || !permission {
		log.Info("Sign Request for key: %s denied", key.Marshal())
		return nil, errors.New("Biometrics not checked")
	}

	var rand = rand.Reader
	log.Info("Sign Request for key: %s %s accepted", ssh.FingerprintSHA256(key), sshKey.Name)
	return signer.Sign(rand, data)
}

func (vaultAgent) Signers() ([]ssh.Signer, error) {

	return []ssh.Signer{}, nil
}

func (vaultAgent) Unlock(passphrase []byte) error {
	return nil
}

type SSHAgentServer struct {
	vault               *vault.Vault
	config              *config.Config
	unlockRequestAction func() bool
}

func (v *SSHAgentServer) SetUnlockRequestAction(action func() bool) {
	v.unlockRequestAction = action
}

func NewVaultAgent(vault *vault.Vault, config *config.Config) SSHAgentServer {
	return SSHAgentServer{
		vault:  vault,
		config: config,
		unlockRequestAction: func() bool {
			log.Info("Unlock Request, but no action defined")
			return false
		},
	}
}

func (v SSHAgentServer) Serve() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	path := home + "/.goldwarden-ssh-agent.sock"

	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			log.Error("Could not remove old socket file: %s", err)
			return
		}
	}
	listener, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}

	log.Info("SSH Agent listening on %s", path)

	for {
		var conn, err = listener.Accept()
		if err != nil {
			panic(err)
		}

		callingContext := sockets.GetCallingContext(conn)

		log.Info("SSH Agent connection from %s>%s>%s \nby user %s", callingContext.GrandParentProcessName, callingContext.ParentProcessName, callingContext.ProcessName, callingContext.UserName)
		log.Info("SSH Agent connection accepted")

		go agent.ServeAgent(vaultAgent{
			vault:               v.vault,
			config:              v.config,
			unlockRequestAction: v.unlockRequestAction,
			context:             callingContext,
		}, conn)
	}
}
