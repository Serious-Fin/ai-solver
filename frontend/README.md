# AI Solver

## Launching

### Launch locally (without Docker containers)

To run the project locally UNIX environment is needed. If you have a Mac or Linux PC, then it will work. Windows users should download [https://learn.microsoft.com/en-us/windows/wsl/install](Windows Subsystem for Linux).

Clone the project locally:

```text
https://github.com/Serious-Fin/ai-solver.git
```

Add the following files to the project (TODO: provide their schemas):

- `.env` in `/ai-solver/frontend/`
- `.env` in `/ai-solver/api/`
- `database.db` in `ai-solver/api/data/`

Launch a shell script ([https://github.com/tmux/tmux/wiki/Installing](tmux) terminal utility is needed)

```text
./workspace.sh
```

Frontend is now accessible via the browser at [http://localhost:5173](http://localhost:5173)


## Cleanup of server

SSH into server:

```text
ssh -i ~/.ssh/droplet-ai-solver root@138.68.76.119
```

Cleanup

```text
# Stop all running containers
docker stop $(docker ps -q)

# Remove all containers
docker rm $(docker ps -aq)

# Remove all images
docker rmi $(docker images -q)
```

## Launching new

Go to `/ai-solver/frontend`:

```text
docker build --platform linux/amd64 -t seriousfin/ai-solver-frontend:latest .
```

Go to `/ai-solver/api`:

```text
docker build --platform linux/amd64 -t seriousfin/ai-solver-api:latest .
```

Go to a directory where images may be saved temporarily and run:

```text
docker save -o api.tar seriousfin/ai-solver-api:latest
docker save -o web.tar seriousfin/ai-solver-frontend:latest
```

Transfer images to server:
```text
scp -i ~/.ssh/droplet-ai-solver web.tar root@138.68.76.119:/root
```
