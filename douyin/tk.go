package douyin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type IApp interface {
	StartRefreshToken()
	GetRoomId(token string) (string, string, string, error)
	StartTask(roomid string, msg_type string) (string, error)
	StopTask(roomid string, msg_type string) (string, error)
}

func NewApp(appid, secret string) IApp {
	return &App{
		appid:  appid,
		secret: secret,
	}
}

type App struct {
	appid             string
	secret            string
	accessToken       string
	refreshTokenChan  chan string
	refreshTokenMutex sync.Mutex
}

func (s *App) StartRefreshToken() {
	go func() {
		s.refreshTokenChan = make(chan string, 1)
		for {
			var err error
			s.accessToken, err = s.getAccessToken()
			if err != nil {
				// handle error
				log.Errorf("Failed to refresh access token: %v", err)
			}
			select {
			case m := <-s.refreshTokenChan:
				log.Info(m)
			// Timeout
			case <-time.After(time.Hour): //
				log.Info("RefreshToken after a hour")
			}

		}
	}()
}

func (s *App) StartTask(roomid string, msg_type string) (string, error) {
	url := "https://webcast.bytedance.com/api/live_data/task/start"
	data := map[string]string{
		"roomid":   roomid,
		"appid":    s.appid,
		"msg_type": msg_type,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("access-token", s.accessToken)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	if result["err_no"].(float64) != 0 {
		if int64(result["err_no"].(float64)) == 40022 {
			s.needRefreshToken(fmt.Sprintf("RefreshToken:StartTask%f%v", result["err_no"].(float64), result["err_msg"].(string)))
		}
		return "", errors.New(result["err_msg"].(string))
	}
	return result["data"].(map[string]interface{})["task_id"].(string), nil
}

func (s *App) StopTask(roomid string, msg_type string) (string, error) {
	url := "https://webcast.bytedance.com/api/live_data/task/stop"
	data := map[string]string{
		"roomid":   roomid,
		"appid":    s.appid,
		"msg_type": msg_type,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("access-token", s.accessToken)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	if result["err_no"].(float64) != 0 {
		if int64(result["err_no"].(float64)) == 40022 {
			s.needRefreshToken(fmt.Sprintf("RefreshToken:StopTask%f%v", result["err_no"].(float64), result["err_msg"].(string)))
		}
		return "", errors.New(result["err_msg"].(string))
	}
	return "", nil
}

func (s *App) SendGiftPostRequest(roomid string, appid string, sec_gift_id_list []string) ([]string, error) {
	url := "https://webcast.bytedance.com/api/gift/top_gift"
	data := map[string]interface{}{
		"room_id":          roomid,
		"app_id":           appid,
		"sec_gift_id_list": sec_gift_id_list,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-token", s.accessToken)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if errcode, ok := result["errcode"]; ok {
		if errcode.(float64) != 0 {
			if int64(errcode.(float64)) == 40022 {
				s.needRefreshToken(fmt.Sprintf("RefreshToken:SendGiftPostRequest%f%v", result["errcode"].(float64), result["err_msg"].(string)))
			}
			return nil, errors.New(result["errmsg"].(string))
		}
	}

	success_top_gift_id_list := result["data"].(map[string]interface{})["success_top_gift_id_list"].([]interface{})
	var res []string
	for _, v := range success_top_gift_id_list {
		res = append(res, v.(string))
	}
	return res, nil
}

func (s *App) GetRoomId(token string) (string, string, string, error) {
	url := "http://webcast.bytedance.com/api/webcastmate/info"
	data := map[string]string{
		"token": token,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("X-Token", s.accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}
	var result = &struct {
		Data struct {
			Info struct {
				RoomId   string `json:"room_id"`
				Uid      string `json:"anchor_open_id"`
				Nickname string `json:"nick_name"`
			} `json:"info"`
		} `json:"data"`
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return "", "", "", err
	}
	if result.ErrCode != 0 {
		if result.ErrCode == 40022 {
			s.needRefreshToken(fmt.Sprintf("RefreshToken:GetRoomId%d%s", result.ErrCode, result.ErrMsg))
		}
		return "", "", "", errors.New(result.ErrMsg)
	}
	if result.Data.Info.RoomId == "" {
		return "", "", "", errors.New("no data in response")
	}

	roomId := result.Data.Info.RoomId
	uid := result.Data.Info.Uid
	nickname := result.Data.Info.Nickname
	return roomId, uid, nickname, nil
}

func (s *App) needRefreshToken(reason string) {
	s.refreshTokenMutex.Lock()
	defer s.refreshTokenMutex.Unlock()
	if len(s.refreshTokenChan) == 0 {
		s.refreshTokenChan <- reason
	}
}

// 新建一个golang func，参数appid secret grant_type  发送以json格式http请求返回asccess_token
func (s *App) getAccessToken() (string, error) {
	url := "https://developer.toutiao.com/api/apps/v2/token"
	data := map[string]string{
		"appid":      s.appid,
		"secret":     s.secret,
		"grant_type": "client_credential",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	if result["err_no"].(float64) != 0 {
		return "", errors.New(result["err_tips"].(string))
	}
	s.accessToken = result["data"].(map[string]interface{})["access_token"].(string)
	return s.accessToken, nil
}

func (s *App) GetAccessToken() string {
	return s.accessToken
}
