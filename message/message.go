package message

import (
	"app/tools"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"strings"
)

var (
	typeMap  map[int32]protoreflect.MessageType // 此map初始化后，所有session都会并发读取
	msgNames map[string]int32                   // 此map初始化后，所有session都会并发读取
)

// 协议enum->message自动解析:
// 1、不区分大小写
// 2、过滤下划线

func InitMessageParser(packageName, msgEnum string) {
	typeMap = make(map[int32]protoreflect.MessageType)
	msgNames = make(map[string]int32)
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
			typeMap[msgId] = tp
			fullName = string(tp.Descriptor().FullName())
		} else if tp, ok := adapter[ln]; ok {
			typeMap[msgId] = tp
			fullName = string(tp.Descriptor().FullName())
			delete(adapter, ln)
		} else {
			if isForce {
				panic(fmt.Errorf("msg format error msgTypeName=%v", msgTypeName))
			} else {
				return
			}
		}
		if _, ok := msgNames[fullName]; ok {
			panic(fmt.Sprintf("msg name repeated name=%s", fullName))
		}
		msgNames[fullName] = msgId
	}
	for i := 0; i < values.Len(); i++ {
		msgTypeName := string(values.Get(i).Name())
		msgId := int32(values.Get(i).Number())
		msgDealFunc(msgTypeName, msgId, true)
	}

	tools.AssertTrue(len(adapter) == 0, "msg not fix")
}
func Marshal(msg proto.Message) (int32, []byte, bool) {
	msgId, ok := MsgToId(msg)
	if !ok {
		return 0, nil, false
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return 0, nil, false
	}
	return msgId, data, true
}

func Unmarshal(msgId int32, data []byte) (proto.Message, bool) {
	tp, ok := typeMap[msgId]
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
func MsgNameToId(msgName string) (msgId int32, ok bool) {
	msgId, ok = msgNames[msgName]
	return
}

func MsgToId(msg proto.Message) (msgId int32, ok bool) {
	return MsgNameToId(MsgName(msg))
}

func convertMsgName(msgName string) (name string) {
	words := strings.Split(msgName, "_")
	for _, word := range words {
		name += strings.ToLower(word)
	}
	return
}

func MsgName(msg proto.Message) string {
	return string(msg.ProtoReflect().Type().Descriptor().FullName())
}

// 之前定义的消息id和消息体不适用自动解析的部分
func adapters() map[string]protoreflect.MessageType {
	return map[string]protoreflect.MessageType{}
}

func ReadMessage(conn *websocket.Conn) (proto.Message, error) {
	_, dataId, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	msgId := binary.BigEndian.Uint32(dataId)
	_, dataMsg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	msg, ok := Unmarshal(int32(msgId), dataMsg)
	if !ok {
		return nil, errors.New("unmarshal err")
	}
	return msg, nil
}

func SendMessage(conn *websocket.Conn, msg proto.Message) error {
	msgId, data, ok := Marshal(msg)
	if !ok {
		return errors.New("marshal err")
	}
	dataId := make([]byte, 4)
	binary.BigEndian.PutUint32(dataId, uint32(msgId))
	if err := conn.WriteMessage(websocket.BinaryMessage, dataId); err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return err
	}
	return nil
}
