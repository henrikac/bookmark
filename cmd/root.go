// Copyright (C) 2022 Henrik A. Christensen
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// configPath is the path to where configurations are stored.
	configPath string
	// storePath is the path to where bookmarks are stored.
	storePath string
	// store is a map of stored bookmarks.
	store   BookmarkStore
	rootCmd = NewRootCmd()
)

// NewRootCmd initializes a new root command.
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bookmark",
		Short: "A CLI Bookmarker",
	}
}

// Execute checks whether the config folder and the config file exists
// and if they does not they will be created. Execute then loads the users
// bookmarks and then execute the root command.
func Execute() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configFolderPath := filepath.Join(configDir, "bookmark")
	if _, err := os.Stat(configFolderPath); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(configFolderPath, 0777)
		if err != nil {
			return err
		}
	}
	configFilePath := filepath.Join(configDir, "bookmark", "config.json")
	if _, err = os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		err = createConfigFile(configFilePath)
		if err != nil {
			return err
		}
	}
	configPath = configFilePath
	store, err = loadBookmarkStore()
	if err != nil {
		return err
	}
	return rootCmd.Execute()
}

// createConfigFile creates a new config file with default values.
func createConfigFile(filename string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	config := Config{
		Store: filepath.Join(homeDir, ".bookmarks.json"),
	}
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0666)
}

// loadBookmarkStore loads the user's bookmarks.
func loadBookmarkStore() (BookmarkStore, error) {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	storePath = config.Store
	if _, err := os.Stat(storePath); errors.Is(err, os.ErrNotExist) {
		return BookmarkStore{}, nil
	}
	store, err := os.ReadFile(storePath)
	if err != nil {
		return nil, err
	}
	var bookmarks BookmarkStore
	err = json.Unmarshal(store, &bookmarks)
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}
