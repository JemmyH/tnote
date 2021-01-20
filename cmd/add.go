/*
Copyright © 2021 JemmyHu <hujm20151021@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/jemmyh/tnote/note"
	"github.com/spf13/cobra"
)

var content string

func init() {
	addCmd.Flags().StringVarP(&content, "content", "c", "", "content of a note")

	rootCmd.AddCommand(addCmd)
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a note to your notebook",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		notebook := note.GetNoteBook(GetUserName())
		newNote := note.NewNote(content)
		userName := GetUserName()
		passwd := GetInputPassword()
		notebook.AddNote(newNote, userName, passwd)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		CheckNotebook(GetUserName())
	},
}
