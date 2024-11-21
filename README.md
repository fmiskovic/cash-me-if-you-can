![go workflow](https://github.com/fmiskovic/cash-me-if-you-can/actions/workflows/test.yml/badge.svg)
![lint workflow](https://github.com/fmiskovic/cash-me-if-you-can/actions/workflows/lint.yml/badge.svg)

# Cash Me If You Can

Technical test for the position of Backend Developer at [VALSEA TECHNOLOGY](www.valsea.com).

## How to run

```bash 
make run
```  

## Project Structure

Top Level Directories

- [api/](api) - http server, handlers and routes.
- [cmd/](cmd) - cli commands like `migrate` and `server`.
- [config/](config) - configuration and loading environment variables.
- [database/](database) - database service, repositories and migration files.
- [internal/](internal) - core logic, `services` as business use cases and `model` as domain entities.
- [pkg/](pkg) - reusable packages.
- [tests/](tests) - e2e tests.

### Generating mocks

We use [gomock](https://github.com/uber-go/mock) to generate mocks.

If you change the interface make sure to always run this command:
```bash
make mocks
```

## Testing

To run the tests locally, run `make test` to run all the unit tests
or run `go ./... -run <test-name>` to run specific unit test.

By default `make test` will run the tests in parallel n-times.
You can also do this manually by running: `go test ./... -parallel -count=5`


## Available commands

Run `make help` to see all available commands.