package utils

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	uuidLock *sync.Mutex
	lastNum  int64
	count    int
)

func init() {
	uuidLock = new(sync.Mutex)
	count = 0
}

// UUID generate unique values
func UUID() string {
	uuidLock.Lock()
	result := time.Now().UnixNano()
	if lastNum == result {
		count++
	} else {
		count = 0
		lastNum = result
	}
	uuidLock.Unlock()
	return MD5(strconv.Itoa(int(lastNum)) + strconv.Itoa(count))
}

func IsDirExist(path string) bool {
	fi, err := os.Stat(path)
	log.Println(err)
	return err == nil && fi.IsDir() || os.IsExist(err)
}

func IsFileExist(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir() || os.IsExist(err)
}

func MD5(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA512(byteArr []byte) (string, error) {
	return SHA512Stream(bytes.NewReader(byteArr))
}

func SHA512Stream(reader io.Reader) (string, error) {
	hashed := sha512.New()
	io.Seeker(reader).Seek(0, 0)
	if _, err := io.Copy(hashed, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hashed.Sum(nil))
}

func Compare(a, b string) int {
	switch {
	case a == b:
		return 0
	case a < b:
		return -1
	default:
		return 1
	}
}

func GenerateRSAkeyPair(bits int) ([]byte, []byte, error) {
	prvKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	prvBytes := x509.MarshalPKCS1PrivateKey(prvKey)
	pubBytes, err := x509.MarshalPKIXPublicKey(&prvKey.PublicKey)

	if err != nil {
		return nil, nil, err
	}

	prvBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: prvBytes}
	pubBlock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes}

	return pem.EncodeToMemory(prvBlock), pem.EncodeToMemory(pubBlock), nil
}

func RSAEncrypt(keyBytes, contentBytes []byte) ([]byte, error) {
	pubKey, err := getPubKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, pubKey, contentBytes)
}

func RSADecrypt(keyBytes, contentBytes []byte) ([]byte, error) {
	prvKey, err := getPrvKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, prvKey, contentBytes)
}

func SHA256Sign(keyBytes, contentBytes []byte) ([]byte, error) {
	prvKey, err := getPrvKey(keyBytes)
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(contentBytes)
	return rsa.SignPKCS1v15(rand.Reader, prvKey, crypto.SHA256, hashed[:])
}

func SHA256Verify(keyBytes, contentBytes, signBytes []byte) error {
	pubKey, err := getPubKey(keyBytes)
	if err != nil {
		return err
	}

	signStr := hex.EncodeToString(signBytes)
	newSignBytes, _ := hex.DecodeString(signStr)
	hashed := sha256.Sum256(contentBytes)
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], newSignBytes)
}

func getPrvKey(prvBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(prvBytes)
	if block == nil {
		return nil, errors.New("Fail to decode private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
func getPubKey(pubBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubBytes)
	if block == nil {
		return nil, errors.New("Fail to decode public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Fail get public key form pulic interface")
	}

	return pubKey, nil
}
