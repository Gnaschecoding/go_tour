package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	//vp.AddConfigPath("2_blog-serie/configs/") //这个地方我保持疑问

	if len(configs) > 0 {
		for _, config := range configs {
			if config != "" {
				vp.AddConfigPath(config)
			}
		}
	} else {
		vp.AddConfigPath("2_blog-serie/configs/")
	}

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	s := &Setting{vp: vp}

	s.WatchSettingChange()

	return s, nil
}

func (s *Setting) WatchSettingChange() {
	go func() {

		//viper.WatchConfig() 可以监听配置文件的变更，并在发生变化时自动重新读取和更新配置。
		s.vp.WatchConfig()
		//OnConfigChange() 注册回调 ,如果上面的 监听起作用了，就执行 回调函数
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}
