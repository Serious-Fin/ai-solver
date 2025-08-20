# AI Solver

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
