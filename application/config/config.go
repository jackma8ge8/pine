package config

// ==========================================
// ServerConfig
// ==========================================
var serverConfig *ServerConfig

// ServerConfig 服务器配置 配置文件
type ServerConfig struct {
	SystemName  string `validate:"required"`
	ID          string `validate:"required"`
	Kind        string `validate:"required"`
	Host        string `validate:"required"`
	Port        uint32 `validate:"gte=1,lte=65535"`
	IsConnector bool
	Token       string `validate:"required"`
	LogType     string `validate:"oneof=Console File"`
	LogLevel    string `validate:"oneof=Debug Info Warn Error"`
}

// SetServerConfig 保存服务器配置
func SetServerConfig(sc *ServerConfig) {
	serverConfig = sc
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *ServerConfig {
	return serverConfig
}

// ==========================================
// ZkConfig
// ==========================================
var zkConfig *ZkConfig

// ZkConfig zk 配置文件
type ZkConfig struct {
	Host string `validate:"required"`
	Port uint32 `validate:"gte=1,lte=65535"`
}

// SetZkConfig 配置zookeeper配置
func SetZkConfig(zc *ZkConfig) {
	zkConfig = zc
}

// GetZkConfig 获取zk配置
func GetZkConfig() *ZkConfig {
	return zkConfig
}
