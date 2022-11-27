package main

import (
	"fmt"
	"os/exec"
)

func main() {
	//cmd := exec.Command("sh", "-c", "kill -SIGHUP "+strconv.Itoa(syscall.Getpid()))
	//cmd := exec.Command("sh", "-c", "ls")
	//kill -SIGTSTP pid    关闭
	//kill -SIGHUP pid     重启
	cmd := exec.Command("/bin/bash", "-c", "kill -SIGHUP 32625")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
