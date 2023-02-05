package configs

import "github.com/spf13/viper"

func InitConfigs() {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}
}
