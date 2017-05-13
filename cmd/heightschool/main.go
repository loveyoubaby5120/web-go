package main

import (
	"fmt"
	"strings"
)

func main() {
	http.HandleFunc("/", sayhelloName)       //设置访问的路由
}

