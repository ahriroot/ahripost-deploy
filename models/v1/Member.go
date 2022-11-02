package model_v1

type Member struct {
	RID        int64   `json:"_id" gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	UserRID    int64   `json:"user_id"`
	User       User    `json:"user" gorm:"foreignKey:UserRID;references:RID"`
	MemberRID  int64   `json:"member_id"`
	Member     User    `json:"member" gorm:"foreignKey:MemberRID;references:RID"`
	ProjectRID int64   `json:"project_id"`
	Project    Project `json:"project" gorm:"foreignKey:ProjectRID;references:RID"`
	Status     int     `json:"status" gorm:"column:status"`
}
