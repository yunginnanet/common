package hash

import (
	"bytes"
	"os"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
	"git.tcp.direct/kayos/common/squish"
)

const kayos = "Kr+6ONDx+cq/WvhHpQE/4LVuJYi9QHz1TztHNTWwa9KJWqHxfTNLKF3YxrcLptA3wO0KHm83Lq7gpBWgCQzPag=="

func TestBlake2bsum(t *testing.T) {
	og := squish.B64d(kayos)
	newc := Blake2bSum([]byte("kayos\n"))
	if !bytes.Equal(newc, og) {
		t.Fatalf("wanted: %v, got %v", kayos, squish.B64e(newc))
	}
	if !BlakeEqual([]byte("kayos\n"), []byte{107, 97, 121, 111, 115, 10}) {
		t.Fatalf("BlakeEqual should have been true. %s should == %s", []byte("kayos\n"), []byte{107, 97, 121, 111, 115, 92, 110})
	}
}

func TestBlakeFileChecksum(t *testing.T) {
	path := t.TempDir() + "/blake2b.dat"
	err := os.WriteFile(path, []byte{107, 97, 121, 111, 115, 10}, os.ModePerm)
	if err != nil {
		t.Errorf("failed to write test fle for TestBlakeFileChecksum: %s", err.Error())
	}
	filecheck, err2 := BlakeFileChecksum(path)
	if err2 != nil {
		t.Errorf("failed to read test fle for TestBlakeFileChecksum: %s", err2.Error())
	}
	if len(filecheck) == 0 {
		t.Errorf("Got nil output from BlakeFileChecksum")
	}
	if !bytes.Equal(filecheck, squish.B64d(kayos)) {
		t.Fatalf("wanted: %v, got %v", kayos, squish.B64e(filecheck))
	}
	badfile, err3 := BlakeFileChecksum(t.TempDir() + "/" + entropy.RandStr(50))
	if err3 == nil {
		t.Errorf("shouldn't have been able to read phony file")
	}
	if len(badfile) != 0 {
		t.Errorf("Got non-nil output from bogus file: %v", badfile)
	}
	if !bytes.Equal(filecheck, squish.B64d(kayos)) {
		t.Fatalf("wanted: %v, got %v", kayos, squish.B64e(filecheck))
	}
}
