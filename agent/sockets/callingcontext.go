package sockets

import (
	"net"
	"os/user"
	"time"

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
}

func GetCallingContext(connection net.Conn) CallingContext {
	creds, err := peercred.Get(connection)
	errorContext := CallingContext{
		UserName:               "unknown user",
		ProcessName:            "unknown process",
		ParentProcessName:      "unknown parent",
		GrandParentProcessName: "unknown grandparent",
		ProcessPid:             time.Now().UTC().Nanosecond(),
		ParentProcessPid:       time.Now().UTC().Nanosecond(),
		GrandParentProcessPid:  time.Now().UTC().Nanosecond(),
	}
	if err != nil {
		return errorContext
	}
	pid, _ := creds.PID()
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

	uid, _ := creds.UserID()
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

	username, err := user.LookupId(uid)
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
	}
}
