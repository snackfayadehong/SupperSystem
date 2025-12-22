package clientDb

import (
	"SupperSystem/configs"
	logger2 "SupperSystem/pkg/logger"
	"fmt"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

func Init(rootPath string) error {
	var DbPwd string
	DbPwd, err := conf.DecryptionPwd(rootPath)
	if err != nil {
		return err
	}
	//log := logger.New(newMyWriter(), logger.Config{LogLevel: logger.Info})
	log := logger2.NewGormCustomLogger()
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&loc=%s&encrypt=disable", userName, password, ipAddr, port, dbName, loc)
	// dsn := "sqlserver://sa:密码@127.0.0.1:1433?database=dbStatus&encrypt=disable"
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=disable", conf.Configs.DBClient.Username, DbPwd,
		conf.Configs.DBClient.IP, conf.Configs.DBClient.DbName)
	DB, _ = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger:      log,
		PrepareStmt: true, // 执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "TB_",
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	sqlDb, err := DB.DB()
	sqlDb.SetMaxOpenConns(20)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(10 * time.Minute)
	sqlDb.SetConnMaxIdleTime(time.Hour)
	return nil
}
