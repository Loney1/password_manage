package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	ErrEmptyContent = errors.New("encrpty plain content empty")
	ErrNotFull      = errors.New("crypto/cipher: input not full blocks")
)

type Aes struct {
	cipher []byte
}

func NewAes(cipher []byte) *Aes {
	return &Aes{
		cipher: cipher,
	}
}

func (aes *Aes) Encrypt(src string) ([]byte, error) {
	return AesEncrypt(src, aes.cipher)
}

func (aes *Aes) Decrypt(crypted []byte) ([]byte, error) {
	return AesDecrypt(crypted, aes.cipher)
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))

	if len(crypted)%blockMode.BlockSize() != 0 {
		return nil, ErrNotFull
	}

	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)

	return origData, nil
}

func AesEncrypt(src string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if src == "" {
		return nil, ErrEmptyContent
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))

	if len(content)%ecb.BlockSize() != 0 {
		return nil, ErrNotFull
	}
	ecb.CryptBlocks(crypted, content)

	return crypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	if padding == 0 {
		return ciphertext
	} else {
		return append(ciphertext, bytes.Repeat([]byte{byte(0)}, padding)...)
	}
}

func PKCS5UnPadding(origData []byte) []byte {
	for i := len(origData) - 1; ; i-- {
		if origData[i] != 0 {
			return origData[:i+1]
		}
	}
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	//if len(src)%x.blockSize != 0 {
	//	panic("crypto/cipher: input not full blocks")
	//}
	//if len(dst) < len(src) {
	//	panic("crypto/cipher: output smaller than input")
	//}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	//if len(src)%x.blockSize != 0 {
	//	panic("crypto/cipher: input not full blocks")
	//}
	//if len(dst) < len(src) {
	//	panic("crypto/cipher: output smaller than input")
	//}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
