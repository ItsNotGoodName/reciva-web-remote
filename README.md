# reciva-web-remote
Control your legacy Reciva based internet radios (Crane, Grace Digital, Tangent, etc.) via REST api or web browser.

## Usage

This program listens on port `8080` by default. It can be changed by setting the `PORT` environment variable. 
```
export PORT=9000
```
It also needs port `8058` for listening for upnp notify requests.

Run the program by executing `./reciva-web-remote` or `.\reciva-web-remote.exe` depending on your platform.

## Build

Install npm packages.
```
make npm
```

Build program.
```
make build
```

## Development

Install npm packages.
```
make npm
```

Run frontend.
```
make dev-frontend
```

Run backend.
```
make dev-backend
```