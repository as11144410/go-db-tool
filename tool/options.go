package tool

//  Options
//  @Description: Define the Options interface
type Options interface {
	Apply(*MysqlConfig)
}

//  funcOptions
//  @Description: Define the funcOptions structure
type funcOptions struct {
	f func(*MysqlConfig)
}

// Apply
//  @Description: Define the Apply method for funcOptions
//  @receiver fo
//  @param mc
func (fo funcOptions) Apply(mc *MysqlConfig) {
	fo.f(mc)
}
