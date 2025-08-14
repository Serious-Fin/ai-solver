# ai-solver

Make AI solve real-life programming problems by providing prompts and guiding it through the process

Keep in mind:

- descriptions: add empty lines after each paragraph
- input/output should be on separate lines
- each example should be quoted instead of code
- then each input word can be bold if example is quoted?

Stage 1:

- add login/sign-up screens
- create login/sign-up functionality
- try and host it on the internet

## Deploying

### API

#### Preparation

Added the following to `.env` file:

```text
GIN_MODE=release
```

TODO: use `SetTrustedProxies()` to let traffic only from frontend IP?

#### Build & Run

Command to build the API docker image:

```text
docker build -f Dockerfile --tag go-api .
```

Command to run image as container:

```text
docker run -p 8080:8080 go-api
```

### Frontend

#### Build & Run

Build frontend via:

```text
docker build -f Dockerfile --tag svelte-frontend .
```

Run image:

```text
docker run -p 3000:3000 svelte-frontend
```
