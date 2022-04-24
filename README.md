# reciva-web-remote

Control your legacy Reciva based internet radios (Crane, Grace Digital, Tangent, etc.) via REST api or web browser.

![Desktop Demo](/assets/desktop-demo.png)

## Features

- Discover radios on your local network via UPnP
- Toggle radio's power
- Change radio's volume
- Play presets that are on the radio
- Host playlists for radios

## Running

Download and extract the zip file from [releases](https://github.com/ItsNotGoodName/reciva-web-remote/releases).

### Windows

Open the extracted folder and double-click on the `start` file to run on port 80.

You may need to press `Enter` after the terminal opens and also press allow when the firewall prompt opens.

## Configuration

This program's web server listens on port `8080` by default. It can be changed by setting the `-port` flag.

```
./reciva-web-remote -port 9000
```

It needs port `8058` for UPnP notify requests. It can be changed by setting the `-cport` flag.

```
./reciva-web-remote -cport 9058
```

It looks for `reciva-web-remote.json` in the current folder by default. It can be changed with the `config` flag.

```
./reciva-web-remote -config test.json
```

### Host Playlists

Playlists can be hosted on the program's web server, ex. `http://example.com/01.m3u`, `http://example.com/02.m3u`. This is only useful if you were able to point your presets to a domain that you own or an IP before Reciva shutdown its services.

Open `reciva-web-remote.json` and add the following content before running the program.

```json
{
  "presets": [
    {
      "url": "http://example.com/01.m3u"
    },
    {
      "url": "http://example.com/02.m3u"
    }
  ]
}
```

The program will host the m3u playlists on `/01.m3u` and `/02.m3u`.

The contents of the playlists can be changed in the web interface.

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

## Todo

- Dark mode
- Volume slider
- Refresh volume on an interval
