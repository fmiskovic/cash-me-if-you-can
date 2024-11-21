![go workflow](https://github.com/fmiskovic/cash-me-if-you-can/actions/workflows/test.yml/badge.svg)
![lint workflow](https://github.com/fmiskovic/cash-me-if-you-can/actions/workflows/lint.yml/badge.svg)

# Cash Me If You Can

Technical test for the position of Backend Developer at [VALSEA TECHNOLOGY](www.valsea.com).

_Although the task requirement mentioned that in-memory storage would suffice, I opted to use Postgres instead. 
<br>This choice was made to avoid using mutexes in services and 
because database transactions and row locking are more appropriate for this type of task._


## How to run the service

Prerequisites:
- _Have Docker and Docker Compose installed and running_

Clone the repository and run the following commands to start the service:


```bash 
make docker-run
``` 

```bash
make migrate-up
```

```bash 
make run
```  

## API Endpoints Curl Examples

### Create New Account

```bash
curl -X POST http://localhost:8080/accounts -d '{"initial_balance":1000.12233121,"owner":"John Snow"}' -H "Content-Type: application/json"
```

### Retrieve Account Details 
Replace <account_id> with the one you got from the previous request.

```bash
curl -X GET http://localhost:8080/accounts/<account_id>
```

### List All Accounts

```bash
curl -X GET http://localhost:8080/accounts
```

### Create Transaction
Replace <account_id> with the one you got from the previous request.


```bash
curl -X POST http://localhost:8080/accounts/<account_id>/transactions -d '{"amount":999.12233121,"type":"deposit"}' -H "Content-Type: application/json"
```

### Retrieve Transactions for an Account
Replace <account_id> with the one you got from the previous request.

```bash
curl -X GET http://localhost:8080/accounts/905b1267-fe4f-4766-bdbd-4c2d9c761af0/transactions
```

### Transfer Between Accounts
Replace <from_account_id> and <to_account_id> with real account ids.

```bash
curl -X POST http://localhost:8080/transfer -d '{"amount":100.12233121,"from_account_id":"<from_account_id>","to_account_id":"<to_account_id>"}' -H "Content-Type: application/json"
```

## Testing

To run the tests locally, run `make test` to run all the unit tests
or run `go ./... -run <test-name>` to run specific unit test.

By default `make test` will run the tests in parallel n-times.
You can also do this manually by running: `go test ./... -parallel -count=5`

## Project Structure

Top Level Directories

- [api/](api) - http server, handlers and routes.
- [cmd/](cmd) - cli commands like `migrate` and `server`.
- [config/](config) - configuration and loading environment variables.
- [database/](database) - database service, repositories and migration files.
- [internal/](internal) - core logic, `services` as business use cases and `model` as domain entities.
- [pkg/](pkg) - reusable packages.
- [tests/](tests) - e2e tests.


## Available commands

Run `make help` to see all available commands.