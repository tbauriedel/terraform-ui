# Configuration

Configuration of `terraform-ui` is written in JSON.

The default configuration file is located under `config.json`.  
The config file is provided as the command line argument `--config <path to file>`.

An example configuration file can be found [here](./examples/config.json).

# Settings

The configuration settings are separated into logical sections. Each section is located at the top level of the JSON
file.

## Logging

Inside the `logging` section, the following settings can be configured.

```json
{
  "logging": {
    "type": "stdout",
    "file": "terraform-ui-core.log",
    "level": "info"
  }
}
```

**Reference**:

| Field   | Type   | Required    | Default            | Description                                                             |
|---------|--------|-------------|--------------------|-------------------------------------------------------------------------|
| `type`  | string | Yes         | `stdout`           | Logging backend. Possible values: `stdout`, `file`.                     |
| `file`  | string | Conditional | `terraform-ui.log` | Path to the log file. Required if `type` is `file`.                     |
| `level` | string | Yes         | `info`             | Log verbosity level. Possible values: `debug`, `info`, `warn`, `error`. |

## Listener

Inside the `listener` section, the following settings can be configured.

```json
{
  "listener": {
    "listen_addr": ":4890",
    "read_timeout": "5s",
    "idle_timeout": "30s",
    "tls_enabled": true,
    "tls_cert_file": "server.crt",
    "tls_key_file": "server.key",
    "tls_skip_verify": false
  }
}
```

**Reference**:

| Field             | Type                   | Required    | Default | Description                                                                                                       |
|-------------------|------------------------|-------------|---------|-------------------------------------------------------------------------------------------------------------------|
| `listen_addr`     | string                 | Yes         | `:4890` | Network address the server listens on. Format: `host:port` (e.g. `:4890`, `127.0.0.1:8080`).                      |
| `read_timeout`    | string (time.Duration) | Yes         | `10s`   | Maximum duration for reading the entire request. Uses Go duration format (e.g. `5s`, `1m`). `0` means no timeout. |
| `idle_timeout`    | string (time.Duration) | Yes         | `120s`  | Maximum amount of time to wait for the next request when keep-alives are enabled. Go duration format.             |
| `tls_enabled`     | boolean                | Yes         | `false` | Enables TLS for the listener.                                                                                     |
| `tls_cert_file`   | string                 | Conditional | —       | Path to the TLS certificate file (PEM). Required if `tls_enabled` is `true`.                                      |
| `tls_key_file`    | string                 | Conditional | —       | Path to the TLS private key file (PEM). Required if `tls_enabled` is `true`.                                      |
| `tls_skip_verify` | boolean                | Conditional | `false` | Disables TLS certificate verification. **Use only for testing.**                                                  |

---

**time.Duration format**:  
Fields of type 'time.Duration' must be specified as strings using Go’s duration format:

- ms – milliseconds
- s – seconds
- m – minutes
- h – hours

Examples:
```
"read_timeout": "10s"
"idle_timeout": "2m"
```