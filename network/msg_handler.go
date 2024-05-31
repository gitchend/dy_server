package network

import (
	"app/tools"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"reflect"
	"sync"
)

/*
按消息/请求的数据类型注册回调
注意:不是区分具体类型的值
*/

func NewMsgTypeHandler() *MsgTypeHandler {
	return &MsgTypeHandler{msgMap: make(map[string]MsgFun)}
}

type MsgFun func(sourceId, msg interface{})

type MsgTypeHandler struct {
	msgMap map[string]MsgFun
}

// 注册消息回调
func (h *MsgTypeHandler) RegistMsg(msg interface{}, f MsgFun) {
	msgName := MsgName(msg)
	_, ok := h.msgMap[msgName]
	tools.AssertTrue(!ok, "regist repeated msg=%v", msgName)
	h.msgMap[msgName] = f
}

func (h *MsgTypeHandler) HandleMsg(sourceId, msg interface{}) bool {
	msgName := MsgName(msg)
	handler := h.msgMap[msgName]
	tools.AssertTrue(handler != nil, "msg=%v not regist handler", msgName)
	handler(sourceId, msg)
	return true
}

var _msgName sync.Map //
func MsgName(msg interface{}) string {
	if msg == nil {
		return "nil"
	}

	if str, ok := msg.(string); ok {
		return str
	}

	tp := reflect.TypeOf(msg)
	//if msgName, ok := _msgName.Load(tp); ok {
	//	return msgName.(string)
	//}

	var msgName string
	if pb, ok := msg.(proto.Message); ok {
		msgName = string(pb.ProtoReflect().Type().Descriptor().FullName())
	} else if tp.Kind() == reflect.Ptr {
		msgName = tp.Elem().PkgPath() + "/" + tp.Elem().Name()
	} else {
		msgName = tp.PkgPath() + "/" + tp.Name()
	}

	_msgName.Store(tp, msgName)
	return msgName
}
func FindMsgByName(msgName string) (protoreflect.MessageType, error) {
	return protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(msgName))
}
