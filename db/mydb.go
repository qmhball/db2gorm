package db

import (
	"gorm.io/gorm/schema"

	//"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql(dsn string) error{
	//如果不声明err而是使用DB, err:=gorm.Open, DB和err都会被作为局部变量...
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
			//DSN: viper.GetString("mysql.dns"), // DSN data source name
			DSN:dsn,
		}), &gorm.Config{
			NamingStrategy:schema.NamingStrategy{
				SingularTable: true,
			},
		})

		if err != nil {
			return err
		}

		return nil
}

