package cipher

import (
	"goproxy/utils"
)

func GetDriver (iniParser utils.IniParser) (driver Driver, err error) {
	key := iniParser.GetString("common", "cipher")
	if key == "" {
		key = "caesar"
	}

	switch key {
	case "caesar":
		cip := Caesar{}
		cip.SetDis(byte(iniParser.GetInt32("caesar", "dis")))
		driver = Driver(&cip)
	default:
		panic("unknow cipher")
	}

	return driver, nil
}
