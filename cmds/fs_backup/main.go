package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/matryer/filedb"
	"log"
	"strings"
)

const (
	list   = "list"
	add    = "add"
	remove = "remove"
)

type path struct {
	Path string
	Hash string
}

func (p path) String() string {
	return fmt.Sprintf("%s [%s]", p.Path, p.Hash)
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatalln(fatalErr)
		}
	}()

	var dbpath = flag.String("db", "./backupdata", "path to database directory")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fatalErr = errors.New("invalid usage; must specify command")
		return
	}

	db, err := filedb.Dial(*dbpath)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()
	col, err := db.C("paths")
	if err != nil {
		fatalErr = err
		return
	}

	firstNonFlagArgument := args[0]
	switch strings.ToLower(firstNonFlagArgument) {
	case list:
		fatalErr = listPaths(col)
		return
	case add:
		paths := args[1:]
		if len(paths) == 0 {
			fatalErr = errors.New("must specify path to add")
			return
		}
		fatalErr = addingPaths(col, paths)
		return
	case remove:
		paths := args[1:]
		if len(paths) == 0 {
			fatalErr = errors.New("must specify path to remove")
			return
		}
		fatalErr = removingPaths(col, paths)
		return
	}
}

func removingPaths(col *filedb.C, paths []string) error {
	var fatalErr error
	var path path
	err := col.RemoveEach(func(i int, data []byte) (removed bool, stop bool) {
		err := json.Unmarshal(data, &path)
		if err != nil {
			fatalErr = err
			return false, true
		}

		for _, p := range paths {
			if path.Path == p {
				fmt.Printf("- %s\n", path)
				return true, false
			}
		}
		return false, false
	})
	if err != nil {
		fatalErr = err
	}
	return fatalErr
}

func addingPaths(col *filedb.C, paths []string) error {
	var fatalErr error
	for _, p := range paths {
		path := &path{Path: p, Hash: "Not yet archived"}
		if err := col.InsertJSON(path); err != nil {
			fatalErr = err
			break
		}
		fmt.Printf("+ %s\n", path)
	}
	return fatalErr
}

func listPaths(col *filedb.C) error {
	var fatalErr error
	var path path
	err := col.ForEach(func(i int, data []byte) bool {
		err := json.Unmarshal(data, &path)
		if err != nil {
			fatalErr = err
			return true
		}
		fmt.Printf("= %s\n", path)
		return false
	})
	if err != nil {
		fatalErr = err
	}
	return fatalErr
}
