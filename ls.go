package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	path := "."
	var rootStat unix.Stat_t
	err := unix.Lstat(path, &rootStat)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootFd, err := unix.Open(path, 0, 777)
	if err != nil {
		fmt.Printf("open error: %v\n", err)
		os.Exit(1)
	}
	buf := make([]byte, 8192)
	numEntries := -1
	bufLen := 0
	currentBuf := buf
	for numEntries != 0 {
		numEntries, err = unix.ReadDirent(rootFd, currentBuf)
		if err != nil {
			fmt.Printf("ReadDirent error: %v\n", err)
			os.Exit(1)
		}
		bufLen += numEntries
		currentBuf = buf[bufLen:]
	}

	buf = buf[:bufLen]
	bytesLeft := bufLen
	for bytesLeft != 0 {
		consumed, _, newnames := unix.ParseDirent(buf, bufLen, nil)
		for _, name := range newnames {
			fmt.Printf("%v\n", name)
		}
		bytesLeft -= consumed
		buf = buf[consumed:]
	}

	err = unix.Close(rootFd)
	if err != nil {
		fmt.Printf("close error: %v\n", err)
		os.Exit(1)
	}
}
