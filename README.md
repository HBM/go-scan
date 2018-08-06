
# HBM scan web app

[HBM scan](https://github.com/HBM/scan-spec) implemented in Go. Results are served via JSON to a frontend written in React. All you need is a web browser.

## Goal

1. A single binary for multiple platforms (Windows, Linux, Mac)
1. No dependencies on target system (no Java, no Docker, no Go, no nothing)
1. No installation required. Simply use the binary for your platform.

## Development

The app is written in Go.

One goroutine starts a simple HTTP web server and serves all files from the `/public` folder. When `mode` is set to production all files from the `/public` folder are embedded into the Go executable using [statik](https://github.com/rakyll/statik). The web server also provides a `/json` endpoint for the frontend.

A second goroutine starts an infinte loop and listens for incoming UDP messages. All raw JSON messages are converted into proper types/struct and stored in an in-memory database. If devices stop announcing themselves they are removed from the database.

Start the app in development mode by using the following two commands. First of all transpile and bundle the JavaScript code.

```sh
npm run watch
```

Then run the backend in `development` mode using the transpiled `/public/main.js` bundle.

```sh
make dev
```

## Release

To build a new release and binaries for each platform simply run the following command.

```sh
make release
```
