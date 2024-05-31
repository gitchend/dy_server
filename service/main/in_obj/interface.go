package in_obj

import "app/service/main/message/pb"

type IGameMgr interface {
	Login(appId string, token string) (string, string, string, pb.ERROR_CODE)
	StartTask(appId string, roomId string, msgType string) pb.ERROR_CODE
	StopTask(appId string, roomId string, msgType string) pb.ERROR_CODE
}
