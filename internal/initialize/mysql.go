package initialize

import (
	"database/sql"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/package/setting"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// checkErrorPanicC handles errors by logging them and triggering a panic
func checkErrorPanicC(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

// initMysqlConnection initializes a MySQL connection based on provided settings
func initMysqlConnection(m setting.MySQLSetting) *sql.DB {
	dsn := "%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, m.Username, m.Password, m.Host, m.Port, m.Dbname)
	db, err := sql.Open("mysql", s)
	checkErrorPanicC(err, fmt.Sprintf("initMysql initialization error for %s", m.Dbname))
	return db
}

// InitMysqlC initializes MySQL Master connection and MySQL Slave connections
func InitMysqlC() {
	// Initialize MySQL Master connection
	global.MdbcHaproxy = initMysqlConnection(global.Config.MySQLHaproxy)
	SetPoolC(global.MdbcHaproxy)

	// Initialize MySQL Slave connections
	//InitMysqlSlaves()

	// Optional: You can add more logic here for handling other slave connections if needed
}

// InitMysqlSlaves initializes all MySQL Slave connections (Slave, Slave2, Slave3)
// func InitMysqlSlaves() {
// 	// Initialize MySQL Slave connection
// 	global.MdbcSlave = initMysqlConnection(global.Config.MySQLSlave)
// 	SetPoolC(global.MdbcSlave)

// 	// Initialize MySQL Slave2 connection
// 	global.MdbcSlave2 = initMysqlConnection(global.Config.MySQLSlave2)
// 	SetPoolC(global.MdbcSlave2)

// 	// Initialize MySQL Slave3 connection
// 	global.MdbcSlave3 = initMysqlConnection(global.Config.MySQLSlave3)
// 	SetPoolC(global.MdbcSlave3)
// }

// SetPoolC configures the connection pool settings for MySQL connections
func SetPoolC(sqlDb *sql.DB) {
	if sqlDb != nil {
		var m *setting.MySQLSetting
		sqlDb.SetConnMaxIdleTime(time.Duration(m.MaxIdleConnes) * time.Second)
		sqlDb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime) * time.Second)
		sqlDb.SetMaxIdleConns(m.MaxOpenConnes)
	}
}
