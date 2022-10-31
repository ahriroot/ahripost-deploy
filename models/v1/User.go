package model_v1

type User struct {
	RID      int64  `json:"_id" gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}

type Token struct {
	RID     int64  `json:"_id" gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	UserRID int64  `json:"user_id"`
	User    User   `json:"user" gorm:"foreignKey:UserRID;references:RID"`
	Token   string `json:"token" gorm:"column:token"`
}
