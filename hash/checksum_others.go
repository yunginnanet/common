//go:build !amd64

package hash

func checksum(data []byte) uint16 {
	dataSize := len(data)
	var sum uint32

	for offset := 0; offset < dataSize-1; offset += 2 {
		r11 := uint16(data[offset]) << 8
		r8 := uint16(data[offset+1])
		r11 |= r8
		sum += uint32(r11)
		if sum > 0xFFFF { // 65535, max unsignd 16 bit integer
			sum = (sum & 0xFFFF) + (sum >> 16)
		}
	}

	if dataSize%2 != 0 {
		r8 := uint32(data[dataSize-1]) << 8
		sum += r8
		if sum > 0xFFFF {
			sum = (sum & 0xFFFF) + (sum >> 16)
		}
	}

	if sum > 0xFFFF {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	return ^uint16(sum)
}
