package models

type User struct {
	FullName string `bson:"full_name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
