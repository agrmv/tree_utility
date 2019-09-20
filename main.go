package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	newLine      = "\n"
	tab          = "\t"
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

func dirSort(dir []os.FileInfo) {
	sort.SliceStable(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})
}

func isIgnore(info os.FileInfo) bool {
	if info.Name() != ".git" && info.Name() != ".idea" && info.Name() != "README.md" {
		return true
	}
	return false
}

func readDir(path string) (err error, files []os.FileInfo) {
	file, err := os.Open(path)

	//Readdir считывает содержимое каталога, связанного с файлом, и возвращает фрагмент до n значений
	//Если n <= 0, Readdir возвращает все FileInfo из каталога
	files, err = file.Readdir(0)

	// Просто file.Close () может вернуть ошибку, но мы не будем об этом знать
	defer func() {
		if fileErr := file.Close(); fileErr != nil {
			err = fileErr
		}
	}()

	return err, files
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	_, dir := readDir(path)
	dirSort(dir)

	var graphicsSymbol strings.Builder
	for range strings.Split(path, "/") {
		graphicsSymbol.WriteString(tab)
	}

	//Доделать форматирование вывода и вынести в отдельную функцию
	for _, node := range dir {
		if node.IsDir() && isIgnore(node) {
			fmt.Fprintf(out, "%s%s\n", graphicsSymbol.String(), node.Name())
			err = dirTree(out, filepath.Join(path, node.Name()), printFiles)
		}
	}

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
