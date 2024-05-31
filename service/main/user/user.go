package user

import (
	"app/network"
	"app/service/main/in_obj"
	"app/service/main/message/pb"
	"app/service/main/redis"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func NewUser(mgr in_obj.IGameMgr) *User {
	return &User{mgr: mgr}
}

type User struct {
	network.BaseNetHandler
	mgr           in_obj.IGameMgr
	appId         string
	roomId        string
	uid           string
	nickName      string
	isDebug       bool
	isRunning     bool
	audienceSet   sync.Map
	msgChan       chan func()
	closeChan     chan struct{}
	pushCloseChan chan struct{}
}

func (s *User) OnSessionCreated() {
	fmt.Println("[SessionCreated]")
	//主协程
	s.msgChan = make(chan func())
	s.closeChan = make(chan struct{})
	s.pushCloseChan = make(chan struct{})
	s.audienceSet = sync.Map{}
	go func() {
		for {
			select {
			case f := <-s.msgChan:
				f()
			case <-s.closeChan:
				return
			}
		}
	}()
}

func (s *User) OnSessionClosed() {
	fmt.Println("[SessionClosed]")
	if !s.isRunning {
		return
	}
	s.setPushActive(false)
	s.isRunning = false
	s.pushCloseChan <- struct{}{}
	s.closeChan <- struct{}{}
}

func (s *User) OnRecv(msgId int32, data interface{}) {
	if msgId == int32(pb.MSG_TYPE__Ping) {
		s.SendMsg(&pb.Pong{
			ServerTime: time.Now().Unix(),
			ClientTime: data.(*pb.Ping).ClientTime,
		})
		return
	}
	s.Log("Recv", data)
	s.msgChan <- func() {
		switch msgId {
		case int32(pb.MSG_TYPE__Login):
			s.OnLogin(data.(*pb.Login))
		case int32(pb.MSG_TYPE__PlayStart):
			s.OnPlayStart(data.(*pb.PlayStart))
		case int32(pb.MSG_TYPE__PlayEnd):
			s.OnPlayEnd(data.(*pb.PlayEnd))
		case int32(pb.MSG_TYPE__Report):
			s.OnReport(data.(*pb.Report))
		case int32(pb.MSG_TYPE__GetRank):
			s.OnGetRank(data.(*pb.GetRank))
		}
	}
}

func (s *User) OnRecvPush(data string) {
	msg := &pb.NotifyAudienceAction{}
	err := json.Unmarshal([]byte(data), msg)
	if err != nil {
		return
	}
	if _, ok := s.audienceSet.Load(msg.OpenId); !ok {
		audienceData := redis.GetAudience(s.appId, msg.OpenId)
		if audienceData != nil {
			s.audienceSet.Store(msg.OpenId, 1)
			audienceMsg := &pb.NotifyNewAudience{
				Audience: audienceData,
			}
			s.msgChan <- func() {
				s.SendUserMsg(audienceMsg)
			}
		}
	}
	s.msgChan <- func() {
		s.SendUserMsg(msg)
	}
}

func (s *User) SendUserMsg(msg interface{}) {
	s.Log("Send", msg)
	s.SendMsg(msg)
}
func (s *User) Log(logType string, msg interface{}) {
	fmt.Printf("[%s][%s](%s,%s,%s)%v\n", logType, network.MsgName(msg), s.roomId, s.uid, s.nickName, msg)
}
func (s *User) OnLogin(msg *pb.Login) {
	var (
		roomId   string
		uid      string
		nickName string
		result   pb.ERROR_CODE
	)
	if !msg.IsDebug {
		roomId, uid, nickName, result = s.mgr.Login(msg.AppId, msg.Token)
	} else {
		roomId, uid, nickName, result = "1111111111111111111", "debug", "debug", pb.ERROR_CODE_SUCCESS
	}
	s.appId = msg.AppId
	s.roomId = roomId
	s.uid = uid
	s.nickName = nickName
	s.isDebug = msg.IsDebug
	s.SendUserMsg(&pb.LoginResult{
		Result:   result,
		RoomId:   roomId,
		UID:      uid,
		NickName: nickName,
	})

	//test
	s.OnPlayStart(nil)
}
func (s *User) OnPlayStart(msg *pb.PlayStart) {
	if s.isRunning {
		s.SendUserMsg(&pb.PlayStartResult{Result: pb.ERROR_CODE_GAME_IS_RUNNING})
		return
	}
	s.isRunning = true
	s.setPushActive(true)
	s.SendUserMsg(&pb.PlayStartResult{Result: pb.ERROR_CODE_SUCCESS})
}
func (s *User) OnPlayEnd(msg *pb.PlayEnd) {
	if !s.isRunning {
		s.SendUserMsg(&pb.PlayEndResult{Result: pb.ERROR_CODE_GAME_IS_STOPPED})
		return
	}
	s.isRunning = false
	s.setPushActive(false)
	s.SendUserMsg(&pb.PlayEndResult{Result: pb.ERROR_CODE_SUCCESS})
}
func (s *User) OnReport(msg *pb.Report) {
	err := redis.UpdateReport(s.appId, msg)
	if err != nil {
		s.SendUserMsg(&pb.ReportResult{Result: pb.ERROR_CODE_FAIL})
		return
	}
	var openIdList []string
	for _, info := range msg.Info {
		openIdList = append(openIdList, info.OpenId)
	}
	result := redis.GetAudienceInfoList(s.appId, openIdList)
	s.SendUserMsg(&pb.ReportResult{
		Result: pb.ERROR_CODE_SUCCESS,
		Info:   result,
	})
}
func (s *User) OnGetRank(msg *pb.GetRank) {
	rankList, err := redis.GetRank(s.appId, msg.TopCount)
	if err != nil {
		s.Log("OnGetRank err", err)
		s.SendUserMsg(&pb.GetRankResult{Result: pb.ERROR_CODE_FAIL})
		return
	}
	s.SendUserMsg(&pb.GetRankResult{Result: pb.ERROR_CODE_SUCCESS, Info: rankList})
}

func (s *User) setPushActive(isActive bool) {
	if isActive {
		s.Log("SetPushActive1", s.appId)
		s.Log("SetPushActive2", s.roomId)
		redis.Subscribe(s.appId, s.roomId, s.OnRecvPush, s.pushCloseChan)
		if !s.isDebug {
			s.mgr.StartTask(s.appId, s.roomId, "1")
			s.mgr.StartTask(s.appId, s.roomId, "2")
			s.mgr.StartTask(s.appId, s.roomId, "3")
		}
	} else {
		s.pushCloseChan <- struct{}{}
		if !s.isDebug {
			s.mgr.StopTask(s.appId, s.roomId, "1")
			s.mgr.StopTask(s.appId, s.roomId, "2")
			s.mgr.StopTask(s.appId, s.roomId, "3")
		}
	}
}
