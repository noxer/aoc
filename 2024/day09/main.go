package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"time"
)

func main() {
	var err error

	//	task1("example.txt")

	if len(os.Args) <= 1 {
		fmt.Println("Missing argument, please specify the task you want to execute (1 or 2).")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "1":
		err = task1(os.Args[2:]...)
	case "2":
		err = task2(os.Args[2:])
	default:
		fmt.Println("Invalid argument, please specify the task you want to execute (1 or 2).")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error executing task %s: %s\n", os.Args[1], err)
		os.Exit(1)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

type File struct {
	ID   int
	Size int
}

type FileSystem struct {
	Files []File
}

func (fs *FileSystem) Defrag() {
	freeIndex := fs.FreeIndex()
	fileIndex := fs.FileIndex()

	for freeIndex < fileIndex && freeIndex >= 0 && fileIndex >= 0 {
		fs.DefragFile(freeIndex, fileIndex)

		freeIndex = fs.FreeIndex()
		fileIndex = fs.FileIndex()
	}
}

func (fs *FileSystem) FreeIndex() int {
	for i, f := range fs.Files {
		if f.ID < 0 && f.Size > 0 {
			return i
		}
	}

	return -1
}

func (fs *FileSystem) FileIndex() int {
	for i := len(fs.Files) - 1; i >= 0; i-- {
		f := fs.Files[i]
		if f.ID >= 0 && f.Size > 0 {
			return i
		}
	}

	return -1
}

func (fs *FileSystem) DefragFile(freeIndex, fileIndex int) {
	free := fs.Files[freeIndex]
	file := fs.Files[fileIndex]

	if free.Size == file.Size {
		fs.Files[freeIndex] = file
		fs.Files[fileIndex] = free
		return
	}

	if free.Size < file.Size {
		fs.Files[freeIndex].ID = file.ID
		fs.Files[fileIndex].Size -= free.Size
		return
	}

	fs.Files[freeIndex].Size -= file.Size
	fs.Files[fileIndex].ID = -1
	fs.Files = slices.Insert(fs.Files, freeIndex, file)
}

func (fs *FileSystem) Checksum() int {
	index := 0
	checksum := 0

	for _, file := range fs.Files {
		if file.Size <= 0 {
			continue
		}

		if file.ID < 0 {
			index += file.Size
			continue
		}

		for range file.Size {
			checksum += index * file.ID
			index++
		}
	}

	return checksum
}

func parseDiskMap(name string) (*FileSystem, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	isFile := true
	id := 0
	fs := &FileSystem{}
	for b, err := r.ReadByte(); err == nil; b, err = r.ReadByte() {
		f := File{
			Size: int(b - '0'),
		}

		if isFile {
			f.ID = id
			id++
		} else {
			f.ID = -1
		}

		isFile = !isFile
		fs.Files = append(fs.Files, f)
	}

	return fs, nil
}

func task1(args ...string) error {
	fs, err := parseDiskMap(args[0])
	if err != nil {
		return err
	}

	start := time.Now()

	fs.Defrag()
	checksum := fs.Checksum()

	elapsed := time.Since(start)

	fmt.Printf("Checksum: %d (%s)\n", checksum, elapsed)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (fs *FileSystem) FreeIndexWithSize(size int) int {
	for i, f := range fs.Files {
		if f.ID < 0 && f.Size >= size {
			return i
		}
	}

	return -1
}

func (fs *FileSystem) FileIndexWithID(id int) int {
	for i, f := range fs.Files {
		if f.ID == id {
			return i
		}
	}

	return -1
}

func (fs *FileSystem) DefragSpace() {
	maxID := fs.Files[len(fs.Files)-1].ID
	for id := maxID; id >= 0; id-- {
		fileIndex := fs.FileIndexWithID(id)
		freeIndex := fs.FreeIndexWithSize(fs.Files[fileIndex].Size)

		if freeIndex < 0 || freeIndex > fileIndex {
			continue
		}

		fs.DefragFile(freeIndex, fileIndex)
	}
}

func task2(args []string) error {
	fs, err := parseDiskMap(args[0])
	if err != nil {
		return err
	}

	start := time.Now()

	fs.DefragSpace()
	checksum := fs.Checksum()

	elapsed := time.Since(start)

	fmt.Printf("Checksum: %d (%s)\n", checksum, elapsed)

	return nil
}
