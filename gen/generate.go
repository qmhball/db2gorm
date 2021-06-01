package gen

import (
	"github.com/qmhball/db2gorm/db"
	"github.com/qmhball/db2gorm/util"
	"github.com/qmhball/db2gorm/tpl"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strings"
	"text/template"
)


type GenConf struct {
	Dsn string	//数据库配置
	WritePath string	//生成文件路径
	Stdout bool	//true时只输出至标准输出
	Overwrite bool //true则覆盖原文件
}

func(g GenConf)isValid() error{
	if g.Dsn == ""{
		return fmt.Errorf("dsn is empty")
	}

	if g.Stdout == false{
		if g.WritePath == ""{
			return fmt.Errorf("when Stdout is false WritePath can't be empty")
		}

		if util.PathExists(g.WritePath) == false{
			return fmt.Errorf("path: %s not exists", g.WritePath)
		}
	}

	return nil
}

func GenerateAll(conf GenConf) (err error){
	err = conf.isValid();
	if err != nil {
		fmt.Printf("GenerateAll err:%s\n", err)
		return err
	}

	err = db.InitMysql(conf.Dsn)
	if err != nil{
		fmt.Printf("mysql init err:%s\n", err)
		return err
	}

	var tbls Tables
	if err = tbls.GetTables(db.DB); err != nil{
		fmt.Printf("GetTables error:%s\n", err)
		return err
	}

	for _, tblName := range tbls{
		err = doGenerateOne(db.DB, conf, tblName)
		if err != nil{
			fmt.Printf("GenerateOne for table %s err:%s\n", tblName, err)
			return err
		}
	}

	fmt.Printf("GenerateAll Succeed\n")
	return nil
}

func GenerateOne(conf GenConf, tblName string) (err error){
	err = conf.isValid();
	if err != nil {
		fmt.Printf("GenerateOne for table %s err:%s\n", tblName, err)
		return err
	}

	err = db.InitMysql(conf.Dsn)
	if err != nil{
		return fmt.Errorf("mysql init err:%s", err)
	}

	err = doGenerateOne(db.DB, conf, tblName)
	if err != nil{
		fmt.Printf("GenerateOne for table %s err:%s\n", tblName, err)
	}

	return nil
}

func doGenerateOne(orm *gorm.DB, conf GenConf, tblName string)(err error){
	info, err := GetTableInfo(orm, tblName)
	if err != nil{
		return err
	}

	tpl, err := template.New("struct").Parse(tpl.StructTpl)
	if err != nil{
		return err
	}

	//只输出至标准输出的情况
	if conf.Stdout {
		err = tpl.Execute(os.Stdout, info)
		if err != nil{
			return err
		}

		fmt.Printf("\n\n")
		return nil
	}

	//输出至文件的情况
	path, err := mkDir(conf.WritePath, info.DirName)
	if err != nil{
		return fmt.Errorf("mkdir %s err:%s", path, err)
	}

	fName := mkFileName(path, info.PackageName)
	if util.PathExists(fName) && conf.Overwrite == false {
		return fmt.Errorf("file :%s is exists", fName)
	}

	fp, err := os.Create(fName)
	if err != nil{
		return err
	}

	err = tpl.Execute(fp, info)
	if err != nil{
		return err
	}

	fmt.Printf("Generate Table:%s to %s\n", info.TableName, fName)

	return nil
}

//拼接路径和文件名
func mkFileName(path string, name string) string{
	path = strings.TrimRight(path, "/")
	return fmt.Sprintf("%s/%s.go", path, name)
}

//创建并返回文件所在目录
func mkDir(base string, tblDir string) (string, error){
	//base = strings.TrimRight(base, "/")
	fullDir := fmt.Sprintf("%s/%s", base, tblDir)

	err := os.MkdirAll(fullDir, 0766)
	return fullDir, err
}