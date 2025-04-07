package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
)

var krokiURL = []byte("https://kroki.io/")

func makeURL(dtype []byte, ftype []byte, code []byte) ([]byte, error) {
	encoded, err := encode(code)
	if err != nil {
		return []byte{}, err
	}
	url := append(krokiURL, dtype...)
	url = append(url, '/')
	url = append(url, ftype...)
	url = append(url, '/')
	return append(url, encoded...), nil
}

func encode(input []byte) ([]byte, error) {
	var compressed bytes.Buffer
	writer, err := zlib.NewWriterLevel(&compressed, 9)
	if err != nil {
		return []byte{}, fmt.Errorf("error creating deflate writer: %v", err)
	}

	_, err = writer.Write([]byte(input))
	writer.Close()
	if err != nil {
		return []byte{}, fmt.Errorf("fail to create payload: %w", err)
	}

	result := make([]byte, base64.URLEncoding.EncodedLen(compressed.Len()))
	base64.URLEncoding.Encode(result, compressed.Bytes())
	return result, nil
}
