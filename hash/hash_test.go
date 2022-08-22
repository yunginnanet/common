package hash

import (
	"bytes"
	"encoding/base64"
	"os"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
	"git.tcp.direct/kayos/common/squish"
)

const (
	kayosBlake2b = "Kr+6ONDx+cq/WvhHpQE/4LVuJYi9QHz1TztHNTWwa9KJWqHxfTNLKF3YxrcLptA3wO0KHm83Lq7gpBWgCQzPag=="
	kayosMD5     = "aoNjwileNitB208vOwpIow=="
	kayosSHA1    = "M23ElC0sQYAK+MaMZVmza2L8mss="
	kayosSHA256  = "BagY0TmoGR3O7t80BGm4K6UHPlqEg6HJirwQmhrPK4U="
	kayosSHA512  = "xiuo2na76acrWXCTTR++O1pPZabOhyj8nbfb5Go3e1pEq9VJYIsOioTXalf2GCuERmFecWkmaL5QI8mIXXWpNA=="
)

func TestBlake2bsum(t *testing.T) {
	og := squish.B64d(kayosBlake2b)
	newc := Blake2bSum([]byte("kayos\n"))
	if !bytes.Equal(newc, og) {
		t.Fatalf("wanted: %v, got %v", kayosBlake2b, squish.B64e(newc))
	}
	if !BlakeEqual([]byte("kayos\n"), []byte{107, 97, 121, 111, 115, 10}) {
		t.Fatalf("BlakeEqual should have been true. %s should == %s", []byte("kayos\n"), []byte{107, 97, 121, 111, 115, 92, 110})
	}
	t.Logf("[blake2bSum] success: %s", kayosBlake2b)
}

func TestBlakeFileChecksum(t *testing.T) {
	path := t.TempDir() + "/blake2b.dat"
	err := os.WriteFile(path, []byte{107, 97, 121, 111, 115, 10}, os.ModePerm)
	if err != nil {
		t.Errorf("[FAIL] failed to write test fle for TestBlakeFileChecksum: %s", err.Error())
	}
	filecheck, err2 := BlakeFileChecksum(path)
	if err2 != nil {
		t.Errorf("[FAIL] failed to read test fle for TestBlakeFileChecksum: %s", err2.Error())
	}
	if len(filecheck) == 0 {
		t.Errorf("[FAIL] got nil output from BlakeFileChecksum")
	}
	if !bytes.Equal(filecheck, squish.B64d(kayosBlake2b)) {
		t.Fatalf("[FAIL] wanted: %v, got %v", kayosBlake2b, squish.B64e(filecheck))
	}
	badfile, err3 := BlakeFileChecksum(t.TempDir() + "/" + entropy.RandStr(50))
	if err3 == nil {
		t.Errorf("[FAIL] shouldn't have been able to read phony file")
	}
	if len(badfile) != 0 {
		t.Errorf("[FAIL] got non-nil output from bogus file: %v", badfile)
	}
	if !bytes.Equal(filecheck, squish.B64d(kayosBlake2b)) {
		t.Fatalf("[FAIL] wanted: %v, got %v", kayosBlake2b, squish.B64e(filecheck))
	}
	err = os.WriteFile(path+".empty", []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("[FAIL] failed to write test fle for TestBlakeFileChecksum: %s", err.Error())
	}
	_, err4 := BlakeFileChecksum(path + ".empty")
	if err4 == nil {
		t.Fatalf("[FAIL] should have failed to read empty file")
	}
}

func TestSum(t *testing.T) {
	if Sum(TypeNull, []byte("yeet")) != nil {
		t.Fatal("Sum(TypeNull, []byte(\"yeet\")) should have returned nil")
	}

	var (
		ogsha1, _   = base64.StdEncoding.DecodeString(kayosSHA1)
		ogsha256, _ = base64.StdEncoding.DecodeString(kayosSHA256)
		ogsha512, _ = base64.StdEncoding.DecodeString(kayosSHA512)
		ogmd5, _    = base64.StdEncoding.DecodeString(kayosMD5)
		newsha1     = Sum(TypeSHA1, []byte("kayos\n"))
		newsha256   = Sum(TypeSHA256, []byte("kayos\n"))
		newsha512   = Sum(TypeSHA512, []byte("kayos\n"))
		newmd5      = Sum(TypeMD5, []byte("kayos\n"))
	)

	if !bytes.Equal(newsha1, ogsha1) {
		t.Fatalf("[sha1] wanted: %v, got %v", ogsha1, newsha1)
	}
	t.Logf("[sha1]   success: %s", kayosSHA1)

	if !bytes.Equal(newsha256, ogsha256) {
		t.Fatalf("[sha256] wanted: %v, got %v", ogsha256, newsha256)
	}
	t.Logf("[sha256] success: %s", kayosSHA256)

	if !bytes.Equal(newsha512, ogsha512) {
		t.Fatalf("[sha512] wanted: %v, got %v", ogsha512, newsha512)
	}
	t.Logf("[sha512] success: %s", kayosSHA512)

	if !bytes.Equal(newmd5, ogmd5) {
		t.Fatalf("[md5] wanted: %v, got %v", ogmd5, newmd5)
	}
	t.Logf("[md5]    success: %s", kayosMD5)
}
