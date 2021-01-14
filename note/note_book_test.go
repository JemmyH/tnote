package note

import (
	"testing"

	"github.com/boltdb/bolt"
)

/*
* @CreateTime: 2021/1/13 11:41
* @Author: hujiaming
* @Description:
 */

var (
	testNoteBook *NoteBook
	testUserName        = "test"
	testPwd             = "123456"
	testPubKey   []byte = GetPubKey(testUserName)
	testPrvKey   []byte = GetPubKey(testUserName)
)

func TestGetNoteBook(t *testing.T) {
	b := GetNoteBook(testUserName)
	_printNoteBook(t, b)
}

func TestNoteBook_AddNote(t *testing.T) {
	b := GetNoteBook(testUserName)
	// newNote := NewNote("this is a test content")
	// b.AddNote(newNote)
	_printNoteBook(t, b)
}

func _printNoteBook(t *testing.T, b *NoteBook) {
	err := b.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetDbBucketName()))

		head := DeserializeNote(b.Get([]byte(GetDbHeadKey())), testPrvKey, testPwd)
		tmp := head
		for {
			t.Log(tmp.String())
			if tmp.NextID == nil {
				break
			}
			tmp = DeserializeNote(b.Get(tmp.NextID), testPrvKey, testPwd)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoteBook_DeleteNotePrefix(t *testing.T) {
	b := GetNoteBook(testUserName)
	b.DeleteNotePrefix(StringToBytes(""), testUserName, testPwd)
	_printNoteBook(t, b)
}
