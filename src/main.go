package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// scan given a path crawls it and its subfolders,
// searching for existing Git repos
func scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repositories := recursiveScanFolder(folder)
	filePath := getDotFilePath()
	addNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\n\nSuccessfully added\n\n")
}

// stats generates a graph of your Git contributions
func stats(email string) {
	print("stats")
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

// recusriveScanFolder starts the recursive search of git repositories
// living in the 'folder' subtree.
func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

// getDotFilePath returns the dot file for the repos list.
//
// Creates it and the enclosing folder if it does not exist.

// main called function
func main() {
	var folder string
	var email string
	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repos")
	flag.StringVar(&email, "email", "your@email.com", "the email to scan")
	flag.Parse()

	if folder != "" {
		scan(folder)
		return
	}

	stats(email)
}
