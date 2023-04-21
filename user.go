package expenses

type User struct{
	Id 				int	   `json:"-"`
	UserName 		string `json:"username"`
	First_Last_Name string `json:"first_last_name"`
	ChatId 			int    `json:"chatID"`
}