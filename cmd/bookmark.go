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
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// BookmarkStore is a collection of the user's bookmarks.
type BookmarkStore = map[string]string

var (
	bookmarkAddCmd    = NewBookmarkAddCmd()
	bookmarkExecCmd   = NewBookmarkExecCmd()
	bookmarkListCmd   = NewBookmarkListCmd()
	bookmarkRemoveCmd = NewBookmarkRemoveCmd()
	bookmarkSearchCmd = NewBookmarkSearchCmd()
)

// NewBookmarkAddCmd initializes a new add command.
func NewBookmarkAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new bookmark",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			bookmark := args[0]
			bookmarkCmd := strings.Join(args[1:], " ")
			store, err := loadBookmarkStore()
			if err != nil {
				return err
			}
			if val, found := store[bookmark]; found {
				fmt.Printf("%s already exists: %s\n", bookmark, val)
				var input string
				fmt.Printf("Do you want to override it (y/N)? ")
				_, _ = fmt.Scanln(&input)
				if strings.ToLower(input) == "y" {
					store[bookmark] = bookmarkCmd
					err := updateStore(store)
					if err != nil {
						return err
					}
					fmt.Printf("Bookmark \"%s\" has been updated successfully!\n", bookmark)
				}
				return nil
			}
			store[bookmark] = bookmarkCmd
			err = updateStore(store)
			if err != nil {
				return err
			}
			fmt.Printf("New bookmark \"%s\" has been added successfully!\n", bookmark)
			return nil
		},
	}
}

// NewBookmarkExecCmd initializes a new exec command.
func NewBookmarkExecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exec",
		Short: "Execute a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := loadBookmarkStore()
			if err != nil {
				return err
			}
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

// NewBookmarkListCmd initializes a new list command.
func NewBookmarkListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your current saved bookmarks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := loadBookmarkStore()
			if err != nil {
				return err
			}
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

// NewBookmarkRemoveCmd initializes a new remove command.
func NewBookmarkRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := loadBookmarkStore()
			if err != nil {
				return err
			}
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
				err := updateStore(store)
				if err != nil {
					return err
				}
				fmt.Printf("\"%s\" was removed successfully!\n", bookmark)
			}
			return nil
		},
	}
}

// NewBookmarkSearchCmd initializes a new search command.
func NewBookmarkSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "seach for a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := loadBookmarkStore()
			if err != nil {
				return err
			}
			if len(store) == 0 {
				fmt.Println("You have no saved bookmarks")
				return nil
			}
			if val, found := store[args[0]]; found {
				fmt.Println(val)
				return nil
			}
			fmt.Printf("Unable to find bookmark: %s\n", args[0])
			return nil
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

// loadBookmarkStore loads the user's bookmarks.
func loadBookmarkStore() (BookmarkStore, error) {
	storePath := viper.GetViper().GetString("storePath")
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

func updateStore(store BookmarkStore) error {
	storePath := viper.GetViper().GetString("storePath")
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
