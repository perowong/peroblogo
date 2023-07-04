# Getting start

Development getting start

```bash
docker compose -f docker-compose.local-mysql.yaml up -d --no-deps
go run main.go
```

Production deploy

```shell
docker compose -f docker-compose.prod.yaml build &&
docker compose -f docker-compose.prod.yaml up -d --no-deps
```
