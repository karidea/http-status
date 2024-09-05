# http-status

A fast command-line http status tool that returns the http response status codes and paths for a given list of url endpoints. Optionally it will resolve domain names to a specific IP address for testing.

## Installation

   ```bash
   git clone https://github.com/karidea/http-status.git
   cd http-status
   go build http-status.go
   ```

## Usage

``` bash
Usage of ./http-status:
  -file string
        Path to the JSON file containing the https URL endpoints (required)
  -ip string
        IP address to resolve the domain to (optional)
```

The JSON file should contain an array of endpoint URLs you want to test. Example:

```json
[
    "https://example.com/healthcheck",
    "https://api.example.com/liveness"
]
```


```bash
./http-status --file=path/to/endpoints.json
```

Optionally, specify an IP address to resolve the domain names to, for testing endpoint liveness from a specific IP:

```bash
./http-status --file=path/to/endpoints.json --ip=203.0.113.5
```

It will output the HTTP status codes and paths for each endpoint

```bash
404 /healthcheck
200 /liveness
```
