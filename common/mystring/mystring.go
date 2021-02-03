package mystring

import (
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	//test dev:q dev

	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter  index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func init() {
}
func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
func ReverseByte(s []byte) []byte {
	var result []byte
	for i := len(s) - 1; i >= 0; i-- {
		result = append(result, s[i])
	}
	return result
}
func IsInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func HtmlUnEscape(content string) string {
	//find and replace utf16
	if content == "" {
		return ""
	}
	content = strings.Replace(content, `%u`, `\u`, -1)
	t2, _ := strconv.Unquote(`"` + content + `"`)
	t3, _ := url.QueryUnescape(t2)

	return t3
}

func Alphabetnumeric(dat string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "abcdefghij"
	}
	return reg.ReplaceAllString(dat, "")
}
