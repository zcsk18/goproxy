package proto

import "goproxy/utils"

func GetDriver(iniParser utils.IniParser) (driver Driver, err error) {
	switch iniParser.GetString("common", "proto") {
	case "socks5" :
		driver = Driver(&Socks5{})
	default:
		panic("unknow proto")
	}

	return driver, nil
}