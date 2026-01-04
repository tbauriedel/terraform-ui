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

type UserGroupReference struct {
	Username  string `json:"username"`
	GroupName string `json:"group_name"`
	UserID    int    `json:"user_id"`
	GroupID   int    `json:"group_id"`
}

type GroupPermissionReference struct {
	GroupName    string `json:"group_name"`
	Permission   string `json:"permission"`
	GroupID      int    `json:"group_id"`
	PermissionID int    `json:"permission_id"`
}
