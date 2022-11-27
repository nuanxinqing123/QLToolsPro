package reload

import (
	"os"
	"syscall"
)

type sighandle map[os.Signal]HandleFunc

// HandleFunc 处理函数
type HandleFunc func(s Service)

var defaultSignals []os.Signal

func init() {
	// 需要监听的信号
	defaultSignals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	}
}
