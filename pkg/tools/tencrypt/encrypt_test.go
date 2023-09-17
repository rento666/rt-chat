package tencrypt

import (
	"fmt"
	"rt-chat/pkg/tools/ttrans"
	"testing"
)

func TestSplit(t *testing.T) {
	s := "81uUpJ97/0uKMoRS4g9Xew==|InLFJtKmExXojgBEZ1RdbGeisWvEk/ToYFKvCHOp7Fs="
	fmt.Println(ttrans.SplitInString(s, "|"))
}

func TestRsa(t *testing.T) {
	//// 生成并保存公钥和私钥
	//err := GenerateAndSaveKeyPair()
	//if err != nil {
	//	fmt.Println("Failed to generate and save key pair:", err)
	//	return
	//}

	// 使用公钥加密数据
	ciphertext, err := RsaEncrypt([]byte("Hello, World!"))
	if err != nil {
		fmt.Println("Failed to encrypt with public key:", err)
		return
	}

	// 使用私钥解密数据
	plaintext, err := RsaDecrypt(ciphertext)
	if err != nil {
		fmt.Println("Failed to decrypt with private key:", err)
		return
	}

	fmt.Println("Plaintext:", string(plaintext))
}

func TestMD5(t *testing.T) {

}
