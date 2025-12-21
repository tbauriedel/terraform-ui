# REST API

One of the main parts of the `resource-nexus-core` is the powerful REST API.  
Every workload done by the `resource-nexus-core` is triggered via this.

## General Information

### Logging

Each request is logged by default. Based on the type of your configured logging instance, messages are saved to stdout or a file.  
An example message looks like this:
`{"time":"2025-12-21T17:41:01.285593+01:00","level":"INFO","msg":"new request: [GET] /foobar [::1]:51952 curl/8.7.1"}`

### Rate Limiting

To prevent the `resource-nexus-core` from being overloaded, a **global** rate limiting is implemented. For now, limits are not "per-user", but rather for all requests.  

The limiter is token-based.  
There is a global "bucket" of tokens that are available. Each request consumes one token.  
Tokens are generated in a configurable interval. The bucket also has a configurable size.

The default values are:
- Bucket size: 25 tokens
- Token generation interval: 5 per second