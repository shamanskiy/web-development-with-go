To build and run web app:

```
go run main.go
```

To start Postgres image:

```
docker compose up
```

To start Postgres image in detached mode:

```
docker compose up -d
```

To start Postgres image use Ctrl+C or for detached mode:

```
docker compose stop
```

To connect to DB with pgsl CLI in Docker:

```
docker exec -it web-development-with-go-db-1 /usr/bin/psql -U baloo -d lenslocked
```
