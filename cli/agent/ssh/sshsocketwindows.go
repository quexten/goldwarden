//go:build windows

package ssh

import (
	"github.com/Microsoft/go-winio"
	"github.com/quexten/goldwarden/cli/agent/sockets"
	"golang.org/x/crypto/ssh/agent"
)

func (v SSHAgentServer) Serve() {
	pipePath := `\\.\pipe\openssh-ssh-agent`

	l, err := winio.ListenPipe(pipePath, nil)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()
	log.Info("Server listening on named pipe %v\n", pipePath)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
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
