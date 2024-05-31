package msg_util

import (
	"app/tools"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"strings"
)

type ProtoParser struct {
	typemap  map[int32]protoreflect.MessageType // 此map初始化后，所有session都会并发读取
	msgNames map[string]int32                   // 此map初始化后，所有session都会并发读取
}

func NewProtoParser(packageName, msgEnum string) *ProtoParser {
	p := &ProtoParser{
		typemap:  make(map[int32]protoreflect.MessageType),
		msgNames: make(map[string]int32),
	}
	p.init(packageName, msgEnum)
	return p
}

// 协议enum->message自动解析:
// 1、不区分大小写
// 2、过滤下划线
func (s *ProtoParser) init(packageName, msgEnum string) {
	lowerNames := map[string]protoreflect.MessageType{}
	protoregistry.GlobalTypes.RangeMessages(func(messageType protoreflect.MessageType) bool {
		if !strings.HasPrefix(string(messageType.Descriptor().FullName()), packageName+".") {
			return true
		}
		name := convertMsgName(string(messageType.Descriptor().Name()))
		tools.AssertTrue(nil == lowerNames[name], "msg name repeated name=%s", messageType.Descriptor().FullName())
		lowerNames[name] = messageType
		return true
	})

	enums, err := protoregistry.GlobalTypes.FindEnumByName(protoreflect.FullName(packageName + "." + msgEnum))
	tools.AssertNil(err)

	values := enums.Descriptor().Values()
	adapter := adapters()
	msgDealFunc := func(msgTypeName string, msgId int32, isForce bool) {
		fullName := ""
		ln := convertMsgName(msgTypeName)
		if tp, ok := lowerNames[ln]; ok {
			s.typemap[msgId] = tp
			fullName = string(tp.Descriptor().FullName())
		} else if tp, ok := adapter[ln]; ok {
			s.typemap[msgId] = tp
			fullName = string(tp.Descriptor().FullName())
			delete(adapter, ln)
		} else {
			if isForce {
				panic(fmt.Errorf("msg format error msgTypeName=%v", msgTypeName))
			} else {
				return
			}
		}
		if _, ok := s.msgNames[fullName]; ok {
			panic(fmt.Sprintf("msg name repeated name=%s", fullName))
		}
		s.msgNames[fullName] = msgId
	}
	for i := 0; i < values.Len(); i++ {
		msgTypeName := string(values.Get(i).Name())
		msgId := int32(values.Get(i).Number())
		msgDealFunc(msgTypeName, msgId, true)
	}

	tools.AssertTrue(len(adapter) == 0, "msg not fix")
}

func (p *ProtoParser) Unmarshal(msgId int32, data []byte) (proto.Message, bool) {
	tp, ok := p.typemap[msgId]
	if !ok {
		return nil, false
	}

	msg := tp.New().Interface()
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, false
	}
	return msg, true
}
func (p *ProtoParser) UnmarshalByName(msgName string, data []byte) (proto.Message, bool) {
	msgId, ok := p.MsgNameToId(msgName)
	if !ok {
		return nil, false
	}
	return p.Unmarshal(msgId, data)
}

func (s *ProtoParser) MsgIdToName(msgId int32) (msgName string, ok bool) {
	if ptype, has := s.typemap[msgId]; has {
		return string(ptype.Descriptor().FullName()), true
	}
	return
}

func (s *ProtoParser) MsgNameToId(msgName string) (msgId int32, ok bool) {
	msgId, ok = s.msgNames[msgName]
	return
}

func (s *ProtoParser) MsgToId(msg proto.Message) (msgId int32, ok bool) {
	return s.MsgNameToId(ProtoFullName(msg))
}

func convertMsgName(msgName string) (name string) {
	words := strings.Split(msgName, "_")
	for _, word := range words {
		name += strings.ToLower(word)
	}
	return
}

func ProtoFullName(msg proto.Message) string {
	return string(msg.ProtoReflect().Type().Descriptor().FullName())
}

// 之前定义的消息id和消息体不适用自动解析的部分
func adapters() map[string]protoreflect.MessageType {
	return map[string]protoreflect.MessageType{}
}
