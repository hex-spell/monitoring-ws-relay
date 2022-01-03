# monitoring-ws-relay
Websockets server for [esp8266-monitoring](https://github.com/Modulariz/esp8266-monitoring), relays its sockets to all the clients of the server

# Environment variables
- `SECRET`: secret to be passed as `a` query string param of the websockets handshake

# Handshake params
- `a`: `SECRET`
- `t`: Client type, for now there is only "listener", it receives all messages from every client, I'm planning to make more to reduce network usage