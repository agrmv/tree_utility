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
		return true
	}
	return false
}

/*func getFileInfo(file os.FileInfo) string {
	if file.Size() == 0 {
		return file.Name() + " (empty)"
	}
	return file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)"
}*/

func readDir(path string) (err error, files []os.FileInfo) {
	file, err := os.Open(path)

	files, err = file.Readdir(-1)
	dirSort(files)

	defer func() {
		if fileErr := file.Close(); fileErr != nil {
			err = fileErr
		}
	}()

	return err, files
}

func getLastElementIndex(files []os.FileInfo, printFiles bool) int {
	lastIndex := len(files) - 1

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
	var nestedLevelItem string

	if isLastElement {
		prefix = lastItem
		nestedLevelItem = graphicsSymbol + tab
	} else {
		prefix = middleItem
		nestedLevelItem = graphicsSymbol + continueItem + tab
	}

	return prefix, nestedLevelItem

}

func printDir(out io.Writer, path string, printFiles bool, graphicsSymbol string) error {

	err, files := readDir(path)

	for i, info := range files {
		isLastElement := i == getLastElementIndex(files, printFiles)
		prefix, nestedLevelItem := getGraphicsSymbol(graphicsSymbol, isLastElement)

		if info.IsDir() && isIgnore(info) {
			fmt.Fprintf(out, "%s%s\n", graphicsSymbol+prefix, info.Name())
			err = printDir(out, filepath.Join(path, info.Name()), printFiles, nestedLevelItem)
		} else if printFiles && isIgnore(info) {
			fmt.Fprintf(out, "%s%s\n", graphicsSymbol+prefix, File{info.Name(), info.Size()})
		}

	}
	return err
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	err = printDir(out, path, printFiles, "")
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
