package fs_backup_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zucchinidev/fs_backup"
	"testing"
)

func TestDirHash(t *testing.T) {

	hash1a, err := fs_backup.DirHash("test/hash1")
	require.NoError(t, err)
	hash1b, err := fs_backup.DirHash("test/hash1")
	require.NoError(t, err)

	require.Equal(t, hash1a, hash1b, "hash1 and hash1b should be identical")

	hash2, err := fs_backup.DirHash("test/hash2")
	require.NoError(t, err)

	require.NotEqual(t, hash1a, hash2, "hash1 and hash2 should not be the same")

}
