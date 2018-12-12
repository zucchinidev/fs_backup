package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/matryer/filedb"
	"github.com/zucchinidev/fs_backup"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type path struct {
	Path string
	Hash string
}

func main() {

	var fatalErr error
	defer func() {
		if fatalErr != nil {
			log.Fatalln(fatalErr)
		}
	}()

	var (
		interval = flag.Duration("interval", 10*time.Second, "Define the interval in between checks")
		archive  = flag.String("archive", "archive", "path to archive location")
		dbPath   = flag.String("dbPath", "./dbPath", "path to filedb database")
	)
	flag.Parse()
	monitor := fs_backup.Monitor{
		Archiver:    fs_backup.ZIP,
		Paths:       make(map[string]string),
		Destination: *archive,
	}
	db, err := filedb.Dial(*dbPath)
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

	var path path
	err = col.ForEach(func(_ int, data []byte) bool {
		if err := json.Unmarshal(data, &path); err != nil {
			fatalErr = err
			return true
		}
		log.Println("path", path)
		monitor.Paths[path.Path] = path.Hash
		return false // carry on
	})

	if err != nil {
		fatalErr = err
		return
	}
	if fatalErr != nil {
		return
	}
	if len(monitor.Paths) < 1 {
		fatalErr = errors.New("no paths - use fs_backup tool to add at least one")
	}

	check(monitor, col)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-time.After(*interval):
			check(monitor, col)
		case <-signalChan:
			fmt.Println()
			fmt.Println("Stopping...")
			return
		}
	}
}

func check(monitor fs_backup.Monitor, col *filedb.C) {
	fmt.Println("Checking...")
	counter, err := monitor.Now()
	if err != nil {
		log.Fatalln("failed to backup:", err)
	}

	if counter > 0 {
		fmt.Printf("    Archived %d directories\n", counter)
		// updated hashes
		var path path
		err := col.SelectEach(func(_ int, data []byte) (include bool, returnedData []byte, stop bool) {
			if err := json.Unmarshal(data, &path); err != nil {
				log.Println("failed to unmarshal data (skipping):", err)
				return true, data, false
			}
			path.Hash, _ = monitor.Paths[path.Path]
			newData, err := json.Marshal(&path)
			if err != nil {
				log.Println("failed to marshal data (skipping):", err)
				return true, data, false
			}
			return true, newData, false
		})
		if err != nil {
			log.Println("failed to select each paths:", err)
		}
	} else {
		log.Println("    No changes")
	}
}
