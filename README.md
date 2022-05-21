# oklookat / site
Hello, it's my site.

Idea: own webpage with blog/projects/cool stuff.

# Lore
[click here for the best experience](https://www.youtube.com/watch?v=S4-YFU8hx7k)

I guess my programming interest started with SAMP (San Andreas Multiplayer): i played it with cheats. It was be fun, but i wanted to develop my own cheat. Of course - no luck with that. I was young, and C++ broke my brain. But i was able to write (CLEO) scripts, it was be fun also. Then i do some experiments with own SAMP server, own Minecraft server, some Python stuff and other.

And then i trying to make my own cool website on Wordpress. Of course - no luck with that, i was still young and didn't even know which way to think. Then I began to write music and somehow forgot about this business (although I continued to do strange computer things like installing macOS on AMD FX, (try to)learn Arduino, Linux, VFX, 3D, animation, etc).

A few years later, in ~2019 i decide to create own blog because i like write text and programming, and plus web dev. has a lower entry threshold than CPP low-level things for example.
But I didn't even know what is HTML. Oh.

First it was Wordpress. It was inconvenient to change something and when i saw in guides ```<?php``` (JUST IN THE MIDDLE OF HTML MARKUP!!1!!) I thought it was written by some alien. I got tired of understanding these strange things and moved on.

Ok, next was MySQL + PHP + jQuery/Bootstrap. Yes, because im cool guy a decide to make CMS. Slowly, I began to understand something, but it was all rather crooked. I broke down when i start to trying make PHP routing. Sad - i not use git, it would be very interesting to see what I wrote in those days. Much has been written. Many times I started to write everything again. Then there was Laravel (a little).

Next: Django + Postgres + Django templates. Ok, I already understood something.

Next: experiments with Angular, React, Vue.

Next: Postgres + Django REST framework + Vue(afaik?).

Next: a lot of work with [Postgres + AdonisJS + Vue](https://github.com/oklookat/my-site/tree/obsolete-1). And i start using git!

There were some intermediate steps and we are here. Own little golang framework, admin panel and main site. Backend seems stable and cool, frontend - not really. Refine and refactoring can be endlessly. I don't know if the site will be ready, or how long it will last. But at least it's a good experience and case for portfolio.

In general, fun.

*There are a lot of 'I's in this text - sorry, my English vocabulary is small (although google translate helped me a little).*

# Development (Docker)
With Docker you don't need to setup nginx, create database, installing languages, etc.

All this will be done automatically in the container.


## Requirements
**All**
- [VSCode](https://code.visualstudio.com)
- [VSCode Docker extension](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker)
- [VSCode Remote Development extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack)

**Windows**
- [WSL2](https://docs.microsoft.com/en-us/windows/wsl/install)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- Optional: [read this guide](https://docs.microsoft.com/en-us/windows/wsl/tutorials/wsl-vscode)

**Linux**
- Docker and Docker Compose


## Prepare (Windows)
We need to make root cert (mkcert) on host and copy them to nginx devcontainer in WSL
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