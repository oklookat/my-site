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
- All actions with this repo will be in WSL

**Linux**
- Docker and Docker Compose


## How
0. Clone this repo
1. Open VSCode
2. (Windows) Attach VSCode to WSL
4. F1 -> type "Open Folder in Container"
5. Any dir with .devcontainer folder can be opened

After choose, Docker start building containers and then you can work.

In the same way, you can open other directories (but without option 0)


## Useful commands
close WSL:
```  
wsl --shutdown
```