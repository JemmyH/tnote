package note

import (
	"testing"
	"time"
)

/*
* @CreateTime: 2021/1/13 11:31
* @Author: JemmyHu <hujm20151021@gmail.com>
* @Description:
 */

var note *Note

func init() {
	note = NewNote("test content")
	note.SetID(time.Now())
}
func TestNewNote(t *testing.T) {
	note = NewNote("test content")
	t.Log(note.String())
}
