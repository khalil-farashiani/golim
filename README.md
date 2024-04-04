# Golim

## Golim is rate golim based on token bucket algorithm

`Golim` is a rate limiter, `Golim` program written in Go that allows you to control the frequency and concurrency of requests to your web service. It uses a token bucket algorithm to regulate the incoming traffic and prevent overload or abuse. You can customize the parameters of `Golim` to suit your needs, such as the bucket size, the refill rate, endpoint customization, etc.

`Golim` is useful for web developers who want to protect their web service from excessive or malicious requests, while ensuring a fair and smooth user experience. Golim is also easy to use and integrate with your web service with any language like C#, PHP, JS, Python, Golang etc. as it only requires a minimal dependencies.

### Dependencies
all dependencies automatically resolve
- sqlc (as query builder)
- redis-go
- ff cli
- robfig-cron

### Usage

first of all export redis server address with
```bash
export REDIS_URI=redis://localhost:6379
```

- #### install golim
```bash
go install github.com/khalil-farashiani/golim@latest
```

- #### initial golim limiter
```bash
golim init -n <limiter id> -d <destination address without scheme like google.com or 8.8.8.8>
```

- #### add a role to specific limiter
```bash
golim add -l <limiter id> -e <endpoint address like this /users/> -b <bucket size> -a <add token rate per minute> -i <initial token>
```

- ####  remove a role from specific limiter
```bash
golim remove -i <role id>
```

- #### list all role from specific limiter
```bash
golim get -l
```

- #### remove a limiter
```bash
golim removel -l <limiter id>
```
- #### remove a role
```bash
golim remove -i <role id>
```

``all flags have alternative``
- -n <kbd>→</kbd> --name  
- -p <kbd>→</kbd> --port
- -l <kbd>→</kbd> --limiter
- -n <kbd>→</kbd> --name
- -d <kbd>→</kbd> --destination
- -e <kbd>→</kbd> --endpoint
- -b <kbd>→</kbd> --bsize
- -a <kbd>→</kbd> --add_token
- -i <kbd>→</kbd> --initial_token

## TODO features
- [ ] add default limiter
- [ ] add regex
- [ ] add ui version
- [ ] make service open failed
