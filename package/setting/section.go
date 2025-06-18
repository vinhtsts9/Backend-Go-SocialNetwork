package setting

type Config struct {
	Logger       LogSetting    `mapstructure:"log"`
	Server       ServerSetting `mapstructure:"server"`
	Redis        RedisSetting  `mapstructure:"redis"`
	JWT          JWTSetting    `mapstructure:"jwt"`
	MySQLHaproxy MySQLSetting  `mapstructure:"mysql_haproxy"`
	// MySQLSlave   MySQLSetting  `mapstructure:"mysql_slave"`
	// MySQLSlave2  MySQLSetting  `mapstructure:"mysql_slave2"`
	// MySQLSlave3  MySQLSetting  `mapstructure:"mysql_slave3"`
	KafkaBroker   Kafka         `mapstructure:"kafka"`
	ElasticSearch ElasticSearch `mapstructure:"elasticsearch"`
}
type Cloudinary struct {
	CloudName string `mapstructure:"CLOUD_NAME"`
	ApiKey    string `mapstructure:"API_KEY"`
	ApiSecret string `mapstructure:"API_SECRET"`
}
type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}
type RedisSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}
type MySQLSetting struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Dbname          string `mapstructure:"dbname"`
	MaxIdleConnes   int    `mapstructure:"maxIdleConnes"`
	MaxOpenConnes   int    `mapstructure:"maxOpenConnes"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
}
type LogSetting struct {
	LogLevel   string `mapstructure:"log_level"`
	FileName   string `mapstructure:"file_log_name"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}
type JWTSetting struct {
	TOKEN_HOUR_LIFESPAN uint   `mapstructure:"TOKEN_HOUR_LIFESPAN"`
	API_SECRET_KEY      string `mapstructure:"API_SECRET_KEY"`
	JWT_EXPIRATION      string `mapstructure:"JWT_EXPIRATION"`
	REFRESH_EXPIRATION  string `mapstructure:"REFRESH_EXPIRATION"`
}
type Kafka struct {
	Brokers string `mapstructure:"brokers"`
}
type ElasticSearch struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
