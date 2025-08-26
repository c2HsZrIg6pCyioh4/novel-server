package models

import "time"

// User 用户表映射
type User struct {
	ID            int64     `db:"id" json:"id"`
	OpenID        string    `db:"openid" json:"openid,omitempty"`
	UnionID       string    `db:"unionid" json:"unionid,omitempty"`
	Nickname      string    `db:"nickname" json:"nickname,omitempty"`
	PhoneNumber   int64     `db:"phonenumber" json:"phonenumber,omitempty"`
	Email         string    `db:"email" json:"email,omitempty"`
	RandNum       string    `db:"randnum" json:"randnum,omitempty"`
	EnabledStatus bool      `db:"enabled_status" json:"enabled_status,omitempty"`
	CreateTime    time.Time `db:"createtime" json:"createtime"`
	UpdateTime    time.Time `db:"updatetime" json:"updatetime"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	Username      string    `db:"username" json:"username,omitempty"`
	Sub           string    `db:"sub" json:"sub,omitempty"`
	AppleSub      string    `db:"apple_sub" json:"apple_sub,omitempty"`
}
