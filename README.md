# OAuth Provider in GO

## Postgres database in docker-compose

A postgres database can be spun-up if one is not already running on your system.

`docker-compose up -d`

## Consider using `gow`

The `gow` package provides file watching so that changes can be reloaded as soon as files are changed.

https://github.com/mitranim/gow

`go install github.com/mitranim/gow@latest`

`gow run .`
