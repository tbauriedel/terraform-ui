# Configuration

Configuration of `resource-nexus-core` is written in JSON.

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
    "file": "resource-nexus-core.log",
    "level": "info"
  }
}
```

**Reference**:

| Field   | Type   | Required    | Default                   | Description                                                             |
|---------|--------|-------------|---------------------------|-------------------------------------------------------------------------|
| `type`  | string | Yes         | `stdout`                  | Logging backend. Possible values: `stdout`, `file`.                     |
| `file`  | string | Conditional | `resource-nexus-core.log` | Path to the log file. Required if `type` is `file`.                     |
| `level` | string | Yes         | `info`                    | Log verbosity level. Possible values: `debug`, `info`, `warn`, `error`. |

## Listener

Inside the `listener` section, the following settings can be configured.

```json
{
  "listener": {
    "listenAddr": ":4890",
    "readTimeout": "5s",
    "idleTimeout": "30s",
    "tlsEnabled": true,
    "tlsCertFile": "server.crt",
    "tlsKeyFile": "server.key",
    "tlsSkipVerify": false,
    "globalRateLimitBucketSize": 25,
    "globalRateLimitGeneration": 5,
    "ipBasedRateLimitBucketSize": 10,
    "ipBasedRateLimitGeneration": 2
  }
}
```

**Reference**:

| Field                        | Type                   | Required    | Default | Description                                                                                                                |
|------------------------------|------------------------|-------------|---------|----------------------------------------------------------------------------------------------------------------------------|
| `listenAddr`                 | string                 | Yes         | `:4890` | Network address the server listens on. Format: `host:port` (e.g. `:4890`, `127.0.0.1:8080`).                               |
| `readTimeout`                | string (time.Duration) | Yes         | `10s`   | Maximum duration for reading the entire request. Uses Go duration format (e.g. `5s`, `1m`). `0` means no timeout.          |
| `idleTimeout`                | string (time.Duration) | Yes         | `120s`  | Maximum amount of time to wait for the next request when keep-alives are enabled. Go duration format.                      |
| `tlsEnabled`                 | boolean                | Yes         | `false` | Enables TLS for the listener.                                                                                              |
| `tlsCertFile`                | string                 | Conditional | —       | Path to the TLS certificate file (PEM). Required if `tlsEnabled` is `true`.                                                |
| `tlsKeyFile`                 | string                 | Conditional | —       | Path to the TLS private key file (PEM). Required if `tlsEnabled` is `true`.                                                |
| `tlsSkipVerify`              | boolean                | Conditional | `false` | Disables TLS certificate verification. **Use only for testing.**                                                           |
| `globalRateLimitBucketSize`  | float64                | Conditional | `25`    | Number of tokens in the bucket. Check the "Rate Limiting" section in [REST-API docs](./20-REST-API.md) for details.        |
| `globalRateLimitGeneration`  | int                    | Conditional | `5`     | Number of tokens generated per second. Check the "Rate Limiting" section in [REST-API docs](./20-REST-API.md) for details. 
| `ipBasedRateLimitBucketSize` | float64                | Conditional | `10`    | Number of tokens in the bucket. Check the "Rate Limiting" section in [REST-API docs](./20-REST-API.md) for details.        |
| `ipBasedRateLimitGeneration` | int                    | Conditional | `2`     | Number of tokens generated per second. Check the "Rate Limiting" section in [REST-API docs](./20-REST-API.md) for details. 

---

**time.Duration format**:  
Fields of type 'time.Duration' must be specified as strings using Go’s duration format:

- ms – milliseconds
- s – seconds
- m – minutes
- h – hours

Examples:

```
"readTimeout": "10s"
"idleTimeout": "2m"
```