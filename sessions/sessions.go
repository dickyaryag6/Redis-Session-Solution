package sessions

var SessionsStore Store

type Session struct {
	Name string `json:"name"`
	UserID int `json:"userID"`
}

type Store interface {
	Get(string) (Session, error)
	Set(string, Session) error
}
