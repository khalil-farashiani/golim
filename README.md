# Golim

## Golim is rate golim based on token bucket algorithm

`Golim` is a rate golim program written in Go that allows you to control the frequency and concurrency of requests to your web service. It uses a token bucket algorithm to regulate the incoming traffic and prevent overload or abuse. You can customize the parameters of `Golim` to suit your needs, such as the bucket size, the refill rate, the timeout duration, etc.

`Golim` is useful for web developers who want to protect their web service from excessive or malicious requests, while ensuring a fair and smooth user experience. Golim is also easy to use and integrate with your existing web service, as it only requires a few lines of code and minimal dependencies.

### Dependencies
all dependencies automatically resolve
- sqlc (as query builder)
- redis-go
- ff cli

## TODO (Golim v1):
- [x] initial the project
- [x] add db configuration
- [x] add roles logic to project
- [ ] implement bucket algorithm
- [ ] add cli version
- [ ] add ui version