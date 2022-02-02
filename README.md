# Let's Auth With Go - From JWT to JSON

Carson Anderson

DevX Engineer, Weave

@carson_ops


## Quickstart

```sh
# Gen a new keypair
openssl genpkey -out auth.ed
openssl pkey -in auth.ed -pubout > auth.ed.pub

## try it with a local issuer
t=$(go run ./cmd/jwt-issue auth.ed)
echo "TOKEN: $t"
go run ./cmd/jwt-validate/ auth.ed.pub $t

# as a one-liner
go run ./cmd/jwt-validate/ auth.ed.pub $(go run ./cmd/jwt-issue auth.ed)

## try it with services
# run the basic auth api with the private key
go run ./cmd/0-auth-api auth.ed

# run a frontend with the public key
go run ./cmd/1-frontend auth.ed.pub
# or try the version with middleware
go run ./cmd/1-frontend-mw auth.ed.pub

# run a backend with the public key
go run ./cmd/2-backend auth.ed.pub
# or try the version with middleware
go run ./cmd/2-backend-with-middleware auth.ed.pub

# do a test request to just get a token and hit the frontend with it
t=$(curl admin:pass@localhost:8081/login); echo $t;curl -H "Authorization: Bearer $t" localhost:8082/
t=$(curl admin:pass@localhost:8081/login); echo $t;curl -H "Authorization: Bearer $t" localhost:8082/claims

# do a test request to just get a token and hit the frontend which calls the backend, passing the token on
t=$(curl admin:pass@localhost:8081/login); echo $t;curl -H "Authorization: Bearer $t" localhost:8082/hello;echo
```

## A note about encryption

To illustrate security best practices; the code here uses [Ed25519](https://ed25519.cr.yp.to/) keys.

These are supported by Go but may not work as easily for other languages.
However, nearly all the code here is the same regardless of JWT singing method
and nothing shown here can't be done with things like RSA or HMAC signing instead.

## Running the presentation

This presentation uses a custom theme and can be run by installing the [`go-present`](https://pkg.go.dev/golang.org/x/tools/present) tool and starting it:

It also does some setup work to fake out the go env to enable "commands" to exec in the presentation against the local machine.

Run it with the run script

```bash
./run
```
