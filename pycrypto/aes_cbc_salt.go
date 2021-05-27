package pycrypto

//参考：https://github.com/funny/crypto/tree/master/aes256cbc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"go.uber.org/zap"
)

const (
	CBC_SALT_LEN = 8
	CBC_IV_LEN   = 16 //BYTES
	CBC_KEY_LEN  = 32
	CBC_CRED_LEN = 48 //CBC_IV_LEN+CBC_KEY_LEN
)

//预先生成PrePadPatterns
var PrePadPatterns [aes.BlockSize + 1][]byte

//fix header
var CbcfixedSaltHeader = []byte("Salted__")

func init() {
	for i := 0; i < len(PrePadPatterns); i++ {
		PrePadPatterns[i] = bytes.Repeat([]byte{byte(i)}, i)
	}
	/*
		[]
		[1]
		[2 2]
		[3 3 3]
		[4 4 4 4]
		[5 5 5 5 5]
		[6 6 6 6 6 6]
		[7 7 7 7 7 7 7]
		[8 8 8 8 8 8 8 8]
		[9 9 9 9 9 9 9 9 9]
		[10 10 10 10 10 10 10 10 10 10]
		[11 11 11 11 11 11 11 11 11 11 11]
		[12 12 12 12 12 12 12 12 12 12 12 12]
		[13 13 13 13 13 13 13 13 13 13 13 13 13]
		[14 14 14 14 14 14 14 14 14 14 14 14 14 14]
		[15 15 15 15 15 15 15 15 15 15 15 15 15 15 15]
		[16 16 16 16 16 16 16 16 16 16 16 16 16 16 16 16]
	*/
}

type Creds [CBC_CRED_LEN]byte

func (c *Creds) Extract(password, salt []byte) (key, iv []byte) {
	m := c[:]
	buf := make([]byte, 0, 16+len(password)+len(salt))
	var prevSum [16]byte
	for i := 0; i < 3; i++ {
		n := 0
		if i > 0 {
			n = 16
		}
		buf = buf[:n+len(password)+len(salt)]
		copy(buf, prevSum[:])
		copy(buf[n:], password)
		copy(buf[n+len(password):], salt)
		prevSum = md5.Sum(buf)
		copy(m[i*16:], prevSum[:])
	}
	return c[:32], c[32:]
}

type Cbc256WithSalt struct {
	SecretPass string
	Salt       string
	Logger     *zap.Logger
}

func NewCbc256WithSalt(logger *zap.Logger, pass string) *Cbc256WithSalt {
	c := Cbc256WithSalt{
		SecretPass: pass,
		Logger:     logger,
	}

	return &c
}

func (c *Cbc256WithSalt) Encrypt(origin_text string) ([]byte, error) {
	var (
		creds Creds
	)
	origin_text_c := []byte(origin_text)
	// Generate random salt
	var salt [CBC_SALT_LEN]byte
	//_, err := io.ReadFull(rand.Reader, salt)	//WRONG cannot use salt (type [8]byte) as type []byte in argument to io.ReadFull
	_, err := io.ReadFull(rand.Reader, salt[:])
	if err != nil {
		c.Logger.Error("generate random error", zap.String("errmsg", err.Error()))
		return nil, err
	}

	/*
		|Salted__(8 byte)|salt(8 byte)|plaintext|
	*/
	data := make([]byte, len(origin_text)+aes.BlockSize /*16*/)
	copy(data[0:], CbcfixedSaltHeader)
	copy(data[8:], salt[:])
	copy(data[aes.BlockSize:], origin_text_c)

	key, iv := creds.Extract([]byte(c.SecretPass), salt[:])
	padded, err := pkcs7Pading(data)
	if err != nil {
		c.Logger.Error("pkcs7Pading error", zap.String("errmsg", err.Error()))
		return nil, err
	}

	cc, err := aes.NewCipher(key)
	if err != nil {
		c.Logger.Error("NewCipher error", zap.String("errmsg", err.Error()))
		return nil, err
	}
	cbc := cipher.NewCBCEncrypter(cc, iv)
	//fmt.Println(padded[aes.BlockSize:])

	// 只从plaintext位置开始加密（上图）
	cbc.CryptBlocks(padded[aes.BlockSize:], padded[aes.BlockSize:])
	return padded, nil
}

//for decrypt
func (c *Cbc256WithSalt) Decrypt(encrypt_str []byte) ([]byte, error) {
	/*
		|Salted__(8 byte)|salt(8 byte)|encrypt_text|
	*/
	if len(encrypt_str) < aes.BlockSize {
		return nil, errors.New("length illegal")
	}
	saltHeader := encrypt_str[:aes.BlockSize]
	if !bytes.Equal(saltHeader[:8], CbcfixedSaltHeader) {
		return nil, errors.New("check cbc fixed header error")
	}
	var creds Creds
	key, iv := creds.Extract([]byte(c.SecretPass), saltHeader[8:])

	if len(encrypt_str) == 0 || len(encrypt_str)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("encrypt_str length illegal: len=%d", len(encrypt_str))
	}
	cc, err := aes.NewCipher(key)
	if err != nil {
		c.Logger.Error("NewCipher error", zap.String("errmsg", err.Error()))
		return nil, err
	}
	cbc := cipher.NewCBCDecrypter(cc, iv)
	cbc.CryptBlocks(encrypt_str[aes.BlockSize:], encrypt_str[aes.BlockSize:])

	//删除加密时候填充的padding
	return pkcs7Unpading(encrypt_str[aes.BlockSize:])
}

func pkcs7Pading(data []byte) ([]byte, error) {
	if len(data)%aes.BlockSize == 0 {
		//no need to padding
		return data, nil
	}
	padlen := 1
	for ((len(data) + padlen) % aes.BlockSize) != 0 {
		padlen = padlen + 1
	}
	return append(data, PrePadPatterns[padlen]...), nil
}

//
func pkcs7Unpading(data []byte) ([]byte, error) {
	if len(data)%aes.BlockSize != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen > aes.BlockSize || padlen == 0 {
		return nil, errors.New("param illegal")
	}
	if !bytes.Equal(PrePadPatterns[padlen], data[len(data)-padlen:]) {
		return nil, errors.New("param illegal")
	}
	return data[:len(data)-padlen], nil
}

func main() {
	logger, _ := zap.NewProduction()
	c := NewCbc256WithSalt(logger, "aaaaaa")
	enc, err := c.Encrypt("ABCEDFG")
	fmt.Println(enc, err)

	dec, err := c.Decrypt(enc)
	fmt.Println(string(dec), err)
}
