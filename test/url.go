package main

import "fmt"

func main() {
	fmt.Println(StringHTTP("http://www.baidu.com:15000/"))
}

func StringHTTP(url string) string {
	s := []byte(url)
	if s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	url = string(s)
	return url
}
