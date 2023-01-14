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
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/henrikac/bookmark/internal/store"
	"github.com/spf13/cobra"
)

// BookmarkStore is a collection of the user's bookmarks.
// type BookmarkStore = map[string]string

var (
	bookmarkStore     = store.NewBookmarkFileStore()
	bookmarkAddCmd    = BookmarkAddCmd(bookmarkStore)
	bookmarkExecCmd   = BookmarkExecCmd(bookmarkStore)
	bookmarkListCmd   = BookmarkListCmd(bookmarkStore)
	bookmarkRemoveCmd = BookmarkRemoveCmd(bookmarkStore)
	bookmarkSearchCmd = BookmarkSearchCmd(bookmarkStore)
)

// BookmarkAddCmd initializes a new add command.
func BookmarkAddCmd(bs store.BookmarkStoreLoadUpdater) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new bookmark",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			bookmarkCmd := strings.Join(args[1:], " ")
			bookmarks, err := bs.Load()
			if err != nil {
				return err
			}
			if val, found := bookmarks[name]; found {
				cmd.Printf("%s already exists: %s\n", name, val)
				var input string
				cmd.Printf("Do you want to override it (y/N)? ")
				_, _ = fmt.Scanln(&input)
				if strings.ToLower(strings.TrimSpace(input)) == "y" {
					bookmarks[name] = bookmarkCmd
					err := bs.Update(bookmarks)
					if err != nil {
						return err
					}
					cmd.Printf("Bookmark \"%s\" has been updated successfully!\n", name)
				}
				return nil
			}
			bookmarks[name] = bookmarkCmd
			err = bs.Update(bookmarks)
			if err != nil {
				return err
			}
			cmd.Printf("New bookmark \"%s\" has been added successfully!\n", name)
			return nil
		},
	}
}

// BookmarkExecCmd initializes a new exec command.
func BookmarkExecCmd(bs store.BookmarkStoreLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "exec",
		Short: "Execute a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bookmarks, err := bs.Load()
			if err != nil {
				return err
			}
			if len(bookmarks) == 0 {
				cmd.Println("You have no saved bookmarks")
				return nil
			}
			name := args[0]
			if _, found := bookmarks[name]; !found {
				cmd.Printf("Unable to find bookmark: \"%s\"\n", name)
				return nil
			}
			bookmarkCmdArr := splitOnSpace(bookmarks[name])
			var command *exec.Cmd
			cmdAndArgs := bookmarkCmdArr[0]
			if len(bookmarkCmdArr) > 1 {
				for _, part := range bookmarkCmdArr[1:] {
					cmdAndArgs += fmt.Sprintf(" %s", part)
				}
			}
			if runtime.GOOS == "windows" {
				command = exec.Command("cmd", "/c", cmdAndArgs)
			} else {
				command = exec.Command("bash", "-c", cmdAndArgs)
			}
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			return command.Run()
		},
	}
}

// BookmarkListCmd initializes a new list command.
func BookmarkListCmd(bs store.BookmarkStoreLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your current saved bookmarks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			bookmarks, err := bs.Load()
			if err != nil {
				return err
			}
			cmd.Println("ID: BOOKMARK: COMMAND")
			counter := 1
			for bm, c := range bookmarks {
				cmd.Printf("%d: %s: %s\n", counter, bm, c)
				counter += 1
			}
			return nil
		},
	}
}

// BookmarkRemoveCmd initializes a new remove command.
func BookmarkRemoveCmd(bs store.BookmarkStoreLoadUpdater) *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bookmarks, err := bs.Load()
			if err != nil {
				return err
			}
			if len(bookmarks) == 0 {
				cmd.Println("You have no saved bookmarks")
				return nil
			}
			name := args[0]
			if _, found := bookmarks[name]; !found {
				cmd.Printf("Unable to find bookmark: \"%s\"\n", name)
				return nil
			}
			var input string
			cmd.Printf("Are you sure you want to remove \"%s\" (y/N)? ", name)
			_, _ = fmt.Scanln(&input)
			if strings.ToLower(strings.TrimSpace(input)) == "y" {
				delete(bookmarks, name)
				err := bs.Update(bookmarks)
				if err != nil {
					return err
				}
				cmd.Printf("\"%s\" was removed successfully!\n", name)
			}
			return nil
		},
	}
}

// BookmarkSearchCmd initializes a new search command.
func BookmarkSearchCmd(bs store.BookmarkStoreLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "seach for a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bookmarks, err := bs.Load()
			if err != nil {
				return err
			}
			if len(bookmarks) == 0 {
				cmd.Println("You have no saved bookmarks")
				return nil
			}
			if val, found := bookmarks[args[0]]; found {
				cmd.Println(val)
				return nil
			}
			cmd.Printf("Unable to find bookmark: \"%s\"\n", args[0])
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
