package gen

import (
	"fmt"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"github.com/qmhball/db2gorm/util"
)

// show full columns 获取的原始数据
type Column struct {
	Field   string  `gorm:"column:Field"`
	Type    string  `gorm:"column:Type"`
	Key     string  `gorm:"column:Key"`
	Desc    string  `gorm:"column:Comment"`
	Null    string  `gorm:"column:Null"`
	Default string `gorm:"column:Default"`
}
/*
	Field:age
	Type:int(11)
	Key:
	Desc:
	Null:NO
	Default:-1
*/

type Columns []Column
/*
 * 将原始的Column信息转为ColumnInfo
 */
func (cs Columns)change2Info() ([]ColumnInfo, error){
	num := len(cs)
	csInfo := make([]ColumnInfo, num)

	for idx, one := range cs{
		if err := csInfo[idx].set(one); err != nil{
			return nil, err
		}
	}

	return csInfo, nil
}


/*获取DB.tblName的表结构*/
func (cs *Columns)getTableColumns(orm *gorm.DB, tblName string) error{
	// Get table annotations.获取表注释
	sql := fmt.Sprintf("show FULL COLUMNS from `%s`", tblName)
	res := orm.Raw(sql).Scan(cs)
	if res.Error != nil{
		return res.Error
	}

	//fmt.Printf("%+v", list)
	return nil
}

//将Column转为生成模板需要的数据
type ColumnInfo struct {
	Field   string
	Type    string
	Default string
}

/*用一个原始的Column填充一个ColumnInfo*/
func (c *ColumnInfo)set(src Column) (err error){
	c.setField(src.Field)
	if err := c.setType(src.Type); err != nil {
		return err
	}

	c.setDefault(src.Default)

	return nil
}

/*将Field字段改成首字母大写的驼峰字串*/
func (c *ColumnInfo)setField(name string){
	tmp := strings.Split(name, "_")
	for _,  v := range tmp{
		if v == "id" {
			c.Field += "ID"
		}else {
			c.Field += util.StrFirstToUpper(v)
		}
	}
}

//sql中的type与go的类型对应
func (c *ColumnInfo)setType(typeName string) error{
	//精确匹配
	if v, ok := TypeMysqlDicMp[typeName]; ok {
		c.Type = v
		return nil
	}

	// Fuzzy Regular Matching.模糊正则匹配
	for _, l := range TypeMysqlMatchList {
		if ok, _ := regexp.MatchString(l.Key, typeName); ok {
			c.Type = l.Value
			return nil
		}
	}

	return fmt.Errorf("no type for src typeName:%s", typeName)
}

//默认值字段
func (c *ColumnInfo)setDefault(defaultVal string){
	if defaultVal != "" {
		c.Default = fmt.Sprintf("`gorm:\"default:%s\"`", defaultVal)
	}else {
		c.Default = ""
	}
}

//获取DB.tblName对应的表模板信息
//外部只需要这个方法
func GetTableColumnsInfo(orm *gorm.DB, tblName string) ([]ColumnInfo, error){
	var columns Columns
	if err := columns.getTableColumns(orm, tblName); err != nil{
		return nil, err
	}

	columnsInfo, err := columns.change2Info()
	if err != nil{
		return nil, err
	}

	return columnsInfo, nil
}