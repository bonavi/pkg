package xorMask

import "encoding/base64"

// Mask применяет XOR-маскировку к данным и возвращает base64-строку.
func Mask(data []byte, key []byte) string {
	masked := make([]byte, len(data))
	for i := range data {
		masked[i] = data[i] ^ key[i%len(key)]
	}
	return base64.URLEncoding.EncodeToString(masked)
}

// Unmask декодирует base64-строку и снимает XOR-маскировку.
func Unmask(s string, key []byte) ([]byte, error) {
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	for i := range data {
		data[i] ^= key[i%len(key)]
	}
	return data, nil
}