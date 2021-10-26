# this is obsolete version, dont use it. But if you want:
## 1. frontend (elvenfront) uses url's from golang api routes, need correct url's to nodejs routes.

# oklookat.ru source (elven) #

## Powered by
### Backend
AdonisJS v5 - https://adonisjs.com
### Frontend
Vue.js v3 - https://vuejs.org/


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
### Elven backend (backend)
- ```node ace migration:run``` = run migration
- ```node ace migration:rollback``` = reset tables
- ```node ace elven:superuser``` = create superuser
- ```node ace serve --watch``` = run in dev
- ```node ace build --production``` = build prod
### Elven frontend (elvenfront)
- ```npm run dev``` = run in dev


## F.A.Q
### Elven?
Elven is my CMS. It's not WP or other super cool things, but anyway it's my and it's cool. That's all I can say. 
