package main

import "os"
import "fmt"
import "github.com/amarburg/go-quicktime"
import "flag"
import "strings"

func main() {

	el := flag.String("eagerload", "", "Atom types to eagerload")

	flag.Parse()

	atoms := strings.Split(*el, ",")
	eagerLoadFunc := func(conf *quicktime.BuildTreeConfig) {
		if len(atoms) > 0 {
			conf.EagerloadTypes = atoms
		}
	}

	// Now handle the remaining args
	if len(flag.Args()) == 0 {
		fmt.Printf("Usage:   dump_mov_tree [MOV files]\n")
		os.Exit(-1)
	}

	for _, filename := range flag.Args() {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Couldn't open file %s: %s\n", filename, err.Error())
			continue
		}

		fileinfo, _ := file.Stat()

		tree, err := quicktime.BuildTree(file, uint64(fileinfo.Size()), eagerLoadFunc)
		if err != nil {
			fmt.Printf("Error loading atom tree for %s: %s\n", filename, err.Error())
		}

		quicktime.DumpTree(tree)
	}
}
