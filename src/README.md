
# Development with Docker #
##### Tested in WSL2 (Ubuntu)


#### How to run
1. In /backend make copy of .env.example and rename it to .env
2. Edit .env as you like (almost)
3. Install Docker and Node.js
4. In /src run "./run.sh"
5. Run migration and create superuser (see [Container Commands](#container-commands))


#### Open ports:
##### Services
- Nginx: localhost:1111 (use for access to backend)
- Database: db (only between containers, not available outside)
- Adminer: localhost:2222
##### Main stuff
- Backend: localhost:3333 (now is closed, because nginx)
- Frontend (admin panel): localhost:4444


## <a name="host-commands"></a>Host Commands:
* do not forget add your user to docker group
- ```docker-compose up --build -d``` = build all stuff (run in /src dir)
- ```docker exec -it backend <command>``` = run command in backend
- ```docker exec -it elvenfront <command>``` = run command in frontend (admin panel)


## <a name="container-commands"></a>Container commands (via docker exec):
- ```node ace migration:run``` = run migration (backend)
- ```node ace elven:superuser``` = create superuser (backend)
- ```node ace migration:rollback``` = reset tables (backend)


##### Other useful commands
- ```docker kill $(docker ps -q)``` = kill all containers


#### Extra:
- In WSL2 before start docker run "service docker start"
- If you want to copy something from Windows to WSL2, do it only from WSL2 console (/mnt/c). If you directly copy files via Explorer to \\wsl$\Ubuntu it break permissions on files, and you need use chmod, what annoying.