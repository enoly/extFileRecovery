package main

import (
	"fmt"
	"os"
	"strings"

	ext3 "github.com/enoly/extFileRecovery/pkg/ext3"
	structure "github.com/enoly/extFileRecovery/pkg/ext3/structure"
)

func main() {
	args := os.Args[1:]
	storage := ""
	directory := "/"
	restoreType := "journal"

	if len(args) > 0 {
		storage = args[0]
	} else {
		fmt.Println("Please specify arguments")
		os.Exit(-1)
	}

	extFs, err := ext3.New(storage)
	if err != nil {
		fmt.Printf("Unable to read ext3 from %v: %v\n", storage, err)
		os.Exit(-1)
	}
	defer extFs.Close()

	if len(args) > 1 {
		restoreTypeArg := args[1]
		switch restoreTypeArg {
		case "indirect":
			restoreType = "indirect"
		case "journal":
			restoreType = "journal"
		case "fragments":
			restoreType = "fragments"
		default:
			{
				fmt.Println("Wrong restore type")
				os.Exit(-1)
			}
		}
	}

	if len(args) > 2 {
		directory = args[2]
	}

	dir, err := readDirectory(extFs, nil, directory)
	if err != nil {
		fmt.Printf("Unable to restore from %v: %v\n", directory, err)
		os.Exit(-1)
	}

	deletedFromDir, err := extFs.GetDeletedFilesFromDir(dir)
	if err != nil {
		fmt.Println("Unable to read deleted files from root directory: ", err)
	}

	for _, entry := range dir.Entries {
		fmt.Printf("%v inode: %v\n", entry.Name, entry.InodePtr)
	}

	if len(deletedFromDir.Entries) > 0 {
		fmt.Println("\ndeleted files")
	}
	for _, entry := range deletedFromDir.Entries {
		fmt.Printf("*%v inode: %v\n", entry.Name, entry.InodePtr)
	}

	switch restoreType {
	case "indirect":
		restoreFromIndirectBlocks(extFs, deletedFromDir)
	case "journal":
		restoreFromJournal(extFs, deletedFromDir)
	case "fragments":
		restoreFragments(extFs)
	default:
		{
			fmt.Println("Wrong restore type")
			os.Exit(-1)
		}
	}
}

func readDirectory(extFs *ext3.Ext3, parent *structure.Directory, path string) (*structure.Directory, error) {
	splittedPath := strings.Split(path, "/")

	if parent == nil && splittedPath[0] == "" {
		root, err := extFs.ReadDirFromInodePtr(2)
		if err != nil {
			return nil, fmt.Errorf("unable to read root directory: %v", err)
		}

		return readDirectory(extFs, root, path[1:])
	}

	if len(splittedPath) == 1 && splittedPath[0] == "" {
		return parent, nil
	}

	for _, entry := range parent.Entries {
		if entry.FileType == structure.Directory_DirEntry_FileTypeEnum__Dir && entry.Name == splittedPath[0] {
			dir, err := extFs.ReadDirFromInodePtr(entry.InodePtr)
			if err != nil {
				return nil, fmt.Errorf("unable to read directory %v: %v", entry.Name, err)
			}

			seek := 0
			if len(splittedPath) > 1 {
				seek = 1
			}

			return readDirectory(extFs, dir, path[len(entry.Name)+seek:])
		}
	}

	return nil, fmt.Errorf("unable to find directory")
}

func restoreFromJournal(ext3fs *ext3.Ext3, deletedFromDir *structure.DeletedFiles) {
	if len(deletedFromDir.Entries) == 0 {
		fmt.Println("Nothing to restore")
		os.Exit(0)
	}

	fmt.Println("Try to restore deleted files from journal...")
	for _, entry := range deletedFromDir.Entries {
		inode, err := ext3fs.GetFreshInodeFromJournal(entry.InodePtr)
		if err != nil || inode == nil {
			fmt.Printf("Inode for %v not found in journal\n", entry.Name)
		}

		file, err := ext3fs.ReadFileFromInode(inode)
		if err != nil {
			fmt.Printf("Unable to restore %v: %v\n", entry.Name, err)
			return
		}

		if err := os.WriteFile(entry.Name, *file, 0777); err != nil {
			fmt.Printf("Unable to store resotred %v: %v\n", entry.Name, err)
		}
	}
}

func restoreFromIndirectBlocks(ext3fs *ext3.Ext3, deletedFromDir *structure.DeletedFiles) {
	if len(deletedFromDir.Entries) == 0 {
		fmt.Println("Nothing to restore")
		os.Exit(0)
	}

	fmt.Println("Try to restore deleted files from indirect blocks...")
	for _, entry := range deletedFromDir.Entries {
		fileBlockGroup := ext3fs.GetInodeBlockGroup(entry.InodePtr)
		blockBitmap, err := ext3fs.GetBlockGroupBitmap(fileBlockGroup)
		if err != nil {
			fmt.Printf("Unable to restore %v from indirect blocks: %v", entry.Name, err)
		}

		inodeTableSize := ext3fs.Superblock.InodesPerGroup * uint32(ext3fs.Superblock.InodeSize) / ext3fs.BlockSize
		blocksStart := ext3fs.GDT[fileBlockGroup].InodeTableBlock + inodeTableSize
		for i := blocksStart; i < (fileBlockGroup+1)*ext3fs.Superblock.BlocksPerGroup && i < ext3fs.Superblock.BlocksCount; i++ {
			if (*blockBitmap)[i-fileBlockGroup*ext3fs.Superblock.BlocksPerGroup] {
				addrBlock, err := readAddrBlock(ext3fs, uint32(i))
				if err != nil || addrBlock == nil {
					continue
				}

				directBlocks := []uint32{}
				for j := addrBlock.BlocksPtrs[0].Ptr - 12; j < addrBlock.BlocksPtrs[0].Ptr; j++ {
					directBlocks = append(directBlocks, j)
				}

				for _, ptr := range addrBlock.BlocksPtrs {
					if ptr.Ptr == 0 {
						break
					}
					directBlocks = append(directBlocks, ptr.Ptr)
				}

				file, err := ext3fs.ReadFileFromPtrs(&directBlocks)
				if err != nil {
					fmt.Println("Unable to restore file from indirect blocks: ", err)
				}

				if err := os.WriteFile(fmt.Sprintf("%dBlock%v", i, entry.Name), *file, 0777); err != nil {
					fmt.Printf("Unable to store resotred %v: %v\n", fmt.Sprintf("fileFromBlock%d", i), err)
				}
			}
		}
	}
}

func restoreFragments(ext3fs *ext3.Ext3) {
	fmt.Println("Try to restore fragments from indirect blocks...")
	blockGroupNum := ext3fs.Superblock.BlocksCount / ext3fs.Superblock.BlocksPerGroup
	if ext3fs.Superblock.BlocksCount%ext3fs.Superblock.BlocksPerGroup != 0 {
		blockGroupNum++
	}

	for blockGroup := uint32(0); blockGroup < blockGroupNum; blockGroup++ {
		blockBitmap, err := ext3fs.GetBlockGroupBitmap(blockGroup)
		if err != nil {
			fmt.Printf("Unable to restore fragments from indirect blocks: %v", err)
		}

		inodeTableSize := ext3fs.Superblock.InodesPerGroup * uint32(ext3fs.Superblock.InodeSize) / ext3fs.BlockSize
		blocksStart := ext3fs.GDT[blockGroup].InodeTableBlock + inodeTableSize
		for i := blocksStart; i < (blockGroup+1)*ext3fs.Superblock.BlocksPerGroup; i++ {
			if (*blockBitmap)[i-blockGroup*ext3fs.Superblock.BlocksPerGroup] {
				innerPtrs, err := recursiveReadPtrs(ext3fs, uint32(i))
				if err != nil || innerPtrs == nil {
					continue
				}

				directBlocks := []uint32{}
				startOfZeroes := 0
				for j := len(*innerPtrs) - 1; j >= 0; j-- {
					if (*innerPtrs)[j] != 0 {
						break
					}

					startOfZeroes = j
				}

				for j, ptr := range *innerPtrs {
					if j == startOfZeroes {
						break
					}
					directBlocks = append(directBlocks, ptr)
				}

				file, err := ext3fs.ReadFileFromPtrs(&directBlocks)
				if err != nil {
					fmt.Println("Unable to restore fragment from indirect blocks: ", err)
				}

				if err := os.WriteFile(fmt.Sprintf("FragmentFromBlock%d", i), *file, 0777); err != nil {
					fmt.Printf("Unable to store resotred %v: %v\n", fmt.Sprintf("FragmentFromBlock%d", i), err)
				}
			}
		}
	}
}

func readAddrBlock(ext3fs *ext3.Ext3, blockPtr uint32) (*structure.IndirectBlock, error) {
	addrBlock, err := ext3fs.ReadIndirectAddrBlock(blockPtr)
	if err != nil {
		return nil, err
	}

	tmp := addrBlock.BlocksPtrs[0].Ptr
	for j := 1; j < len(addrBlock.BlocksPtrs); j++ {
		if addrBlock.BlocksPtrs[j].Ptr != 0 {
			if j < 4 && addrBlock.BlocksPtrs[j].Ptr <= tmp {
				return nil, nil
			}

			if addrBlock.BlocksPtrs[j].Ptr > ext3fs.Superblock.BlocksCount {
				return nil, nil
			}

			tmp = addrBlock.BlocksPtrs[j].Ptr
		}
	}

	if addrBlock.BlocksPtrs[0].Ptr == 0 {
		return nil, nil
	}

	return addrBlock, nil
}

func recursiveReadPtrs(ext3fs *ext3.Ext3, blockPtr uint32) (*[]uint32, error) {
	ptrs := []uint32{}
	addrBlock, err := ext3fs.ReadIndirectAddrBlock(blockPtr)
	if err != nil {
		return nil, err
	}

	tmp := addrBlock.BlocksPtrs[0].Ptr
	for j := 1; j < len(addrBlock.BlocksPtrs); j++ {
		if addrBlock.BlocksPtrs[j].Ptr != 0 {
			if addrBlock.BlocksPtrs[j].Ptr <= tmp {
				return nil, nil
			}

			if addrBlock.BlocksPtrs[j].Ptr > ext3fs.Superblock.BlocksCount {
				return nil, nil
			}

			tmp = addrBlock.BlocksPtrs[j].Ptr
		}
	}

	if addrBlock.BlocksPtrs[0].Ptr == 0 {
		return nil, nil
	}

	errFlag := false
	for _, ptr := range addrBlock.BlocksPtrs {
		innerPtrs, err := recursiveReadPtrs(ext3fs, ptr.Ptr)
		if err != nil || innerPtrs == nil {
			errFlag = true
			break
		}

		ptrs = append(ptrs, *innerPtrs...)
	}

	if errFlag {
		ptrs = []uint32{}
		for _, ptr := range addrBlock.BlocksPtrs {
			ptrs = append(ptrs, ptr.Ptr)
		}
	}

	return &ptrs, nil
}
