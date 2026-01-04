# API Resources

Documentation of the API resources provided by `resource-nexus-core`.

### /system/health

Necessary permission: `system:health:get`

`GET /system/health`: Prints health information about the `resource-nexus-core` instance.

Example response:
```json
{
  "databaseStatus": true,
  "databaseStatusMessage": "OK",
  "version": "0.1.0"
}
```

### /auth/user/add

Necessary permission: `auth:user:add`

`POST /auth/user/add -d '{"name":"hercules","password_hash":"argon2id-hash","is_admin":false}'`: Creates a new user.  

Body:
- `name`: Name of the user
- `password_hash`: argon2id password hash
- `is_admin`: Boolean value indicating whether the user is an admin. 

Example response:
```json
{
  "message":"entity created successfully"
}
```

### /auth/group/add

Necessary permission: `auth:group:add`

`POST /auth/group/add -d '{"name":"default-users"}`: Creates a new group.

Body:
- `name`: Name of the group

Example response:
```json
{
  "message": "entity created successfully"
}
```

### /auth/usergroup/add

Necessary permission: `auth:usergroup:add`

`POST /auth/usergroup/add -d '{"username":"hercules","group_name":"default-users"}`: Create a new user greoup reference. Adds the user to the group.

Body:
- `username`: ID of the user
- `group_name`: ID of the group

Example response:
```json
{
  "message": "user group reference added"
}
```

### /auth/grouppermission/add

Necessary permission: `auth:grouppermission:add`

`POST /auth/grouppermission/add -d '{"group_name":"default-users","permission":"auth:user:add"}`: Adds a permission to a group.

Body:
- `group_name`: Name of the group
- `permission`: Permission to add

Example response:
```json
{
  "message": "permission group reference added"
}
```