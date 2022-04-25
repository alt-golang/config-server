package gin

import (
	"fmt"
	"github.com/alt-golang/config-server/service"
	"github.com/alt-golang/logger"
	g "github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type Server struct {
	Logger        logger.Logger
	Port          int
	Context       string
	Mode          string
	Engine        *g.Engine
	ConfigService service.ConfigService
}

func (Server Server) Init() {

	Server.Engine.Use(
		g.LoggerWithFormatter(func(param g.LogFormatterParams) string {
			log := ""
			consoleLogger := Server.Logger.(logger.ConsoleLogger)
			if consoleLogger.IsInfoEnabled() {
				log = consoleLogger.Formatter.Format(param.TimeStamp, consoleLogger.Config.Category, consoleLogger.Config.Levels.GetNameForValue(logger.INFO), param.ErrorMessage,
					&struct {
						ClientIP   string
						Path       string
						Method     string
						StatusCode string
						Latency    string
						UserAgent  string
					}{
						ClientIP:   param.ClientIP,
						Path:       param.Path,
						Method:     param.Method,
						StatusCode: strconv.Itoa(param.StatusCode),
						Latency:    fmt.Sprint(param.Latency),
						UserAgent:  param.Request.UserAgent(),
					})
			}

			return log + "\n"
		}))
	Server.Engine.Use(g.Recovery())
	Server.Engine.Use(func(context *g.Context) {
		Server.Logger.Info("middleware /*:" + fmt.Sprint(context.Request.Body))
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		context.Header("Access-Control-Allow-Header", "Content-Type, Authorization")
		context.Header("Content-Type", "application/json")
		context.Next()
	})

	Server.Engine.OPTIONS("*all", func(context *g.Context) {
		Server.Logger.Info("OPTIONS /*: " + fmt.Sprint(context.Request.Body))
		context.Status(200)
	})

	Server.Engine.GET("/", func(context *g.Context) {
		Server.Logger.Info("GET /:" + fmt.Sprint(context.Request.Body))
		result, err := Server.ConfigService.Get("", "", "", context.Query("path"))
		Server.Response(context, result, err)

	})
	Server.Engine.GET("/:env", func(context *g.Context) {
		Server.Logger.Info("GET /:" + fmt.Sprint(context.Request.Body))
		result, err := Server.ConfigService.Get(context.Param("env"), "", "", context.Query("path"))
		Server.Response(context, result, err)
	})
	Server.Engine.GET("/:env/:instance", func(context *g.Context) {
		Server.Logger.Info("GET /:" + fmt.Sprint(context.Request.Body))
		result, err := Server.ConfigService.Get(context.Param("env"), context.Param("instance"), "", context.Query("path"))
		Server.Response(context, result, err)
	})
	Server.Engine.GET("/:env/:instance/*profiles", func(context *g.Context) {
		Server.Logger.Info("GET /:" + fmt.Sprint(context.Request.Body))
		result, err := Server.ConfigService.Get(context.Param("env"), context.Param("instance"), strings.Replace(context.Param("profiles"), "/", ",", 0), context.Query("path"))
		Server.Response(context, result, err)
	})
}

func (Server Server) Response(context *g.Context, result interface{}, err error) {
	fmt.Println(result)
	if err == nil {
		if result == nil {
			if context.Query("default") != "" {
				context.IndentedJSON(200, fmt.Sprint(context.Query("default")))
			} else {
				context.IndentedJSON(404, "Not Found")
			}
		} else {
			context.IndentedJSON(200, result)
		}

	} else {
		context.IndentedJSON(500, err)
	}
}

func (Server Server) Run() {
	err := Server.Engine.Run(fmt.Sprintf("127.0.0.1:%d", Server.Port))
	Server.Logger.Fatal("Server (gin) failed to start: " + fmt.Sprint(err))
}
