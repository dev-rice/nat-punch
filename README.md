### Simple UDP NAT Punchthrough
Implementation taken from https://en.wikipedia.org/wiki/UDP_hole_punching#Flow

### Running the example
On some server with known, public IP:
```
go run cmd/introducer/main.go --listen-address :9001
```

On client A (can be anywhere)
```
go run cmd/client/main.go --introducer-address <introducer-ip>:9001
```

On client B (can be anywhere)
```
go run cmd/client/main.go --introducer-address <introducer-ip>:9001
```

After both clients connect to the introducer server, they will be able to send packets to each other without the need for any port forwarding rules. In this example, they just send their "secret number" to each other.

### Example output:
#### Introducer
```
listening on :9001
read 16 bytes from 35.238.218.17:43155: hello introducer
read 16 bytes from 67.123.133.53:62917: hello introducer
sending peer info to 35.238.218.17:43155
sending peer info to 67.123.133.53:62917
```
#### Client A
```
wrote 16 bytes: hello introducer
read 19 bytes from 35.225.210.119:9001: 67.123.133.53:64097
introducer introduced me to peer '67.123.133.53:64097'
my secret number: 1844085407
sending secret number to 67.123.133.53:64097
read 31 bytes from 67.123.133.53:64097: hello. secret number: 554376187
sending secret number to 67.123.133.53:64097
read 31 bytes from 67.123.133.53:64097: hello. secret number: 554376187
sending secret number to 67.123.133.53:64097
read 31 bytes from 67.123.133.53:64097: hello. secret number: 554376187
sending secret number to 67.123.133.53:64097
read 31 bytes from 67.123.133.53:64097: hello. secret number: 554376187
```

#### Client B
```
wrote 16 bytes: hello introducer
read 19 bytes from 35.225.210.119:9001: 35.238.218.17:44576
introducer introduced me to peer '35.238.218.17:44576'
my secret number: 554376187
sending secret number to 35.238.218.17:44576
read 32 bytes from 35.238.218.17:44576: hello. secret number: 1844085407
sending secret number to 35.238.218.17:44576
read 32 bytes from 35.238.218.17:44576: hello. secret number: 1844085407
sending secret number to 35.238.218.17:44576
read 32 bytes from 35.238.218.17:44576: hello. secret number: 1844085407
sending secret number to 35.238.218.17:44576
read 32 bytes from 35.238.218.17:44576: hello. secret number: 1844085407
```
