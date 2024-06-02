package hash

import "testing"

func TestChecksum(t *testing.T) {
	data := []byte("hello")
	expected := uint16(48173)
	actual := checksum(data)
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}
