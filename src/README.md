# Development with Docker #
### With Docker you dont care about nodejs, postgres, deps, and other borning configuration. All you need is just type some commands for build containers and run them. 


#### How to run
1. Install Docker and Node.js
2. In /backend and /frontend run ```npm install```
3. In /backend make copy of .env.example and rename it to .env
4. Edit .env
5. In /src run ```docker-compose build```
6. In /src run ```docker-compose up```
7. Run migration
8. Create superuser
9. Instruction may be changed and etc


#### Commands:
* use sudo if your not added account to Docker group. In Windows commands little different.
- ```sudo docker-compose up --build -d``` = build all stuff (run in /src dir)
- ```sudo docker exec -it backend <command>``` = run command in backend
- ```sudo docker exec -it elvenfront <command>``` = run command in frontend (admin panel)


## Useful commands (via docker exec)
- ```node ace migration:run``` = run migration (backend)
- ```node ace migration:rollback``` = reset tables (backend)
- ```node ace elven:superuser``` = create superuser (backend)


#### Open ports:
1. Nginx: localhost:1111 (use for access to backend)
2. Adminer: localhost:2222
3. Backend: localhost:3333 (now is closed, because nginx)
4. Frontend (admin panel): localhost:4444


#### Other information:
1. For local development we need node_modules installed on host. The entrypoint.sh file copy node_modules from container to our host. But if you need edit this .sh file in Windows and then run all stuff, you get error like ".sh not found". For fix this error you need convert your editied .sh to Linux file (because Windows file has bad line endings, or idk).
2. Run "service docker start" before run docker (do it once on WSL started. Or you can add docker to systemctl or something)
3. If permission denied in entrypoint.sh: cd to this .sh on host machine, and run chmod +x entrypoint.sh