# Description

The following repo is a small investigation into what is sent over the wire when a TLS connection is established between a client and a server by inititating an HTTP/S request and then revert to use of the raw TCP stream.

The server can operate with either TLS on/off. Make the switch by changing the constant at the top of `hijack.go`.
The client will switch to a TLS client if TLS is enabled.

The server will wait for a request on an endpoint `/foo` and the client will establish a TCP connection and send an HTTP GET request to the server's `/foo` endpoint. The server will immediately grab the raw TCP stream and reply with a message. The client will use its raw TCP stream and read the response. Messages will be sent back and forth continuously.

# Instructions

To run the server you first need a server certificate. Run the following commands to generate a cert and key.

```
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

Next if you want to inspect the packets, run the following in a separate terminal.

```
sudo tcpdump  -s0 -vv -n port 8081 -i lo -A
```

Finally, run the server.

```
go run hijack.go
```

If TLS is turned off you should be able to see the messages "Hijack server" and "Hello world" repeatedly. With TLS turned on you will see cipher text/random characters.

This shows that TLS is maintained when no longer using HTTP.