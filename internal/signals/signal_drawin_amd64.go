// +build amd64,darwin

package signals

import (
	"os"
	"syscall"
)

const (
	SIGABRT   = syscall.Signal(0x6)
	SIGALRM   = syscall.Signal(0xe)
	SIGBUS    = syscall.Signal(0xa)
	SIGCHLD   = syscall.Signal(0x14)
	SIGCONT   = syscall.Signal(0x13)
	SIGEMT    = syscall.Signal(0x7)
	SIGFPE    = syscall.Signal(0x8)
	SIGHUP    = syscall.Signal(0x1)
	SIGILL    = syscall.Signal(0x4)
	SIGINFO   = syscall.Signal(0x1d)
	SIGINT    = syscall.Signal(0x2)
	SIGIO     = syscall.Signal(0x17)
	SIGIOT    = syscall.Signal(0x6)
	SIGKILL   = syscall.Signal(0x9)
	SIGPIPE   = syscall.Signal(0xd)
	SIGPROF   = syscall.Signal(0x1b)
	SIGQUIT   = syscall.Signal(0x3)
	SIGSEGV   = syscall.Signal(0xb)
	SIGSTOP   = syscall.Signal(0x11)
	SIGSYS    = syscall.Signal(0xc)
	SIGTERM   = syscall.Signal(0xf)
	SIGTRAP   = syscall.Signal(0x5)
	SIGTSTP   = syscall.Signal(0x12)
	SIGTTIN   = syscall.Signal(0x15)
	SIGTTOU   = syscall.Signal(0x16)
	SIGURG    = syscall.Signal(0x10)
	SIGUSR1   = syscall.Signal(0x1e)
	SIGUSR2   = syscall.Signal(0x1f)
	SIGVTALRM = syscall.Signal(0x1a)
	SIGWINCH  = syscall.Signal(0x1c)
	SIGXCPU   = syscall.Signal(0x18)
	SIGXFSZ   = syscall.Signal(0x19)
)

var AllSignals = []os.Signal{
	SIGABRT,
	SIGALRM,
	SIGBUS,
	SIGCHLD,
	SIGCONT,
	SIGEMT,
	SIGFPE,
	SIGHUP,
	SIGILL,
	SIGINFO,
	SIGINT,
	SIGIO,
	SIGIOT,
	SIGKILL,
	SIGPIPE,
	SIGPROF,
	SIGQUIT,
	SIGSEGV,
	SIGSTOP,
	SIGSYS,
	SIGTERM,
	SIGTRAP,
	SIGTSTP,
	SIGTTIN,
	SIGTTOU,
	SIGURG,
	SIGUSR1,
	SIGUSR2,
	SIGVTALRM,
	SIGWINCH,
	SIGXCPU,
	SIGXFSZ,
}
