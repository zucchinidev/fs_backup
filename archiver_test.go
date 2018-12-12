package fs_backup_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zucchinidev/fs_backup"
	"os"
	"testing"
)

func setup(t *testing.T) {
	_ = os.MkdirAll("test/output", 0777)
}

func teardown(t *testing.T) {
	_ = os.RemoveAll("/test/output")
}

func TestZipper_Archive(t *testing.T) {
	setup(t)
	defer teardown(t)
	err := fs_backup.ZIP.Archive("test/hash1", "test/output/1.zip")
	require.NoError(t, err)
}

type call struct {
	Src  string
	Dest string
}

type TestArchiver struct {
	Archives []*call
}

func (a *TestArchiver) DestFmt() string {
	return "%d.zip"
}

func (a *TestArchiver) Archive(src, dest string) error {
	a.Archives = append(a.Archives, &call{Src: src, Dest: dest})
	return nil
}

var _ fs_backup.Archiver = (*TestArchiver)(nil)
