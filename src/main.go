package main

import (
	"flag"
)

// scan given a path crawls it and its subfolders,
// searching for existing Git repos
func scan(path string) {
	print("scan")
}

// stats generates a graph of your Git contributions
func stats(email string) {
	print("stats")
}

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
