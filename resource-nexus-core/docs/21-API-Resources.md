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

Necessary permission: `auth:user:create`

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

Necessary permission: `auth:group:create`

`POST /auth/group/add -d '{"name":"default-users"}`: Creates a new group.

Body:
- `name`: Name of the group

Example response:
```json
{
  "message": "entity created successfully"
}
```