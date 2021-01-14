package note

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/big"
	"reflect"
	"time"
	"unsafe"
)

/*
* @CreateTime: 2021/1/12 18:34
* @Author: Jemmy@hujm20151021@gmail.com
* @Description: utils
 */

// Data a helper for encoding and decoding.
type Data struct {
	NewData string
}

// EncryptString encodes `source` with gob.
func EncryptString(source string) []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(Data{NewData: source})
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// DecryptString decodes `source` with gob.
func DecryptString(source []byte) string {
	var data Data
	decoder := gob.NewDecoder(bytes.NewReader(source))
	err := decoder.Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	return data.NewData
}

// -----------------------------------------------------------------------
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58Encode encodes a byte array to Base58
func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	reverseBytes(result)
	for b := range input {
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result
}

// Base58Decode decodes Base58-encoded data
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}

// ReverseBytes reverses a byte array
func reverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// FormatTimestamp transform timestamp to formatted time string.
func FormatTimestamp(t int64) string {
	if t == 0 {
		return ""
	}
	return time.Unix(0, t).Format("2006-01-02 15:04:05")
}

// StringToBytes change str to bytes without malloc.
func StringToBytes(str string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bytesHeader := &reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(bytesHeader))
}
