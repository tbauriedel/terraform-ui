package database

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	IsAdmin      bool   `json:"is_admin"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Permission struct {
	ID       int
	Category string
	Resource string
	Action   string
}

type UserGroup struct {
	UserID  int
	GroupID int
}

type GroupPermission struct {
	GroupID      int
	PermissionID int
}
