package clientDb

import (
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB       *gorm.DB
	userName = "sa"
	password = "20010126"
	ipAddr   = "127.0.0.1"
	// port     = 3306
	dbName = "WZ_MY"
	loc    = "local"
)

func InitDb() error {
	var err error
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&loc=%s&encrypt=disable", userName, password, ipAddr, port, dbName, loc)
	// dsn := "sqlserver://sa:密码@127.0.0.1:1433?database=dbStatus&encrypt=disable"
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=disable", userName, password, ipAddr, dbName)
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "TB_",
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	return err
}
