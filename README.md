# go-quicktime

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/3372baadfefe4b77b53eb7b44a18e7a9)](https://www.codacy.com/app/amarburg/go-quicktime?utm_source=github.com&utm_medium=referral&utm_content=amarburg/go-quicktime&utm_campaign=badger)
[![GoDoc](https://godoc.org/github.com/amarburg/go-quicktime?status.svg)](https://godoc.org/github.com/amarburg/go-quicktime)
[![wercker status](https://app.wercker.com/status/6a3ea4d08d325030c0988bd0f099e563/s/master "wercker status")](https://app.wercker.com/project/byKey/6a3ea4d08d325030c0988bd0f099e563)

quicktime is a simple Quicktime container parser written in Go.

Uses the `io.ReaderAt` interface for random access into the Quicktime file.   Originally written for use with [LazyFS](https://github.com/amarburg/go-lazyfs)


The core unit is an `Atom` which is a generic structure storing the the type, size, and location of each Atom.    Each Atom can also optionally store its data as a byte slice.   Atoms are arranged in a tree structure which reflects the structure of the file.

When the data is loaded, the Atom can be parsed into a named Atom type (e.g. `MOOVAtom`, `FTYPAtom`, etc).   These structures break out the actual fields within the Atom.   In some cases these named types do understand the Atom hierarchy.

Forked from [here](https://github.com/chunkedswarm/go-quicktime).

## TODO

- [ ] Write [GoDoc](https://blog.golang.org/godoc-documenting-go-code) and get it published on [GoDoc.org](https://godoc.org/)

## License
This library is under the [MIT License](http://opensource.org/licenses/MIT)
