package sender

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
)

type Wechat struct {
	CorpID         string
	AgentID        int
	Secret         string
	EncodingAESKey string
	TokenCache     *cache.Cache
	Timeout        int
}

func NewWechat(CorpID string, AgentID int, Secret, EncodingAESKey string) *Wechat {
	tc := cache.New(6000*time.Second, 5*time.Second)
	return &Wechat{CorpID: CorpID, AgentID: AgentID,
		Secret: Secret, EncodingAESKey: EncodingAESKey, TokenCache: tc}
}

//发送信息
type Content struct {
	Content string `json:"content"`
}

type MsgPost struct {
	ToUser  string  `json:"touser"`
	MsgType string  `json:"msgtype"`
	AgentID int     `json:"agentid"`
	Text    Content `json:"text"`
}

type GroupMsg struct {
	ChatId  string  `json:"chatid"`
	MsgType string  `json:"msgtype"`
	Text    Content `json:"text"`
}

func (wechat Wechat) CreateChatGroup(name, owner, chatId string, userList []string) error {
	token, found := wechat.TokenCache.Get("token")
	if !found {
		log.Printf("token获取失败!")
		return errors.New("token获取失败!")
	}
	accessToken, ok := token.(AccessToken)
	if !ok {
		return errors.New("token解析失败!")
	}

	url := "https://qyapi.weixin.qq.com/cgi-bin/appchat/create?access_token=" + accessToken.AccessToken

	data := make(map[string]interface{})
	data["name"] = name
	data["owner"] = owner
	data["userlist"] = userList
	data["chatid"] = chatId
	result, err := wechat.WxPost(url, data)
	if err != nil {
		log.Printf("请求微信失败: %v", err)
	}

	log.Printf("创建讨论组: %s, 所有者: %s, ChatID: %s, 微信返回结果: %v", name, owner, chatId, result)
	return nil
}

func (wechat Wechat) SendMsgToGroup(chatId, content string) error {
	text := Content{}
	text.Content = content

	msg := GroupMsg{
		ChatId:  chatId,
		MsgType: "text",
		Text:    text,
	}

	token, found := wechat.TokenCache.Get("token")
	if !found {
		log.Printf("token获取失败!")
		return errors.New("token获取失败!")
	}
	accessToken, ok := token.(AccessToken)
	if !ok {
		return errors.New("token解析失败!")
	}

	url := "https://qyapi.weixin.qq.com/cgi-bin/appchat/send?access_token=" + accessToken.AccessToken

	result, err := wechat.WxPost(url, msg)
	if err != nil {
		log.Printf("请求微信失败: %v", err)
	}
	log.Printf("发送信息给%s, 信息内容: %s, 微信返回结果: %v", chatId, content, result)
	return nil
}

func (wechat Wechat) SendMsg(toUser, content string) error {
	if userList := strings.Split(toUser, ","); len(userList) > 1 {
		toUser = strings.Join(userList, "|")
	}

	text := Content{}
	text.Content = content

	msg := MsgPost{
		ToUser:  toUser,
		MsgType: "text",
		AgentID: wechat.AgentID,
		Text:    text,
	}

	token, found := wechat.TokenCache.Get("token")
	if !found {
		log.Printf("token获取失败!")
		return errors.New("token获取失败!")
	}
	accessToken, ok := token.(AccessToken)
	if !ok {
		return errors.New("token解析失败!")
	}

	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + accessToken.AccessToken

	result, err := wechat.WxPost(url, msg)
	if err != nil {
		log.Printf("请求微信失败: %v", err)
	}
	log.Printf("发送信息给%s, 信息内容: %s, 微信返回结果: %v", toUser, content, result)
	return nil
}

//开启回调模式验证
func (wechat Wechat) WxAuth(echostr string) error {
	if echostr == "" {
		return errors.New("无法获取请求参数, echostr 为空")
	}

	wByte, err := base64.StdEncoding.DecodeString(echostr)
	if err != nil {
		return errors.New("接受微信请求参数 echostr base64解码失败(" + err.Error() + ")")
	}
	key, err := base64.StdEncoding.DecodeString(wechat.EncodingAESKey + "=")
	if err != nil {
		return errors.New("配置 EncodingAESKey base64解码失败(" + err.Error() + "), 请检查配置文件内 EncodingAESKey 是否和微信后台提供一致")
	}

	keyByte := []byte(key)
	x, err := AesDecrypt(wByte, keyByte)
	if err != nil {
		return errors.New("aes 解码失败(" + err.Error() + "), 请检查配置文件内 EncodingAESKey 是否和微信后台提供一致")
	}

	buf := bytes.NewBuffer(x[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)

	//验证返回数据ID是否正确
	appIDstart := 20 + length
	if len(x) < int(appIDstart) {
		return errors.New("获取数据错误, 请检查 EncodingAESKey 配置")
	}
	id := x[appIDstart : int(appIDstart)+len(wechat.CorpID)]
	if string(id) == wechat.CorpID {
		return nil
	}
	return errors.New("微信验证appID错误, 微信请求值: " + string(id) + ", 配置文件内配置为: " + wechat.CorpID)
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

//从微信获取 AccessToken
func (wechat Wechat) GetAccessTokenFromWeixin() {

	for {
		if wechat.CorpID == "" || wechat.Secret == "" {
			log.Printf("corpId或者secret 获取失败, 请检查配置文件")
			return
		}

		WxAccessTokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + wechat.CorpID + "&corpsecret=" + wechat.Secret

		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		}
		client := &http.Client{Transport: tr, Timeout: time.Duration(wechat.Timeout) * time.Second}
		result, err := client.Get(WxAccessTokenUrl)
		if err != nil {
			log.Printf("获取微信 Token 返回数据错误: %v, 10秒后重试!", err)
			time.Sleep(10 * time.Second)
			continue
		}

		res, err := ioutil.ReadAll(result.Body)

		if err != nil {
			log.Printf("获取微信 Token 返回数据错误: %v, 10秒后重试!", err)
			time.Sleep(10 * time.Second)
			continue
		}
		newAccess := AccessToken{}
		err = json.Unmarshal(res, &newAccess)
		if err != nil {
			log.Printf("获取微信 Token 返回数据解析 Json 错误: %v, 10秒后重试!", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if newAccess.ExpiresIn == 0 || newAccess.AccessToken == "" {
			log.Printf("获取微信错误代码: %v, 错误信息: %v, 10秒后重试!", newAccess.ErrCode, newAccess.ErrMsg)
			time.Sleep(10 * time.Second)
			continue
		}

		//延迟时间
		wechat.TokenCache.Set("token", newAccess, time.Duration(newAccess.ExpiresIn)*time.Second)
		log.Printf("微信 Token 更新成功: %s,有效时间: %v", newAccess.AccessToken, newAccess.ExpiresIn)
		time.Sleep(time.Duration(newAccess.ExpiresIn-1000) * time.Second)

	}

}

//微信请求数据
func (wechat Wechat) WxPost(url string, data interface{}) (string, error) {
	jsonBody, err := encodeJson(data)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: time.Duration(wechat.Timeout) * time.Second}
	r, err := client.Post(url, "application/json;charset=utf-8", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

//获取当前运行路径
func GetWorkPath() string {
	if file, err := exec.LookPath(os.Args[0]); err == nil {
		return filepath.Dir(file) + "/"
	}
	return "./"
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("aes解密失败: %v", err)
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//string 类型转 int
func StringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("agent 类型转换失败, 请检查配置文件中 agentid 配置是否为纯数字(%v)", err)
		return 0
	}
	return n
}

//json序列化(禁止 html 符号转义)
func encodeJson(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
