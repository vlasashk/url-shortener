# URL Shortener
## Build
#### Prerequisites
- docker

1. Clone project:
```
git clone git@github.com:vlasashk/url-shortener.git
cd url-shortener
```
2. Run:
```
docker compose up --build
```
3. Test:
```
go test -v ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html
```
## Project information

### Restrictions/Peculiarities
- Alias имеет ограниченный срок жизни (1 месяц)
- Каждые 100 посещений обновляют срок жизни alias на 1 месяц
- Для очистки старых записей, у которых закончился срок жизни, используется cronjob сервис. В конфиге можно контролировать частоту выполнения запроса на удаление не актуальных записей раз в день/неделю/месяц

### Tools used
- PostgreSQL as database
- [jackc/pgx](https://pkg.go.dev/github.com/jackc/pgx) package as toolkit for PostgreSQL
- [go-chi/chi](https://pkg.go.dev/github.com/go-chi/chi) package as router for building HTTP service
- [rs/zerolog](https://github.com/rs/zerolog) package for logging
- [stretchr/testify](https://github.com/stretchr/testify) package for testing
- [vektra/mockery](https://github.com/vektra/mockery) package for mock generation
- Docker for deployment

### Functionality
#### URL manipulation
- {POST} /alias - Создание alias
    ```
    {
        "original": "https://test.com",
    }
    ```
- {GET} /{alias} - Получение оригинальной ссылки по alias (редирект)