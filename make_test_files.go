// Build a directory structure with the required number of files in
//
// Run with make_test_files.go [flag] <directory>
package main

import (
	cryptrand "crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var (
	// Flags
	numberOfFiles            = flag.Int("n", 1000, "Number of files to create")
	zeroes                   = flag.Bool("z", false, "Fill files with zeroes instead of random data")
	verbose                  = flag.Bool("v", false, "Be more verbose")
	doSync                   = flag.Bool("sync", false, "Fsync each file")
	loop                     = flag.Bool("loop", false, "Loop forever")
	averageFilesPerDirectory = flag.Int("files-per-directory", 10, "Average number of files per directory")
	maxDepth                 = flag.Int("max-depth", 10, "Maximum depth of directory heirachy")
	minFileSize              = flag.Int64("min-size", 0, "Minimum size of file to create")
	maxFileSize              = flag.Int64("max-size", 100, "Maximum size of files to create")
	minFileNameLength        = flag.Int("min-name-length", 4, "Minimum size of file to create")
	maxFileNameLength        = flag.Int("max-name-length", 12, "Maximum size of files to create")
	defaultSeed              = flag.Int64("seed", 1, "Seed for the random number generator")

	directoriesToCreate int
	totalDirectories    int
	totalSize           int64
	fileNames           = map[string]struct{}{} // keep a note of which file name we've used already
)

// randomString create a random string for test purposes
func randomString(n int) string {
	const (
		vowel     = "aeiou"
		consonant = "bcdfghjklmnpqrstvwxyz"
		digit     = "0123456789"
	)
	pattern := []string{consonant, vowel, consonant, vowel, consonant, vowel, consonant, digit}
	out := make([]byte, n)
	p := 0
	for i := range out {
		source := pattern[p]
		p = (p + 1) % len(pattern)
		out[i] = source[rand.Intn(len(source))]
	}
	return string(out)
}

// fileName creates a unique random file or directory name
func fileName() (name string) {
	for {
		length := rand.Intn(*maxFileNameLength-*minFileNameLength) + *minFileNameLength
		name = randomString(length)
		if _, found := fileNames[name]; !found {
			break
		}
	}
	fileNames[name] = struct{}{}
	return name
}

// dir is a directory in the directory heirachy being built up
type dir struct {
	name     string
	depth    int
	children []*dir
	parent   *dir
}

// Create a random directory heirachy under d
func (d *dir) createDirectories() {
	for totalDirectories < directoriesToCreate {
		newDir := &dir{
			name:   fileName(),
			depth:  d.depth + 1,
			parent: d,
		}
		d.children = append(d.children, newDir)
		totalDirectories++
		switch rand.Intn(4) {
		case 0:
			if d.depth < *maxDepth {
				newDir.createDirectories()
			}
		case 1:
			return
		}
	}
	return
}

// list the directory heirachy
func (d *dir) list(path string, output []string) []string {
	dirPath := filepath.Join(path, d.name)
	output = append(output, dirPath)
	for _, subDir := range d.children {
		output = subDir.list(dirPath, output)
	}
	return output
}

type zeroReader struct{}

func (z zeroReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
		n++
	}
	return n, nil
}

// writeFile writes a random file at dir/name
func writeFile(dir, name string) {
	size := rand.Int63n(*maxFileSize-*minFileSize) + *minFileSize
	totalSize += size
	if *verbose {
		log.Printf("Making %s/%s size %d", dir, name, size)
	}
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatalf("Failed to make directory %q: %v", dir, err)
	}
	path := filepath.Join(dir, name)
	fd, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to open file %q: %v", path, err)
	}
	var in io.Reader = cryptrand.Reader
	if *zeroes {
		in = zeroReader{}
	}
	_, err = io.CopyN(fd, in, size)
	if err != nil {
		log.Fatalf("Failed to write %v bytes to file %q: %v", size, path, err)
	}
	if *doSync {
		err = fd.Sync()
		if err != nil {
			log.Fatalf("Failed to sync file %q: %v", path, err)
		}
	}
	err = fd.Close()
	if err != nil {
		log.Fatalf("Failed to close file %q: %v", path, err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [flags] <directory>

This command makes a random directory structure with random files in
<directory>.  The options can be used to control exactly which files
get made.

The file names and sizes will be identical each time the command is
run with the same parameters.  -seed can be used to change what is
created.

Options:
`, os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("Require 1 directory argument")
	}
	outputDirectory := args[0]

	for {
		// create identical files each time round the loop
		rand.Seed(*defaultSeed)
		fileNames = map[string]struct{}{}
		totalSize = 0

		log.Printf("Storing files in %v", outputDirectory)
		start := time.Now()

		directoriesToCreate = *numberOfFiles / *averageFilesPerDirectory
		if *verbose {
			log.Printf("Creating %d directories", directoriesToCreate)
		}
		totalDirectories = 0

		root := &dir{name: outputDirectory, depth: 1}
		for totalDirectories < directoriesToCreate {
			root.createDirectories()
		}
		dirs := root.list("", []string{})
		for i := 0; i < *numberOfFiles; i++ {
			dir := dirs[rand.Intn(len(dirs))]
			writeFile(dir, fileName())
		}

		end := time.Now()
		dt := end.Sub(start)
		log.Printf("That took %v to write %d files, total size %d, @ %.3fMiBytes/s", dt, *numberOfFiles, totalSize, float64(totalSize)/float64(dt)*float64(time.Second)/1024/1024)

		if !*loop {
			break
		}
	}
}
