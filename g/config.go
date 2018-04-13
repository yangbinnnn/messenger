package g

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type HttpConfig struct {
	Listen string `json:"listen"`
	Token  string `json:"token"`
}

type SmtpConfig struct {
	Enable   bool   `json:"enable"`
	TLS      bool   `json:"tls"`
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	Timeout  int    `json:"timeout"`
}

type WechatConfig struct {
	Enable         bool   `json:"enable"`
	CorpID         string `json:"corpid"`
	AgentId        int    `json:"agentid"`
	Secret         string `json:secret`
	EncodingAESKEY string `json:aeskey`
	Timeout        int    `json:"timeout"`
}

type GlobalConfig struct {
	Debug  bool          `json:"debug"`
	Http   *HttpConfig   `json:"http"`
	Smtp   *SmtpConfig   `json:"smtp"`
	Wechat *WechatConfig `json:"wechat"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func FileString(path string) (string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !FileIsExist(cfg) {
		return fmt.Errorf("configuration file %s is nonexistent", cfg)
	}

	ConfigFile = cfg

	configContent, err := FileString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	configLock.Lock()
	defer configLock.Unlock()
	UseEnvConfig(&c)
	config = &c

	log.Println("load configuration file", cfg, "successfully")
	return nil
}

func UseEnvConfig(c *GlobalConfig) {
	enable := os.Getenv("USE_ENV_CONFIG")
	if enable != "true" {
		return
	}
	log.Println("load env configuration")
}
