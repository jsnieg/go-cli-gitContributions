package main

// this automatically updates upon writing code, pwetty cool...
// i said pwetty and not pretty cause i can't pronounce R's!
import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// getDotFilePath returns the dot file for the repos list.
//
// Creates it and the enclosing folder if it does not exist.
func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "/.gogitlocalstats"
	return dotFile
}

// openFile opens the file located at 'filePath'.
//
// Creates it if not existing.
func openFile(filePath string) *os.File {
	// Tries to open our dotfile like '.git'
	// f, is the return and err is error code
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0755)

	// If there's an error and if does not exist create the file to fill it with repos scanned.
	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			// other error
			panic(err)
		}
	}

	return f
}

// parseFileLinesToSlice given a filepath string, gets the content
// of each line and parses it to a slice of strings.
func parseFileLinesToSlice(filePath string) []string {
	f := openFile(filePath)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	return lines
}

// sliceContains returns true if 'slice' contains 'value'
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// joinSlices adds the element of the 'new' slice
// into the 'existing' slice, only if not already there.
func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

// dumpStringsSliceToFile writes content to the file in path 'filePath' (overwriting existing content)
func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	// function deprecated as of Go 1.16 [refactoring needed]
	ioutil.WriteFile(filePath, []byte(content), 0755) 
}

// addNewSliceElementsToFile given a slice of strings representing paths stores them
// to the filesystem.
func addNewSliceElementsToFile(filePath string, newRepos []string) {
	existingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringsSliceToFile(repos, filePath)
}

// recusriveScanFolder starts the recursive search of git repositories
// living in the 'folder' subtree.
func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

// scan given a path crawls it and its subfolders,
// searching for existing Git repos
func scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repositories := recursiveScanFolder(folder)
	filePath := getDotFilePath()
	addNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\n\nSuccessfully added\n\n")
}

// scanGitFolders returns a list of subfolders of 'folder' ending with '.git'.
//
// Returns the base foler of the repo, the .git folder parent.
//
// Recursively searches in the subfolders by passing an existing 'folders' slice.
func scanGitFolders(folders []string, folder string) []string {
	// trim the last '/'
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.ReadDir(-1) // this highlights that 'Readdir' is less efficient than 'ReadDir' per docs
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()

			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}

			// explicitly avoid folders that may persist of large data
			// for this purpose these are useless
			if file.Name() == "vendor" || file.Name() = "node_modules" {
				continue
			}

			// recursive scan
			folders = scanGitFolders(folders, path)
		}
	}

	return folders
}