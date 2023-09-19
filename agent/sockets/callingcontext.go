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
}

func GetCallingContext(connection net.Conn) CallingContext {
	creds, err := peercred.Get(connection)
	if err != nil {
		panic(err)
	}
	pid, _ := creds.PID()
	process, err := gops.FindProcess(pid)
	if process.Executable() == "ssh-keygen" {
		process, err = gops.FindProcess(process.PPid())
		pid = process.Pid()
	}

	uid, _ := creds.UserID()
	ppid := process.PPid()
	if err != nil {
		panic(err)
	}

	parentProcess, err := gops.FindProcess(ppid)
	if err != nil {
		panic(err)
	}

	parentParentProcess, err := gops.FindProcess(parentProcess.PPid())
	if err != nil {
		panic(err)
	}

	username, err := user.LookupId(uid)
	if err != nil {
		panic(err)
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
