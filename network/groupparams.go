package network

//request command type
type GroupCommandType string

const (
	GroupADD    GroupCommandType = "add"
	GroupDELETE GroupCommandType = "delete"
	GroupGetAll GroupCommandType = "getall"
	GroupGet    GroupCommandType = "get"
	//EXIST  CommandType = "exist"
	//PASS CommandType = "pass"
)
