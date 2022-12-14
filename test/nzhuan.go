package main

import (
	"fmt"
	"strings"
)

func main() {
	//var a = []byte("1\n2\n3\n4\n5\n6\n7\n8\n9")
	var a = []byte("12dsf345;6; 7a89")
	fmt.Println(a)
	fmt.Println("1\\n2\\n3\\n4\\n5\\n6\\n7\\n8\\n9")
	//fmt.Println(fmt.Sprintf("%s", a))
	ss := fmt.Sprintf("%s", a)
	sss := strings.Split(ss, "\n")
	//fmt.Println(sss)
	fmt.Println(sss)
	fmt.Println(len(sss))
}
