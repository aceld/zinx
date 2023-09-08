/**
 * User: coder.sdp@gmail.com
 * Date: 2023/9/8
 * Time: 14:18
 */

package zconf

import (
	"os"
	"path/filepath"
)

const (
	// EnvConfigFilePathKey (Set configuration file path export ZINX_CONFIG_FILE_PATH = xxxxxxzinx.json)
	// (设置配置文件路径 export ZINX_CONFIG_FILE_PATH = xxx/xxx/zinx.json)
	EnvConfigFilePathKey     = "ZINX_CONFIG_FILE_PATH"
	EnvDefaultConfigFilePath = "/conf/zinx.json"
)

var env = new(zEnv)

type zEnv struct {
	configFilePath string
}

func init() {
	configFilePath := os.Getenv(EnvConfigFilePathKey)
	if configFilePath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		configFilePath = filepath.Join(pwd, EnvDefaultConfigFilePath)
	}
	var err error
	configFilePath, err = filepath.Abs(configFilePath)
	if err != nil {
		panic(err)
	}
	env.configFilePath = configFilePath
}

func GetConfigFilePath() string {
	return env.configFilePath
}
