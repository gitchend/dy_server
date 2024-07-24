package game_mgr

import (
	"app/message/pb"
	"app/redis"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"sort"
	"strings"
)

func (s *HttpMgr) OnDouyinPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func (s *HttpMgr) OnDouyinDataPush(c *gin.Context) {
	var err error
	defer func() {
		c.JSON(200, gin.H{
			"message": "ok",
		})
		if err != nil {
			fmt.Println("OnDouyinDataPush err", err)
		}
		return
	}()
	appId := c.Param("appId")
	//appSecret, ok := APP_TOKEN_MAP[appId]
	//if !ok {
	//	return
	//}
	roomId := c.GetHeader("x-roomid")
	msgType := c.GetHeader("x-msg-type")
	//headerCheck := map[string]string{
	//	"x-nonce-str": c.GetHeader("x-nonce-str"),
	//	"x-timestamp": c.GetHeader("x-timestamp"),
	//	"x-msg-type":  msgType,
	//	"x-roomid":    roomId,
	//}
	bodyData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	//signatureRemote := c.GetHeader("x-signature")
	//signatureLocal := signature(headerCheck, string(bodyData), appSecret)
	//if  signatureRemote != signatureLocal {
	//	return
	//}

	var audienceData []*pb.AudienceBasic
	var notifyData []*pb.NotifyAudienceAction

	switch msgType {
	case "live_comment":
		var rawData []*struct {
			AvatarUrl string `json:"avatar_url"`
			Content   string `json:"content"`
			MsgID     string `json:"msg_id"`
			NickName  string `json:"nickname"`
			SecOpenid string `json:"sec_openid"`
			Timestamp int64  `json:"timestamp"`
		}

		err = json.Unmarshal(bodyData, &rawData)
		if err != nil {
			return
		}
		for _, data := range rawData {
			audienceData = append(audienceData, &pb.AudienceBasic{
				OpenId:    data.SecOpenid,
				NickName:  data.NickName,
				AvatarUrl: data.AvatarUrl,
			})
			notifyData = append(notifyData, &pb.NotifyAudienceAction{
				OpenId:  data.SecOpenid,
				Content: data.Content,
			})
		}

	case "live_gift":
		var rawData []*struct {
			MsgID     string `json:"msg_id"`
			SecOpenid string `json:"sec_openid"`
			SecGiftID string `json:"sec_gift_id"`
			GiftNum   int    `json:"gift_num"`
			GiftValue int    `json:"gift_value"`
			AvatarUrl string `json:"avatar_url"`
			NickName  string `json:"nickname"`
			Timestamp int    `json:"timestamp"`
			Test      bool   `json:"test"`
		}

		err = json.Unmarshal(bodyData, &rawData)
		if err != nil {
			return
		}
		for _, data := range rawData {
			audienceData = append(audienceData, &pb.AudienceBasic{
				OpenId:    data.SecOpenid,
				NickName:  data.NickName,
				AvatarUrl: data.AvatarUrl,
			})
			notifyData = append(notifyData, &pb.NotifyAudienceAction{
				OpenId:  data.SecOpenid,
				GiftId:  data.SecGiftID,
				GiftNum: int32(data.GiftNum),
			})
		}
	case "live_like":
		var rawData []*struct {
			MsgID     string `json:"msg_id"`
			SecOpenid string `json:"sec_openid"`
			LikeNum   int    `json:"like_num"`
			AvatarUrl string `json:"avatar_url"`
			NickName  string `json:"nickname"`
			Timestamp int    `json:"timestamp"`
		}

		err = json.Unmarshal(bodyData, &rawData)
		if err != nil {
			return
		}
		for _, data := range rawData {
			audienceData = append(audienceData, &pb.AudienceBasic{
				OpenId:    data.SecOpenid,
				NickName:  data.NickName,
				AvatarUrl: data.AvatarUrl,
			})
			notifyData = append(notifyData, &pb.NotifyAudienceAction{
				OpenId:  data.SecOpenid,
				LikeNum: int32(data.LikeNum),
			})
		}
	}

	for _, audience := range audienceData {
		redis.SetAudienceBasic(appId, audience)
	}
	for _, notify := range notifyData {
		redis.Publish(appId, roomId, notify)
	}
}

func signature(header map[string]string, bodyStr, secret string) string {
	keyList := make([]string, 0, 4)
	for key, _ := range header {
		keyList = append(keyList, key)
	}
	sort.Slice(keyList, func(i, j int) bool {
		return keyList[i] < keyList[j]
	})
	kvList := make([]string, 0, 4)
	for _, key := range keyList {
		kvList = append(kvList, key+"="+header[key])
	}
	urlParams := strings.Join(kvList, "&")
	rawData := urlParams + bodyStr + secret
	md5Result := md5.Sum([]byte(rawData))
	return base64.StdEncoding.EncodeToString(md5Result[:])
}
