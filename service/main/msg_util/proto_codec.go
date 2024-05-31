package msg_util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"google.golang.org/protobuf/proto"

	"app/network"
	"app/tools"
)

func NewProtoCodec(parser *ProtoParser, maxDecode int, isRaw bool) *ProtoCodec {
	return &ProtoCodec{
		parser:    parser,
		maxDecode: maxDecode,
		isRaw:     isRaw,
	}
}

type ProtoCodec struct {
	maxDecode int
	msglen    uint32
	context   bytes.Buffer
	parser    *ProtoParser
	isRaw     bool //是否需要解析
}

const STREAM_HEADLEN = 4
const STREAM_MSGID_LEN = 2 //uint16

var ErrUnmarshal = errors.New("message unmarshal failed")

func (s *ProtoCodec) Decode(data []byte) (msgIds []int32, ret []interface{}, err error) {
	s.context.Write(data)

	for s.context.Len() >= STREAM_HEADLEN+STREAM_MSGID_LEN {
		if s.msglen == 0 {
			d := s.context.Bytes()
			s.msglen = binary.BigEndian.Uint32(d[:STREAM_HEADLEN]) - STREAM_HEADLEN //客户端headlen也算入长度
			if s.msglen < STREAM_MSGID_LEN {
				err = errors.New("data is too small")
				return
			}
			if s.maxDecode > 0 && int(s.msglen) > s.maxDecode {
				err = network.ErrRecvLen
				return
			}
		}

		if int(s.msglen)+STREAM_HEADLEN > s.context.Len() {
			break
		}

		d := make([]byte, s.msglen+STREAM_HEADLEN)
		if n, err := s.context.Read(d); n != int(s.msglen)+STREAM_HEADLEN || err != nil {
			s.msglen = 0
			continue
		}

		msgId := int32(binary.BigEndian.Uint16(d[STREAM_HEADLEN : STREAM_HEADLEN+STREAM_MSGID_LEN]))
		var msg interface{}
		if s.isRaw {
			msg = d[STREAM_HEADLEN+STREAM_MSGID_LEN:]
		} else {
			var ok bool
			msg, ok = s.parser.Unmarshal(msgId, d[STREAM_HEADLEN+STREAM_MSGID_LEN:])
			if !ok {
				err = ErrUnmarshal
				return
			}
		}

		s.msglen = 0
		ret = append(ret, msg)
		msgIds = append(msgIds, msgId)
	}
	return
}

func (s *ProtoCodec) Encode(data interface{}) []byte {
	if pb, ok := data.(proto.Message); ok {
		msgId, ok := s.parser.MsgToId(pb)
		tools.AssertTrue(ok, "msgId=%v not found", msgId, pb)

		pbBuf, err := proto.Marshal(pb)
		tools.AssertNil(err)
		buf := make([]byte, STREAM_HEADLEN+STREAM_MSGID_LEN+len(pbBuf))
		binary.BigEndian.PutUint32(buf, uint32(len(buf)))
		binary.BigEndian.PutUint16(buf[STREAM_HEADLEN:], uint16(msgId))
		copy(buf[STREAM_HEADLEN+STREAM_MSGID_LEN:], pbBuf)
		return buf
	} else if byteArray, ok := data.([]byte); ok {
		return byteArray
	} else {
		return nil
	}
}
