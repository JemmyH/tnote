package note

import (
	"fmt"
	"testing"
)

/*
* @CreateTime: 2021/1/13 19:46
* @Author: JemmyHu <hujm20151021@gmail.com>
* @Description:
 */

func TestGenRsaKey(t *testing.T) {
	sourceData := "好的代码本身就是最好的说明文档"
	pwd := "123456"
	// 创建公私钥
	pubKey, prvKey, err := RSAGenKeyWithPwd(2048, pwd)
	if err != nil {
		panic(err)
	}
	fmt.Println("gen pubKey and prvKey ok!")
	fmt.Printf("before encrypt: %s\n", sourceData)
	// 使用公钥加密
	encryptData, err := RSAEncrypt([]byte(sourceData), pubKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("after encrypt: %v\n", encryptData)
	// 使用私钥解密
	decryptData, err := RSADecryptWithPwd(encryptData, prvKey, pwd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("after decrypt: %s\n", string(decryptData))
	fmt.Printf("equal? %v \n", string(decryptData) == sourceData)
}

func TestRSAEncrypt(t *testing.T) {
	sourceData := "我的头发长，天下我为王"
	// 创建公私钥
	pubKey, prvKey, err := RSAGenKey(2048)
	if err != nil {
		panic(err)
	}
	fmt.Println("gen pubKey and prvKey ok!")
	fmt.Printf("before encrypt: %s\n", sourceData)
	// 使用公钥加密
	encryptData, err := RSAEncrypt([]byte(sourceData), pubKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("after encrypt: %v\n", encryptData)
	// 使用私钥解密
	decryptData, err := RSADecrypt(encryptData, prvKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("after decrypt: %s\n", string(decryptData))
	fmt.Printf("equal? %v \n", string(decryptData) == sourceData)
}
