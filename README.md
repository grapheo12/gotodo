# GOTODO

Simple REST API for a todo application, written in Go.

Done for learning Backend Development in Go.

### Stack

```
net/http
gorilla/mux
jinzhu/gorm
dgrijalva/jwt-go
```

### Reusable Features

Directory structure has a close resemblance with Django directory structure.
So this can be used as a template Webapp.

Also the auth module provides JWT authentication (Uses `Bearer` header).

A handy authentication middleware called `LoginRequired` is built to protect the endpoints.