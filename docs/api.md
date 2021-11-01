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

The uuid GET parameter is optional. Only differential state updates are sent.

```
GET ws://localhost:8080/v1/radio/ws?uuid=:uuid
```

### Modify radio state

Not all JSON parameters have to be sent. Only what you want changed on the radio.

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

Retrieves volume from radio and update radio state.

```
POST http://localhost:8080/v1/radio/:uuid/volume
```

### Renew UPnP subscription to radio event publisher

```
POST http://localhost:8080/v1/radio/:uuid/renew
```
