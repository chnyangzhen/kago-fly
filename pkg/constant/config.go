package constant

// 环境变量的配置
const (
	LogConfigName        = "log_config_name"
	DefaultLogConfigName = "./conf.d/logger.yml"
	EnvLogConfigName     = "LOG_CONFIG_NAME"

	ConfigNames       = "config_names"
	EnvConfigNames    = "CONFIG_NAMES"
	DefaultConfigPath = "./conf.d"
)

// 配置文件
const (
	Application = "application"
	Logger      = "logger"
	Hystrix     = "hystrix"
	Transporter = "transporter"
	Discovery   = "discovery"
	Go2Sky      = "go2sky"
	Consumer    = "consumer"
)

// viper配置
const (
	DynamicTypeStaticValue  = 0 << iota
	DynamicTypeLookupConfig = 1
	DynamicTypeLookupEnv    = 2
)

const (
	ApolloConfig = "apollo_config"
)

const (
	DubboConsumerConfig = "dubbo-consumer-config"
)
