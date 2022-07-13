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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// A Config describes a user's bookmark configurations.
type Config struct {
	// Store specifies the path to where the user's bookmarks are stored.
	Store string `json:"store"`
}

var (
	configCmd     = NewConfigCmd()
	configListCmd = NewConfigListCmd()
)

func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Handle bookmark configurations",
	}
}

func NewConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all your bookmark configurations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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

func init() {
	configCmd.AddCommand(configListCmd)
	rootCmd.AddCommand(configCmd)
}
