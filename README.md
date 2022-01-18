Location History Server
=======================

`HISTORY_SERVER_LISTEN_ADDR`: default is `8080`

## HTTP endpoints

* `PUT /location/{order_id}`
* `GET /location/{order_id}?max=<N>`
* `DELETE /location/{order_id}`

### `PUT /location/{order_id}`

append a location to the history for the specified order.

Example interaction:
```
PUT /location/def456
{
	"lat": 12.34,
	"lng": 56.78
}

200 OK
```
### `GET /location/{order_id}?max=<N>`

Retrieve at most N items of history for the specified order. The most recent locations (in chronological order of insertion) should be returned first, if history is truncated by the `max` parameter.

Example interaction:
```
GET /location/abc123?max=2

200 OK
{
	"order_id": "abc123",
	"history": [
		{"lat": 12.34, "lng": 56.78},
		{"lat": 12.34, "lng": 56.79}
	]
}
```
The `max` query parameter may or may not be present. If it is not present, the entire history should be returned.

### `DELETE /location/{order_id}`

delete history for the specified order. Example interaction:
```
DELETE /location/xyz987

200 OK
```