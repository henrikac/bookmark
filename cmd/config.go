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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// A Config describes a user's bookmark configurations.
type Config struct {
	// StorePath specifies the path to where the user's bookmarks are stored.
	StorePath string `json:"storePath"`
}

var (
	configCmd     = NewConfigCmd()
	configListCmd = NewConfigListCmd()
	configSetCmd  = NewConfigSetCmd()
)

// NewConfigCmd initializes a new config command.
func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Handle bookmark configurations",
	}
}

// NewConfigListCmd initializes a new config list command.
func NewConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all your bookmark configurations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := viper.GetViper().ConfigFileUsed()
			data, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}
			dst := &bytes.Buffer{}
			err = json.Indent(dst, data, "", "    ")
			if err != nil {
				return err
			}
			fmt.Println(dst.String())
			return nil
		},
	}
}

func NewConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Sets configuration <config> to the given <value>",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := viper.GetViper().GetString(args[0])
			if config == "" {
				return fmt.Errorf("unable to find the given config: \"%s\"", args[0])
			}
			switch args[0] {
			case "storePath":
				newStorePath := args[1]
				if !filepath.IsAbs(newStorePath) {
					absPath, err := filepath.Abs(newStorePath)
					if err != nil {
						return err
					}
					newStorePath = absPath
				}
				dir, file := filepath.Split(newStorePath)
				if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
					err = os.MkdirAll(dir, 0750)
					if err != nil {
						return err
					}
				}
				if file == "" {
					file = "bookmarks.json"
				}
				if _, err := os.Stat(config); errors.Is(err, os.ErrNotExist) {
					b, err := json.Marshal(make(map[string]string))
					if err != nil {
						return err
					}
					err = os.WriteFile(config, b, 0666)
					if err != nil {
						return err
					}
				}
				err := os.Rename(config, filepath.Join(dir, file))
				if err != nil {
					return err
				}
				viper.GetViper().Set(args[0], newStorePath)
				return viper.GetViper().WriteConfig()
			default:
				return fmt.Errorf("unable to find the given config: \"%s\"", args[0])
			}
		},
	}
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
