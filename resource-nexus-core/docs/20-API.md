# REST API

One of the main parts of the `resource-nexus-core` is the powerful REST API.  
Every workload done by the `resource-nexus-core` is triggered via this.

### Authentication

The REST API is secured via basic authentication.  
To handle the authentication and validate credentials, users and encoded password hashes are stored inside the database.
The used algorithm to hash passwords is `argon2id`. Plaintext passwords are never stored or printed!

An authentication process works like this:

- User sends username and password to the REST API as part of the header (base64 encoded)
- Check if the user exists
- Take the provided password and hash it with the stored hash algorithm
- Compare the hashed password with the stored hash

If any of these steps fails, the authentication fails, and the request is rejected with a `401 Unauthorized`. The user
itself will only get an http error and no proper error message. More details why the authentication failed can be found
in the logs.

#### Permissions

Permissions are granted with using the implement RBAC. Permissions are never granted to the user directly but to a
group.  
A user can be part of multiple groups. All permissions granted by the groups are merged into one set of permissions.

Permissions are structured like this: `category:resource:action`. For example: `security:user:get`

Available permissions are:

- `system:health:get`: Get general health status of resource-nexus-core

Admin users are flagged inside the database and have **all permissions**. Here, the permission is set on the user itself
and no group is used.

### Logging

Each request is logged by default. Based on the type of your configured logging instance, messages are saved to stdout
or a file.  
An example message looks like this:  
`{"time":"2025-12-21T17:41:01.285593+01:00","level":"INFO","msg":"new request: [GET] /foobar [::1]:51952 curl/8.7.1"}`

### Rate Limiting

To prevent the `resource-nexus-core` from being overloaded, global and ip-based rate limiting is implemented.

The rate limiters are token-based, where one request needs one token.
To store available tokens, "buckets" are used. There is one "global" bucket for global rate limiting and one bucket per
ip address.
When no token is left inside the desired bucket, the request is rejected with a `429 Too Many Requests` response.

Tokens are generated in a configurable interval. The bucket also has configurable sizes. See
the [configuration reference](./10-Config.md) for more information.

First, the global rate limiter is applied. When the global limit is reached, it will reject the request.  
If the global limit is not reached, the ip-based rate limiter is applied.

The default values for the rate limiters can be found inside the [configuration reference](./10-Config.md).
