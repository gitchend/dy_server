package game_mgr

import (
	"app/game_user"
	"app/in_obj"
	"app/message"
	"app/message/pb"
	"app/redis"
	"app/sdk/douyin"
	"encoding/json"
	"os"
)

var APP_TOKEN_MAP = map[string]string{
	"debug": "default",
}

type GameMgr struct {
	appMap  map[string]douyin.IApp
	httpMgr *HttpMgr
}

func NewGameMgr() *GameMgr {
	appInfo := os.Getenv("APP_INFO")
	if appInfo != "" {
		appInfoMap := make(map[string]string)
		err := json.Unmarshal([]byte(appInfo), &appInfoMap)
		if err == nil {
			for appId, appSecret := range appInfoMap {
				APP_TOKEN_MAP[appId] = appSecret
			}
		}
	}
	return &GameMgr{}
}
func (s *GameMgr) Run() error {
	s.appMap = make(map[string]douyin.IApp)
	s.httpMgr = NewHttpMgr(s)

	message.InitMessageParser("message", "MSG_TYPE")
	redis.Init()

	for appId, appSecret := range APP_TOKEN_MAP {
		newApp := douyin.NewApp(appId, appSecret)
		s.appMap[appId] = newApp
	}

	return s.httpMgr.Run()
}

func (s *GameMgr) NewUser() in_obj.IGameUser {
	return game_user.NewUser(s)
}

func (s *GameMgr) SdkLogin(appId string, token string) (string, string, string, pb.ERROR_CODE) {
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

func (s *GameMgr) SdkStartTask(appId string, roomId string, msgType string) pb.ERROR_CODE {
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

func (s *GameMgr) SdkStopTask(appId string, roomId string, msgType string) pb.ERROR_CODE {
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
