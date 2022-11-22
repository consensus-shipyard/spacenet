# Faucet

## Development
To run the service even in the development mode, you must provide an X509 certificate.

The easiest way to do that is to use [mkcert](https://github.com/FiloSottile/mkcert)
tool and `make cert` command.

The run `make all` to ensure that tests pass and `make demo` to run a demo accessible on localhost.