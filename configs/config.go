package configs

import (
	"encoding/json"
	"io/ioutil"
)

type Conf struct {
	//websocket config
	Name          string `json:"name"`
	IpVersion            string `json:"ipVersion"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	HeartBeatTime int    `json:"heartBeatTime"`
	InChanSize    int    `json:"inChanSize"`
	OutChanSize   int    `json:"outChanSize"`

	//redis config
	RedisAddr string `json:"redisAddr"`
	RedisPort int    `json:"redisPort"`
	RedisPw   string `json:"redisPw"`

	//db config
	DbAddr     string `json:"dbAddr"`
	DbPort     int    `json:"dbPort"`
	DbDatabase string `json:"dbDatabase"`
	DbUserName string `json:"dbUserName"`
	DbPw       string `json:"dbPw"`

	WorkerPoolSize uint64 `json:"workerPoolSize"`
	MaxWorkTaskLen uint32 `json:"maxWorkTaskLen"`
	MaxConn int `json:"maxConn"`

	GetTeacherInfoUrl string `json:"getTeacherInfoUrl"`
}

var GConf *Conf

func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	conf := Conf{}
	err = json.Unmarshal(content, &conf)

	if err != nil {
		return err
	}
	GConf = &conf
	return nil
}
