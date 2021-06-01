package db

import (
	"gorm.io/gorm/schema"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql(dsn string) (err error){
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

