package cmd_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/henrikac/bookmark/cmd"
	"github.com/henrikac/bookmark/internal/store"
	"github.com/spf13/cobra"
)

type memoryBookmarkStore struct {
	Bookmarks map[string]string
}

func (s *memoryBookmarkStore) Load() (store.BookmarkContainer, error) {
	return s.Bookmarks, nil
}

func (s *memoryBookmarkStore) Update(store store.BookmarkContainer) error {
	return nil
}

func newMemoryBookmarkStore() *memoryBookmarkStore {
	return &memoryBookmarkStore{
		Bookmarks: store.BookmarkContainer{},
	}
}

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buff := new(bytes.Buffer)
	cmd.SetOut(buff)
	cmd.SetErr(buff)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buff.String(), err
}

func userInput(input string) *os.File {
	content := []byte(input)
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
	return tmpfile
}

// capture os.Stdout
func capture() func() (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	done := make(chan error, 1)

	save := os.Stdout
	os.Stdout = w

	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		w.Close()
		err := <-done
		return buf.String(), err
	}
}

func TestBookmarkAddCmd(t *testing.T) {
	s := newMemoryBookmarkStore()
	root := cmd.NewRootCmd()
	addCmd := cmd.BookmarkAddCmd(s)
	bookmarkName := "hello"
	bookmarkCmd := "echo \"Hello World\""
	root.AddCommand(addCmd)
	output, err := executeCommand(root, "add", bookmarkName, bookmarkCmd)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if output != "New bookmark \"hello\" has been added successfully!\n" {
		t.Errorf(
			"Expected: New bookmark \"hello\" has been added successfully!\nGot: %s",
			output,
		)
	}
	if len(s.Bookmarks) != 1 {
		t.Errorf("Expected 1 stored bookmark\nFound: %d", len(s.Bookmarks))
	}
	val, found := s.Bookmarks[bookmarkName]
	if !found {
		t.Errorf("Expected to find bookmark named \"%s\"", bookmarkName)
	}
	if val != bookmarkCmd {
		t.Errorf("Expected to find: %s\nFound: %s", bookmarkCmd, val)
	}
}

func TestBookmarkListCmdWithNoBookmarks(t *testing.T) {
	s := newMemoryBookmarkStore()
	root := cmd.NewRootCmd()
	listCmd := cmd.BookmarkListCmd(s)
	root.AddCommand(listCmd)
	output, err := executeCommand(root, "list")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	expected := "ID: BOOKMARK: COMMAND\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
}

func TestBookmarkListCmd(t *testing.T) {
	s := newMemoryBookmarkStore()
	s.Bookmarks["hello"] = "echo \"Hello world\""
	s.Bookmarks["list"] = "ls"
	root := cmd.NewRootCmd()
	listCmd := cmd.BookmarkListCmd(s)
	root.AddCommand(listCmd)
	output, err := executeCommand(root, "list")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if output == "" {
		t.Error("Expected output but got none")
	}
	expectedOutput := `ID: BOOKMARK: COMMAND
1: hello: echo "Hello world"
2: list: ls
`
	if output != expectedOutput {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestBookmarkRemoveCmdWithNoBookmarks(t *testing.T) {
	s := newMemoryBookmarkStore()
	root := cmd.NewRootCmd()
	removeCmd := cmd.BookmarkRemoveCmd(s)
	root.AddCommand(removeCmd)
	output, err := executeCommand(root, "remove", "test")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	expected := "You have no saved bookmarks\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
}

func TestBookmarkRemoveCmdUnknownCmd(t *testing.T) {
	s := newMemoryBookmarkStore()
	s.Bookmarks["test"] = "bad command"
	root := cmd.NewRootCmd()
	removeCmd := cmd.BookmarkRemoveCmd(s)
	root.AddCommand(removeCmd)
	output, err := executeCommand(root, "remove", "unknown")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	expected := "Unable to find bookmark: \"unknown\"\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
}

func TestBookmarkRemoveCmd(t *testing.T) {
	input := userInput("y")
	defer os.Remove(input.Name())
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = input

	s := newMemoryBookmarkStore()
	s.Bookmarks["test"] = "bad command"
	root := cmd.NewRootCmd()
	removeCmd := cmd.BookmarkRemoveCmd(s)
	root.AddCommand(removeCmd)
	output, err := executeCommand(root, "remove", "test")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	output = strings.TrimPrefix(
		output,
		"Are you sure you want to remove \"test\" (y/N)? ",
	)
	expected := "\"test\" was removed successfully!\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
	if len(s.Bookmarks) != 0 {
		t.Error("Failed to remove bookmark")
	}
}

func TestBookmarkExecCmdWithNoBookmarks(t *testing.T) {
	s := newMemoryBookmarkStore()
	root := cmd.NewRootCmd()
	execCmd := cmd.BookmarkExecCmd(s)
	root.AddCommand(execCmd)
	output, err := executeCommand(root, "exec", "test")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	expected := "You have no saved bookmarks\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
}

func TestBookmarkExecCmdWithUnknownCommand(t *testing.T) {
	s := newMemoryBookmarkStore()
	s.Bookmarks["hello"] = "echo \"Hello world\""
	root := cmd.NewRootCmd()
	execCmd := cmd.BookmarkExecCmd(s)
	root.AddCommand(execCmd)
	output, err := executeCommand(root, "exec", "test")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	expected := "Unable to find bookmark: \"test\"\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
}

func TestBookmarkExecCmd(t *testing.T) {
	s := newMemoryBookmarkStore()
	s.Bookmarks["hello"] = "echo \"Hello world\""
	root := cmd.NewRootCmd()
	execCmd := cmd.BookmarkExecCmd(s)
	root.AddCommand(execCmd)
	done := capture()
	_, err := executeCommand(root, "exec", "hello")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	output, err := done()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	expected := "Hello world\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
}

func TestBookmarkSearchCmdWithNoBookmarks(t *testing.T) {
	s := newMemoryBookmarkStore()
	root := cmd.NewRootCmd()
	execCmd := cmd.BookmarkSearchCmd(s)
	root.AddCommand(execCmd)
	output, err := executeCommand(root, "search", "test")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	expected := "You have no saved bookmarks\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
}

func TestBookmarkSearchCmdUnknownBookmark(t *testing.T) {
	s := newMemoryBookmarkStore()
	s.Bookmarks["hello"] = "echo \"Hello world\""
	root := cmd.NewRootCmd()
	execCmd := cmd.BookmarkSearchCmd(s)
	root.AddCommand(execCmd)
	output, err := executeCommand(root, "search", "test")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	expected := "Unable to find bookmark: \"test\"\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
}

func TestBookmarkSearchCmd(t *testing.T) {
	s := newMemoryBookmarkStore()
	s.Bookmarks["hello"] = "echo \"Hello world\""
	root := cmd.NewRootCmd()
	execCmd := cmd.BookmarkSearchCmd(s)
	root.AddCommand(execCmd)
	output, err := executeCommand(root, "search", "hello")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	expected := "echo \"Hello world\"\n"
	if output != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, output)
	}
}
