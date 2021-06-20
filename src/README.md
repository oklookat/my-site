# Development #

#### How to run
1. Install Node.js
2. In /backend and /frontend run ```npm install```
3. In /backend make copy of .env.example and rename it to .env
4. Edit .env as your need
5. Run migration
6. Create superuser


## Useful commands
- ```node ace migration:run``` = run migration (backend)
- ```node ace migration:rollback``` = reset tables (backend)
- ```node ace elven:superuser``` = create superuser (backend)


#### Open ports:
1. Nginx: ?
2. pgAdmin & DB: localhost:5432
3. Backend API: localhost:3333
4. Frontend (admin panel): localhost:8080


#### Other information:
1. ...
