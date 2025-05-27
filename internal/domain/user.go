package domain

import "time"

type User struct {
	Id       int64
	Uid      uint64
	Email    string
	Password string

	NickName string
	Phone    string
	AboutMe  string
	Ctime    time.Time
	Birthday time.Time
}
