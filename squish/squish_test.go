package squish

import (
	"bytes"
	"testing"
)

const lip string = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed a ante sit amet purus blandit auctor. Nullam ornare enim sed nibh consequat molestie. Duis est lectus, vestibulum vel felis vel, convallis cursus ex. Morbi nec placerat orci. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Praesent a erat sit amet libero convallis ornare a venenatis dolor. Pellentesque euismod risus et metus porttitor, vel consectetur lacus tempus. Integer elit arcu, condimentum quis nisi eget, dapibus imperdiet nulla. Cras sit amet ante in urna varius tempus. Integer tristique sagittis nunc vel tincidunt.

Integer non suscipit ligula, et fermentum sem. Duis id odio lorem. Sed id placerat urna, eu vehicula risus. Duis porttitor hendrerit risus. Curabitur id tellus ac arcu aliquet finibus. Pellentesque et nisl ante. Mauris sapien nisl, pretium in ligula tempus, posuere mattis turpis.

Proin et tempus enim. Nullam at diam est. Vivamus ut lectus hendrerit, interdum ex id, ultricies sapien. Praesent rhoncus turpis dolor, quis lobortis tortor pellentesque id. Pellentesque eget nisi laoreet, fringilla augue eu, cursus risus. Integer consectetur ornare laoreet. Praesent ligula sem, tincidunt at ligula at, condimentum venenatis tortor.

Nam laoreet enim leo, sed finibus lorem egestas vel. Maecenas varius a leo non placerat. Donec scelerisque, risus vel finibus ornare, arcu ligula interdum justo, in ultricies urna mi et neque. Curabitur sed sem dui. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Maecenas eget laoreet nisi. Nam rhoncus sapien ac interdum sagittis. Nulla fermentum sem nec tellus dignissim lacinia. Curabitur ornare lectus non dictum laoreet. Praesent tempor risus at tortor tempor finibus. Cras id dolor mi.

Mauris ut mi quis est vehicula molestie. Mauris eu varius urna. Integer sodales nunc at risus rutrum eleifend. In sed bibendum lectus. Morbi ipsum sapien, blandit in dignissim eu, ultrices non odio. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Nulla eget volutpat ligula, at elementum dui. Aliquam sed enim scelerisque, facilisis magna vitae, dignissim enim. Pellentesque non ultricies urna. Proin fermentum erat semper efficitur auctor. Vestibulum posuere non tortor vitae tincidunt. 
`

func TestGzip(t *testing.T) {
	gsUp, err := Gzip([]byte(lip))
	if err != nil {
		t.Fatalf("Gzip compression failed: %e", err)
	}

	if bytes.Equal(gsUp, []byte(lip)) {
		t.Fatalf("Gzip didn't change the data at all despite being error free...")
	}

	if len(gsUp) == len([]byte(lip)) {
		t.Fatalf("Gzip didn't change the sise of the data at all despite being error free...")
	}

	profit := len([]byte(lip)) - len(gsUp)
	t.Logf("[PASS] Gzip compress succeeded, squished %d bytes.", profit)

	hosDown, err := Gunzip(gsUp)

	if err != nil {
		t.Fatalf("Gzip decompression failed: %e", err)
	}

	if !bytes.Equal(hosDown, []byte(lip)) {
		t.Fatalf("Gzip decompression failed, data does not appear to be the same after decompression")
	}

	if len(hosDown) != len([]byte(lip)) {
		t.Fatalf("Gzip decompression failed, data [%d] does not appear to be the same [%d] length after decompression", hosDown, len([]byte(lip)))
	}

	t.Logf("[PASS] Gzip decompress succeeded, restored %d bytes.", profit)
}
