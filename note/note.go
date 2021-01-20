package note

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

/*
* @CreateTime: 2021/1/12 18:28
* @Author: Jemmy@hujm20150121@gmail.com
* @Description: note is a message in notebook
 */

// Note means a diary note
type Note struct {
	ID        []byte // 某一天20191214
	PrevID    []byte // 上一条ID
	NextID    []byte // 下一条ID
	Content   []byte // 日记内容
	IsSecret  bool   // TODO: 是否加密，如果加密，则需要密码访问
	Timestamp int64  // 时间戳
}

// NewNote returns a `Note` object with `content` and `timestamp`
func NewNote(content string) *Note {
	return &Note{
		Content:   EncryptString(content),
		Timestamp: time.Now().UnixNano(),
	}
}

func (n *Note) String() string {
	if n == nil {
		return ""
	}
	var lines []string
	lines = append(lines, fmt.Sprintf("* ----------- Note ------------------"))
	lines = append(lines, fmt.Sprintf("|  ID:       %s", n.ID))
	lines = append(lines, fmt.Sprintf("|  PrevID:   %s", n.PrevID))
	lines = append(lines, fmt.Sprintf("|  NextID:   %s", n.NextID))
	lines = append(lines, fmt.Sprintf("|  DateTime: %s", FormatTimestamp(n.Timestamp)))
	lines = append(lines, fmt.Sprintf("|  Content:  %s", DecryptString(n.Content)))
	lines = append(lines, fmt.Sprintf("* -----------------------------------\n"))
	return strings.Join(lines, "\n")
}

func (n *Note) SimpleString() string {
	if n == nil {
		return ""
	}
	return fmt.Sprintf("%s    %s", FormatTimestamp(n.Timestamp), DecryptString(n.Content))
}

// Serialize encodes a Note to bytes.
func (n *Note) Serialize(pubKey []byte) []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(n)
	if err != nil {
		log.Panic(err)
	}
	encBytes, err := RSAEncrypt(encoded.Bytes(), pubKey)
	if err != nil {
		log.Panic(err)
	}
	return []byte(base64.StdEncoding.EncodeToString(encBytes))
}

// DeserializeNote decodes bytes to Note.
func DeserializeNote(data []byte, prvKey []byte, pwd string) *Note {
	var note Note

	d, err := base64.StdEncoding.DecodeString(string(data))
	source, err := RSADecryptWithPwd(d, prvKey, pwd)
	if err != nil {
		if err == x509.IncorrectPasswordError {
			fmt.Println("incorrect password")
			os.Exit(1)
		}
		log.Panic(err)
	}
	decoder := gob.NewDecoder(bytes.NewReader(source))

	err = decoder.Decode(&note)
	if err != nil {
		log.Panic(err)
	}
	return &note
}

// 设置当前日记的ID(t为创建时间的时间戳)
func (n *Note) SetID(t time.Time) {
	n.ID = []byte(t.Format(idLayout))
}

// GetHeadNote returns a Note with default `headContent`
func GetHeadNote() *Note {
	return generateNoteFromContent(GetHeadDefaultContent(), []byte(GetDbHeadKey()))
}

// GetTailNote return a Note with default `tailContent`
func GetTailNote() *Note {
	return generateNoteFromContent(GetTailDefaultContent(), []byte(GetDbTailKey()))
}

func generateNoteFromContent(content string, id []byte) *Note {
	return &Note{
		ID:        id,
		Content:   EncryptString(content),
		Timestamp: time.Now().UnixNano(),
	}
}
