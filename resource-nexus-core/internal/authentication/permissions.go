package authentication

// permissions return the permissions map.
func permissions() map[string]string {
	return map[string]string{
		"/system/health":  "system:health:get",
		"/auth/user/add":  "auth:user:create",
		"/auth/group/add": "auth:group:create",
	}
}

// GetPermissionForPath returns the permission for the given path / url.
func GetPermissionForPath(url string) (string, bool) {
	perm, ok := permissions()[url]
	if !ok {
		return "", false
	}

	return perm, true
}

// BuildPermissionString builds a permission string from category, action and resource.
//
// Format: category:action:resource.
func BuildPermissionString(category, resource, action string) string {
	return category + ":" + resource + ":" + action
}
