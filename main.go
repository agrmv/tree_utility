package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	tab          = "\t"
	middleItem   = "├───"
	continueItem = "│"
	lastItem     = "└───"
)

type File struct {
	name string
	size int64
}

func (file File) String() string {
	if file.size == 0 {
		return file.name + " (empty)"
	}
	return file.name + " (" + strconv.FormatInt(file.size, 10) + "b)"
}

func dirSort(dir []os.FileInfo) {
	sort.SliceStable(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})
}

func isIgnore(info os.FileInfo) bool {
	if info.Name() != ".git" && info.Name() != ".idea" && info.Name() != "README.md" && info.Name() != "dockerfile" {
		return false
	}
	return true
}

func readDir(path string) (error, []os.FileInfo) {
	file, err := os.Open(path)

	dir, err := file.Readdir(-1)

	defer func() {
		if fileErr := file.Close(); fileErr != nil {
			err = fileErr
		}
	}()

	return err, dir
}

func getLastElementIndex(files []os.FileInfo, printFiles bool) int {
	lastIndex := len(files) - 1

	//Считает только папки, игнорируя файлы
	if !printFiles {
		for i := lastIndex; i >= 0; i-- {
			if files[i].IsDir() {
				return i
			}
		}
	}

	return lastIndex
}

func getGraphicsSymbol(graphicsSymbol string, isLastElement bool) (string, string) {
	var prefix string
	var nestedLevelGraphicsSymbol string

	if isLastElement {
		prefix = lastItem
		nestedLevelGraphicsSymbol = graphicsSymbol + tab
	} else {
		prefix = middleItem
		nestedLevelGraphicsSymbol = graphicsSymbol + continueItem + tab
	}

	return prefix, nestedLevelGraphicsSymbol
}

func printDirTree(out io.Writer, path string, printFiles bool, graphicsSymbol string) error {

	err, dir := readDir(path)
	dirSort(dir)

	for i, info := range dir {
		isLastElement := i == getLastElementIndex(dir, printFiles)
		prefix, nestedLevelGraphicsSymbol := getGraphicsSymbol(graphicsSymbol, isLastElement)
		if !isIgnore(info) {
			if info.IsDir() {
				fmt.Fprintf(out, "%s%s\n", graphicsSymbol+prefix, info.Name())
				err = printDirTree(out, filepath.Join(path, info.Name()), printFiles, nestedLevelGraphicsSymbol)
			} else if printFiles {
				fmt.Fprintf(out, "%s%s\n", graphicsSymbol+prefix, File{info.Name(), info.Size()})
			}
		}

	}
	return err
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := printDirTree(out, path, printFiles, "")
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
