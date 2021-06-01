package util

import (
	"fmt"
	"os"
	"strings"
)

func StrFirstToUpper(str string) string{
	if len(str) < 1{
		return ""
	}

	strArray := []rune(str)

	if strArray[0] >= 97 && strArray[0] <= 122 {
		strArray[0] -= 32
	}

	return string(strArray)
}

//将下划线连接的字串改为驼峰
//abc_def --> AbcDef
func StrCamel(str string) string{
	tmp := strings.Split(str, "_")
	var result string
	for _,  v := range tmp{
		result += StrFirstToUpper(v)
	}

	return result
}

//判断文件或路径是否存在
func PathExists(path string) bool{
	if _, err := os.Stat(path); os.IsNotExist(err){
		return false
	}

	return true
}

//获取tpl所在目录
func GetTplPath() string{
	dir, _ := os.Getwd()

	return fmt.Sprintf("%s../tpl", dir)
}