package utils

import (
	"github.com/fsnotify/fsnotify"

	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var Conf *viper.Viper

func init() {
	Conf = viper.New()

	Conf.SetConfigType("toml")
	Conf.SetConfigName("torch-client")
	Conf.AddConfigPath(`.`)

	Conf.SetDefault("Supervisor", false)

	Conf.SetDefault("URL", "")
	Conf.SetDefault("Frequency", "")

	Conf.SetDefault("Instance.Type", "")
	Conf.SetDefault("Token", "")

	replacer := strings.NewReplacer(".", "_")
	Conf.SetEnvKeyReplacer(replacer)
	err := Conf.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		_, err := os.Create("./torch-client.toml")
		if err != nil {
			panic(err)
		}

		err = Conf.WriteConfig()
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	Conf.WatchConfig()
	Conf.OnConfigChange(func(in fsnotify.Event) {
		err := Conf.ReadInConfig()
		if err != nil {
			log.Println(err)
		}
	})
}
