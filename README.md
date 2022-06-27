# Web Client

this is a http client wrapper to facilitate building rest clients.

## How to test

* all of the unit tests

```sh
go test -race ./...
```

* specific unit test

```sh
go test -timeout 2s -count 1 -run ^TestGetRequest$ github.com/fernandoocampo/webclient
```

