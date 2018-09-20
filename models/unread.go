package models

type Unread struct {
	UID     string `gorm:"column:uid" json:"uid"`
	GID     int64  `gorm:"column:gid" json:"gid"`
	LastMID string `gorm:"column:last_mid" json:"last_mid"`
}

func (Unread) TableName() string {
	return "unread"
}

func (u *Unread) Save() {
}
