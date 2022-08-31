package domain

const (
	AuthCookie = "auth"
)

type Cookie struct {
	Name   string
	Token  string
	MaxAge int
}
