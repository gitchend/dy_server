package game_user

import (
	"app/in_obj"
	"app/message/pb"
	"app/redis"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

func NewUser(mgr in_obj.IGameMgr) in_obj.IGameUser {
	return &GameUser{mgr: mgr}
}

type GameUser struct {
	mgr           in_obj.IGameMgr
	appId         string
	roomId        string
	uid           string
	nickName      string
	isDebug       bool
	isRunning     bool
	audienceSet   sync.Map
	sendChan      chan proto.Message
	msgChan       chan func()
	closeChan     chan struct{}
	pushCloseChan chan struct{}
}

func (s *GameUser) OnSessionCreated(sendChan chan proto.Message) {
	fmt.Println("[SessionCreated]")
	//主协程
	s.sendChan = sendChan
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

func (s *GameUser) OnSessionClosed() {
	fmt.Println("[SessionClosed]")
	if !s.isRunning {
		return
	}
	s.setPushActive(false)
	s.isRunning = false
	s.pushCloseChan <- struct{}{}
	s.closeChan <- struct{}{}
}

func (s *GameUser) OnRecv(msg proto.Message) {
	s.Log("Recv", msg)
	s.msgChan <- func() {
		switch v := msg.(type) {
		case *pb.Ping:
			s.OnPing(v)
		case *pb.Login:
			s.OnLogin(v)
		case *pb.PlayStart:
			s.OnPlayStart(v)
		case *pb.PlayEnd:
			s.OnPlayEnd(v)
		case *pb.Report:
			s.OnReport(v)
		case *pb.GetScoreRank:
			s.OnGetScoreRank(v)
		case *pb.GetMonthScoreRank:
			s.OnGetMonthScoreRank(v)
		default:
			s.Log("Recv unknown", msg)
		}
	}
}

func (s *GameUser) OnRecvPush(data string) {
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
				s.sendUserMsg(audienceMsg)
			}
		}
	}
	s.msgChan <- func() {
		s.sendUserMsg(msg)
	}
}

func (s *GameUser) sendUserMsg(msg proto.Message) {
	s.Log("Send", msg)
	s.sendMsg(msg)
}

func (s *GameUser) sendMsg(msg proto.Message) {
	s.sendChan <- msg
}

func (s *GameUser) Log(logType string, data ...interface{}) {
	fmt.Printf("[%s](%s,%s,%s)%v\n", logType, s.roomId, s.uid, s.nickName, data)
}
func (s *GameUser) OnPing(msg *pb.Ping) {
	s.sendMsg(&pb.Pong{
		ServerTime: time.Now().Unix(),
		ClientTime: msg.ClientTime,
	})
}
func (s *GameUser) OnLogin(msg *pb.Login) {
	var (
		roomId   string
		uid      string
		nickName string
		result   pb.ERROR_CODE
	)
	if !msg.IsDebug {
		roomId, uid, nickName, result = s.mgr.SdkLogin(msg.AppId, msg.Token)
	} else {
		roomId, uid, nickName, result = "1111111111111111111", "debug", "debug", pb.ERROR_CODE_SUCCESS
	}
	s.appId = msg.AppId
	s.roomId = roomId
	s.uid = uid
	s.nickName = nickName
	s.isDebug = msg.IsDebug
	s.sendUserMsg(&pb.LoginResult{
		Result:   result,
		RoomId:   roomId,
		UID:      uid,
		NickName: nickName,
	})

	//test
	s.OnPlayStart(nil)
}
func (s *GameUser) OnPlayStart(msg *pb.PlayStart) {
	if s.isRunning {
		s.sendUserMsg(&pb.PlayStartResult{Result: pb.ERROR_CODE_GAME_IS_RUNNING})
		return
	}
	s.isRunning = true
	s.setPushActive(true)
	s.sendUserMsg(&pb.PlayStartResult{Result: pb.ERROR_CODE_SUCCESS})
}
func (s *GameUser) OnPlayEnd(msg *pb.PlayEnd) {
	if !s.isRunning {
		s.sendUserMsg(&pb.PlayEndResult{Result: pb.ERROR_CODE_GAME_IS_STOPPED})
		return
	}
	s.isRunning = false
	s.setPushActive(false)
	s.sendUserMsg(&pb.PlayEndResult{Result: pb.ERROR_CODE_SUCCESS})
}
func (s *GameUser) OnReport(msg *pb.Report) {
	err := redis.UpdateReport(s.appId, msg)
	if err != nil {
		s.sendUserMsg(&pb.ReportResult{Result: pb.ERROR_CODE_FAIL})
		return
	}
	var openIdList []string
	for _, info := range msg.Info {
		openIdList = append(openIdList, info.OpenId)
	}
	result := redis.GetAudienceInfoList(s.appId, openIdList)
	s.sendUserMsg(&pb.ReportResult{
		Result: pb.ERROR_CODE_SUCCESS,
		Info:   result,
	})
}
func (s *GameUser) OnGetScoreRank(msg *pb.GetScoreRank) {
	rankList, err := redis.GetScoreRank(s.appId, msg.TopCount)
	if err != nil {
		s.Log("OnGetScoreRank err", err)
		s.sendUserMsg(&pb.GetScoreRankResult{Result: pb.ERROR_CODE_FAIL})
		return
	}
	s.sendUserMsg(&pb.GetScoreRankResult{Result: pb.ERROR_CODE_SUCCESS, Info: rankList})
}

func (s *GameUser) OnGetMonthScoreRank(msg *pb.GetMonthScoreRank) {
	rankList, err := redis.GetMonthScoreRank(s.appId, msg.TopCount)
	if err != nil {
		s.Log("OnGetMonthScoreRank} err", err)
		s.sendUserMsg(&pb.GetMonthScoreRankResult{Result: pb.ERROR_CODE_FAIL})
		return
	}
	s.sendUserMsg(&pb.GetMonthScoreRankResult{Result: pb.ERROR_CODE_SUCCESS, Info: rankList})
}

func (s *GameUser) setPushActive(isActive bool) {
	if isActive {
		redis.Subscribe(s.appId, s.roomId, s.OnRecvPush, s.pushCloseChan)
		if !s.isDebug {
			s.mgr.SdkStartTask(s.appId, s.roomId, "live_comment")
			s.mgr.SdkStartTask(s.appId, s.roomId, "live_gift")
			s.mgr.SdkStartTask(s.appId, s.roomId, "live_like")
		}
	} else {
		s.pushCloseChan <- struct{}{}
		if !s.isDebug {
			s.mgr.SdkStopTask(s.appId, s.roomId, "live_comment")
			s.mgr.SdkStopTask(s.appId, s.roomId, "live_gift")
			s.mgr.SdkStopTask(s.appId, s.roomId, "live_like")
		}
	}
}
