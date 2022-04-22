package context

import (
	"fmt"
	"github.com/alt-golang/config"
	"github.com/alt-golang/config-server/service"
	gin "github.com/alt-golang/config-server/web/bindings/gin"
	"github.com/alt-golang/logger"
	g "github.com/gin-gonic/gin"
	"strconv"
)

var cfg = config.GetConfigFromDir("config/internal")
var dir, _ = config.GetWithDefault("config.dir", "config")
var ConfigService = service.ConfigService{
	Logger: logger.GetLogger("github.com/alt-golang/config-server/service/ConfigService"),
	Dir:    dir.(string),
}

var port, _ = config.Get("server.port")
var portInt, _ = strconv.Atoi(fmt.Sprint(port))
var mode, _ = config.Get("server.mode")

var Server = gin.Server{
	Logger:        logger.GetLogger("github.com/alt-golang/config-server/web/bindings/gin/Server"),
	Port:          portInt,
	Context:       "",
	Mode:          mode.(string),
	ConfigService: ConfigService,
}

func Start() {
	g.SetMode(mode.(string))
	Server.Engine = g.New()
	Server.Init()
	Server.Run()
}
