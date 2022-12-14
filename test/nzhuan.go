package main

import "fmt"

func main() {
	var a = []byte("123\n456")
	fmt.Println(a)
	fmt.Println(fmt.Sprintf("%s", a))
}
