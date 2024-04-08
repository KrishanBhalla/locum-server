package models

import badger "github.com/dgraph-io/badger/v4"

type DbCloser interface {
	CloseDB() error
}

func lookupByKey(db *badger.DB, key string, data []byte) ([]byte, error) {
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			// This func with val would only be called if item.Value encounters no error.
			data = append(data, val...)
			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}
	return data, nil
}

func lookupByPrefix(db *badger.DB, prefix []byte, data [][]byte) ([][]byte, error) {
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				data = append(data, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return data, nil
}

func appendToSliceWithoutDuplicates(oldElements []string, newSlice ...string) []string {
	seen := make(map[string]bool)
	for _, v := range oldElements {
		seen[v] = true
	}
	newElements := make([]string, 0, len(newSlice))
	for _, v := range newSlice {
		if _, ok := seen[v]; !ok {
			newElements = append(newElements, v)
			seen[v] = true
		}
	}
	return append(oldElements, newElements...)
}

func removeFromSlice(original []string, key string) []string {
	for i, v := range original {
		if v == key {
			if i < len(original) {
				return append(original[:i], original[i+1:]...)
			} else {
				return original[:len(original)-1]
			}
		}
	}
	return original
}
