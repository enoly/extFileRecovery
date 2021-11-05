package main

import (
	"bytes"
	"fmt"
	"syscall"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

func main() {
	disk := "\\\\.\\PHYSICALDRIVE1"
	var numRead int
	var err error

	fd, err := syscall.Open(disk, syscall.O_RDONLY, 0777)

	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	syscall.Seek(fd, 2048, 0)
	buffer := make([]byte, 1024)

	numRead, err = syscall.Read(fd, buffer)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Read %v byte\n", numRead)
	fmt.Println(string(buffer))

	err = syscall.Close(fd)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	sb := NewExt3Superblock()
	err = sb.Read(kaitai.NewStream(bytes.NewReader(buffer)), sb, sb)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}
