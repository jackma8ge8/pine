package application

import (
	"github.com/jackma8ge8/pine/application/config"
	"github.com/jackma8ge8/pine/connector"
)

// AsConnector 作为Connector 启动
func (app Application) AsConnector(authFunc func(uid string, token string, sessionData map[string]string) error) {
	config.GetServerConfig().IsConnector = true
	connector.RegisteAuth(authFunc)
}
