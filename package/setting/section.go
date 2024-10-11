package setting

type Config struct {
	Mysql  MySQLSetting  `mapstructure:"mysql"`
	Logger LogSetting    `mapstructure:"logger"`
	Server ServerSetting `mapstructure:"server"`
	Redis  RedisSetting  `mapstructure:"redis"`
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
	FileName   string `mapstruture:"file_name"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}
