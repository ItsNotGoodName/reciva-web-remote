# reciva-web-remote
Control your legacy Reciva based internet radios (Crane, Grace Digital, Tangent, etc.) via REST api or web browser.

## Usage

This program's web server listens on port `8080` by default. It can be changed by setting the `-port` flag.
```
./reciva-web-remote -port=9000
```
It also needs port `8058` for UPnP notify requests. It can be changed by setting the `-cport` flag.
```
./reciva-web-remote -cport=9058
```

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