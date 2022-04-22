package service

import (
	"fmt"
	"github.com/alt-golang/config"
	"github.com/alt-golang/logger"
	"os"
)

type ConfigService struct {
	Logger logger.Logger
	Dir    string
}

func (configService ConfigService) Get(environment string, instance string, profile string, path string) interface{} {

	configService.Logger.Info(fmt.Sprintf("Fetching config for environment:%s, instance:%s, profile:%s , path: %s", fmt.Sprint(environment), fmt.Sprint(instance), fmt.Sprint(profile), fmt.Sprint(path)))

	os.Setenv("GO_ENV", environment)
	os.Setenv("GO_APP_INSTANCE", instance)
	os.Setenv("GO_PROFILES_ACTIVE", profile)

	conf := config.GetServiceConfigFromDir(configService.Dir)
	var result interface{}
	result, _ = conf.Get("")
	if path != "" {
		result, _ = conf.Get(path)
	}

	configService.Logger.Info(fmt.Sprintf("Sending config for environment:%s, instance:%s, profile:%s , path: %s as:", environment, instance, profile, path) + fmt.Sprint(result))

	return result
}
