package db

type (
	Users struct {
				Id   	int64  `db:"id"`
				Name	string `db:"name"`
				Age		int	   `db:"age"`
	}
	
	RC_User struct {
					Id 		int 	`db:"id"`
					Name 	string 	`db:"name"`
	}
)
