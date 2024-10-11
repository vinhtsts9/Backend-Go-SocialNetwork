package initialize

import (
	"database/sql"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"time"

	"go.uber.org/zap"
)

func checkErrorPanicC(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}

}
func InitMysqlC() {
	m := global.Config.Mysql

	dsn := "%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, m.Username, m.Password, m.Host, m.Port, m.Dbname)
	db, err := sql.Open("mysql", s)
	checkErrorPanicC(err, "initMysql initialization error")

	global.Mdbc = db

	SetPoolC()
	// migrationTables()
}

func SetPoolC() {
	m := global.Config.Mysql
	sqlDb, err := global.Mdb.DB()
	if err != nil {
		fmt.Printf("mysql error: %s::", err)
	}
	sqlDb.SetConnMaxIdleTime(time.Duration(m.MaxIdleConnes))
	sqlDb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))
	sqlDb.SetMaxIdleConns(m.MaxOpenConnes)
}

// func migrationTables() {
// 	err := global.Mdb.AutoMigrate(
// 		&po.User{},
// 		&po.Role{},
// 	)
// 	if err != nil {
// 		fmt.Println("Migrating tables error:", err)
// 	}
// }
