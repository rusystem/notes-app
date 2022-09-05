package domain

import "time"

const (
	USER    = "user"
	NOTE    = "note"
	SIGN_IN = "signIn"
	SIGN_UP = "signUp"
	CREATE  = "create"
	GET     = "get"
	UPDATE  = "update"
	DELETE  = "delete"
)

type LogItem struct {
	Entity    string    `json:"entity"`
	Action    string    `json:"action"`
	EntityID  int64     `json:"entityID"`
	Timestamp time.Time `json:"timestamp"`
}
