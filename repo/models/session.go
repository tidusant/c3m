package models

import (
	"time"
)

type Session struct {
	Time     time.Time
	UserID   string
	ShopID   string
	UserName string
}
