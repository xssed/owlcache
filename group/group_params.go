package group

//request command type
type GroupCommandType string

const (
	GroupADD    GroupCommandType = "add"
	GroupDELETE GroupCommandType = "delete"
	GroupGetAll GroupCommandType = "getall"
	GroupGet    GroupCommandType = "get"
)

//response status type
type ResStatus int

//response 状态
const (
	SUCCESS         ResStatus = 200
	ERROR           ResStatus = 500
	NOT_FOUND       ResStatus = 404
	UNKNOWN_COMMAND ResStatus = 501
	NOT_PASS        ResStatus = 401
)

//status to string
func ResStatusToString(resstatus ResStatus) string {

	var s string
	switch resstatus {
	case SUCCESS:
		s = "SUCCESS"
	case ERROR:
		s = "ERROR"
	case NOT_FOUND:
		s = "NOT_FOUND"
	case UNKNOWN_COMMAND:
		s = "UNKNOWN_COMMAND"
	case NOT_PASS:
		s = "NOT_PASS"
	}

	return s
}
