package filesystem_watcher

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// be responsible for archiving the source folder and storing it in the destination path.
type Archiver interface {
	DestFmt() string
	Archive(src, dest string) error
}

type zipper struct{}

func (z *zipper) DestFmt() string {
	return "%d.zip"
}

func (z *zipper) Archive(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()
	return filepath.Walk(src, zipperExecutor(w))
}

func zipperExecutor(zipWriter *zip.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil // Skip
		}

		if err != nil {
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		fileWriter, err := zipWriter.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(fileWriter, in)
		if err != nil {
			return err
		}
		return nil
	}
}

// This curious snippet of Go voodoo is actually a very interesting way of exposing the
// intent to the compiler without using any memory (literally 0 bytes).
// Zip is an Archiver that zips and unzips files.
var ZIP Archiver = (*zipper)(nil)
