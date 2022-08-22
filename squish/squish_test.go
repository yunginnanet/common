package squish

import (
	"bytes"
	"encoding/base64"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
)

const lip string = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed a ante sit amet purus blandit auctor. Nullam ornare enim sed nibh consequat molestie. Duis est lectus, vestibulum vel felis vel, convallis cursus ex. Morbi nec placerat orci. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Praesent a erat sit amet libero convallis ornare a venenatis dolor. Pellentesque euismod risus et metus porttitor, vel consectetur lacus tempus. Integer elit arcu, condimentum quis nisi eget, dapibus imperdiet nulla. Cras sit amet ante in urna varius tempus. Integer tristique sagittis nunc vel tincidunt.

Integer non suscipit ligula, et fermentum sem. Duis id odio lorem. Sed id placerat urna, eu vehicula risus. Duis porttitor hendrerit risus. Curabitur id tellus ac arcu aliquet finibus. Pellentesque et nisl ante. Mauris sapien nisl, pretium in ligula tempus, posuere mattis turpis.

Proin et tempus enim. Nullam at diam est. Vivamus ut lectus hendrerit, interdum ex id, ultricies sapien. Praesent rhoncus turpis dolor, quis lobortis tortor pellentesque id. Pellentesque eget nisi laoreet, fringilla augue eu, cursus risus. Integer consectetur ornare laoreet. Praesent ligula sem, tincidunt at ligula at, condimentum venenatis tortor.

Nam laoreet enim leo, sed finibus lorem egestas vel. Maecenas varius a leo non placerat. Donec scelerisque, risus vel finibus ornare, arcu ligula interdum justo, in ultricies urna mi et neque. Curabitur sed sem dui. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Maecenas eget laoreet nisi. Nam rhoncus sapien ac interdum sagittis. Nulla fermentum sem nec tellus dignissim lacinia. Curabitur ornare lectus non dictum laoreet. Praesent tempor risus at tortor tempor finibus. Cras id dolor mi.

Mauris ut mi quis est vehicula molestie. Mauris eu varius urna. Integer sodales nunc at risus rutrum eleifend. In sed bibendum lectus. Morbi ipsum sapien, blandit in dignissim eu, ultrices non odio. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Nulla eget volutpat ligula, at elementum dui. Aliquam sed enim scelerisque, facilisis magna vitae, dignissim enim. Pellentesque non ultricies urna. Proin fermentum erat semper efficitur auctor. Vestibulum posuere non tortor vitae tincidunt. 
`

func TestGzip(t *testing.T) {
	gsUp := Gzip([]byte(lip))
	if bytes.Equal(gsUp, []byte(lip)) {
		t.Fatalf("[FAIL] Gzip didn't change the data at all despite being error free...")
	}
	if len(gsUp) == len([]byte(lip)) || len(gsUp) > len([]byte(lip)) {
		t.Fatalf("[FAIL] Gzip didn't change the sise of the data at all (or it grew)...")
	}
	if len(gsUp) == 0 {
		t.Fatalf("[FAIL] ended up with 0 bytes after compression...")
	}
	profit := len([]byte(lip)) - len(gsUp)
	t.Logf("[PASS] Gzip compress succeeded, squished %d bytes.", profit)
	hosDown, err := Gunzip(gsUp)
	if err != nil {
		t.Fatalf("Gzip decompression failed: %s", err.Error())
	}
	if !bytes.Equal(hosDown, []byte(lip)) {
		t.Fatalf("[FAIL] Gzip decompression failed, data does not appear to be the same after decompression")
	}
	if len(hosDown) != len([]byte(lip)) {
		t.Fatalf("[FAIL] Gzip decompression failed, data [%d] does not appear to be the same [%d] length after decompression", hosDown, len([]byte(lip)))
	}
	t.Logf("[PASS] Gzip decompress succeeded, restored %d bytes.", profit)
	_, err = Gunzip(nil)
	if err == nil {
		t.Fatalf("[FAIL] Gunzip didn't fail on nil input")
	}
}

func TestGunzipMustFails(t *testing.T) {
	blank := ""
	_, err := Gunzip([]byte(blank))
	if err == nil {
		t.Fatalf("[FAIL] Gunzip didn't fail on empty input")
	}
	_, err = UnpackStr(blank)
	if err == nil {
		t.Fatalf("[FAIL] UnpackStr didn't fail on empty input")
	}
	junk := "junk"
	_, err = Gunzip([]byte(junk))
	if err == nil {
		t.Fatalf("[FAIL] Gunzip didn't fail on junk input")
	}
	_, err = UnpackStr(junk)
	if err == nil {
		t.Fatalf("[FAIL] UnpackStr didn't fail on junk input")
	}
}

func TestGzipEntropic(t *testing.T) {
	for i := 0; i < 50; i++ {
		dat := []byte(entropy.RandStr(entropy.RNG(55) * 1024))
		for len(dat) < 1024 {
			dat = []byte(entropy.RandStr(entropy.RNG(55) * 1024))
		}
		gzTest(dat, t)
	}
}

func gzTest(dat []byte, t *testing.T) {
	t.Logf("Testing Gzip on %d bytes of data", len(dat))
	gsUp := Gzip(dat)
	if bytes.Equal(gsUp, dat) {
		t.Fatalf("[FAIL] Gzip didn't change the data at all despite being error free...")
	}
	if len(gsUp) == len(dat) || len(gsUp) > len(dat) {
		t.Fatalf("[FAIL] Gzip didn't change the sise of the data at all (or it grew)... before: %d  after: %d",
			len(dat), len(gsUp))
	}
	if len(gsUp) == 0 {
		t.Fatalf("[FAIL] ended up with 0 bytes after compression...")
	}
	profit := len(dat) - len(gsUp)
	t.Logf("[PASS] Gzip compress succeeded, squished %d bytes.", profit)
	hosDown, err := Gunzip(gsUp)
	if err != nil {
		t.Fatalf("Gzip decompression failed: %s", err.Error())
	}
	if !bytes.Equal(hosDown, dat) {
		t.Fatalf("[FAIL] Gzip decompression failed, data does not appear to be the same after decompression")
	}
	if len(hosDown) != len(dat) {
		t.Fatalf("[FAIL] Gzip decompression failed, data [%d] does not appear to be the same [%d] length after decompression", hosDown, len(dat))
	}
	t.Logf("[PASS] Gzip decompress succeeded, restored %d bytes.", profit)
}

func TestGzipDeterministic(t *testing.T) {
	packed := Gzip([]byte(lip))
	for n := 0; n < 10; n++ {
		again := Gzip([]byte(lip))
		if !bytes.Equal(again, packed) {
			t.Fatalf("[FAIL] Gzip is not deterministic")
		}
	}
}

func TestUnpackStr(t *testing.T) { //nolint:cyclop
	gzd := Gzip([]byte(lip))
	if len(gzd) == 0 {
		t.Fatalf("[FAIL] Gzip failed to compress data")
	}
	gzdSanity, gzdErr := Gunzip(gzd)
	if gzdErr != nil {
		t.Fatalf("Gzip failed: %s", gzdErr.Error())
	}
	if !bytes.Equal(gzdSanity, []byte(lip)) {
		t.Fatalf("Bytes not equal after Gzip: %v != %v", gzdSanity, []byte(lip))
	}
	packed := B64e(gzd)
	if len(packed) == 0 {
		t.Fatalf("[FAIL] B64e failed to encode data")
	}
	t.Logf("Packed: %s", packed)
	sanity1, err1 := base64.StdEncoding.DecodeString(packed)
	if err1 != nil {
		t.Fatalf("b64 failed: %s", err1.Error())
	}
	if !bytes.Equal(sanity1, gzd) {
		t.Fatalf("Bytes not equal after b64: %v != %v", sanity1, gzd)
	}
	sanity2, err2 := Gunzip(sanity1)
	if err2 != nil {
		t.Fatalf("Gzip failed: %s", err2.Error())
	}
	if !bytes.Equal(sanity2, []byte(lip)) {
		t.Fatalf("Bytes not equal after Gzip: %v != %v", sanity2, []byte(lip))
	}
	unpacked, err := UnpackStr(packed)
	switch {
	case err != nil:
		t.Errorf("[FAIL] %s", err.Error())
	case unpacked != lip:
		t.Fatalf("[FAIL] unpackstr decided to not work, who knows why. If you see this than I have already become a janitor.\n"+
			"unpacked: %s != packed: %s", unpacked, lip)
	default:
		t.Logf("[PASS] TestUnpackStr")
	}
	_, nilerr := UnpackStr("")
	if nilerr == nil {
		t.Fatalf("[FAIL] unpackstr didn't fail on empty input")
	}

}
