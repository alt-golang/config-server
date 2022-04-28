package context

import (
	"fmt"
	"github.com/alt-golang/config"
	"github.com/alt-golang/config-server/service"
	"github.com/alt-golang/config-server/web/bindings/gin"
	"github.com/alt-golang/logger"
	g "github.com/gin-gonic/gin"
	"os"
	"strconv"
)

func getCfgFlag() string {
	var cfgFlag = ""
	if len(os.Args) >= 3 {
		cfgFlag = os.Args[2]
	}
	if cfgFlag == "" {
		cfgFlag = "config" + fmt.Sprint(string(os.PathSeparator)) + "internal"
	}
	return cfgFlag
}

var cfg = config.GetConfigFromDir(getCfgFlag())
var dir, _ = config.GetWithDefault("config.dir", "config")
var gitUrl, _ = config.GetWithDefault("git.url", "")
var gitBranch, _ = config.GetWithDefault("git.branch", "main")
var gitUser, _ = config.GetWithDefault("git.username", "")
var gitToken, _ = config.GetWithDefault("git.token", "")

var ConfigService = service.ConfigService{
	Logger:      logger.GetLogger("github.com/alt-golang/config-server/service/ConfigService"),
	Dir:         dir.(string),
	GitUrl:      gitUrl.(string),
	GitBranch:   gitBranch.(string),
	GitUsername: gitUser.(string),
	GitToken:    gitToken.(string),
}

var host, _ = config.GetWithDefault("server.host", "0.0.0.0")
var hostStr = fmt.Sprint(host)
var port, _ = config.GetWithDefault("server.port", "80")
var portInt, _ = strconv.Atoi(fmt.Sprint(port))
var mode, _ = config.GetWithDefault("server.mode", "debug")

var Server = gin.Server{
	Logger:        logger.GetLogger("github.com/alt-golang/config-server/web/bindings/gin/Server"),
	Host:          hostStr,
	Port:          portInt,
	Context:       "",
	Mode:          mode.(string),
	ConfigService: ConfigService,
}

func Start() {
	ConfigService.Init()
	g.SetMode(mode.(string))
	Server.Engine = g.New()
	Server.Init()
	Server.Run()
}
