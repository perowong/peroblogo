package conf

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/perowong/peroblogo/utils"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	App struct {
		Name string
		Port string
	}
	Mysql map[string]struct {
		Dsn string
	}
}

var (
	Env string
	C   Conf
)

func loadLocalConfig(env string) {
	pwd, _ := os.Getwd()
	configName := fmt.Sprintf("conf.%s.yaml", env)

	data, err := utils.ReadFile(filepath.Join(pwd, "conf", configName))
	if err != nil {
		log.Fatalln(err.Error())
	}

	C = Conf{}
	err = yaml.Unmarshal(data, &C)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func init() {
	var env string
	flag.StringVar(&env, "env", "dev", "set env")
	flag.Parse()

	Env = env
	loadLocalConfig(env)
}
