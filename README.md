# AI Solver

Make AI solve real-life programming problems by providing prompts and guiding it through the process

## Launching

### Launch locally (without Docker containers)

This is used for local development only. Convenient because of hot reloading feature.

To run the project locally UNIX environment is needed. If you have a Mac or Linux PC, then it will work. Windows users should download [Windows Subsystem for Linux](https://learn.microsoft.com/en-us/windows/wsl/install).

Clone the project locally:

```text
https://github.com/Serious-Fin/ai-solver.git
```

Add the following files to the project (TODO: provide their schemas):

- `.env` in `/ai-solver/frontend/`
- `.env` in `/ai-solver/api/`
- `database.db` in `ai-solver/api/data/`

Launch a shell script ([tmux](https://github.com/tmux/tmux/wiki/Installing) terminal utility is needed)

```text
./workspace.sh
```

Website is now accessible via the browser at [http://localhost:5173](http://localhost:5173)

> **Note:** Logging in does not work locally, as OAuth clients are configured to work on original host name.

### Launch locally (with Docker containers)

This is used to test out locally how the program will behave in the cloud. Convenient to check if web to api communication works as intended.

Download [Docker Desktop](https://docs.docker.com/desktop/)

Clone the project locally:

```text
https://github.com/Serious-Fin/ai-solver.git
```

Add the following files to the project (TODO: provide their schemas):

- `.env` in `/ai-solver/frontend/`
- `.env` in `/ai-solver/api/`
- `database.db` in `ai-solver/api/data/`

Run the development docker compose file:

```text
docker compose -f compose.dev.yaml up --build
```

Website is now accessible via the browser at [http://localhost](http://localhost)

---

Keep in mind:

- descriptions: add empty lines after each paragraph
- input/output should be on separate lines
- each example should be quoted instead of code
- then each input word can be bold if example is quoted?

## Deploying

### API

#### Preparation

Added the following to `.env` file:

```text
GIN_MODE=release
```

TODO: use `SetTrustedProxies()` to let traffic only from frontend IP?
