**oklookat / API**

powered by [servus](./core/)
# Requirements

- [go 1.18+](https://go.dev/dl)

- [GCC](https://gcc.gnu.org) or [MinGW](https://www.mingw-w64.org/downloads)


# Commands & flags
see: [help.go](./apps/elven/cmd/help.go)


# Routes (/elven)
*information here is rarely updated

## authorization
### /auth - manage auth tokens
- **POST "/login"** = generate and send auth token by username and password
- **POST "/logout"** = delete auth token by auth token

## users
### /users/me - current user
- **GET "/"** = get current user info by auth token
- **POST "/change"** = change username or password

## articles
### /article/articles - manage articles
- **GET "/"** = get paginated
- **POST "/"** = create
- **GET "/id"** = get one
- **PUT "/id"** = full update
- **PATCH "/id"** = update specific fields (json-patch)
- **DELETE "/id"** = delete

## files
### /files - manage files
- **GET "/"** = get paginated
- **POST "/"** = upload
- **DELETE "/id"** = delete
