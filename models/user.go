package models

import "time"

type User struct {
	ID        int
	Phone     string
	Password  string
	CreatedAt time.Time
	LastIp    string
	Email     string
	Status    int
}

func (u User) TableName() string {
	return "user"
}
