package main

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestGetExtensionFromFilePath(t *testing.T) {
	ret := getExtensionFromFilePath("test.txt")
	if ret != ".txt" {
		t.Errorf("got %s; want .txt", ret)
	}

	ret = getExtensionFromFilePath("test without extension")
	if ret != "" {
		t.Errorf("got %s; want \"\"", ret)
	}
}

func BenchmarkGetExtensionFromFilePath(b *testing.B) {
	const numberOfTexts = 1000

	texts := make([]string, numberOfTexts)
	for i := 0; i < numberOfTexts; i++ {
		texts[i] = RandStringBytesMaskImprSrcSB(i * 2)
		if i%2 == 0 {
			texts[i] = texts[i] + ".txt"
		}
		if i%4 == 0 {
			texts[i] = texts[i] + ".svg"
		}
	}
	b.ResetTimer()
	counter := 0
	for i := 0; i < b.N; i++ {
		if counter >= numberOfTexts {
			counter = 0
		}
		_ = getExtensionFromFilePath(texts[counter])
		counter++
	}
}

//see: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ./*-+\\"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
