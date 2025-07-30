package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// aes uses symmetric key block cipher
func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	Block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	GCM, err := cipher.NewGCM(Block)
	if err != nil {
		return nil, err
	}

	nonceByte := make([]byte, GCM.NonceSize())

	if _, err := rand.Read(nonceByte); err != nil {
		return nil, err
	}

	ciphertext := GCM.Seal(nonceByte, nonceByte, plaintext, nil)
	return ciphertext, nil
}

func decrypt(decrypText []byte, key []byte) ([]byte, error) {
	Block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	GCM, err := cipher.NewGCM(Block)
	if err != nil {
		return nil, err
	}

	gcmNonce := GCM.NonceSize()
	nonce, ciphertext := decrypText[:gcmNonce], decrypText[gcmNonce:]

	plainText, err := GCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

/**
// Doesn't have utilize AES-ni instructions, also symmetric key block cipher
func weakEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := chacha20poly1305.New(key)
	checkError(err)

	chachaNonce := make([]byte, block.NonceSize())

	ciphertext := block.Seal(chachaNonce, chachaNonce, plaintext, nil)
	return ciphertext, nil
}

func weakDecrypt(decryptText []byte, key []byte) ([]byte, error) {
	block, err := chacha20poly1305.New(key)
	checkError(err)

	chachaNonce := block.NonceSize()

	nonce, decrypText := decryptText[:chachaNonce], decryptText[chachaNonce:]
	plaintext, err := block.Open(nil, nonce, decrypText, nil)
	return plaintext, err
}
*/
