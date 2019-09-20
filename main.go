package main

import (
	"fmt"
	_ "io"
	"os"
	_ "path/filepath"
	"sort"
	_ "strings"
)

func dirSort(dir []os.FileInfo) {
	sort.SliceStable(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})
}

func dirTree(out interface{}, path string, printFiles bool) (err error) {

	file, err := os.Open(path)

	//Readdir считывает содержимое каталога, связанного с файлом, и возвращает фрагмент до n значений
	//Если n <= 0, Readdir возвращает все FileInfo из каталога
	dir, err := file.Readdir(0)

	// Просто file.Close () может вернуть ошибку, но мы не будем об этом знать
	defer func() {
		if fileErr := file.Close(); fileErr != nil {
			err = fileErr
		}
	}()

	dirSort(dir)

	for _, info := range dir {
		fmt.Println(info.Name())
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
