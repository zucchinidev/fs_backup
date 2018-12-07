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
		fatalErr = err
		return
	case add:
		if len(args[1:]) == 0 {
			fatalErr = errors.New("must specify path to add")
			return
		}
		for _, p := range args[1:] {
			path := &path{Path: p, Hash: "Not yet archived"}
			if err := col.InsertJSON(path); err != nil {
				fatalErr = err
				return
			}
			fmt.Printf("+ %s\n", path)
		}
	case remove:

	}
}
