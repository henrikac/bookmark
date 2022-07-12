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

package cmd_test

import (
	"bytes"

	"github.com/henrikac/cli-bookmark/cmd"
)

func ExampleNewVersionCmd() {
	root := cmd.NewRootCmd()
	version := cmd.NewVersionCmd()
	root.AddCommand(version)

	buff := bytes.NewBufferString("")

	root.SetOut(buff)
	root.SetArgs([]string{"version"})
	root.Execute()

	// Output:
	// v1.0.0
}
