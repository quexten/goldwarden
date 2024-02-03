//go:build !windows

package ssh

import (
	"net"
	"os"

	"github.com/quexten/goldwarden/agent/sockets"
	"golang.org/x/crypto/ssh/agent"
)

func (v SSHAgentServer) Serve() {
	path := v.runtimeConfig.SSHAgentSocketPath
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
