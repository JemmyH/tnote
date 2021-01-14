package note

import (
	"log"

	"github.com/boltdb/bolt"
)

/*
* @CreateTime: 2021/1/12 20:49
* @Author: hujiaming
* @Description: a iterator for note
 */

type NoteIter struct {
	CurKey []byte
	Db     *bolt.DB
	PubKey []byte
	PrvKey []byte
	Passwd string
}

func NewNoteIter(curKey string, db *bolt.DB, userName, pwd string) *NoteIter {
	return &NoteIter{
		CurKey: []byte(curKey),
		Db:     db,
		PubKey: GetPubKey(userName),
		PrvKey: GetPrvKey(userName),
		Passwd: pwd,
	}
}

func (i *NoteIter) Next() *Note {
	var note *Note
	err := i.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetDbBucketName()))
		if b == nil {
			return nil
		}
		noteBytes := b.Get(i.CurKey)
		if noteBytes == nil {
			return nil
		}
		note = DeserializeNote(noteBytes, i.PrvKey, i.Passwd)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	if note != nil {
		i.CurKey = note.NextID
	}

	return note
}

func (i *NoteIter) Prev() *Note {
	var note *Note
	err := i.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetDbBucketName()))
		if b == nil {
			return nil
		}
		noteBytes := b.Get(i.CurKey)
		if noteBytes == nil {
			return nil
		}
		note = DeserializeNote(noteBytes, i.PrvKey, i.Passwd)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	if note != nil {
		i.CurKey = note.PrevID
	}
	return note
}
