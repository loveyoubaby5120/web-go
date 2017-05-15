package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
)

func main() {

	listener, _ := net.Listen("tcp", ":8000")
	for {
		c, _ := listener.Accept()
		go connHandler(c)
	}
}

func connHandler(c net.Conn) {
	defer c.Close()
	e := reflect.TypeOf(c).Elem()
	n := e.NumField()
	for i := 0; i < n; i++ {
		fmt.Fprintf(os.Stdout, "e.Field(%d) = %v           %v\n", i, e.Field(i).Name, e.Field(i))
	}

	fmt.Fprintf(os.Stdout, "====================\n")

	v := reflect.ValueOf(c)
	fmt.Fprintf(os.Stdout, "reflect.ValueOf(c) %v\n", v)

	v = v.Elem()
	fmt.Fprintf(os.Stdout, "reflect.ValueOf(c).Elem() %v\n", v)
	fmt.Fprintf(os.Stdout, "====================\n")
	for i := 0; i < v.NumField(); i++ {
		fmt.Fprintf(os.Stdout, "v.Field(%d) = %v   %v\n", i, v.Field(i), v.Field(i).Type())
		switch v.Field(i).Type().Kind() {
		case reflect.Int:
			fmt.Fprintf(os.Stdout, "%d\n", v.Field(i).Int())
		}
	}
	fmt.Fprintf(os.Stdout, "====================\n")

	fmt.Fprintf(os.Stdout, "reflect.ValueOf(c).Elem().FieldByName(\"fd\") %v\n", v.FieldByName("fd"))
	fmt.Fprintf(os.Stdout, "reflect.ValueOf(c).Elem().FieldByName(\"fd\").Elem() %v\n", reflect.ValueOf(c).Elem().FieldByName("fd").Elem())
	v = reflect.ValueOf(c).Elem().FieldByName("fd").Elem()

	fmt.Fprintf(os.Stdout, "====================\n")
	for i := 0; i < v.NumField(); i++ {
		fmt.Fprintf(os.Stdout, "v.Field(%d) = %v   %v  %s\n", i, v.Field(i), v.Field(i).Type(), v.Field(i).String())
		switch v.Field(i).Type().Kind() {
		case reflect.Int:
			fmt.Fprintf(os.Stdout, "%d\n", v.Field(i).Int())
		}
	}
	fmt.Fprintf(os.Stdout, "====================\n")

	sysfd := v.FieldByName("sysfd")
	fmt.Fprintf(os.Stdout, "v.FieldByName(\"sysfd\") = %v\n", sysfd)

	fd := uintptr(sysfd.Int())
	fmt.Fprintf(os.Stdout, "fd %v\n", fd)

	f := os.NewFile(fd, "listen socket")
	fmt.Fprintf(os.Stdout, "f = %v\n", f)
	conn, err := net.FileConn(f)

	if err != nil {
		fmt.Fprintf(os.Stdout, "err %v\n", err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "c %v\n", conn)
	conn.Write([]byte("hello world ~!\n"))
}
