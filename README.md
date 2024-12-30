# Advent of swole

## Tools

- [Postgresql](https://www.postgresql.org/) (Database)
- [SQLC](https://sqlc.dev/) (SQL)
- [Goose](https://github.com/pressly/goose) (database migrations)
- [CompileDaemon](https://github.com/githubnemo/CompileDaemon) (run + watch)
- [Testify](https://github.com/stretchr/testify) (Better testing)
- [TailwindCSS](https://tailwindcss.com/) (Styling)

## Commands

### Run and watch for go changes

```
CompileDaemon --command="./swole" --include="*.html"
```

### Run and watch for css changes

```
npm run css
```

## Database management

```
goose -dir ./db/migrations postgres "$DATABASE_URL" up
```

```
goose -dir ./db/migrations postgres "$TEST_DATABASE_URL" up
```

## SQLC generation

```
sqlc generate
```

## Run tests

```
go test -v ./...
```
