# dtt-game-backend

This is the backend service repository of MatchMe app. DTT stands for [Discrete Trial Training](https://www.autismspeaks.org/expert-opinion/what-discrete-trial-training).

To get started with this project, you'll need the following prerequisites:

- [MongoDB](https://www.mongodb.com/docs/manual/installation/) (tested with v3.6 and v7.0)
- A [Replicate](https://replicate.com/) API key
- A writable directory to store user-uploaded and AI-generated photos

Moreover, you'll need both backend and [frontend](https://github.com/chinese-slacking-party/dtt-game-frontend) to make the app work. After you've deployed the two ends, you'll use [`nginx`](http://nginx.org/en/docs/beginners_guide.html) to map the services onto the same port (80 or 443). See appendix for sample `nginx` configuration snippet.

You can choose to **run a binary release** or **build from source code**.

## Running a binary release

See the [Releases](https://github.com/chinese-slacking-party/dtt-game-backend/releases) page. We currently provide binaries for Windows (x64) and Linux (x64). Those executables are completely standalone (no external library required). You can run it directly in your shell with:

```bash
# Linux
export REPLICATE_API_KEY=xxxx
./linux_server
```

```powershell
# Windows PowerShell
$env:REPLICATE_API_KEY=xxxx
.\windows_server.exe
```

## Building from source code

You'll need [`go`](https://go.dev/) and [`git`](https://git-scm.com/) to build the project.

After cloning this repository, `cd` into it and run:

```bash
go build ./cmd/server
```

As long as you have a working Internet connection (or `GOPROXY` properly set), `go` will automatically download all the dependencies and generate the executable.

## Appendix 1: Sample `nginx` configuration

```nginx
# Add this in your `server` block:
    location /api { # Backend
        proxy_pass http://localhost:8080;
    }
    location / { # Frontend
        proxy_pass http://localhost:3000;
    }
```

## Appendix 2: Customizing your build

As of v0.1.0, MatchMe backend uses in-code [constants](/config/constants.go) for configuration. Before you build, make sure to change those constants into desired values:

- `DBAddr` and `DBName` for your MongoDB configuration
- `PhotoDir` for your writable UGC directory
- `OurAddr` for your IP or domain

If you have your own deployment on Replicate, you should change the `CreatePredictionWithDeployment` invocation in [this file](/handlers/album/album.go) to use yours as well.

More flexible settings are planned for the next release.
