package xerrors

import (
	"encoding/gob"

	"git.tcp.direct/kayos/common/pool"
)

type Codec struct {
	Encoder    gob.GobEncoder
	Decoder    gob.GobDecoder
	BufferPool pool.Pool[pool.ByteBuffer]
}
