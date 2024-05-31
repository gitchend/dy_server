package main

import (
	"app/network"
	"app/service/main/message/pb"
	"app/service/main/msg_util"
	"app/service/robot/user"
	"math/rand"
	"sync"
	"testing"
)

func TestName(t *testing.T) {
	client := network.NewTcpClient(
		"129.211.8.84:8899",
		newParser,
		newHandler)
	err := client.Connect(false)
	if err != nil {
		panic(err)
	}
	//登陆
	login := &pb.Login{
		AppId:   "debug",
		IsDebug: true,
	}
	err = client.SendMsg(login)
	if err != nil {
		panic(err)
	}

	//上传分数
	openIdList := []string{
		"7374686837987841059",
		"7374698021348676648",
		"7372845558429537319",
		"7374677100777002003",
		"7374698021348660264",
		"7374673415892489242",
		"7374683012379776011",
		"7374686060594746368",
		"7374687164753499177",
		"7372850431053173801",
	}
	report := &pb.Report{}
	for _, openId := range openIdList {
		report.Info = append(report.Info, &pb.ReportInfo{
			OpenId: openId,
			Score:  rand.Int31n(100),
			IsWin:  rand.Float32() > 0.5,
		})
	}
	err = client.SendMsg(report)
	if err != nil {
		panic(err)
	}
	//拉排行榜
	rank := &pb.GetRank{
		TopCount: 5,
	}
	err = client.SendMsg(rank)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func newParser() network.ICodec {
	return msg_util.NewProtoCodec(msg_util.NewProtoParser("message", "MSG_TYPE"), 1024*64, false)
}

func newHandler() network.INetHandler {
	return user.NewUser()
}
