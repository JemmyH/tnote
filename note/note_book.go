package note

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	user2 "os/user"
	"time"

	"github.com/boltdb/bolt"
)

/*
* @CreateTime: 2021/1/12 18:30
* @Author: Jemmy@hujm20150121@gmail.com
* @Description: notebook
 */

// NoteBook represents a diary book that contains many diary notes.
type NoteBook struct {
	Head []byte
	Tail []byte
	DB   *bolt.DB
}

// CreateNewNoteBook creates a NoteBook for userName.
func CreateNewNoteBook(userName, passwd string) *NoteBook {
	dbFile := GetDbFilePath(userName)
	fmt.Println(dbFile)
	if isFileExists(dbFile) {
		fmt.Println("notebook already exists")
		os.Exit(1)
	}

	// create rsa key for user
	GenKeyForUserWithPasswd(userName, passwd)

	var headBytes, tailBytes []byte
	headNote := GetHeadNote()
	tailNote := GetTailNote()
	headNote.NextID = tailNote.ID
	tailNote.PrevID = headNote.ID

	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		var (
			b      *bolt.Bucket
			err    error
			pubKey []byte = GetPubKey(userName)
		)
		b = tx.Bucket([]byte(GetDbBucketName()))
		if b == nil {
			b, err = tx.CreateBucket([]byte(GetDbBucketName()))
			if err != nil {
				return err
			}
		}
		// add head to db
		err = b.Put([]byte(GetDbHeadKey()), headNote.Serialize(pubKey))
		if err != nil {
			return err
		}
		// add tail to db
		err = b.Put([]byte(GetDbTailKey()), tailNote.Serialize(pubKey))
		if err != nil {
			return err
		}
		headBytes = []byte(GetDbHeadKey())
		tailBytes = []byte(GetDbTailKey())

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &NoteBook{
		Head: headBytes,
		Tail: tailBytes,
		DB:   db,
	}
}

// GetNoteBook return userName's NoteBook
func GetNoteBook(userName string) *NoteBook {
	dbFile := GetDbFilePath(userName)
	if isFileExists(dbFile) == false {
		return nil
	}
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.Panic(err)
	}

	return &NoteBook{
		Head: []byte(GetDbHeadKey()),
		Tail: []byte(GetDbTailKey()),
		DB:   db,
	}
}

// AddNote add a Note to NoteBook.
func (b *NoteBook) AddNote(note *Note, userName string, pwd string) {
	if note.ID == nil {
		note.SetID(time.Now())
	}
	err := b.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GetDbBucketName()))

		// if exists, just return
		noteInDb := b.Get(note.ID)
		if noteInDb != nil {
			return nil
		}
		pubKey := GetPubKey(userName)
		prvKey := GetPrvKey(userName)

		// get tail note
		tail := DeserializeNote(b.Get([]byte(GetDbTailKey())), prvKey, pwd)
		// get tail's prev note
		tailPrev := DeserializeNote(b.Get(tail.PrevID), prvKey, pwd)

		// insert note between tailPrev and tail
		tailPrev.NextID = note.ID
		note.PrevID = tailPrev.ID
		note.NextID = tail.ID
		tail.PrevID = note.ID

		// update tail's timestamp to now
		tail.Timestamp = time.Now().UnixNano()

		// save tailPrev
		if err := b.Put(tailPrev.ID, tailPrev.Serialize(pubKey)); err != nil {
			return err
		}
		// save current note
		if err := b.Put(note.ID, note.Serialize(pubKey)); err != nil {
			return err
		}
		// save tail
		if err := b.Put(tail.ID, tail.Serialize(pubKey)); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("save note ok with id: %s\n", note.ID)
}

// DeleteNotePrefix delete notes that match prefix of `id`.
func (b *NoteBook) DeleteNotePrefix(id []byte, userName, pwd string) {
	// you can not delete `head` or `tail` key
	if len(id) == 0 || bytes.Equal(id, []byte(GetDbHeadKey())) || bytes.Equal(id, []byte(GetDbTailKey())) {
		return
	}
	err := b.DB.Update(func(tx *bolt.Tx) error {
		cnt := 0
		b := tx.Bucket([]byte(GetDbBucketName()))
		c := b.Cursor()
		pubKey := GetPubKey(userName)
		prvKey := GetPrvKey(userName)
		for k, v := c.Seek(id); k != nil && bytes.HasPrefix(k, id); k, v = c.Next() {
			note := DeserializeNote(v, prvKey, pwd)
			prevNote := DeserializeNote(b.Get(note.PrevID), prvKey, pwd)
			nextNote := DeserializeNote(b.Get(note.NextID), prvKey, pwd)
			prevNote.NextID = nextNote.ID
			nextNote.PrevID = prevNote.ID

			// save prevNote
			if err := b.Put(prevNote.ID, prevNote.Serialize(pubKey)); err != nil {
				return err
			}
			// save nextNote
			if err := b.Put(nextNote.ID, nextNote.Serialize(pubKey)); err != nil {
				return err
			}

			note = nil // help gc
			cnt++
		}
		fmt.Printf("delete %d notes\n", cnt)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// GetDbFilePath returns the notebook path for userName.
func GetDbFilePath(userName string) string {
	// ~/.terminal_note/terminal_diary_xxx.db
	folderName := GetAppFolderPath()
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err := os.Mkdir(folderName, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
	}
	return folderName + "/" + fmt.Sprintf(dbFileNameTemplate, userName)
}

func GenKeyForUserWithPasswd(userName, passwd string) {
	pubKey, prvKey, err := RSAGenKeyWithPwd(rsaBits, passwd)
	if err != nil {
		log.Panic(err)
	}
	// save pubKey
	pubFile, err := os.OpenFile(GetPublicKeyPath(userName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Panic(err)
	}
	defer pubFile.Close()
	if _, err := pubFile.Write(pubKey); err != nil {
		log.Panic(err)
	}

	// save private key
	privateFile, err := os.OpenFile(GetPrivateKeyPath(userName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Panic(err)
	}
	defer privateFile.Close()
	if _, err := privateFile.Write(prvKey); err != nil {
		log.Panic(err)
	}
}

// isFileExists return true if `fileName` exists
func isFileExists(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

func CheckDbFileExist(userName string) bool {
	return isFileExists(GetDbFilePath(userName))
}

func GetAppFolderPath() string {
	//  ~/.terminal_note/
	user, _ := user2.Current()
	return user.HomeDir + "/" + GetAppHome()
}

func GetPublicKeyPath(userName string) string {
	return GetAppFolderPath() + "/" + fmt.Sprintf(GetPubKeyName(), userName)
}

func GetPrivateKeyPath(userName string) string {
	return GetAppFolderPath() + "/" + fmt.Sprintf(GetPrivateKeyName(), userName)
}

// GetPubKey returns userName's public key
func GetPubKey(userName string) []byte {
	res, err := ioutil.ReadFile(GetPublicKeyPath(userName))
	if err != nil {
		log.Panic(err)
	}
	return res
}

// GetPrvKey returns userName's private key
func GetPrvKey(userName string) []byte {
	res, err := ioutil.ReadFile(GetPrivateKeyPath(userName))
	if err != nil {
		log.Panic(err)
	}
	return res
}
