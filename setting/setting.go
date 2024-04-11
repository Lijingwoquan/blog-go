package setting

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	//加载配置文件
	viper.SetConfigFile("./config/config.yaml")
	if err = viper.ReadInConfig(); err != nil {
		fmt.Printf("viper.ReadInConfig() failed,err:%v", err)
		return
	}
	//实时监控配置文件
	viper.WatchConfig()

	//配置文件修改之后的回调函数
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件修改了,name:%v,op:%v\n", e.Name, e.Op)
	})
	return
}
