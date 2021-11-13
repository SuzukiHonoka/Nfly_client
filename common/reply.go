package common

type Reply struct {
	Action string      `json:"action"`
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}