# reciva-web-remote

Control your legacy Reciva based internet radios (Crane, Grace Digital, Tangent, etc.) via REST api or web browser.

![Desktop Demo](/assets/desktop-demo.png)

## Features
- Discover radios on the local network via UPnP
- Change power, volume, and presets on the radios via UPnP
- Host playlists for presets
## Running

Download and extract the zip file from [releases](https://github.com/ItsNotGoodName/reciva-web-remote/releases). Open the extracted folder and run the executable. On Windows, you may need to press `Enter` after the terminal opens.

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

Playlists can be hosted on it's web server, ex. `http://example.com/01.m3u`, `http://example.com/02.m3u`. This is only useful if you were able to point your presets to a domain that you own or an IP before Reciva shutdown it's services.

Create a `reciva-web-remote.json` file with the following content and then run the program.

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

Now the program will host playlists on `/01.m3u` and `/02.m3u`. The contents of the playlists can be changed in the web interface.

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
