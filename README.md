# 1. 功能
根据数据库表生成gorm需要的struct。支持指定单表生成，也可以全库生成。

比如有如下数据表：
```mysql
Table: user
Create Table: CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `age` int(11) NOT NULL DEFAULT '-1',
  `is_admin` tinyint(1) NOT NULL DEFAULT '0',
  `is_valid` tinyint(1) NOT NULL DEFAULT '1',
  `login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8 COMMENT='用户表'
```

db2gorm可以在指定的目录下生成 user/user.go，内容如下：

```
package user

type User struct{
    ID uint32 
    Name string 
    Age int `gorm:"default:-1"`
    IsAdmin bool `gorm:"default:0"`
    IsValid bool `gorm:"default:1"`
    LoginTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
```

# 2. 使用
## 2.1 指定单表生成文件
```go
package main

import (
    "github.com/qmhball/db2gorm/gen"
)

dsn := "root:user237@tcp(10.10.10.237:3306)/mydb?charset=utf8&parseTime=true&loc=Local"
        
//生成指定单表
tblName := "User"
gen.GenerateOne(gen.GenConf{
    Dsn:       dsn,
    WritePath: "./model",
    Stdout:    false,
    Overwrite: true,
}, tblName)
```

gen.GenConf的说明如下：
- Dsn：数据库配置
- WritePath：指定文件写入的目录。生成前目录一定要存在。上例中会在model目录下生成user/user.go
- Stdout：是否输出至标准输出。如果Stdout为true，则生成的struct不会写入文件。
- Overwrite：当原文件存在时，是否进行覆盖。true为覆盖。


## 2.2 全库生成
```
    gen.GenerateAll(gen.GenConf{
        Dsn:       dsn,
        WritePath: "./model",
        Stdout: false,
        Overwrite: true,
    })
```

# 3. 推荐的数据库配置
建议使用gorm v1.2以上版本(v1.1*的版本和v1.2差别比较大)

建议gorm.Open时指定SingularTable为true，即使用单数表名。这样就不必在struct上定义TableName方法指定表名了。

```go
        gorm.Open(mysql.New(mysql.Config{
			DSN:dsn,
		}), &gorm.Config{
			NamingStrategy:schema.NamingStrategy{
			    //单数表名
				SingularTable: true,
			},
		})
```

# 4. 生成规则
## 4.1 表名
表名对应的大驼峰命名做为struct名，全小写表名做为目录名，文件名和包名。
比如表名为demo_test, 则：
- struct名：DemoTest
- 目录名：demotest
- 文件名：demotest.go
- 包名：demotest


## 4.2 字段名
表名对应的大驼峰命名做为struct字段名。比如is_admin对应IsAdmin。但id字段除名，id固定对应ID。



