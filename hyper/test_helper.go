package hyper

import (
	"fmt"
	"os"

	"github.com/aalda/trees/storage/badger"
)

func openBadgerStore(path string) (*badger.BadgerStore, func()) {
	store := badger.NewBadgerStore(path)
	return store, func() {
		store.Close()
		deleteFile(path)
	}
}

func deleteFile(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Unable to remove db file %s", err)
	}
}
