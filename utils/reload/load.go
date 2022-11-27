package reload

import (
	"flag"
	"log"
	"net"
	"os"
	"syscall"
)

var (
	reLoad uint
)

func init() {
	flag.UintVar(&reLoad, "reLoad", 0, "重启次数")
}

// Service 服务
type Service interface {
	SetCanReLoad(count uint)                        // 设置可重启次数
	CanReLoad() bool                                // 是否可以重启
	IsChild() bool                                  // 是否子进程
	SetSigHandle(sig os.Signal, f HandleFunc) error // 设置信号处理
	SetNotifySigs(sigs []os.Signal) error           // 设置需要监控的信号量
	Start()                                         // 开启监控 注意：必须以这个来堵塞
	Reload() (err error)
	Shutdown()
	Wait()
	Add(delat int)
	Done()
}

// NewService 初始化
func NewService(l net.Listener) Service {
	if flag.Parsed() == false {
		panic(" You Must Run Flag.Parse() at Main Pack ! ")
	}

	var s = new(s)
	s.L = l
	s.reLoad = reLoad
	s.sigs = defaultSignals
	s.SigHandle = make(sighandle)
	s.sigChan = make(chan os.Signal)
	s.stopChan = make(chan struct{})

	// 设置默认处理
	s.setDefaultHandle()

	return s
}

// GetListener 取得连接
func GetListener(laddr string) (l net.Listener, err error) {
	if reLoad > 0 {
		// 子进程用文字描述符来接收数据
		f := os.NewFile(3, "")
		l, err = net.FileListener(f)
		if err != nil {
			log.Printf("net.FileListener error:%v\n", err)
			return
		}
		//log.Printf("laddr : %v ,listener: %v \n", laddr, l)

		// 如果是子进程就杀掉父进程
		syscall.Kill(syscall.Getppid(), syscall.SIGTSTP) //干掉父进程
	} else {
		l, err = net.Listen("tcp", laddr)
		if err != nil {
			log.Printf("net.Listen error: %v", err)
			return
		}
	}
	return
}
