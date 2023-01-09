package server

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	b64 "encoding/base64"
	"fmt"
	"io"
	"log"
)

var PrivKey *rsa.PrivateKey

func init() {
	raw, err := b64.StdEncoding.DecodeString("MIIEpQIBAAKCAQEApvxieEC28whOn2Qi++uxVDM6m2spQHWoAZBz2WiqauAs9tpgeLuQgKBxpK4W+qm0oVksXXMKE5aDYdmvVz4X++S4KkElTJutg7VJzXVSlIOlCIRR6UyPmPQYmdGWx+gf+LBd6HbhAKbSzJBDQ7j4MW4dqjbuKBxqbJQcg13vnvefgIog9XD3k5eA+dkLk9DuMSZDd24ACcBtiwEj2ATy6ccEY4VXGk+AUWA6VowlcUTTNC7Yp+19QgpVpWIpJ8gcIptdBJ9TDbP0/ST48niqQQDWLeXsa0Wg+ZR+/uTV36aMqYQuf4xHrm3hhroK2diOasJEDrD1ml7Afty717YNMQIDAQABAoIBAQClRWawXly0baRjXVjCvaPlEk8PRCCwC8McyTvgEheZcAcQy1JwLDP5GtNfim5z6UM97mRamWF/wZiHYEyKrIpQZS9hotin2e0CToudLmFtXF4a79uibIQzfmRa2XXCpZv/J4/KZN6NJo+8p4vrm0cKpVH3BibwzC8JCA6wdmiTjCiB7inG0bs44Eu/3UOIE5vtz3EDrlPYKP13Lw3GXS7bCqS52qrjD3/b90+7kRjY+lZlLef38ZM6os+Q4cX0PJ0mRJ70RXfvNXOinvEIDFmTW5hv1hxh/02cjbSQjXcOXax0FFQyIVGqCxwakybM4Mg3Em8p3yUzVOYvkVnnLcoBAoGBAMkAsNOKDAdSAJukqdOKupHueZGFV1+kSLdXT7GbfHhxNy4b7XAgbKVQerm2kQzkLgFiQ5ZpbfJJWf3i3NBE2afFwhJRticyYX3hyaTaINxVysFeBL/jKIUIZcLCExxxH3WU8Gl+rXZFtxBbqYdBRx2mQz+dWJ2lA0NvSmBW0jtFAoGBANSs+RUiqoTdE4ffsctmlKO84h9eMxONQ1cZ4SqK60iddzXeA1xQKvVamc44y2bwyVcQKlq451AyMT/wTCv3AJKJK949UA23hx56otw2qgQZuwX9inzCmvazka7bcBF5z/uasw7cSUUM6Q7nb7wtD6b68096AML70luTNFtHxjL9AoGBAMa4iRnK/JNsLi+yvzfmiwfl5ojJdJWZHU8t4hts5sVI4U4TzE4zsFZMV9kttwAww48YsEuPlmSYwoDwfnDl8O4e5P0pjdX4yEwlIy95fE16AEfmhPmVQqUrpTfEmhJfgMPF6V3TIPmyeQeSJ+wRzJZynz/QdyD8WFqeN8FBdP2lAoGBALBehe3mIr0uTX0XoG1Ks6eaA3f5+aeUNa0s9BMAw6AjnfHZHLZYcVepe/WOGfhTZNVDvawgvQs/pKIemDy5iQr8oJmcBSBq+63mC9tNpe7im7ubCFbwV+yQ/BajOivz9ev03dtMCaMu0rOecIYAZIOYh2B4j9sjVM7Go9uzCQXNAoGAfJlUeZhjz0CDO4VfcYs0S73Dq775KeBq/OsYfDEkYWOAheaWxOrTAOkOIznqiLoqH+u+R8rnUbYki6qlrHMr8vh/EXI4HPMja9nuQw3Owxqtg2eWL4h6eX46BvORVhZHq00l8nGIlEM0VIL/bL6BwQslLbmu22NbCOpwzovQzLc=")
	if err != nil {
		log.Fatal("something went horribly wrong")
	}
	PrivKey, err = x509.ParsePKCS1PrivateKey(raw)
}

// Encrypts text with the passphrase
func Encrypt(text string, passphrase string) string {
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err.Error())
	}

	key, iv := __DeriveKeyAndIv(passphrase, string(salt))

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	pad := __PKCS5Padding([]byte(text), block.BlockSize())
	ecb := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(pad))
	ecb.CryptBlocks(encrypted, pad)

	return b64.StdEncoding.EncodeToString([]byte("Salted__" + string(salt) + string(encrypted)))
}

func EncryptWithSalt(text string, passphrase string, salt string) string {
	key, iv := __DeriveKeyAndIv(passphrase, salt)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	pad := __PKCS5Padding([]byte(text), block.BlockSize())
	ecb := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(pad))
	ecb.CryptBlocks(encrypted, pad)

	return b64.StdEncoding.EncodeToString([]byte(string(encrypted)))
}

// Decrypts encrypted text with the passphrase
func Decrypt(encrypted string, passphrase string) string {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			_, ok = r.(error)
			if !ok {
				_ = fmt.Errorf("pkg: %v", r)
			}
		}
	}()
	ct, _ := b64.StdEncoding.DecodeString(encrypted)
	if len(ct) < 16 || string(ct[:8]) != "Salted__" {
		return ""
	}

	salt := ct[8:16]
	ct = ct[16:]
	key, iv := __DeriveKeyAndIv(passphrase, string(salt))
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCDecrypter(block, []byte(iv))
	dst := make([]byte, len(ct))
	cbc.CryptBlocks(dst, ct)

	return string(__PKCS5Trimming(dst))
}

func DecryptWithSalt(encrypted string, passphrase string, salt string) string {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			_, ok = r.(error)
			if !ok {
				_ = fmt.Errorf("pkg: %v", r)
			}
		}
	}()
	ct, _ := b64.StdEncoding.DecodeString(encrypted)
	key, iv := __DeriveKeyAndIv(passphrase, salt)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCDecrypter(block, []byte(iv))
	dst := make([]byte, len(ct))
	cbc.CryptBlocks(dst, ct)

	return string(__PKCS5Trimming(dst))
}

func __PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func __PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func __DeriveKeyAndIv(passphrase string, salt string) (string, string) {
	salted := ""
	dI := ""

	for len(salted) < 48 {
		md := md5.New()
		md.Write([]byte(dI + passphrase + salt))
		dM := md.Sum(nil)
		dI = string(dM[:16])
		salted = salted + dI
	}

	key := salted[0:32]
	iv := salted[32:48]

	return key, iv
}

func RSADecrypt(data string) string {
	dataBytes, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	decryptedBytes, err := PrivKey.Decrypt(nil, dataBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return ""
	}
	return string(decryptedBytes)
}