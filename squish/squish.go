package squish

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"sync"
)

var (
	bufPool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	gzipPool = &sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	}
)

// Gzip compresses as slice of bytes using gzip compression.
func Gzip(data []byte) []byte {
	buf := bufPool.Get().(*bytes.Buffer)
	gz := gzipPool.Get().(*gzip.Writer)
	buf.Reset()
	r, w := io.Pipe()
	gz.Reset(w)
	go func() {
		_, _ = gz.Write(data)
		_ = gz.Close()
		_ = w.Close()
	}()
	n, _ := buf.ReadFrom(r)
	buf.Truncate(int(n))
	_ = r.Close()
	res, _ := io.ReadAll(buf)
	bufPool.Put(buf)
	gzipPool.Put(gz)
	return res
}

// Gunzip decompresses a gzip compressed slice of bytes.
func Gunzip(data []byte) (out []byte, err error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	var n int64
	n, _ = buf.ReadFrom(gz)
	err = gz.Close()
	buf.Truncate(int(n))
	res, _ := io.ReadAll(buf)
	bufPool.Put(buf)
	return res, err
}

// B64e encodes the given slice of bytes into base64 standard encoding.
func B64e(in []byte) (out string) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Grow(base64.StdEncoding.EncodedLen(len(in)))
	b64 := base64.NewEncoder(base64.StdEncoding, buf)
	_, _ = b64.Write(in)
	_ = b64.Close()
	res := buf.Bytes()
	bufPool.Put(buf)
	return string(res)
}

// B64d decodes the given string into the original slice of bytes.
// Do note that this is for non critical tasks, it has no error handling for purposes of clean code.
func B64d(str string) (data []byte) {
	if len(str) == 0 {
		return nil
	}
	data, _ = base64.StdEncoding.DecodeString(str)
	return data
}

// UnpackStr UNsafely unpacks (usually banners) that have been base64'd and then gzip'd.
func UnpackStr(encoded string) (string, error) {
	one := B64d(encoded)
	if len(one) == 0 {
		return "", errors.New("0 length base64 decoding result")
	}
	dcytes, err := Gunzip(one)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return string(dcytes), nil
}
