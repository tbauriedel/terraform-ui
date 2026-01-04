# Getting Started

## Installation

TODO

## Initial Admin Setup

The initial admin user is not created automatically.  
The user needs to be created manually, by using the `resource-nexus-admin` command.

Username is `admin`.

Flags that are needed:
- `-config`: Path to the config file. (Needed for the database connection)
- `-admin-password`: The password for the admin user

Example: `resource-nexus-admin -config"test/testdata/config/config_test_tls.json" -admin-password="admin"`
```
{"time":"2026-01-04T14:33:07.5019+01:00","level":"INFO","msg":"loading config from file","file":"test/testdata/config/config_test_tls.json"}
{"time":"2026-01-04T14:33:07.502515+01:00","level":"INFO","msg":"initializing database connection"}
{"time":"2026-01-04T14:33:07.529891+01:00","level":"INFO","msg":"database connection established and tested successfully"}
{"time":"2026-01-04T14:33:07.80169+01:00","level":"INFO","msg":"admin user created. username: admin"}
{"time":"2026-01-04T14:33:07.801724+01:00","level":"INFO","msg":"closing database connection"}
```