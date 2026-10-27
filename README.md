# WASA Photo

WASA Photo is a photo-sharing social web application developed for the Web and Software Architecture course in the Bachelor's degree in Computer Science at Sapienza University of Rome. Users can publish photos, follow other users, and interact with photos through likes and comments.

## Features

- Create or return to an account by choosing a username.
- View a feed containing photos from followed users.
- Publish and delete photos, and manage a profile and profile picture.
- Search for users, view profiles, and follow or unfollow users.
- Like and comment on photos.
- Ban and unban users.

## Technology

- **Frontend:** Vue 3, Vue Router, Vite, and Axios
- **Backend:** Go HTTP API
- **Database:** SQLite

## Project layout

| Path | Description |
| --- | --- |
| `cmd/webapi/` | Backend server entry point and configuration |
| `service/api/` | HTTP handlers and API routes |
| `service/database/` | SQLite database and data-access logic |
| `webui/` | Vue frontend |
| `doc/api.yaml` | OpenAPI description of the API |
| `demo/config.yml` | Example service configuration |

## Run locally

Install Go 1.17 or later and Node.js with npm. The SQLite Go driver may also require a C compiler, depending on your platform.

1. Start the API from the repository root:

   ```sh
   go run ./cmd/webapi
   ```

   By default, the API listens on `http://localhost:3000` and stores its SQLite database at `/tmp/decaf.db`. The server accepts configuration through command-line flags, `CFG_` environment variables, and a YAML configuration file; run `go run ./cmd/webapi --help` to see the available options.

2. In another terminal, start the frontend:

   ```sh
   cd webui
   npm ci
   npm run dev
   ```

   Open the local URL printed by Vite. The frontend is configured to send API requests to `http://localhost:3000`.

The first time you log in with a username, the API creates the account. Logging in again with that username returns the existing account.

## Build

Build the frontend for standalone hosting with:

```sh
cd webui
npm ci
npm run build-prod
```

To embed the frontend in the Go server, build its dashboard assets first and then compile the backend with the `webui` build tag:

```sh
cd webui
npm ci
npm run build-embed
cd ..
go build -tags webui -o webapi ./cmd/webapi
```

The embedded dashboard is served at `/dashboard/`. The repository also includes `Dockerfile.backend` and `Dockerfile.frontend` for container builds.

## API

The API is described in [`doc/api.yaml`](doc/api.yaml). The server exposes a liveness endpoint at `GET /liveness`.
