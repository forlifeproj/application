package conf

var GConf = struct {
	DB struct {
		User     string
		PassWord string
		StrIp    string
		Port     int
		Database string
		MaxOpen  int
		MaxIdol  int
	}
}{}
