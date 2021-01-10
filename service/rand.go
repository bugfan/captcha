package service

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

const seedOffset = 34052689

func init() {
	rand.Seed(time.Now().UnixNano() & seedOffset)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var letterHex = []rune("abcdef0123456789")

func RandStringHex(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterHex[rand.Intn(len(letterHex))]
	}
	return string(b)
}

var encrytToken = "Ty"

const encryptIndex = 10

func EncodeStr(str string) string {
	sEnc := base64.StdEncoding.EncodeToString([]byte(str))
	left := sEnc[0:encryptIndex]
	right := sEnc[encryptIndex:]
	return left + encrytToken + right
}
func DecodeStr(str string) string {
	defer func() {
		recover()
	}()
	left := str[0:encryptIndex]
	right := str[encryptIndex+len(encrytToken):]
	sDec, err := base64.StdEncoding.DecodeString(left + right)
	if err != nil {
		fmt.Printf("Error decoding string:%s ,source str is:%s", err.Error(), str)
		return str
	}
	return string(sDec)
}
