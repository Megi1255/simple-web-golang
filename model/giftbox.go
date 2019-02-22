package model

import "time"

type GiftBox struct {
	GiftId         int64     `json:"gift_id"`
	SenderId       int64     `json:"sender_id"`
	AmountCurrency int       `json:"amount_currency"`
	SendTime       time.Time `json:"send_time"`
	Content        string    `json:"content"`
}

type Follow struct {
	UserId     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	Created    time.Time `json:"created"`
	State      int       `json:"state"`
}
