package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func main() {
	flag.Parse()
	// determine initial directories
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	now := time.Now()

	var nfiles, nbytes int64
	for _, root := range roots {
		nf, nb := walkDir(root)
		nfiles += nf
		nbytes += nb
	}

	fmt.Println("Total time taken: ", time.Since(now))
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir recursively walks the file tree rooted at dir
// and returns the size of each found file and the number of files.
func walkDir(dir string) (numFiles, size int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			nf, fs := walkDir(subdir)
			numFiles += nf
			size += fs
		} else {
			numFiles++
			fs, _ := entry.Info()
			size += fs.Size()
		}
	}
	return
}

// dirents returns the entries of directory dir.
func dirents(dir string) []fs.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}

	return entries
}
