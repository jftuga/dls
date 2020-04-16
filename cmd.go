/*

dls.go - docker ls
-John Taylor
Apr-16-2020

List files within a running docker container

To statically compile the program on Linux:
go build -tags netgo -ldflags '-extldflags "-static" -s -w'

*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const version = "1.0.0"

type stats struct {
	fileCount int
	dirCount  int
	errCount  int
}

func outputTable(colHeaders []string, tblData [][]string) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(colHeaders)
	table.SetAutoWrapText(false)
	table.AppendBulk(tblData)
	table.Render()
}

func getMetadata(dirName string) ([][]string, [][]string, stats) {
	var allEntries [][]string
	var allErrors [][]string
	errCount := 0
	fileCount := 0
	dirCount := 0
	fileSize := ""
	modTime := ""

	filepath.Walk(dirName,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				allErrors = append(allErrors, []string{fmt.Sprintf("%v", err)})
				errCount++
				return nil
			}

			fileSize = fmt.Sprintf("%9d", info.Size())
			modTime = fmt.Sprintf("%v", info.ModTime())[:19]
			allEntries = append(allEntries, []string{fileSize, modTime, path})
			if info.IsDir() {
				dirCount++
			} else {
				fileCount++
			}

			return nil
		})

	myStats := new(stats)
	myStats.fileCount = fileCount
	myStats.dirCount = dirCount
	myStats.errCount = errCount

	return allEntries, allErrors, *myStats
}

func main() {

	fmt.Printf("args: %v\n\n", os.Args)

	argsShowErrors := flag.Bool("e", false, "show file/directory errors")
	argsVersion := flag.Bool("v", false, "show version and then exit")

	flag.Usage = func() {
		pgmName := os.Args[0]
		if strings.HasPrefix(os.Args[0], "./") {
			pgmName = os.Args[0][2:]
		}
		fmt.Fprintf(os.Stderr, "\n%s: Get file info for a list of multiple directories\n", pgmName)
		fmt.Fprintf(os.Stderr, "usage: %s [directory ...]\n", pgmName)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()
	if *argsVersion {
		fmt.Fprintf(os.Stderr, "version %s\n", version)
		os.Exit(1)
	}

	args := flag.Args()
	var allEntries, allErrors [][]string
	var myStats stats
	if len(args) == 0 {
		allEntries, allErrors, myStats = getMetadata(".")
	}
	//var allFilenames []string

	var colHeader []string
	if *argsShowErrors && len(allErrors) > 0 {
		colHeader = []string{fmt.Sprintf("Errors: %d", myStats.errCount)}
		outputTable(colHeader, allErrors)
	}

	colHeader = []string{"Size", "Mod Time", fmt.Sprintf("Name: Files:%d Dirs:%d", myStats.fileCount, myStats.dirCount)}
	outputTable(colHeader, allEntries)
}
