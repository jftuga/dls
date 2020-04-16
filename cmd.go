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
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const version = "1.0.0"

type stats struct {
	fileCount     int
	dirCount      int
	errCount      int
	totalFileSize int64
}

var excludeDirs = map[string]int{
	"dev":  0,
	"proc": 0,
	"sys":  0,
	".git": 0,
}

func outputTable(colHeaders []string, tblData [][]string) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(colHeaders)
	table.SetAutoWrapText(false)
	table.AppendBulk(tblData)
	table.Render()
}

func matchesExclude(entry string) bool {
	prefix := ""
	for key, _ := range excludeDirs {

		max := len(key)
		if max > len(entry) {
			max = len(entry)
		}
		prefix = entry[:max]
		if key == prefix {
			return true
		}
	}
	return false
}

func getMetadata(dirName string, showAll bool) ([][]string, [][]string, stats) {
	var allEntries, allErrors [][]string
	var errCount, fileCount, dirCount int = 0, 0, 0
	var totalFileSize, size int64 = 0, 0
	var fileSize, modTime string = "", ""

	filepath.Walk(dirName,
		func(path string, info os.FileInfo, err error) error {
			if !showAll && (matchesExclude(path) || "." == path) {
				return nil
			}
			if err != nil {
				allErrors = append(allErrors, []string{fmt.Sprintf("%v", err)})
				errCount++
				return nil
			}

			size = info.Size()
			totalFileSize += size
			fileSize = fmt.Sprintf("%9d", size)
			modTime = fmt.Sprintf("%v", info.ModTime())[:19]
			objType := "F"
			if info.IsDir() {
				path += string(os.PathSeparator)
				objType = "D"
				dirCount++
			} else {
				fileCount++
			}
			allEntries = append(allEntries, []string{fileSize, modTime, objType, path})

			return nil
		})

	myStats := new(stats)
	myStats.fileCount = fileCount
	myStats.dirCount = dirCount
	myStats.errCount = errCount
	myStats.totalFileSize = totalFileSize

	return allEntries, allErrors, *myStats
}

func main() {
	argsShowAll := flag.Bool("a", false, "show all files, including .git, dev, proc, and sys")
	argsShowErrors := flag.Bool("e", false, "show file/directory errors")
	argsVersion := flag.Bool("v", false, "show version and then exit")
	argsShowTotal := flag.Bool("t", false, "show total file size of all files")
	flag.Usage = func() {
		pgmName := os.Args[0]
		if strings.HasPrefix(os.Args[0], "./") {
			pgmName = os.Args[0][2:]
		}
		fmt.Fprintf(os.Stderr, "\n%s: Get file info for a directory\n", pgmName)
		fmt.Fprintf(os.Stderr, "usage: %s [directory]\n", pgmName)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nCurrent directory is the default if no other directory is given on cmd-line\n")
	}
	flag.Parse()
	if *argsVersion {
		fmt.Fprintf(os.Stderr, "version %s\n", version)
		os.Exit(1)
	}

	var allEntries, allErrors [][]string
	var myStats stats
	var myDir = "."
	args := flag.Args()
	if len(args) > 0 {
		myDir = args[0]
		fi, err := os.Stat(myDir)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if !fi.Mode().IsDir() {
			log.Printf("'%s' is not a directory\n", myDir)
			os.Exit(1)
		}
	}
	allEntries, allErrors, myStats = getMetadata(myDir, *argsShowAll)

	var colHeader []string
	if *argsShowErrors && len(allErrors) > 0 {
		colHeader = []string{fmt.Sprintf("Errors: %d", myStats.errCount)}
		outputTable(colHeader, allErrors)
	}

	if *argsShowTotal {
		allEntries = append(allEntries, []string{"", fmt.Sprintf("%9d", myStats.totalFileSize), "", "(total size)"})
		var MBytes, GBytes float64 = 0, 0
		MBytes = float64(myStats.totalFileSize) / (1024 * 1024)
		GBytes = MBytes / 1024
		if MBytes > 0.02 {
			allEntries = append(allEntries, []string{"", fmt.Sprintf("%9.2f", MBytes), "", "(MB total size)"})
		}
		if GBytes > 0.02 {
			allEntries = append(allEntries, []string{"", fmt.Sprintf("%9.2f", GBytes), "", "(GB total size)"})
		}
	}
	colHeader = []string{"Size", "Mod Time", "Type", fmt.Sprintf("Name (Files:%d Dirs:%d)", myStats.fileCount, myStats.dirCount)}
	outputTable(colHeader, allEntries)
}
