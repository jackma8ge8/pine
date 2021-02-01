package application

import (
	"math/rand"
	"time"

	"github.com/jackma8ge8/pine/application/config"
	"github.com/jackma8ge8/pine/logger"
	RpcServer "github.com/jackma8ge8/pine/rpc/server"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	config.ParseConfig()
	logger.SetLogMode(config.GetServerConfig().LogType)
}

// Application app
type Application struct {
}

// Start start application
func (app Application) Start() {
	RpcServer.Start()
}

// App pine application instance
var App *Application

// CreateApp 创建app
func CreateApp() *Application {

	if App != nil {
		return App
	}

	App = &Application{}

	return App
}
