package ContextKey

type contextKey string
type contextKeyInt int

const(
	UserId contextKey="user_id"
	UserRole contextKeyInt=iota
	UserPassword contextKey="user_password"
)