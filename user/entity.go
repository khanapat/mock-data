package user

import (
	"context"
	"time"
)

type User struct {
	ID        *int       `db:"id" json:"id" exmpale:"1"`
	FirstName *string    `db:"first_name" json:"firstName" example:"Khanapat"`
	LastName  *string    `db:"last_name" json:"lastName" example:"Apiwattanawong"`
	Phone     *string    `db:"phone" json:"phone" example:"0859223735"`
	Email     *string    `db:"email" json:"email" example:"k.apiwattanawong@gmail.com"`
	Balance   *float64   `db:"balance" json:"balance" example:"0.5"`
	DateTime  *time.Time `db:"date_time" json:"datetime" example:"2021-01-02 12:13:14"`
}

type UserRepository interface {
	QueryUser(context.Context, map[string]interface{}) (*[]User, error)
}
