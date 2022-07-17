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
	"github.com/spf13/viper"
)

var (
	rootCmd = NewRootCmd()
)

// NewRootCmd initializes a new root command.
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bookmark",
		Short: "A CLI Bookmarker",
	}
}

// Execute executes the root command.
func Execute() error {
	err := initConfig()
	if err != nil {
		return err
	}
	return rootCmd.Execute()
}

// initConfig checks whether the config folder and the config file exists
// and if they does not they will be created. initConfig the loads in the
// configurations.
func initConfig() error {
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
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(configFolderPath)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

// createConfigFile creates a new config file with default values.
func createConfigFile(filename string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	config := Config{
		StorePath: filepath.Join(homeDir, ".bookmarks.json"),
	}
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0666)
}
