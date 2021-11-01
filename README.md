# reciva-web-remote

Control your legacy Reciva based internet radios (Crane, Grace Digital, Tangent, etc.) via REST api or web browser.

## Running

Download and extract the zip file from [releases](https://github.com/ItsNotGoodName/reciva-web-remote/releases). Open the extracted folder and run the executable. On Windows you may need to press `Enter` after running the program.

## Settings

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
