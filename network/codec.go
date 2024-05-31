package network

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type ICodec interface {
	Decode([]byte) ([]int32, []interface{}, error)
	Encode(interface{}) []byte
}

type CodecChain func(l int, paloay interface{})

var ErrRecvLen = errors.New("data is too long")

type StreamCodec struct {
	MaxDecode int
	msglen    uint32
	context   bytes.Buffer
}

const STREAM_HEADLEN = 4

func (s *StreamCodec) Decode(data []byte) ([]int32, []interface{}, error) {
	s.context.Write(data)

	var ret []interface{} = nil
	for s.context.Len() >= STREAM_HEADLEN {
		if s.msglen == 0 {
			d := s.context.Bytes()
			s.msglen = binary.BigEndian.Uint32(d[:STREAM_HEADLEN])
			if s.MaxDecode > 0 && int(s.msglen) > s.MaxDecode {
				return nil, nil, ErrRecvLen
			}
		}

		if int(s.msglen)+STREAM_HEADLEN > s.context.Len() {
			break
		}

		d := make([]byte, s.msglen+STREAM_HEADLEN)
		n, err := s.context.Read(d)
		if n != int(s.msglen)+STREAM_HEADLEN || err != nil {
			s.msglen = 0
			continue
		}
		s.msglen = 0

		each := d[STREAM_HEADLEN:]
		ret = append(ret, each)
	}
	return nil, ret, nil
}

func (s *StreamCodec) Encode(data interface{}) []byte {
	d := make([]byte, STREAM_HEADLEN)
	binary.BigEndian.PutUint32(d, uint32(len(data.([]byte))))
	return append(d, data.([]byte)...)
}
