package squish

import (
        "bytes"
        "compress/gzip"
        "encoding/base64"
        "io"
)

// Gzip compresses as slice of bytes using gzip compression.
func Gzip(data []byte) ([]byte, error) {
        var b bytes.Buffer
        gz := gzip.NewWriter(&b)
        if _, err := gz.Write(data); err != nil {
                return data, err
        }
        if err := gz.Close(); err != nil {
                return data, err
        }
        return b.Bytes(), nil
}


// Gunzip decompresses a gzip compressed slice of bytes.
func Gunzip(data []byte) ([]byte, error) {
        gz, err := gzip.NewReader(bytes.NewReader(data))
        if err != nil {
                return nil, err
        }
        out, err := io.ReadAll(gz)
        if err != nil {
                return nil, err
        }
        return out, err
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
func UnpackStr(encoded string) string {
        dcytes, _ := Gunzip(B64d(encoded))
        return string(dcytes)
}
