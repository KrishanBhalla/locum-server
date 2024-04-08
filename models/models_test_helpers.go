package models

import (
	"os"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/require"
)

func removeDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
}

// Opens a badger db and runs a a test on it.
func runBadgerTest(t *testing.T, test func(t *testing.T, db *badger.DB)) {
	dir, err := os.MkdirTemp("", "models-test")
	require.NoError(t, err)
	defer removeDir(dir)
	opts := new(badger.Options)
	*opts = getTestOptions(dir)

	if opts.InMemory {
		opts.Dir = ""
		opts.ValueDir = ""
	}

	db, err := badger.Open(*opts)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, db.Close())
	}()
	test(t, db)
}

func getTestOptions(dir string) badger.Options {
	opt := badger.DefaultOptions(dir).
		WithSyncWrites(false).
		WithLoggingLevel(badger.WARNING).
		WithInMemory(true)
	return opt
}
