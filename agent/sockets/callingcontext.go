package sockets

import (
	"net"
	"os/user"

	gops "github.com/mitchellh/go-ps"
	"inet.af/peercred"
)

type CallingContext struct {
	UserName               string
	ProcessName            string
	ParentProcessName      string
	GrandParentProcessName string
	ProcessPid             int
	ParentProcessPid       int
	GrandParentProcessPid  int
	Error                  bool
	Authenticated          bool
}

func GetCallingContext(connection net.Conn) CallingContext {
	creds, err := peercred.Get(connection)
	errorContext := CallingContext{
		UserName:               "unknown",
		ProcessName:            "unknown",
		ParentProcessName:      "unknown",
		GrandParentProcessName: "unknown",
		ProcessPid:             0,
		ParentProcessPid:       0,
		GrandParentProcessPid:  0,
		Error:                  true,
		Authenticated:          false,
	}
	if err != nil {
		return errorContext
	}
	uid, _ := creds.UserID()
	username, err := user.LookupId(uid)
	if err != nil {
		return errorContext
	}
	errorContext.UserName = username.Username

	pid, ok := creds.PID()
	if !ok {
		return errorContext
	}

	process, err := gops.FindProcess(pid)
	if err != nil {
		return errorContext
	}
	if process == nil {
		return errorContext
	}

	// git is epheremal and spawns ssh-keygen and ssh so we need to anchor to git
	if process.Executable() == "ssh-keygen" || process.Executable() == "ssh" {
		p, e := gops.FindProcess(process.PPid())
		if p.Executable() == "git" && e == nil {
			process, err = p, e
			pid = process.Pid()
		}
	}

	ppid := process.PPid()
	if err != nil {
		return errorContext
	}

	parentProcess, err := gops.FindProcess(ppid)
	if err != nil {
		return errorContext
	}

	parentParentProcess, err := gops.FindProcess(parentProcess.PPid())
	if err != nil {
		return errorContext
	}

	return CallingContext{
		UserName:               username.Username,
		ProcessName:            process.Executable(),
		ParentProcessName:      parentProcess.Executable(),
		GrandParentProcessName: parentParentProcess.Executable(),
		ProcessPid:             pid,
		ParentProcessPid:       ppid,
		GrandParentProcessPid:  parentParentProcess.PPid(),
		Error:                  false,
	}
}
