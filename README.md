# SwiftLink

Imaginary URL shortener service written in Go.

## Run the server

Change the `SWIFTLY_PREFIX` and `SWIFTLY_PORT` environment variables to your liking.
THe SWIFTLY_PREFIX is the prefix of the shortened link, and the SWIFTLY_PORT is the port the server will run on.

```bash
SWIFTLY_PREFIX="http://localhost:8080/" SWIFTLY_PORT="8080" go run main.go
```

## Examples Query using cURL

```bash
# Shorten the long link
curl -X POST http://localhost:8080/shorten -d "{\"url\":\"https://google.com/longlink\"}"
# Take the shortened link and try out the redirect
curl http://localhost:8080/<shortened-link>
```
