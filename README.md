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

To stop Postgres image use Ctrl+C or for detached mode:

```
docker compose stop
```

To stop and kill Postgres image use

```
DON'T USE IT UNLESS YOU HAVE TO!
docker compose down
```

To connect to DB with pgsl CLI in Docker:

```
docker exec -it web-development-with-go-db-1 /usr/bin/psql -U baloo -d lenslocked
```

To run goose migration (up/down/status/reset)

```
goose postgres \
"host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable" \
status
```

```
goose fix
```

Fill in .env.template file using your credentials from https://mailtrap.io
and put .env in the project root.

Icons https://heroicons.com

Run production docker compose:

```
docker compose -f docker-compose.yml -f docker-compose.production.yml up -d
```
