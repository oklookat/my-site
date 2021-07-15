# oklookat.ru (src) #

## Requirements a.k.a tested with
- Node.js 16.5.0
- PostgreSQL 13.3


## How to run
1. In /backend make copy of .env.example and rename it to .env
2. Edit .env as you like (almost)
3. Run ```npm install``` in folders
4. Create database with name (see PG_DB_NAME in backend .env)
5. Run migration and create superuser (see [Main commands](#main-commands))
6. Run in dev (see [Main commands](#main-commands))


## Open ports:
- Database: localhost:5432
- Backend: localhost:3333
- Frontend: localhost:8080 (admin panel)


## <a name="main-commands"></a>Main commands:
### Backend
- ```node ace migration:run``` = run migration
- ```node ace migration:rollback``` = reset tables
- ```node ace elven:superuser``` = create superuser
- ```node ace serve --watch``` = run in dev
- ```node ace build --production``` = build prod
### Elvenfront
- ```npm run serve``` = run in dev