package hash

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
)

const (
	kayosBlake2b = "Kr+6ONDx+cq/WvhHpQE/4LVuJYi9QHz1TztHNTWwa9KJWqHxfTNLKF3YxrcLptA3wO0KHm83Lq7gpBWgCQzPag=="
	kayosMD5     = "aoNjwileNitB208vOwpIow=="
	kayosSHA1    = "M23ElC0sQYAK+MaMZVmza2L8mss="
	kayosSHA256  = "BagY0TmoGR3O7t80BGm4K6UHPlqEg6HJirwQmhrPK4U="
	kayosSHA512  = "xiuo2na76acrWXCTTR++O1pPZabOhyj8nbfb5Go3e1pEq9VJYIsOioTXalf2GCuERmFecWkmaL5QI8mIXXWpNA=="
	kayosCRC32   = "xtig5w=="
)

var kayosByteSlice = []byte{107, 97, 121, 111, 115, 10}

var (
	ogsha1, _    = base64.StdEncoding.DecodeString(kayosSHA1)
	ogsha256, _  = base64.StdEncoding.DecodeString(kayosSHA256)
	ogsha512, _  = base64.StdEncoding.DecodeString(kayosSHA512)
	ogmd5, _     = base64.StdEncoding.DecodeString(kayosMD5)
	ogBlake2b, _ = base64.StdEncoding.DecodeString(kayosBlake2b)
	ogCRC32, _   = base64.StdEncoding.DecodeString(kayosCRC32)
	valids       = map[Type][]byte{
		TypeSHA1:    ogsha1,
		TypeSHA256:  ogsha256,
		TypeSHA512:  ogsha512,
		TypeMD5:     ogmd5,
		TypeCRC32:   ogCRC32,
		TypeBlake2b: ogBlake2b,
	}
)

func TestSum(t *testing.T) {
	t.Parallel()
	if Sum(TypeNull, []byte("yeet")) != nil {
		t.Fatal("Sum(TypeNull, []byte(\"yeet\")) should have returned nil")
	}

	for k, v := range valids {
		typeToTest := k
		valueToTest := v
		t.Run(typeToTest.String()+"/string_check", func(t *testing.T) {
			t.Parallel()
			if !strings.EqualFold(typeToTest.String(), StringToType(typeToTest.String()).String()) {
				t.Errorf("[FAIL] %s: wanted %s, got %s", typeToTest.String(), typeToTest.String(), StringToType(typeToTest.String()).String())
			}
		})
		t.Run(typeToTest.String()+"/static_check", func(t *testing.T) {
			t.Parallel()
			mySum := Sum(typeToTest, kayosByteSlice)
			if !bytes.Equal(mySum, valueToTest) {
				t.Errorf("[FAIL] %s: wanted %v, got %v", typeToTest.String(), valueToTest, mySum)
			}
		})
		t.Run(typeToTest.String()+"/file_check", func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), typeToTest.String()) // for coverage
			if err := os.WriteFile(path, kayosByteSlice, os.ModePerm); err != nil {
				t.Fatalf("[FAIL] failed to write test fle for TestSum: %s", err.Error())
			}
			res, err := SumFile(typeToTest, path)
			if err != nil {
				t.Fatalf("[FAIL] failed to read test fle for TestSum: %s", err.Error())
			}
			if !bytes.Equal(res, valueToTest) {
				t.Errorf("[FAIL] %s: wanted %v, got %v", typeToTest.String(), valueToTest, res)
			}
		})
	}
	t.Run("bad file", func(t *testing.T) {
		t.Parallel()
		_, err := SumFile(TypeSHA1, "/dev/null")
		if err == nil {
			t.Fatal("SumFile should have returned an error")
		}
		if _, err = SumFile(TypeSHA1, entropy.RandStrWithUpper(500)); err == nil {
			t.Fatal("SumFile should have returned an error")
		}
	})
	t.Run("unknown type", func(t *testing.T) {
		t.Parallel()
		if Type(uint8(94)).String() != "unknown" {
			t.Fatal("Type(uint(9453543)).String() should have returned \"unknown\"")
		}
		if StringToType(entropy.RandStr(10)) != TypeNull {
			t.Fatal("bogus string should have returned TypeNull")
		}
	})
}

var benchData = []byte(entropy.RandStrWithUpper(5000))

func BenchmarkSum(b *testing.B) {
	for sumType := range valids {
		b.Run(sumType.String(), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.SetBytes(int64(len(benchData)))
				Sum(sumType, benchData)
			}
		})
	}
}
