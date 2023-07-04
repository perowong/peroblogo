package conf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/perowong/peroblogo/utils"
	"gopkg.in/yaml.v3"

	logger "github.com/agent-chatee/zap-logger"
)

type Conf struct {
	App struct {
		Name string
		Addr string
	}
	Consul struct {
		Addr string
	}
	Mysql map[string]struct {
		Dsn string
	}
	Github struct {
		ClientID     string `yaml:"client-id"`
		ClientSecret string `yaml:"client-secret"`
	}
	LoggerConfig logger.Config     `yaml:"loggerConfig"`
	ServerNames  map[string]string `yaml:"serverNames"`
}

type EnvType string

const (
	Development EnvType = "dev"
	Test        EnvType = "test"
	Production  EnvType = "prod"
)

var (
	Env EnvType
	C   Conf
)

func getConfDir() string {
	dir := "conf"
	for i := 0; i < 5; i++ {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			break
		}
		dir = filepath.Join("..", dir)
	}

	return dir
}

func loadLocalConfig(env EnvType) {
	configName := fmt.Sprintf("conf.%s.yaml", env)

	data, err := utils.ReadFile(filepath.Join(getConfDir(), configName))
	if err != nil {
		log.Fatal(err.Error())
	}

	C = Conf{}
	err = yaml.Unmarshal(data, &C)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func initLogger() {
	fConfig := C.LoggerConfig.FileConfig
	lConfig := &logger.Config{FileConfig: &logger.FileConfig{
		LogFilePath:   fConfig.LogFilePath,
		MaxSize:       fConfig.MaxSize,
		MaxBackups:    fConfig.MaxBackups,
		MaxAge:        fConfig.MaxAge,
		Console:       fConfig.Console,
		LevelString:   fConfig.LevelString,
		RotationTime:  fConfig.RotationTime,
		FilePathDepth: fConfig.FilePathDepth,
	}}
	logger.Init(lConfig, C.App.Name)
}

func InitConf(env EnvType) {
	Env = env
	loadLocalConfig(Env)
	initLogger()
}
