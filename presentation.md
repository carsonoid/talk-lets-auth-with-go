# Let's Auth With Go
From JWT to gRPC

Carson Anderson
DevX-0, Weave
@carson_ops

https://github.com/carsonoid/talk-lets-auth-with-go

## Why JWT?

.image img/gopher.png 500 _
.caption _Gopher_ by [[https://github.com/MariaLetta/free-gophers-pack][Maria Letta]]

## JWT Benefits

* Fully distributed
  * No DB required
* Issuers are not in the serving path
* Common standard
* Easy to make and use!

## JWT Basics: What is a JWT?

All JWTs (Pronounced JOT) have the same basic 3 part structure:

---

`HEADER`.`PAYLOAD`.`SIGNATURE`

* `HEADER` is a base64 encoded json object
* `PAYLOAD` can be anything but is almost always a base64 encoded JSON object
* `SIGNATURE` can change depending on the signing method used and may be omitted entirely

## Example JWT - JSON data

Header

.code examples/minimal-jwt.json /START HEADER OMIT/,/END HEADER OMIT/

Payload

.code examples/minimal-jwt.json /START PAYLOAD OMIT/,/END PAYLOAD OMIT/

## Example JWT - Base64 Data

Header

.code examples/minimal-jwt.json /START HEADERBASE64 OMIT/,/END HEADERBASE64 OMIT/

Payload

.code examples/minimal-jwt.json /START PAYLOADBASE64 OMIT/,/END PAYLOADBASE64 OMIT/

Then sign it!

## So how do we sign a token?

## Signature Pseduo-code

.code examples/minimal-jwt.json /START SIGNATURE OMIT/,/END SIGNATURE OMIT/

.code examples/minimal-jwt.json /START SIGNATURE EXAMPLE OMIT/,/END SIGNATURE EXAMPLE OMIT/

## Signature Pseduo-code

.code examples/minimal-jwt.json /START SIGNATURE OMIT/,/END SIGNATURE OMIT/

.code examples/minimal-jwt.json /START SIGNATURE EXAMPLE OMIT/,/END SIGNATURE EXAMPLE OMIT/

Any number of signature methods can be used. Common Methods are:
  * symmetric - HMAC
  * asymmetric - RSA, DSA, ED25519

## Signature Pseduo-code

.code examples/minimal-jwt.json /START SIGNATURE OMIT/,/END SIGNATURE OMIT/

.code examples/minimal-jwt.json /START SIGNATURE EXAMPLE OMIT/,/END SIGNATURE EXAMPLE OMIT/


Any number of signature methods can be used. Common Methods are:
  * symmetric - HMAC
  * asymmetric - RSA, DSA, ED25519

Be wary:

* The the result might look different for different methods
* The result is is **not** just base64 text.

## Final Result

`HeaderB64`.`PayloadB64`.`Signature`

.code examples/minimal-jwt.json /START FINAL OMIT/,/END FINAL OMIT/

Resulting tokens are always just long strings. But they have the telltale two `.` chars and 3 part format

> Remember that the header and payload are both just base64 of the contents!
> JWTs **do not encrypt** the payload or header. So don't use them to store secrets!

## So how do we validate a token?

## Validation Pseduo-code

.code examples/minimal-jwt.json /START VALIDATION OMIT/,/END VALIDATION OMIT/

> These 3 operations are often handled in one helper func. Or they may need to be done
> manually. Check your package usage for details.

## What we will be building

* Simple command line token issuer and validator
* A basic auth-api that issues JWT tokens
* A frontend HTTP service that requires a token from our auth server
* A backend gRPC service that requires a token from our auth server via the frontend

Along the way we will build a simple JWT helper library to wrap all
the common things our 3 services will do.

## First thing: Generate Keypair

Asymmetric singing requires a private and public key. 
* We will be exclusively using ed25519 for signing
* All modern versions of openssl support the ed25519 method.

Commands

.code examples/minimal-jwt.json /START KEYGEN OMIT/,/END KEYGEN OMIT/

Example keys

.code examples/auth.ed

.code examples/auth.ed.pub

## Helper Package: simplejwt

This package will do opinionated issuing and validation and will
also be the foundation for all the middleware we will write later.

It will also help us separate server code from the auth concepts.

---

JWT parsing/signing packages:

* We will use: `https://github.com/golang-jwt/jwt`
  * Simple
  * I've used it the most
* Alternative: `github.com/lestrrat-go/jwx`
  * Is more full-featured but more complex
  * NOTE: Does not do signature verification or field validation
    unless you specifically tell it to in the `Parse` function

## Build a jwt issuer in simplejwt

## simplejwt - Issuer type

.code pkg/simplejwt/issuer.go /golang-jwt/,/golang-jwt/

.code pkg/simplejwt/issuer.go /START ISSUER OMIT/,/END ISSUER OMIT/

##

.code pkg/simplejwt/issuer.go /IssueToken/,/^}/

## cmd/jwt-issue

.code cmd/jwt-issue/main.go

## cmd/jwt-issue run

.play examples/jwt-issue.go /START EXEC OMIT/,/END EXEC OMIT/

## Build a jwt validator in simplejwt

## simplejwt - Validator type

.code pkg/simplejwt/validator.go /golang-jwt/,/golang-jwt/

.code pkg/simplejwt/validator.go /START VALIDATOR OMIT/,/END VALIDATOR OMIT/

##

.code pkg/simplejwt/validator.go /GetToken/,/^}/

## cmd/jwt-validate

.code cmd/jwt-validate/main.go

## cmd/jwt-validate run

.play -edit examples/jwt-validate.go /START EXEC OMIT/,/END EXEC OMIT/

## The current architecture

These commands represent a minimal two-party identity and claims exchange.

Note no direct communication was required and no DB was required. With JWT the
only thing you need to make available to consumers is the public key.

.image img/cmd-flow.png

The only thing we will do different from here is adding HTTP/gRPC

## The new architecture

A big advantage of JWT is that the issuer is not in the critical path and there is no DB.

.image img/flow.png _ 800

The big catch is that you need to distribute the public key. We will do that by reading files
but there are other ways like the JWK standard and https fetching from public endpoints.

## Build a basic auth service

##

.code cmd/0-auth-api/main.go /\/\/ AuthService/,/^}/
.code cmd/0-auth-api/main.go /\/\/ NewAuthService/,/^}/

##

.code cmd/0-auth-api/main.go /^func .* HandleLogin/,/^}/

##

.code cmd/0-auth-api/main.go /^func main/,/^}/

## We just built this

.image img/flow-1.png _ 800

##

.play -edit examples/auth.go /START EXEC OMIT/,/END EXEC OMIT/

## Build a basic frontend HTTP service

##

.code cmd/1-frontend/main.go /^type Frontend/,/^}/
.code cmd/1-frontend/main.go /^func NewFrontend/,/^}/

##

Helper func to parse a token from the header

.code cmd/1-frontend/main.go /^func .* getHeaderToken/,/^}/

##

The claims endpoint gets the token and prints the claims

.code cmd/1-frontend/main.go /^func .* ClaimsHandler/,/^}/

##

Without middleware, all endpoints **must** get the token to do auth
even if it's not used

.code cmd/1-frontend/main.go /^func .* RootHandler/,/^}/

##

Now build the main to start our frontend


.code cmd/1-frontend/main.go  /simplejwt/


.code cmd/1-frontend/main.go  /^func main/,/GRPC OMIT/


##

Build a frontend and register handlers

.code cmd/1-frontend/main.go /frontend, err/,/^}/

> Note: We will come back to the backend client and handler after we talk about the backend

##

.play -edit examples/auth-frontend.go /START EXEC OMIT/,/END EXEC OMIT/

## We just built this

.image img/flow-2.png _ 800

## Before Backend - gRPC Foundations

## How do we transfer the token over gRPC in go?

In HTTP requests we have HTTP Headers:

`Authorization: Bearer XXXXXX`

What about in gRPC?

## Use the metadata package

`"google.golang.org/grpc/metadata"`

---

GPC "headers" == `metadata` package + `context`

The metadata package has standardized helpers for setting k/v pairs
and sending it over the wire.

* SET: `metadata.NewOutgoingContext(ctx, map[string]string)`
* GET: `metadata.FromIncomingContext(ctx)`

## Metadata set and get examples

.code cmd/1-frontend/main.go /add the auth token /,/SETTER END/

.code cmd/2-backend/main.go /rip the token /,/tokenString/

## Our GRPC Proto

.code examples/helloworld.proto

## Build a basic backend GRPC service

##

Import

.code cmd/2-backend/main.go /helloworld/,/helloworld/

Server type

.code cmd/2-backend/main.go /^type Backend/,/^}/
.code cmd/2-backend/main.go /^func NewBackend/,/^}/


##

.code cmd/2-backend/main.go /^\/\/ SayHello/,/^}/

## Build the backend main

.code cmd/2-backend/main.go /^func main/,/PG2/

## Build the backend main

.code cmd/2-backend/main.go /PG2/,/^}/

## Use the backend in the frontend

Create the client and set it in the frontend type

.code cmd/1-frontend/main.go /GRPC OMIT/,/NewFrontend/

.code cmd/1-frontend/main.go /^func NewFrontend/,/^}/

## Use the backendClient in the handler

.code cmd/1-frontend/main.go /func .* HelloHandler/,/^}/

## Use the backendClient in the handler

.code cmd/1-frontend/main.go /add the auth token/,/^}/

## We just built this

.image img/flow-3.png _ 800

## Test the gRPC auth with a cli

.code -edit cmd/grpc-local/main.go /func main/,/^}/

## Backend auth fail demo

.play -edit examples/auth-backend.go /START EXEC OMIT/,/END EXEC OMIT/

## Current Architecture

.image img/flow.png _ 800

## Full Stack Demo

.play -edit examples/auth-frontend-backend.go /START EXEC OMIT/,/END EXEC OMIT/

## Current Architecture

.image img/flow.png _ 800

## Too much toil; needs middleware

.image img/flow-mw.png _ 800

## First: A note on context.Context

`ctx` is everywhere. It's in both http and grpc implementations in go.
This makes it an ideal place to keep tokens for individual requests.

## simplejwt - Context Helpers

Context best practices encourage the use of a custom type when setting
keys. This keeps packages from overwriting each other if they use bare values. It is also advised to use helper functions rather than directly manipulate ctx in your end code

.code pkg/simplejwt/ctx.go /middlewareContextKey/,/"token"/

## The Setter Func is simple

.code pkg/simplejwt/ctx.go /ContextWithToken/,/^}/

## Context Value Getter

.code pkg/simplejwt/ctx.go /ContextGetToken/,/^}/

## Context Value "MustGet"

The reason for this function will make sense later

.code pkg/simplejwt/ctx.go /MustContextGetToken/,/^}/

## Build HTTP Server Middleware in simplejwt

##

.code pkg/simplejwt/http-middleware.go /Middleware/,/^}/

.code pkg/simplejwt/http-middleware.go /NewMiddleware/,/^}/

##

.code pkg/simplejwt/http-middleware.go /HandleHTTP/,/^}/

## Frontend: Create Middleware

.code cmd/1-frontend-mw/main.go /func main/,/END HTTP OMIT/

## Frontend: Use Middleware

.code cmd/1-frontend-mw/main.go /business/,/add handlers here/
.code cmd/1-frontend-mw/main.go /root mux/,/^}/

## Frontend: No need for the validator

We can now assume validated tokens in the request context!

.code cmd/1-frontend-mw/main.go /^type Frontend/,/^}/
.code cmd/1-frontend-mw/main.go /^func NewFrontend/,/^}/

## Frontend: Use Simple Handlers!

##
.code cmd/1-frontend-mw/main.go /^func .* ClaimsHandler/,/^}/

##

The root handler is so simple we don't even need the method

.code cmd/1-frontend-mw/main.go /NewServeMux/,/\}\)/

> We will come back to the 'Hello' handler that uses the backend

## We just built this

.image img/flow-mw-http.png _ 800

## Build GRPC Server Middleware in simplejwt

##

Anything that fits an expected signature can be a gRPC server middleware

All we need to do is extend our current middleware with the right method and then use it.

We just need to add a method on this:

.code pkg/simplejwt/http-middleware.go /Middleware/,/^}/

##

.code pkg/simplejwt/grpc-middleware-server.go /Middleware/,/^}/

## Backend: Create And Use Middleware

.code cmd/2-backend-mw/main.go /create middleware/,/^}/

## Backend: No need for the validator

We can now assume validated tokens in the incoming context!

.code cmd/2-backend-mw/main.go /^type Backend/,/^}/
.code cmd/2-backend-mw/main.go /^func NewBackend/,/^}/

## Backend: Cleanup Handler

We can just assume auth all the time and use the context helpers if
we need the token for data

.code cmd/2-backend-mw/main.go /SayHello/,/^}/

## We just built this

.image img/flow-mw-grpc-server.png _ 800

## Build GRPC Client Middleware in simplejwt

##

Anything that fits an expected signature can be a gRPC client middleware

All we need to do is extend our current middleware with the right method and then use it.

We just need to add a method on this:

.code pkg/simplejwt/http-middleware.go /Middleware/,/^}/

##

.code pkg/simplejwt/grpc-middleware-client.go /Middleware/,/^}/

## Frontend: Add middleware to backend client

.code cmd/1-frontend-mw/main.go /create middleware/,/defer/

## Frontend: HelloHandler becomes trivial

* We don't have to check the incoming http handler for a valid token
* We don't have to set the token in the outgoing client calls

.code cmd/1-frontend-mw/main.go /^func .* HelloHandler/,/^}/

## We just built this

.image img/flow-mw-grpc-client.png _ 800

## Full stack demo

.play -edit examples/auth-frontend-backend-mw.go /START EXEC OMIT/,/END EXEC OMIT/

## Final architecture

.image img/flow-mw-all.png _ 800

## Final code view

The final 3 services are incredibly simple:

.link https://github.com/carsonoid/talk-lets-auth-with-go/blob/main/cmd/0-auth-api/main.go Auth API
.link https://github.com/carsonoid/talk-lets-auth-with-go/blob/main/cmd/1-frontend-mw/main.go HTTP Frontend
.link https://github.com/carsonoid/talk-lets-auth-with-go/blob/main/cmd/2-backend-mw/main.go GRPC Backend

`simplejwt` wraps the issuer, validator, and middlewares

.link https://github.com/carsonoid/talk-lets-auth-with-go/tree/main/pkg/simplejwt simplejwt package
