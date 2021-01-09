package service

import (
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
