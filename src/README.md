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
1. For local development we need node_modules installed on host. The entrypoint.sh file copy node_modules from container to our host. But if you need edit this .sh file in Windows and then run all stuff, you get error like ".sh not found". For fix this error you need convert your editied .sh to Linux file (because Windows file has bad line endings, or idk).
2. Run "service docker start" before run docker (do it once on WSL started. Or you can add docker to systemctl or something)
3. If permission denied in entrypoint.sh: cd to this .sh on host machine, and run chmod +x entrypoint.sh