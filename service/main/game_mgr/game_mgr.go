package game_mgr

import (
	"app/douyin"
	"app/network"
	"app/service/main/message/pb"
	"app/service/main/msg_util"
	"app/service/main/user"
)

var APP_TOKEN_MAP = map[string]string{
	"debug": "default",
}

type GameMgr struct {
	appMap map[string]douyin.IApp
}

func NewGameMgr() *GameMgr {
	return &GameMgr{}
}

func (s *GameMgr) Run() error {
	s.appMap = make(map[string]douyin.IApp)

	//for appId, appSecret := range APP_TOKEN_MAP {
	//	newApp := douyin.NewApp(appId, appSecret)
	//	s.appMap[appId] = newApp
	//	newApp.StartRefreshToken()
	//}

	listener := network.NewTcpListener(
		":8899",
		s.newParser,
		s.newHandler)

	err := listener.Listen()
	return err
}

func (s *GameMgr) newParser() network.ICodec {
	return msg_util.NewProtoCodec(msg_util.NewProtoParser("message", "MSG_TYPE"), 1024*64, false)
}

func (s *GameMgr) newHandler() network.INetHandler {
	return user.NewUser(s)
}

func (s *GameMgr) Login(appId string, token string) (string, string, string, pb.ERROR_CODE) {
	app, ok := s.appMap[appId]
	if !ok {
		return "", "", "", pb.ERROR_CODE_INVALID_APPID
	}
	roomId, uid, nickName, err := app.GetRoomId(token)
	if err != nil {
		return "", "", "", pb.ERROR_CODE_INVALID_TOKEN
	}
	return roomId, uid, nickName, pb.ERROR_CODE_SUCCESS
}

func (s *GameMgr) StartTask(appId string, roomId string, msgType string) pb.ERROR_CODE {
	app, ok := s.appMap[appId]
	if !ok {
		return pb.ERROR_CODE_INVALID_APPID
	}

	_, err := app.StartTask(roomId, msgType)
	if err != nil {
		return pb.ERROR_CODE_FAIL
	}
	return pb.ERROR_CODE_SUCCESS
}

func (s *GameMgr) StopTask(appId string, roomId string, msgType string) pb.ERROR_CODE {
	app, ok := s.appMap[appId]
	if !ok {
		return pb.ERROR_CODE_INVALID_APPID
	}

	_, err := app.StopTask(roomId, msgType)
	if err != nil {
		return pb.ERROR_CODE_FAIL
	}
	return pb.ERROR_CODE_SUCCESS
}
