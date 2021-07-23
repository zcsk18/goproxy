package proto

import "goproxy/utils"

func GetDriver(iniParser *utils.IniParser) (driver Driver, err error) {
	key := iniParser.GetString("common", "proto")
	if key == "" {
		key = "socks5"
	}

	switch key {
	case "socks5" :
		driver = Driver(&Socks5{})
	case "http" :
		driver = Driver(&Http{})
	default:
		panic("unknow proto")
	}

	return driver, nil
}