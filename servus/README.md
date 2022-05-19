**oklookat.ru / REST API backend**
**powered by servus little framework**
# Requirements

- **go 1.18+**

- **GCC**

# Useful

**Install all deps:**
```go get ./...```


# Commands & flags

## run args
**path to config:**
```-config=/dir/dir2/config.json```

**create database tables by .sql:**
```-sql=/dir/to/file.sql el:mg```

**delete databaes tables:**
```el:rb```

**create user (-die = delete user if exists):**
```el:tu -username=NAME -password=PASS -die```

**create superuser:**
```same as create user, but el:su instead el:tu```


# Routes (/elven)

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
- **PUT "/id"** = update full
- **PATCH "/id"** = update specific fields
- **DELETE "/id"** = delete

## files
### /files - manage files
- **GET "/"** = get paginated
- **POST "/"** = upload
- **DELETE "/id"** = delete
