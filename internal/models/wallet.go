package models

type Wallet struct {
	Id       int     `db:"id" json:"id"`
	UserId   int     `db:"user_id" json:"userId"`
	Balance  float64 `db:"balance" json:"balance"`
	Currency string  `db:"currency" json:"currency"`
}
