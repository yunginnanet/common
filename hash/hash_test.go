package hash

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
)

const (
	kayosBlake2b      = "Kr+6ONDx+cq/WvhHpQE/4LVuJYi9QHz1TztHNTWwa9KJWqHxfTNLKF3YxrcLptA3wO0KHm83Lq7gpBWgCQzPag=="
	kayosBlake2bHex   = "2abfba38d0f1f9cabf5af847a5013fe0b56e2588bd407cf54f3b473535b06bd2895aa1f17d334b285dd8c6b70ba6d037c0ed0a1e6f372eaee0a415a0090ccf6a"
	kayosMD5          = "aoNjwileNitB208vOwpIow=="
	kayosMD5Hex       = "6a8363c2295e362b41db4f2f3b0a48a3"
	kayosSHA1         = "M23ElC0sQYAK+MaMZVmza2L8mss="
	kayosSHA1Hex      = "336dc4942d2c41800af8c68c6559b36b62fc9acb"
	kayosSHA256       = "BagY0TmoGR3O7t80BGm4K6UHPlqEg6HJirwQmhrPK4U="
	kayosSHA256Hex    = "05a818d139a8191dceeedf340469b82ba5073e5a8483a1c98abc109a1acf2b85"
	kayosSHA512       = "xiuo2na76acrWXCTTR++O1pPZabOhyj8nbfb5Go3e1pEq9VJYIsOioTXalf2GCuERmFecWkmaL5QI8mIXXWpNA=="
	kayosSHA512Hex    = "c62ba8da76bbe9a72b5970934d1fbe3b5a4f65a6ce8728fc9db7dbe46a377b5a44abd549608b0e8a84d76a57f6182b8446615e71692668be5023c9885d75a934"
	kayosCRC32        = "xtig5w=="
	kayosCRC32Hex     = "c6d8a0e7"
	kayosCRC64ISO     = "YVx8IpQawAA="
	kayosCRC64ISOHex  = "615c7c22941ac000"
	kayosCRC64ECMA    = "Nn5+vneo4j4="
	kayosCRC64ECMAHex = "367e7ebe77a8e23e"
)

var kayosByteSlice = []byte{107, 97, 121, 111, 115, 10}

var (
	ogsha1, _      = base64.StdEncoding.DecodeString(kayosSHA1)
	ogsha256, _    = base64.StdEncoding.DecodeString(kayosSHA256)
	ogsha512, _    = base64.StdEncoding.DecodeString(kayosSHA512)
	ogmd5, _       = base64.StdEncoding.DecodeString(kayosMD5)
	ogBlake2b, _   = base64.StdEncoding.DecodeString(kayosBlake2b)
	ogCRC32, _     = base64.StdEncoding.DecodeString(kayosCRC32)
	ogCRC64ISO, _  = base64.StdEncoding.DecodeString(kayosCRC64ISO)
	ogCRC64ECMA, _ = base64.StdEncoding.DecodeString(kayosCRC64ECMA)
	valids         = map[Type][]byte{
		TypeSHA1:      ogsha1,
		TypeSHA256:    ogsha256,
		TypeSHA512:    ogsha512,
		TypeMD5:       ogmd5,
		TypeCRC32:     ogCRC32,
		TypeCRC64ISO:  ogCRC64ISO,
		TypeCRC64ECMA: ogCRC64ECMA,
		TypeBlake2b:   ogBlake2b,
	}
	validHexes = map[Type]string{
		TypeSHA1:      kayosSHA1Hex,
		TypeSHA256:    kayosSHA256Hex,
		TypeSHA512:    kayosSHA512Hex,
		TypeMD5:       kayosMD5Hex,
		TypeCRC32:     kayosCRC32Hex,
		TypeCRC64ISO:  kayosCRC64ISOHex,
		TypeCRC64ECMA: kayosCRC64ECMAHex,
		TypeBlake2b:   kayosBlake2bHex,
	}
)

func init() {
	for k, v := range valids {
		if len(v) == 0 {
			panic("invalid test data")
		}
		if hexBytes, hexErr := hex.DecodeString(validHexes[k]); hexErr != nil {
			panic("invalid test data")
		} else if !bytes.Equal(v, hexBytes) {
			panic("invalid test data")
		}
	}
}

func TestSum(t *testing.T) {
	t.Parallel()
	if Sum(TypeNull, []byte("yeet")) != nil {
		t.Fatal("Sum(TypeNull, []byte(\"yeet\")) should have returned nil")
	}

	for k, v := range valids {
		typeToTest := k
		valueToTest := v
		t.Run(typeToTest.String()+"/enum_check", func(t *testing.T) {
			t.Parallel()
			if !strings.EqualFold(typeToTest.String(), StringToType(typeToTest.String()).String()) {
				t.Errorf("[FAIL] %s: wanted %s, got %s",
					typeToTest.String(), typeToTest.String(), StringToType(typeToTest.String()).String(),
				)
			}
		})
		t.Run(typeToTest.String()+"/static_check", func(t *testing.T) {
			t.Parallel()
			mySum := Sum(typeToTest, kayosByteSlice)
			if !bytes.Equal(mySum, valueToTest) {
				t.Errorf("[FAIL] %s: wanted %v, got %v", typeToTest.String(), valueToTest, mySum)
			}
			RecycleChecksum(mySum)
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

func TestSumHex(t *testing.T) {
	t.Parallel()
	if SumHex(TypeNull, []byte("yeet")) != "" {
		t.Fatal("SumHex(TypeNull, []byte(\"yeet\")) should have returned an empty string")
	}

	for k, v := range validHexes {
		typeToTest := k
		valueToTest := v
		t.Run(typeToTest.String()+"/static_check", func(t *testing.T) {
			t.Parallel()
			mySum := SumHex(typeToTest, kayosByteSlice)
			if mySum != valueToTest {
				t.Errorf("[FAIL] %s: wanted %v, got %v", typeToTest.String(), valueToTest, mySum)
			}
		})
	}
}

func TestSumFile(t *testing.T) {
	t.Parallel()
	if _, err := SumFile(TypeNull, "/dev/null"); err == nil {
		t.Fatal("SumFile(TypeNull, \"/dev/null\") should have returned an error")
	}
	if _, err := SumFile(TypeNull, entropy.RandStrWithUpper(500)); err == nil {
		t.Fatal("SumFile(TypeNull, entropy.RandStrWithUpper(500)) should have returned an error")
	}
}

func BenchmarkSumHex(b *testing.B) {
	runIt := func(length int) {
		dat := []byte(entropy.RandStrWithUpper(length))
		b.Run(strconv.Itoa(length)+"char", func(b *testing.B) {
			for sumType := range valids {
				b.Run(sumType.String(), func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						b.SetBytes(int64(len(dat)))
						SumHex(sumType, dat)
					}
				})
			}
		})
	}
	for i := 0; i != 5; i++ {
		mult := 5
		if i == 0 {
			runIt(50)
			continue
		}
		if i > 1 {
			mult = (mult * 10) * i
		}
		if i > 3 {
			mult = (mult * 100) * i
		}
		runIt(i * mult)
	}

}

func BenchmarkSum(b *testing.B) {
	runIt := func(length int) {
		dat := []byte(entropy.RandStrWithUpper(length))
		b.Run(strconv.Itoa(length)+"char", func(b *testing.B) {
			for sumType := range valids {
				b.Run(sumType.String(), func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						b.SetBytes(int64(len(dat)))
						Sum(sumType, dat)
					}
				})
			}
		})
	}
	for i := 0; i != 5; i++ {
		mult := 5
		if i == 0 {
			runIt(50)
			continue
		}
		if i > 1 {
			mult = (mult * 10) * i
		}
		if i > 3 {
			mult = (mult * 100) * i
		}
		runIt(i * mult)
	}
}
