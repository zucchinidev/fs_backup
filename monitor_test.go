package fs_backup_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zucchinidev/fs_backup"
	"strings"
	"testing"
)

var archiver *TestArchiver

func TestMonitor(t *testing.T) {
	archiver = &TestArchiver{}

	monitor := &fs_backup.Monitor{
		Destination: "test/archive",
		Paths: map[string]string{
			"test/hash1": "abc",
			"test/hash2": "def",
		},
		Archiver: archiver,
	}

	numOfCompressedFiles, err := monitor.Now()
	require.NoError(t, err)
	require.Equal(t, 2, numOfCompressedFiles)
	require.Equal(t, 2, len(archiver.Archives))
	for _, call := range archiver.Archives {
		require.True(t, strings.HasPrefix(call.Dest, monitor.Destination))
		require.True(t, strings.HasSuffix(call.Dest, ".zip"))
	}
}
