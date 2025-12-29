package authentication

type Permission struct {
	Action   string
	Resource string
}

func BuildPermissionString(action, resource string) string {
	return action + ":" + resource
}
