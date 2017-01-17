package main

import "os"
import "fmt"
import "github.com/amarburg/go-quicktime"


func main() {

  if len(os.Args) < 2 {
    fmt.Printf("Usage:   dump_mov_tree [MOV files]\n")
    os.Exit(-1)
  }

  for _,filename := range os.Args[1:] {
    file,err := os.Open( filename )
    if err != nil {
      fmt.Printf("Couldn't open file %s: %s\n", filename, err.Error() )
      continue
    }

    fileinfo,_ := file.Stat()

    tree,err := quicktime.BuildTree( file, fileinfo.Size() )
    if err != nil {
      fmt.Printf("Error loading atom tree for %s: %s\n", filename, err.Error() )
    }

    quicktime.DumpTree( tree )
  }
}
