package tool

// NewFuncTableToStructOption
//  @Description: return a funcOptions
//  @param f
//  @return *funcOptions
func NewFuncTableToStructOption(f func(*MysqlConfig)) *funcOptions {
	return &funcOptions{
		f: f,
	}
}

// WithDsn
//  @Description: Specify dsn (required)
//  @param dsn
//  @return Options
func WithDsn(dsn string) Options {
	return NewFuncTableToStructOption(func(mc *MysqlConfig) {
		mc.dsn = dsn
	})
}

// WithTable
//  @Description: Specify a table name (required)
//  @param table
//  @return Options
func WithTable(table string) Options {
	return NewFuncTableToStructOption(func(mc *MysqlConfig) {
		mc.table = table
	})
}

// WithTagKey
//  @Description: Tag tag key, default to gorm tag (not necessary)
//  @param tagKey
//  @return Options
func WithTagKey(tagKey string) Options {
	return NewFuncTableToStructOption(func(mc *MysqlConfig) {
		mc.tagKey = tagKey
	})
}

// WithSavePath
//  @Description: Write file path, do not write and print to the console (not necessary)
//  @param savePath
//  @return Options
func WithSavePath(savePath string) Options {
	return NewFuncTableToStructOption(func(mc *MysqlConfig) {
		mc.savePath = savePath
	})
}

// WithRealTableName
//  @Description: Obtain the real table name (not necessary)
//  @param tableName
//  @return Options
func WithRealTableName(tableName string) Options {
	return NewFuncTableToStructOption(func(mc *MysqlConfig) {
		mc.realTableName = tableName
	})
}
