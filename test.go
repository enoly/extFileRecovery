package main

import (
	"bytes"
	"fmt"
	"syscall"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

func main() {
	disk := "/dev/sdb1"
	var err error

	fd, err := syscall.Open(disk, syscall.O_RDONLY, 0777)

	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	syscall.Seek(fd, 1024, 0)
	rawSuperblock := make([]byte, 1024)

	_, err = syscall.Read(fd, rawSuperblock)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	sb := NewExt3Superblock()
	err = sb.Read(kaitai.NewStream(bytes.NewReader(rawSuperblock)), sb, sb)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	var blockSize uint32 = 1024 << sb.LogBlockSize
	var GDTSeek uint64 = uint64((sb.FirstDataBlock + 1) * blockSize)

	var blockGroupNum uint32 = sb.BlocksCount / sb.BlocksPerGroup
	if sb.BlocksCount%sb.BlocksPerGroup != 0 {
		blockGroupNum++
	}

	var GDTSize uint32 = blockGroupNum * 32
	fmt.Printf("First data block is %v\n", sb.FirstDataBlock)
	fmt.Printf("GDT starts on %v byte from disk start\n", int64(GDTSeek))
	syscall.Seek(fd, int64(GDTSeek), 0)

	rawGDT := make([]byte, GDTSize)
	_, err = syscall.Read(fd, rawGDT)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	GDT := make([]Ext3GdtRecord, 0)
	for i := 0; i < int(blockGroupNum); i++ {
		GDTRecord := NewExt3GdtRecord()
		GDTRecord.Read(kaitai.NewStream(bytes.NewReader(rawGDT[i*32:(i+1)*32])), GDTRecord, GDTRecord)
		GDT = append(GDT, *GDTRecord)
	}

	fmt.Printf("First block group descriptor %+v\n", GDT[0])
	fmt.Printf("Inode size %v\n", sb.InodeSize)

	var inodeTableSeek uint64 = uint64(GDT[0].InodeTableBlock) * uint64(blockSize)
	syscall.Seek(fd, int64(inodeTableSeek), 0)

	rawInodeTable := make([]byte, sb.InodesPerGroup*uint32(sb.InodeSize))
	_, err = syscall.Read(fd, rawInodeTable)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	inodeTable := make([]Ext3Inode, 0)
	for i := 0; i < int(sb.InodesPerGroup); i++ {
		inode := NewExt3Inode()
		inode.Read(kaitai.NewStream(bytes.NewReader(rawInodeTable[i*int(sb.InodeSize):(i+1)*int(sb.InodeSize)])), inode, inode)
		inodeTable = append(inodeTable, *inode)
	}

	for i := 0; i < 11; i++ {
		fmt.Printf("Inode %v: %+v\n", i+1, inodeTable[i])
	}

	err = syscall.Close(fd)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}
