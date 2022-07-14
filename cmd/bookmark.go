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
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type BookmarkStore = map[string]string

var (
	bookmarkAddCmd    = NewBookmarkAddCmd()
	bookmarkExecCmd   = NewBookmarkExecCmd()
	bookmarkListCmd   = NewBookmarkListCmd()
	bookmarkRemoveCmd = NewBookmarkRemoveCmd()
	bookmarkSearchCmd = NewBookmarkSearchCmd()
)

func NewBookmarkAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new bookmark",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			bookmark := args[0]
			bookmarkCmd := strings.Join(args[1:], " ")
			if val, found := store[bookmark]; found {
				fmt.Printf("%s already exists: %s\n", bookmark, val)
				var input string
				fmt.Printf("Do you want to override it (y/N)? ")
				_, _ = fmt.Scanln(&input)
				if strings.ToLower(input) == "y" {
					store[bookmark] = bookmarkCmd
					err := updateStore()
					if err != nil {
						return err
					}
					fmt.Printf("Bookmark \"%s\" has been updated successfully!\n", bookmark)
				}
				return nil
			}
			store[bookmark] = bookmarkCmd
			err := updateStore()
			if err != nil {
				return err
			}
			fmt.Printf("New bookmark \"%s\" has been added successfully!\n", bookmark)
			return nil
		},
	}
}

func NewBookmarkExecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exec",
		Short: "Execute a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(store) == 0 {
				fmt.Println("You have no saved bookmarks")
				return nil
			}
			bookmark := args[0]
			if _, found := store[bookmark]; !found {
				fmt.Printf("Unable to find bookmark: \"%s\"\n", bookmark)
				return nil
			}
			bookmarkCmdArr := splitOnSpace(store[bookmark])
			fmt.Println(bookmarkCmdArr)
			fmt.Println(len(bookmarkCmdArr))
			for i, item := range bookmarkCmdArr {
				fmt.Printf("%d - %s\n", i, item)
			}
			var command *exec.Cmd
			if len(bookmarkCmdArr) == 1 {
				command = exec.Command(bookmarkCmdArr[0])
			} else {
				command = exec.Command(bookmarkCmdArr[0], bookmarkCmdArr[1:]...)
			}
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			return command.Run()
		},
	}
}

func NewBookmarkListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your current saved bookmarks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(store) == 0 {
				fmt.Println("You have no saved bookmarks")
				return nil
			}
			fmt.Println("ID: BOOKMARK: COMMAND")
			counter := 1
			for bm, cmd := range store {
				fmt.Printf("%d: %s: %s\n", counter, bm, cmd)
				counter += 1
			}
			return nil
		},
	}
}

func NewBookmarkRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(store) == 0 {
				fmt.Println("You have no saved bookmarks")
				return nil
			}
			bookmark := args[0]
			if _, found := store[bookmark]; !found {
				fmt.Printf("Unable to find bookmark: %s\n", bookmark)
				return nil
			}
			var input string
			fmt.Printf("Are you sure you want to remove \"%s\" (y/N)? ", bookmark)
			_, _ = fmt.Scanln(&input)
			if strings.ToLower(input) == "y" {
				delete(store, bookmark)
				err := updateStore()
				if err != nil {
					return err
				}
				fmt.Printf("\"%s\" was removed successfully!\n", bookmark)
			}
			return nil
		},
	}
}

func NewBookmarkSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "seach for a bookmark",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(store) == 0 {
				fmt.Println("You have no saved bookmarks")
				return
			}
			if val, found := store[args[0]]; found {
				fmt.Println(val)
				return
			}
			fmt.Printf("Unable to find bookmark: %s\n", args[0])
		},
	}
}

func init() {
	rootCmd.AddCommand(bookmarkAddCmd)
	rootCmd.AddCommand(bookmarkExecCmd)
	rootCmd.AddCommand(bookmarkListCmd)
	rootCmd.AddCommand(bookmarkRemoveCmd)
	rootCmd.AddCommand(bookmarkSearchCmd)
}

func updateStore() error {
	b, err := json.Marshal(store)
	if err != nil {
		return err
	}
	return os.WriteFile(storePath, b, 0666)
}

func splitOnSpace(s string) []string {
	res := []string{}
	var beg int
	var inString bool
	var quote byte

	for i := 0; i < len(s); i++ {
		if s[i] == ' ' && !inString {
			res = append(res, s[beg:i])
			beg = i + 1
		} else if s[i] == '"' || s[i] == '\'' {
			if !inString {
				quote = s[i]
				inString = true
			} else if (i > 0 && s[i-1] != '\\') && s[i] == quote {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}
