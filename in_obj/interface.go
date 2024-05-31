package in_obj

import (
	"app/message/pb"
	"google.golang.org/protobuf/proto"
)

type IGameMgr interface {
	NewUser() IGameUser
	SdkLogin(appId string, token string) (string, string, string, pb.ERROR_CODE)
	SdkStartTask(appId string, roomId string, msgType string) pb.ERROR_CODE
	SdkStopTask(appId string, roomId string, msgType string) pb.ERROR_CODE
}

type IGameUser interface {
	OnSessionCreated(chan proto.Message)
	OnSessionClosed()
	OnRecv(proto.Message)
	Log(logType string, data ...interface{})
}
