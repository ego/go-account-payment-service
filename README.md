# gopayment WEB service API for account and payment features
## (Account payment service in Go)

## Service feature

* Send payment from one account to another
* Send payments only with same currency
* View one payment by ID
* View one account by ID


## System design

* [Go](https://www.python.org/ "Go")
* [GoKit](https://gokit.io/ "gokit")
* [PostgreSQL 12](https://www.postgresql.org/docs/12/index.html "PostgreSQL 12")
* [Nginx 1.17.6](https://nginx.org "Nginx 1.17.6")


## System requirements


``` bash
make --version
>= GNU Make 3.81
```

```bash
docker --version
>= Docker version 19.03.5, build 633a0ea
```

```
docker-compose --version
>= docker-compose version 1.24.1, build 4667896b
```


## REST API

Nginx reverse proxy on `80`
http://localhost/

Dev app runs on `8888`
http://localhost:8888/


## Development

Make and Makefile for commands
```bash
make help
```

Docker compose files:

* `docker-compose.yml` **production** build
* `docker-compose.dev.yml` **development** build
  Override `ENV` vars and allow databases log queries


`ENV` vars:

* `app.env`
* `postgres.env`


**Run dev application**

```bash
make
```

**Run prod application**

```bash
make up-prod
make ps-prod
make down-prod
```

`make up-prod` run's reverse proxy so you can scale your application.
Variable `SCALE` in `Makefile`.


Coding

```bash
make refactor
```

What did not have time to do:

* DRY
* Move logic to logic layer
* Refactor code, more `right` Go, language specific paradigms
* Refactor code according to community standards and guadilances
* Test =)
* Hot reload for dev
* ...


```bash
curl -XPOST -d '{"name": "Test1", "balance": 100, "currency": "USD"}' -H "Content-Type: application/json" http://localhost:8888/accounts

curl -XGET localhost:8888/accounts/1

curl -XPOST -d '{"name": "Test2", "balance": 100, "currency": "USD"}' -H "Content-Type: application/json" http://localhost:8888/accounts

curl -XGET localhost:8888/accounts/2

curl -XPOST -d '{"name": "Test3", "balance": 100, "currency": "UAH"}' -H "Content-Type: application/json" http://localhost:8888/accounts

curl -XPOST -d '{"account_id":1, "to_account_id": 2, "amount": 100, "direction": "outgoing"}' -H "Content-Type: application/json" http://localhost:8888/payments

curl -XPOST -d '{"account_id":1, "to_account_id": 2, "amount": 100, "direction": "incoming"}' -H "Content-Type: application/json" http://localhost:8888/payments

curl -XPOST -d '{"account_id":1, "to_account_id": 3, "amount": 100, "direction": "outgoing"}' -H "Content-Type: application/json" http://localhost:8888/payments

curl -XGET localhost:8888/payments/1
```
