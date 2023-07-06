package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/korovkin/limiter"
)

func linktest(path string, cycles int) {
	// Create a temp directory
	dir, err := os.MkdirTemp(path, "linktest")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < cycles; i++ {
		orig_file := filepath.Join(dir, fmt.Sprintf("orig%d", i))
		link_file := filepath.Join(dir, fmt.Sprintf("link%d", i))

		// Create a file
		f, err := os.Create(orig_file)
		if err != nil {
			log.Fatal(err)
		}

		// write something to the original file
		f.WriteString(orig_file)
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}

		// hard link file
		err = os.Link(orig_file, link_file)
		if err != nil {
			log.Fatal(err)
		}

		// stat both files
		_, err = os.Stat(orig_file)
		if err != nil {
			log.Fatal(err)
		}
		_, err = os.Stat(link_file)
		if err != nil {
			log.Fatal(err)
		}

		// remove original file
		err = os.Remove(orig_file)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	parallel := flag.Int("p", 10, "number of parallel operations")
	cycles := flag.Int("n", 10, "number of file link cycles per operation")
	flag.Parse()

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	limiter := limiter.NewConcurrencyLimiter(*parallel)
	defer limiter.WaitAndClose()

	for p := 0; p < *parallel; p++ {
		limiter.Execute(func() {
			linktest(path, *cycles)
		})
	}
}
