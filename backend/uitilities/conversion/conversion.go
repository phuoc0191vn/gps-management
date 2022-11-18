package conversion

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

func ToInterfaces(list []string) []interface{} {
	result := make([]interface{}, len(list))
	for i := range list {
		result[i] = list[i]
	}
	return result
}

func ToJson(obj interface{}) string {
	out, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(out)
}

func HashSHA256(message string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(message []byte) string {
	hash := sha256.New()
	hash.Write(message)
	return hex.EncodeToString(hash.Sum(nil))
}

func HashMd5(message string) string {
	hash := md5.New()
	hash.Write([]byte(message))
	return hex.EncodeToString(hash.Sum(nil))
}

func TotalHash(buffer []byte) (hashmd5 string, hashsha256 string) {
	md5Hash := md5.New()
	md5Hash.Write(buffer)
	hashmd5 = hex.EncodeToString(md5Hash.Sum(nil))

	sha256Hash := sha256.New()
	sha256Hash.Write(buffer)
	hashsha256 = hex.EncodeToString(sha256Hash.Sum(nil))
	return
}

// StringToSequence convert binary string to int64 sequence
func StringToSequence(source string) (seq int64, err error) {
	bits := 8
	if len(source) != bits {
		return -1, fmt.Errorf("Wrong source format (%s)", source)
	}
	for i := 0; i < bits; i++ {
		seq += int64(source[i]) << uint((bits-i-1)*bits)
	}
	return
}

// SequenceToString convert sequence number to binary string
func SequenceToString(seq int64) (seqStr string, err error) {
	bits := 8
	source := make([]byte, bits)
	for i := 0; i < bits; i++ {
		source[i] = byte((seq >> uint((bits-i-1)*bits)) & 255)
	}
	return string(source), nil
}

// StringToSlice splits a string into a slice by a separator
func StringToSlice(input, separator string) []string {
	result := make([]string, 0)
	for _, value := range strings.Split(input, separator) {
		result = append(result, strings.TrimSpace(value))
	}
	return result
}

func DecodeBase64Url(value string) string {
	buffer, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return ""
	}

	return string(buffer)
}

func DecodeStdBase64(value string) string {
	buffer, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return ""
	}

	return string(buffer)
}

func EncodeStdBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
