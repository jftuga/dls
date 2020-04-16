/*

dls.go
-John Taylor
Apr-16-2020

List files within a running docker container

Adopted from:
https://yourbasic.org/golang/list-files-in-directory/

To statically compile the program on Linux:
go build -tags netgo -ldflags '-extldflags "-static" -s -w'

*/

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

func main() {
	table := tablewriter.NewWriter(os.Stdout)
	tableErrors := tablewriter.NewWriter(os.Stdout)
	errorCount := 0
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				modTime := fmt.Sprintf("%v", info.ModTime())
				modTime = modTime[:19]
				entry := []string{fmt.Sprintf("%9d", info.Size()), modTime, path}
				table.Append(entry)
			}

			return nil
		})
	if err != nil {
		tableErrors.Append([]string{fmt.Sprintf("%v", err)})
		errorCount++
	}

	if errorCount > 0 {
		tableErrors.SetHeader([]string{fmt.Sprintf("Errors: %d", errorCount)})
		tableErrors.Render()
	}

	table.SetHeader([]string{"Size", "Mod Time", "Name"})
	table.Render()
}
