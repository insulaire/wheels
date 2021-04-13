package main

import "wheels/pkg/tcp"

func main() {
	_ = tcp.NewClient("chen")
	select {}
}
