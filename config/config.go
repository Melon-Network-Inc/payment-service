package config

import "github.com/spf13/viper"

type ServiceConfig struct {
	ServiceName		string
	Version	   		string
	ServerPort		int
	DatabaseUrl  	string
	CacheUrl    	string
	EmailApiKey     string
}

func BuildServerConfig(ServiceName string, Version string) ServiceConfig {
	sc := ServiceConfig{}
	
	viper.SetConfigFile("./pkg/envs/.env")

	if err := viper.ReadInConfig(); err == nil {
		sc = ServiceConfig{
			ServiceName: ServiceName,
			Version: 	 Version,
			ServerPort: viper.Get("PORT").(int),
			DatabaseUrl: viper.Get("DB_URL").(string),
			CacheUrl: viper.Get("CACHE_URL").(string),
			EmailApiKey: viper.Get("SEND_GRID").(string),
		}
	} else {
		sc = ServiceConfig{
			ServiceName: 	ServiceName,
			Version: 	 	Version,
			ServerPort: 	7000,
			DatabaseUrl: 	"postgres://postgres:postgres@localhost:5432/melon_service",
			CacheUrl: 		"localhost:6379",
			EmailApiKey: 	"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",
		}
	}
	return sc
}