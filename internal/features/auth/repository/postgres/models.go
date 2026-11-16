package auth_postgres_repository

import "time"

type UserIdentityModel struct {
	ID        int
	UserID    int
	Provider  string
	Subject   string
	Email     string
	CreatedAt time.Time
}
