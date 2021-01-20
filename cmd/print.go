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
	"bytes"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/jemmyh/tnote/note"
	"github.com/spf13/cobra"
)

var (
	// tn print -o Jemmy 			// 从过去到现在 打印 Jemmy 的所有 notes
	// tn print -o Jemmy -n 10      // 从过去到现在 打印前 10 条
	// tn print -o Jemmy -r -n 10   // 从现在到过去 打印前 10 条
	// tn print -o Jemmy -t range --from 2020-01-06 --to 2020-02-18		 // 从过去到现在 打印2020-01-06~2020-02-18 的所有notes
	// tn print -o Jemmy -t range --from 2020-01-06 --to 2020-02-18 -r	 // 从现在到过去 打印2020-01-06~2020-02-18 的所有notes
	myP *printStruct
)

func init() {
	rootCmd.AddCommand(printCmd)

	myP = &printStruct{}

	// Here you will define your flags and configuration settings.
	printCmd.Flags().StringVarP(&myP.printType, "type", "t", "normal", "print type. Default is `normal`. For range printing, use `range`")
	printCmd.Flags().IntVarP(&myP.number, "number", "n", 0, "number of notes to print")
	printCmd.Flags().BoolVarP(&myP.reverse, "reverse", "r", false, "print order. If `true`, print from past to now.")
	printCmd.Flags().StringVar(&myP.rangeFrom, "from", "", "start time for type=range. Format: 2020-01-02")
	printCmd.Flags().StringVar(&myP.rangeTo, "to", "", "end time for type=range. Format: 2020-10-10")
	printCmd.Flags().BoolVarP(&myP.verbose, "verbose", "v", false, "If true, print note's detail")
}

type printStruct struct {
	printType string
	number    int
	reverse   bool
	rangeFrom string
	rangeTo   string
	verbose   bool
}

func (p *printStruct) getNoteIter(db *bolt.DB, userName, passwd string) *note.NoteIter {
	if p.reverse {
		return note.NewNoteIter(note.GetDbTailKey(), db, userName, passwd)
	}
	return note.NewNoteIter(note.GetDbHeadKey(), db, userName, passwd)
}

func (p *printStruct) getNote(iter *note.NoteIter) *note.Note {
	if p.reverse {
		return iter.Prev()
	}
	return iter.Next()
}

func (p *printStruct) checkContinue(id []byte) bool {
	if p.reverse {
		if bytes.Equal(id, []byte(note.GetDbTailKey())) {
			return true
		}
	} else {
		if bytes.Equal(id, []byte(note.GetDbHeadKey())) {
			return true
		}
	}
	return false
}

func (p *printStruct) checkBreak(id []byte) bool {
	if p.reverse {
		if bytes.Equal(id, []byte(note.GetDbHeadKey())) {
			return true
		}
	} else {
		if bytes.Equal(id, []byte(note.GetDbTailKey())) {
			return true
		}
	}
	return false
}

func (p *printStruct) printNote(n *note.Note) {
	if p.verbose {
		fmt.Println(n.String())
		// TODO: use arrow 
		// if p.reverse && n.PrevID != nil {
		// 	fmt.Println("⇧")
		// 	fmt.Println("|")
		// }else if !p.reverse && n.NextID != nil{
		// 	fmt.Println("|")
		// 	fmt.Println("↓")
		// }
	} else {
		fmt.Println(n.SimpleString())
	}
}

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print notes in your notebook",
	PreRun: func(_ *cobra.Command, _ []string) {
		CheckNotebook(GetUserName())
	},
	Run: func(_ *cobra.Command, _ []string) {
		userName := GetUserName()
		passwd := GetInputPassword()
		switch myP.printType {
		case "normal":
			noteBook := note.GetNoteBook(GetUserName())
			iter := myP.getNoteIter(noteBook.DB, userName, passwd)
			number := myP.number
			for {
				note := myP.getNote(iter)
				if myP.checkContinue(note.ID) {
					continue
				}
				if note == nil || myP.checkBreak(note.ID) {
					break
				}
				myP.printNote(note)
				if myP.number > 0 {
					number--
					if number <= 0 {
						break
					}
				}
			}
		case "range":
			// TODO: range print
		default:
			fmt.Println("invalid type: " + myP.printType)
			return
		}
	},
}
