package main

import (
	"fmt"
	"os"

	ext4 "github.com/enoly/extFileRecovery/pkg/ext4"
)

func main() {
	args := os.Args[1:]
	storage := ""

	if len(args) > 0 {
		storage = args[0]
	} else {
		fmt.Println("Please specify arguments!")
		os.Exit(-1)
	}

	ext4fs, err := ext4.New(storage)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer ext4fs.Close()
	// fmt.Printf("Superblock:\n%+v\n", ext4fs.Superblock)
	// fmt.Println("Groups per flex_bg: ", ext4fs.GroupsPerFlex)

	// fmt.Println("\nGDT:")
	// for _, record := range ext4fs.GDT {
	// 	fmt.Printf("%+v\n", record)
	// }

	// firstGroupBitmap, err := ext4fs.ReadBlockGroupBitmap(0)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }

	// lastUsedBlockPtr := uint32(0)
	// for i := ext4fs.Superblock.BlocksPerGroup - 1; i > 0; i-- {
	// 	if firstGroupBitmap.Flags[i] {
	// 		lastUsedBlockPtr = i
	// 	}
	// }

	// lastUsedBlock, err := ext4fs.ReadBlock(lastUsedBlockPtr)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }
	// fmt.Printf("Last used block from first group:\n%+v\n", lastUsedBlock)

	// inode, err := ext4fs.ReadInode(15)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }
	// fmt.Printf("Inode 2:\n%+v\n", inode)

	// extent, err := ext4fs.GetExtentFromInode(inode)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }
	// fmt.Printf("Inode 2 extent:\n%+v\n", extent)
	// fmt.Println("Inode 2 extent internal nodes:")
	// for _, node := range extent.InternalNodes {
	// 	fmt.Printf("Inode 2 internal node: %+v\n", node)
	// }
	// fmt.Println("Inode 2 extent leaf nodes:")
	// for _, node := range extent.LeafNodes {
	// 	fmt.Printf("Inode 2 leaf node: %+v\n", node)
	// }

	// extent, err := ext4fs.GetExtentFromBlock(33247)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }
	// fmt.Printf("Inode 2 extent:\n%+v\n", extent)
	// fmt.Println("Inode 2 extent internal nodes:")
	// for _, node := range extent.InternalNodes {
	// 	fmt.Printf("Inode 2 internal node: %+v\n", node)
	// }
	// fmt.Println("Inode 2 extent leaf nodes:")
	// for _, node := range extent.LeafNodes {
	// 	fmt.Printf("Inode 2 leaf node: %+v\n", node)
	// }

	// blockUsedFlag, err := ext4fs.IsBlockUsed(33247)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }
	// fmt.Printf("Is block 33247 used: %v\n", blockUsedFlag)

	restoreFromExtents(ext4fs)
}

func restoreFromExtents(ext4fs *ext4.Ext4) {
	for i := uint32(1); i < ext4fs.Superblock.BlocksCount; i++ {
		if ext4fs.CheckExtentMagic(i) {
			extent, err := ext4fs.GetExtentFromBlock(i)
			if err != nil || extent == nil {
				continue
			}

			fmt.Printf("Found indirect extent, try to restore...\n")
			file, err := ext4fs.ReadFileFromExtent(extent)
			if err != nil {
				fmt.Println("unable to restore file from extent")
				continue
			}

			if err := os.WriteFile(fmt.Sprintf("fileFromExtent%d", i), *file, 0777); err != nil {
				fmt.Printf("Unable to store resotred %v: %v\n", fmt.Sprintf("fileFromExtent%d", i), err)
			}
		}
	}
}
