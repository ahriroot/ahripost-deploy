package model_v1

type Project struct {
	RID      int64  `json:"_id" gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	ID       int64  `json:"id" gorm:"column:id"`
	UserRID  int64  `json:"user_id"`
	User     User   `json:"user" gorm:"foreignKey:UserRID;references:RID;constraint:OnDelete:SET NULL"`
	Key      string `json:"key" gorm:"column:key;uniqueIndex"`
	Name     string `json:"name" gorm:"column:name"`
	CreateAt int64  `json:"create_at" gorm:"column:create_at"`
	Public   bool   `json:"public" gorm:"column:public;default:false"`
}

type Item struct {
	RID        int64   `json:"_id" gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	ID         int64   `json:"id" gorm:"column:id"`
	Key        string  `json:"key" gorm:"column:key"`
	Label      string  `json:"label" gorm:"column:label"`
	Type       string  `json:"type" gorm:"column:type"`
	From       string  `json:"from" gorm:"column:from"`
	ProjectRID string  `json:"project_id"`
	Project    Project `json:"project" gorm:"foreignKey:ProjectRID;references:Key"`
	UserRID    int64   `json:"user_id"`
	User       User    `json:"user" gorm:"foreignKey:UserRID;references:RID;constraint:OnDelete:SET NULL"`
	Tag        bool    `json:"tag" gorm:"column:tag"`
	Client     string  `json:"client" gorm:"column:client"`
	Parent     string  `json:"parent" gorm:"column:parent;default:''"`
	LastSync   int64   `json:"last_sync" gorm:"column:last_sync;default:0"`
	LastUpdate int64   `json:"last_update" gorm:"column:last_update"`
	Template   string  `json:"template" gorm:"column:template"`
	Request    string  `json:"request" gorm:"column:request"`
	Response   string  `json:"response" gorm:"column:response"`
}
