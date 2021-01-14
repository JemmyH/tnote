package note

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
* @CreateTime: 2021/1/13 11:31
* @Author: hujiaming
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

func TestNote_Serialize(t *testing.T) {
	data := note.Serialize()
	t.Logf("data: %s", data)

	newNote := DeserializeNote(data)
	t.Log(newNote.String())
	assert.Equal(t, note.String(), newNote.String())
}

func TestGetHeadNote(t *testing.T) {
	head := GetHeadNote()
	t.Log(head.String())
}
