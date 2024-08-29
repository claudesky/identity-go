package models

type User struct {
	Id       string `db:"id"`
	Password string `db:"password" json:"-"`
	Name     string `db:"name"`
	Email    string `db:"email"`
}
