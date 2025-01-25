package utils

/**
 * @Author nick
 * @Blog http://www.lampnick.com
 * @Email nick@lampnick.com
 */
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}
type ecbEncrypter ecb
type ecbDecrypter ecb

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func Base64DecodeAndDecrypt(cipherText string, aesKey []byte) ([]byte, error) {
	decodeCipher, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return []byte{}, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return []byte{}, err
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(decodeCipher))
	blockMode.CryptBlocks(origData, decodeCipher)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func EncryptAndBase64(src []byte, key []byte) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("aesRandomKey error", err)
	}
	ecb := NewECBEncrypter(block)
	src = PKCS5Padding(src, block.BlockSize())
	crypt := make([]byte, len(src))
	ecb.CryptBlocks(crypt, src)
	return base64.StdEncoding.EncodeToString(crypt)
}

//func PKCS5Padding(cipherText []byte, blockSize int) []byte {
//	padding := blockSize - len(cipherText)%blockSize
//	padText := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(cipherText, padText...)
//}
//
//func PKCS5UnPadding(origData []byte) []byte {
//	length := len(origData)
//	// remove the last byte
//	unPadding := int(origData[length-1])
//	return origData[:(length - unPadding)]
//}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (x *ecbEncrypter) BlockSize() int {
	return x.blockSize
}
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func (x *ecbDecrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

// 加密
func CreateDes(value string) string {
	_authSlice := []string{
		ToString(time.Now().Second()),
		ToString(value),
	}
	_desStr, _ := DesEncrypt([]byte(strings.Join(_authSlice, "_")), []byte("20140924"))
	return hex.EncodeToString(_desStr)
}

// 解密
func DecrDes(value string) ([]string, error) {
	_authCode, _err := hex.DecodeString(value)
	if _err != nil {
		return nil, _err
	}
	_authStr, _err := DesDecrypt(_authCode, []byte("20140924"))
	if _err != nil {
		return nil, _err
	}
	return strings.Split(string(_authStr), "_"), nil
}

// RSA 默认公钥私钥
var _publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8RqfMVcTx9LIILsIK0yQ/p4OT
t3GlQh4Nyh5QcJvQ3ZaV7tOzXYXAlMuVcb4oKO3CG0m5TOWnuOsatAkHP2Y0HAVm
B4frRMbokBSrSrhVjGqaFv/EG86t1jdv8oQWvJTafJ5/LmoB09DubYQCw87Ar8jM
rONPNFT1SeCrIwZyvwIDAQAB
-----END PUBLIC KEY-----
`)
var _privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQC8RqfMVcTx9LIILsIK0yQ/p4OTt3GlQh4Nyh5QcJvQ3ZaV7tOz
XYXAlMuVcb4oKO3CG0m5TOWnuOsatAkHP2Y0HAVmB4frRMbokBSrSrhVjGqaFv/E
G86t1jdv8oQWvJTafJ5/LmoB09DubYQCw87Ar8jMrONPNFT1SeCrIwZyvwIDAQAB
AoGAGHHgBXK5YGTR3Kgdf4RMd4tLVRmDQt6jwkyUxQLp6CNtEshwaiBhZlCrYhrj
gplVzVb3qnxmcPFcbRok9fDwVuhO5NRCYkKKoBQeoPlHA3U0BEdfzPDxbmeDKVQ0
MQsjJoLef9Nitl0nA+AcL/5HQcRK2FyCkmKMh3iGqLmsKbECQQD2Fni9N5Uv1A14
iaOsKQJYl90uBgyDxlcmfcm0dP6/e01CVSFqMCydmxEXwUxevHO39B2Z7C370Cr1
4z/diEXnAkEAw9wPGiPalPeXOeE5BeeqNDrpGCXIijZE6Eyu9NCCAEh3S2ezserl
MP6Bgmfwxv5DHQ8AJQTSbCYE/mRVwAkhaQJAEuMqpSss8hzOY9/8hewn1/Df8vZX
441HhxbEcmtAWiX2ig7Kn8HOytHp/+7AE81W/FlqJDQyW09g3LpyXmhlJQJBAIN9
U247L93estoayEuckfnqqt6ZTx7q/Cvwf2zAJubFv8ER5+PETQYtdwjzewQ9YxU5
IuG3cQVGKQgYmDEKcDECQACbmmFH02brHydC7Vc9FEQ3Rj4UT5sWbGzh8c7FaPtS
SvNYbSFZuSjgva26Jkf3Q2MRVP/znXcn/1vNps6x+nM=
-----END RSA PRIVATE KEY-----
`)

// 加密
func RsaEncrypt(publicKey, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext string) ([]byte, error) {
	_value, _err := base64.StdEncoding.DecodeString(ciphertext)
	if _err != nil {
		return nil, _err
	}
	block, _ := pem.Decode(_privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, _value)
}
