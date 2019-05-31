package stream

import (
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetDefault("stream_root_dir", "/tmp/dhara")
}
