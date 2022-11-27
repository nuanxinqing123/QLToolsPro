package reload

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

// S 服务
type s struct {
	L         net.Listener
	reLoad    uint // 初始化次数
	canReLoad uint // 允许重启次数，0为不限制
	SigHandle sighandle
	sigChan   chan os.Signal
	sigs      []os.Signal
	isRun     bool
	stopChan  chan struct{}
	sync.WaitGroup
	sync.RWMutex
}

// Start 开始监听信号量
func (s *s) Start() {
	s.Lock()
	s.isRun = true
	s.Unlock()

	signal.Notify(
		s.sigChan,
		s.sigs...,
	)

	go func() {
		var sig os.Signal
		pid := syscall.Getpid()
		for {
			sig = <-s.sigChan
			log.Println(pid, "Received SIG.", sig)
			if _, ok := s.SigHandle[sig]; ok == true {
				s.SigHandle[sig](s)
			}
		}
	}()

	select {
	case <-s.stopChan:
		return
	}
}

// CanReLoad 是否可以重启
func (s *s) CanReLoad() bool {
	s.RLock()
	defer s.RUnlock()
	if s.canReLoad == 0 {
		return true
	} else if s.reLoad >= s.canReLoad {
		return false
	}
	return true
}

// SetCanReLoad 设置可重启次数
func (s *s) SetCanReLoad(count uint) {
	s.Lock()
	defer s.Unlock()
	s.canReLoad = count
}

// IsChild 是否子进程
func (s *s) IsChild() bool {
	s.RLock()
	defer s.RUnlock()
	return s.reLoad > 0
}

// SetSigHandle 设置信号处理
func (s *s) SetSigHandle(sig os.Signal, f HandleFunc) error {
	if s.isRun == true {
		return fmt.Errorf("不可以在运行中修改")
	}

	s.Lock()
	s.SigHandle[sig] = f
	s.Unlock()

	return nil
}

// SetNotifySigs 设置需要监听的信号量
func (s *s) SetNotifySigs(sigs []os.Signal) error {
	if s.isRun == true {
		return fmt.Errorf("不可以在运行中修改监听信号量")
	}

	s.Lock()
	s.sigs = sigs
	s.Unlock()

	return nil
}

// Shutdown 停止
func (s *s) Shutdown() {
	s.Lock()
	defer s.Unlock()

	s.Wait()

	s.stopChan <- struct{}{}
}

// setDefaultHandle 默认处理
func (s *s) setDefaultHandle() {
	if err := s.SetSigHandle(syscall.SIGHUP, func(s Service) {
		if err := s.Reload(); err != nil {
			fmt.Println(err)
		}
	}); err != nil {
		panic(err)
	}

	// 设置会退出的信号量
	var stopSigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP}
	for _, sig := range stopSigs {
		if err := s.SetSigHandle(sig, func(s Service) {
			s.Shutdown()
		}); err != nil {
			panic(err)
		}
	}

}

// Reload 重启
func (s *s) Reload() (err error) {
	if s.CanReLoad() == false {
		return fmt.Errorf("forked count is %d More The %d ", s.reLoad, s.canReLoad)
	}

	s.Lock()
	defer s.Unlock()
	s.reLoad++

	log.Println("Restart: forked Start....")

	tl := s.L.(*net.TCPListener)
	fl, _ := tl.File()

	path := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			tag := strings.Split(arg, "=")
			if tag[0] == "-reLoad" {
				break
			}
			args = append(args, arg)
		}
	}
	args = append(args, fmt.Sprintf("-reLoad=%d", s.reLoad))

	log.Println(path, args)
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{fl}

	err = cmd.Start()
	if err != nil {
		log.Printf("Restart: Failed to launch, error: %v", err)
	}

	return
}
