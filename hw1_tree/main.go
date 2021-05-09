package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	//"strings"
)

func countFilesInDir(files []string, path string, printFiles bool) []string {
	if printFiles {
		return files
	} else {
		var folders []string

		for _, val := range files {
			newPath := path + "/" + val
			tmpF, _ := os.Stat(newPath)

			if tmpF.IsDir() {
				folders = append(folders, val)
			}
		}

		return folders
	}
}

func recursiveScan(out io.Writer, path string, modifier string, printFiles bool) {

	f, _ := os.Open(path)
	defer f.Close()

	files, _ := f.Readdirnames(-1)
	files = countFilesInDir(files, path, printFiles)
	filesAmount := len(files)

	sort.Strings(files)

	for i, val := range files {
		newPath := path + "/" + val
		tmpF, _ := os.Stat(newPath)

		switch tmpF.IsDir() {
		case true:
			if i == filesAmount-1 {
				_, _ = fmt.Fprintln(out, modifier+"└───"+val)
				recursiveScan(out, newPath, modifier+"\t", printFiles)
			} else {
				_, _ = fmt.Fprintln(out, modifier+"├───"+val)
				recursiveScan(out, newPath, modifier+"│\t", printFiles)
			}

		case false:
			var fSizeString string

			if tmpF.Size() == 0 {
				fSizeString = " (empty)"
			} else {
				fSizeString = fmt.Sprintf(" (%db)", tmpF.Size())
			}

			if i == filesAmount-1 {
				if printFiles {
					_, _ = fmt.Fprintln(out, modifier+"└───"+val+fSizeString)
				}
			} else {
				if printFiles {
					_, _ = fmt.Fprintln(out, modifier+"├───"+val+fSizeString)
				}
			}
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	recursiveScan(out, path, "", printFiles)
	return nil
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
