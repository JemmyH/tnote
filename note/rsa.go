package note

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

/*
* @CreateTime: 2021/1/13 18:58
* @Author: JemmyHu <hujm20151021@gmail.com>
* @Description: rsa
1. GenRsaKey returns publicKey and privateKey
2. RsaEncrypt input data needed to encrypt and publicKey, output data that has been encrypted
3. RsaDecrypt input encrypted data and privateKey, output the decrypted data
*/

// RSAGenKey generate rsa key pair.
func RSAGenKey(bits int) (pubKey, prvKey []byte, err error) {
	/*
		generate privateKey
	*/
	// 1、get a privateKey 
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	// 2、Marshal the privateKey 
	privateStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// 3、put the marshaled privateKey into a Block
	block1 := &pem.Block{
		Type:  "private key",
		Bytes: privateStream,
	}
	prvKey = pem.EncodeToMemory(block1)

	/*
		genarate publicKey from privateKey
	*/
	publicKey := privateKey.PublicKey
	publicStream, err := x509.MarshalPKIXPublicKey(&publicKey)
	block2 := &pem.Block{
		Type:  "public key",
		Bytes: publicStream,
	}
	pubKey = pem.EncodeToMemory(block2)
	return pubKey, prvKey, nil
}

// RSAEncrypt use publicKey to encrypt source data.
func RSAEncrypt(src []byte, pubKey []byte) (res []byte, err error) {
	block, _ := pem.Decode(pubKey)

	// unmarshal publicKey
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	publicKey := keyInit.(*rsa.PublicKey)
	// encrypt data with publicKey
	res, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, src)
	return
}

// RSADecrypt decrypt encrypted data with private key.
func RSADecrypt(encryptedSrc []byte, prvKey []byte) (res []byte, err error) {
	// decode the privateKey
	block, _ := pem.Decode(prvKey)
	blockBytes := block.Bytes
	privateKey, err := x509.ParsePKCS1PrivateKey(blockBytes)
	// decrypt by privateKey
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedSrc)
	return
}

// --------------------------------
// RSAGenKeyWithPwd generate rsa pair key with specified password
func RSAGenKeyWithPwd(bits int, pwd string) (pubKey, prvKey []byte, err error) {
	/*
		generate privateKey
	*/
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	privateStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block1 := &pem.Block{
		Type:  "private key",
		Bytes: privateStream,
	}
	// use optional password
	if pwd != "" {
		block1, err = x509.EncryptPEMBlock(rand.Reader, block1.Type, block1.Bytes, []byte(pwd), x509.PEMCipherAES256)
		if err != nil {
			return nil, nil, err
		}
	}
	prvKey = pem.EncodeToMemory(block1)

	/*
		generate publuicKey from privateKey
	*/
	publicKey := privateKey.PublicKey
	publicStream, err := x509.MarshalPKIXPublicKey(&publicKey)
	block2 := &pem.Block{
		Type:  "public key",
		Bytes: publicStream,
	}
	pubKey = pem.EncodeToMemory(block2)
	return pubKey, prvKey, nil
}

// RSADecryptWithPwd decrypt src with private key and password
func RSADecryptWithPwd(src []byte, prvKey []byte, pwd string) (res []byte, err error) {
	block, _ := pem.Decode(prvKey)
	blockBytes := block.Bytes
	if pwd != "" {
		blockBytes, err = x509.DecryptPEMBlock(block, []byte(pwd))
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(blockBytes)
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	return
}
