# reciva-web-remote

[![GitHub](https://img.shields.io/github/license/itsnotgoodname/reciva-web-remote)](./LICENSE)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/itsnotgoodname/reciva-web-remote)](https://github.com/ItsNotGoodName/reciva-web-remote/tags)
[![GitHub last commit](https://img.shields.io/github/last-commit/itsnotgoodname/reciva-web-remote)](https://github.com/ItsNotGoodName/reciva-web-remote)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/itsnotgoodname/reciva-web-remote)](./go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/ItsNotGoodName/reciva-web-remote)](https://goreportcard.com/report/github.com/ItsNotGoodName/reciva-web-remote)

Control your legacy Reciva based internet radios (Crane, Grace Digital, Tangent, etc.) via web browser or REST API.

![Desktop Demo](/assets/desktop-demo.png)

# Features

- Toggle radio power
- Change radio volume
- Change radio audio source (unreliable)
- Play radio presets (make sure audio source is `Internet radio` or else it will hang)
- [Host playlists](#host-playlists) for radios

# Running

Download and extract the zip/tar.gz file from [releases](https://github.com/ItsNotGoodName/reciva-web-remote/releases).

## Windows

Open the extracted folder and double-click `start.bat` to run the server on port `80`.

You may need to press enter in the terminal and also press allow on the firewall prompt.

You may also want use [winsw](https://github.com/winsw/winsw) to run the program as a Windows service.

# Flags

The web server listens on port `8080` by default. It can be changed by setting the `-port` flag.

```
./reciva-web-remote -port 9000
```

The UPnP control point listens on port `8058` by default. It can be changed by setting the `-cport` flag.

```
./reciva-web-remote -cport 9058
```

The program looks for `reciva-web-remote.json` file in the current folder by default. It can be changed by setting the `-config` flag.

```
./reciva-web-remote -config example.json
```

# Host Playlists

Credit to this [article](https://swling.com/blog/2021/03/how-to-give-your-reciva-wifi-radio-a-second-life-before-the-service-closes-on-april-30-2021/).

Playlists can be hosted on the web server.
This is only useful if you were able to point your presets to a domain or an IP that you own before Reciva shutdown it's service.

Take for example you have a radio where preset 1 and 2 point to `http://192.168.1.2:9000/01.m3u` and `http://192.168.1.2:9000/02.m3u` respectively.

You will have `reciva-web-remote.json` file with the following content.

```json
{
  "presets": [
    {
      "url": "http://192.168.1.2:9000/01.m3u"
    },
    {
      "url": "http://192.168.1.2:9000/02.m3u"
    }
  ]
}
```

Then you will run the program on a machine that has the IP address `192.168.1.2` and with port `9000` available.

```
./reciva-web-remote -port 9000
```

The program will host the playlists on `/01.m3u` and `/02.m3u`.

The contents of the playlists can be changed in the web interface.

![Edit Demo](/assets/desktop-edit-demo.png)

# API

REST API is documented using [Swagger](https://swagger.io/).
Swagger file is stored in [swagger.json](./docs/swagger/swagger.json) and can be viewed using
[Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/ItsNotGoodName/reciva-web-remote/master/docs/swagger/swagger.json).

# Build

Install npm packages.

```
make npm
```

Build program.

```
make build
```

Binary will be in `./bin`.

# Development

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

# To Do

- Readd toasts in web interface
- Add volume slider in web interface
- Discover radios on a timer
- Document WebSocket API
