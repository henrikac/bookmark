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
	rootCmd    = NewRootCmd()
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bookmark",
		Short: "A CLI Bookmarker",
	}
}

func Execute() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath = filepath.Join(configDir, "bookmark")
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(configPath, 0777)
		if err != nil {
			return err
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		config := Config{
			Store: filepath.Join(homeDir, "bookmarks.json"),
		}
		b, err := json.Marshal(config)
		if err != nil {
			return err
		}
		configPath = filepath.Join(configPath, "config.json")
		err = os.WriteFile(configPath, b, 0666)
		if err != nil {
			return err
		}
	} else {
		configPath = filepath.Join(configPath, "config.json")
	}
	return rootCmd.Execute()
}
