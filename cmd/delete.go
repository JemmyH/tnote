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
	"github.com/jemmyh/terminal_note/note"
	"github.com/spf13/cobra"
)

var (
	idPrefix string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete notes prefixed with id",
	Run: func(cmd *cobra.Command, args []string) {
		notebook := note.GetNoteBook(GetUserName())
		userName := GetUserName()
		passwd := GetInputPassword()
		notebook.DeleteNotePrefix(note.StringToBytes(idPrefix), userName, passwd)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		CheckNotebook(GetUserName())
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&idPrefix, "id", "", "delete notes with prefix of id")
}
