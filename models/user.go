package models

type (
	User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age int 	`json:"age"`
	}
)

var (
	Users = map[int]*User{}
	Seq = 1
)
