# Getting start

Start locally for development
```shell
go run main.go
```

Deploy in test env
```shell
docker compose -f docker-compose.test.yaml build && docker compose -f docker-compose.test.yaml up -d
```
