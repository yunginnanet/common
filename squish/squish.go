package squish

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
)

// Gzip compresses as slice of bytes using gzip compression.
func Gzip(data []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	// In theory this should never fail, and I don't know how to make the gzip buffered reader fail in testing.
	_, _ = gz.Write(data)
	_ = gz.Close()
	return b.Bytes()
}

// Gunzip decompresses a gzip compressed slice of bytes.
func Gunzip(data []byte) (out []byte, err error) {
	var gz *gzip.Reader
	gz, err = gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return
	}
	return io.ReadAll(gz)
}

// B64e encodes the given slice of bytes into base64 standard encoding.
func B64e(cytes []byte) (data string) {
	data = base64.StdEncoding.EncodeToString(cytes)
	return
}

// B64d decodes the given string into the original slice of bytes.
// Do note that this is for non critical tasks, it has no error handling for purposes of clean code.
func B64d(str string) (data []byte) {
	data, _ = base64.StdEncoding.DecodeString(str)
	return data
}

// UnpackStr UNsafely unpacks (usually banners) that have been base64'd and then gzip'd.
func UnpackStr(encoded string) (string, error) {
	dcytes, err := Gunzip(B64d(encoded))
	if err != nil {
		return "", err
	}
	return string(dcytes), nil
}
