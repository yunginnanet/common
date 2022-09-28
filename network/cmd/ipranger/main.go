package main

import (
	"os"

	"git.tcp.direct/kayos/common/network"
)

func main() {
	if len(os.Args) < 2 {
		println("mising input")
		return
	}
	iter := network.IterateNetRange(os.Args[1])
	if iter == nil {
		println("invalid input")
		return
	}
	for item := range iter {
		println(item.String())
	}
}
