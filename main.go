package main

import (
	"fmt"
	"os"

	ext3 "github.com/enoly/extFileRecovery/pkg/ext3"
)

func main() {
	ext3fs, err := ext3.New("/dev/sdb1")
	if err != nil {
		fmt.Printf("Unable to read ext3 from /dev/sdb1: %v\n", err)
		os.Exit(0)
	}

	rootDir, err := ext3fs.ReadDirFromInodePtr(2)
	if err != nil {
		fmt.Println("Unable to read root directory: ", err)
	}

	deletedFromRootDir, err := ext3fs.GetDeletedFilesFromDir(rootDir)
	if err != nil {
		fmt.Println("Unable to read deleted files from root directory: ", err)
	}

	for _, entry := range rootDir.Entries {
		fmt.Printf("%v\n", entry.Name)
	}

	fmt.Println()
	for _, entry := range deletedFromRootDir.Entries {
		fmt.Printf("*%v inode: %v\n", entry.Name, entry.InodePtr)
	}

	if len(deletedFromRootDir.Entries) == 0 {
		fmt.Println("Nothing to restore")
		os.Exit(0)
	}

	fmt.Println("Try to restore deleted files from journal...")
	for _, entry := range deletedFromRootDir.Entries {
		inode, err := ext3fs.GetFreshInodeFromJournal(entry.InodePtr)
		if err != nil || inode == nil {
			fmt.Printf("Inode for %v not found in journal\n", entry.Name)
		}

		file, err := ext3fs.ReadFileFromInode(inode)
		if err != nil {
			fmt.Printf("Unable to restore %v: %v\n", entry.Name, err)
		}

		if err := os.WriteFile(entry.Name, *file, 0777); err != nil {
			fmt.Printf("Unable to store resotred %v: %v\n", entry.Name, err)
		}
	}
}
