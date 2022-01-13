# API

## Example Responses

```json
{
  "ok": true,
  "code": 200
}
```

```json
{
  "ok": true,
  "code": 200,
  "result": {
    "uuid": "fb38cb47-a74b-42d5-bb9c-89dcc6a3d960"
  }
}
```

```json
{
  "ok": false,
  "code": 404,
  "error": "radio not found"
}
```

## Radio API

### Discover radios

```
POST http://localhost:8080/v1/radios
```

### Get all radio states

```
GET http://localhost:8080/v1/radios
```

### Get radio state

```
GET http://localhost:8080/v1/radio/:uuid
```

### Get radio state via websocket

Client must send `uuid` after they connect or specify it in the `GET` parameter in order to receive state. The first message sent to client is always the full state. After that, only state changes are sent to the client. The `uuid` is always sent to the client. The connection will terminate if the client sends an invalid `uuid`.

```
GET ws://localhost:8080/v1/radio/ws?uuid=:uuid
```

### Modify radio state

Send only what you want changed on the radio.

```
PATCH http://localhost:8080/v1/radio/:uuid
content-type: application/json

{
	"power": false
	"preset": 10
	"volume": 30
}
```

### Refresh radio volume

Gets volume from radio and updates radio state.

```
POST http://localhost:8080/v1/radio/:uuid/volume
```

### Refresh UPnP subscription to radio's UPnP event publisher

```
POST http://localhost:8080/v1/radio/:uuid
```

## Preset API

### Get presets

```
GET http://localhost:8080/v1/presets
```

### Get preset by url

```
GET http://localhost:8080/v1/preset?url=:url
```

### Update preset

```
POST http://localhost:8080/v1/preset
content-type: application/json

{
	"url: "http://example.com/01.m3u"
	"newName": "Good Music"
	"newUrl" : "http://differentm3u.example.com/08.m3u"
}
```
