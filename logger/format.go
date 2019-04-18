package logger

import (
	"encoding/json"
	"time"
)

type Format map[string]interface{}

func NewFormat(v ...interface{}) string {
	format := make(Format)
	format["time"] = time.Now()

	// switch logger.level {
	// case 0:
	// 	format["level"] = "INFO"
	// case 1:
	// 	format["level"] = "WARN"
	// case 2:
	// 	format["level"] = "ERROR"
	// default:
	// 	format["level"] = "INFO"
	// }

	format["message"] = v

	mjson, _ := json.Marshal(&format)
	mString := string(mjson)

	return mString
}
