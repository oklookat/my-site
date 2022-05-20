# oklookat-site #

# Development (Docker)
*With Docker you don't need to setup nginx, create database, installing languages and more.*

*All this will be done automatically in the container.*

## Requirements


**All**
- [VSCode](https://code.visualstudio.com)
- [VSCode Docker extension](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker)
- [VSCode Remote Development extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack)


**Windows**
- [WSL2](https://docs.microsoft.com/en-us/windows/wsl/install)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- Optional: [read this guide](https://docs.microsoft.com/en-us/windows/wsl/tutorials/wsl-vscode)


**Linux**
- Docker and Docker Compose


## Prepare (Windows)
**We need to make root cert (mkcert) on host and copy them to nginx devcontainer in WSL**
- Clone repo (host)
- *How pt.1*
- Clone repo (WSL)
- Copy ```./nginx/.devcontainer/certs``` from host to WSL
- Start from *How pt.2*


## How
0. Clone repo
1. Run bootstrap and execute: 
```setup dev certs``` | ```generate & copy``` | ```add hosts```
2. Open VSCode
3. (Windows) Attach VSCode to WSL
4. F1 -> type "Open Folder in Container"

- Any dir thats contains .devcontainer folder can be opened

- After choose, Docker start building containers and then you can work

- In the same way, you can open other directories (but start from *pt.2*)

- If you change something in Dockerfile or docker-compose, run ```Rebuild and Reopen in Container```


## Tips and tricks
After you finished, to save your RAM on Windows, close VSCode; Docker Desktop and run
```  
wsl --shutdown
```