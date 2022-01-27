
Location History Server
Task Description
Your task is to implement a toy in-memory location history server.

Clients should be able to speak JSON over HTTP to the server. The three endpoints it should support are:

PUT /location/{order_id}
GET /location/{order_id}?max=<N>
DELETE /location/{order_id}
Details about the endpoints:

PUT /location/{order_id} - append a location to the history for the specified order. Example interaction:

PUT /location/def456
{
	"lat": 12.34,
	"lng": 56.78
}

200 OK
GET /location/{order_id}?max=<N> - Retrieve at most N items of history for the specified order. The most recent locations (in chronological order of insertion) should be returned first, if history is truncated by the max parameter. Example interaction:

GET /location/abc123?max=2

200 OK
{
	"order_id": "abc123",
	"history": [
		{"lat": 12.34, "lng": 56.78},
		{"lat": 12.34, "lng": 56.79}
	]
}
The max query parameter may or may not be present. If it is not present, the entire history should be returned.

DELETE /location/{order_id} - delete history for the specified order. Example interaction:

DELETE /location/xyz987

200 OK
Submission guidelines and notes
Create a Go module which implements the server described above. Running go build in the root directory of this module should produce a working binary. The server should serve HTTP at the address specified in the environment variable HISTORY_SERVER_LISTEN_ADDR. If this environment variable is not set, the listen address should default to :8080.

If the task statement is unclear about something, please feel free to reach out. That being said, decisions about how to treat certain edge cases that may come up are at your discretion.

Good luck, and enjoy! :)
