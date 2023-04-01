package tool

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

//  MysqlConfig
//  @Description: Define MySQL related configurations
type MysqlConfig struct {
	dsn           string
	db            *sql.DB
	table         string
	tagKey        string
	savePath      string
	realTableName string
	err           error
}

//  column
//  @Description: Define mysql fields
type column struct {
	ColumnName    string
	DataType      string
	IsNullable    string
	TableName     string
	ColumnComment string
	Tag           string
}

// NewMysqlTableGoStruct
//  @Description: return MysqlConfig
//  @param opts
//  @return *MysqlConfig
func NewMysqlTableGoStruct(opts ...Options) *MysqlConfig {
	tts := &MysqlConfig{}
	for _, opt := range opts {
		opt.Apply(tts)
	}

	return tts
}

// MysqlConnect
//  @Description: connect mysql
//  @receiver mg
func (mg *MysqlConfig) MysqlConnect() {
	if mg.dsn == "" {
		mg.err = errors.New("Mysql.dsn cannot be empty")
		return
	}

	if mg.db == nil {
		mg.db, mg.err = sql.Open("mysql", mg.dsn)
	}

	return
}

// GetMysqlColumns
//  @Description: get mysql columns
//  @receiver mg
//  @return []*column
//  @return error
func (mg *MysqlConfig) GetMysqlColumns() ([]*column, error) {
	var tableColumnList []*column
	sqlStr := "SELECT COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE table_schema = DATABASE()"

	if mg.table == "" {
		return nil, fmt.Errorf("Table name not specified")
	}

	sqlStr += fmt.Sprintf(" AND TABLE_NAME  = '%s' ORDER BY TABLE_NAME ASC,ORDINAL_POSITION ASC", mg.table)

	rows, err := mg.db.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to query table information :%v", err.Error())
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Table structure not found")
	}

	for rows.Next() {
		c := &column{}

		err = rows.Scan(&c.ColumnName, &c.DataType, &c.IsNullable, &c.TableName, &c.ColumnComment)
		if err != nil {
			return nil, fmt.Errorf("Failed to get column fields :%v", err.Error())
		}

		if mg.tagKey != "" {
			c.Tag = fmt.Sprintf("`%v:\"%v\"`", mg.tagKey, c.ColumnName)
		} else {
			c.Tag = fmt.Sprintf("`gorm:\"%v\"`", c.ColumnName)
		}

		c.ColumnName = mg.camelCase(c.ColumnName)
		c.DataType = typeConvert[c.DataType]

		tableColumnList = append(tableColumnList, c)
	}

	return tableColumnList, nil
}

// Execute
//  @Description: execute
//  @receiver mg
//  @return error
func (mg *MysqlConfig) Execute() error {
	mg.MysqlConnect()

	if mg.err != nil {
		return mg.err
	}

	tableColumnList, err := mg.GetMysqlColumns()
	if err != nil {
		return err
	}

	var structContent string
	var tableName string
	for _, tableColumn := range tableColumnList {
		var columnComment string

		if tableColumn.ColumnComment != "" {
			columnComment = fmt.Sprintf(" // %v", tableColumn.ColumnComment)
		}

		if len(tableName) == 0 {
			tableName = tableColumn.TableName
		}

		structContent += fmt.Sprintf("    %s %s %s%s\n", tableColumn.ColumnName, tableColumn.DataType, tableColumn.Tag, columnComment)
	}

	structName := mg.camelCase(tableName)
	structContent = fmt.Sprintf("type %v struct {\n"+structContent+"}\n\n", structName)

	if mg.realTableName != "" {
		structContent += fmt.Sprintf("func (*%v) TableName() string {\n    return \"%v\"\n}\n", structName, tableName)
	}

	if mg.savePath == "" {
		fmt.Println(structContent)
		return nil
	}

	if err := mg.WriteString(structContent); err != nil {
		return err
	}

	return nil
}

// camelCase
//  @Description: Field Conversion
//  @receiver mg
//  @param str
//  @return string
func (mg *MysqlConfig) camelCase(str string) string {
	var text string

	for _, s := range strings.Split(str, "_") {
		text += strings.ToUpper(s[0:1]) + s[1:]
	}

	return text
}

// FileIsExist
//  @Description: Does the file exist
//  @receiver mg
//  @param filename
//  @return bool
func (mg *MysqlConfig) FileIsExist(filename string) bool {
	exist := true

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

// WriteString
//  @Description:
//  @receiver mg
//  @param content
//  @return error
func (mg *MysqlConfig) WriteString(content string) error {
	var (
		file    *os.File
		fileErr error
	)

	if mg.FileIsExist(mg.savePath) {
		file, fileErr = os.OpenFile(mg.savePath, os.O_RDWR|os.O_TRUNC, 0o666)
	} else {
		file, fileErr = os.Create(mg.savePath)
	}

	if fileErr != nil {
		return fmt.Errorf("File operation failed: %v", fileErr)
	}
	defer file.Close()

	_, fileErr = file.WriteString(content)
	if fileErr != nil {
		return fmt.Errorf("File write failed: %v", fileErr)
	}

	return nil
}
