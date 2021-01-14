package note

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

/*
* @CreateTime: 2021/1/13 18:58
* @Author: hujiaming
* @Description: rsa
1. GenRsaKey 得到公私钥
2. RsaEncrypt 传入要加密的数据 和 公钥，进行加密，得到加密后的数据
3. RsaDecrypt 传入加密后的数据和私钥解密，得到解密后的数据
*/

// RSAGenKey generate rsa key pair.
func RSAGenKey(bits int) (pubKey, prvKey []byte, err error) {
	/*
		生成私钥
	*/
	// 1、使用RSA中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	// 2、通过X509标准将得到的RAS私钥序列化为：ASN.1 的DER编码字符串
	privateStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// 3、将私钥字符串设置到pem格式块中
	block1 := &pem.Block{
		Type:  "private key",
		Bytes: privateStream,
	}
	prvKey = pem.EncodeToMemory(block1)

	/*
		生成公钥
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

	// 使用X509将解码之后的数据 解析出来
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	publicKey := keyInit.(*rsa.PublicKey)
	// 使用公钥加密数据
	res, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, src)
	return
}

// RSADecrypt decrypt encrypted data with private key.
func RSADecrypt(encryptedSrc []byte, prvKey []byte) (res []byte, err error) {
	// 解码
	block, _ := pem.Decode(prvKey)
	blockBytes := block.Bytes
	privateKey, err := x509.ParsePKCS1PrivateKey(blockBytes)
	// 还原数据
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedSrc)
	return
}

// --------------------------------
// RSAGenKeyWithPwd generate rsa pair key with specified password
func RSAGenKeyWithPwd(bits int, pwd string) (pubKey, prvKey []byte, err error) {
	/*
		生成私钥
	*/
	// 1、使用RSA中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	// 2、通过X509标准将得到的RAS私钥序列化为：ASN.1 的DER编码字符串
	privateStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// 3、将私钥字符串设置到pem格式块中
	block1 := &pem.Block{
		Type:  "private key",
		Bytes: privateStream,
	}
	// 通过自定义密码加密
	if pwd != "" {
		block1, err = x509.EncryptPEMBlock(rand.Reader, block1.Type, block1.Bytes, []byte(pwd), x509.PEMCipherAES256)
		if err != nil {
			return nil, nil, err
		}
	}
	prvKey = pem.EncodeToMemory(block1)

	/*
		生成公钥
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
	// 解码
	block, _ := pem.Decode(prvKey)
	blockBytes := block.Bytes
	if pwd != "" {
		blockBytes, err = x509.DecryptPEMBlock(block, []byte(pwd))
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(blockBytes)
	// 还原数据
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	return
}
