package store

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/spf13/viper"
)

// BookmarkStoreLoader is the interface that wraps the Load method.
type BookmarkStoreLoader interface {
	Load() (bookmarkContainer, error)
}

// BookmarkStoreUpdater is the interface that wraps the Update method.
type BookmarkStoreUpdater interface {
	Update(bookmarkContainer) error
}

// BookmarkStoreLoadUpdater is the interface that wraps the Load
// and Update methods.
type BookmarkStoreLoadUpdater interface {
	BookmarkStoreLoader
	BookmarkStoreUpdater
}

// bookmarkContainer
type bookmarkContainer = map[string]string

// BookmarkFileStore
type BookmarkFileStore struct{}

// Load implements the BookmarkStoreLoader interface.
// It loads the user's bookmarks from a json file.
func (s BookmarkFileStore) Load() (bookmarkContainer, error) {
	storePath := viper.GetViper().GetString("storePath")
	if _, err := os.Stat(storePath); errors.Is(err, os.ErrNotExist) {
		return bookmarkContainer{}, nil
	}
	store, err := os.ReadFile(storePath)
	if err != nil {
		return nil, err
	}
	var bc bookmarkContainer
	err = json.Unmarshal(store, &bc)
	if err != nil {
		return nil, err
	}
	return bc, nil
}

// Update implements the BookmarkStoreUpdater interface.
// It writes the user's bookmarks to a json file.
func (s BookmarkFileStore) Update(store bookmarkContainer) error {
	storePath := viper.GetViper().GetString("storePath")
	b, err := json.Marshal(store)
	if err != nil {
		return err
	}
	return os.WriteFile(storePath, b, 0666)
}

// NewBookmarkFileStore initializes a new FileStore.
func NewBookmarkFileStore() *BookmarkFileStore {
	return &BookmarkFileStore{}
}
