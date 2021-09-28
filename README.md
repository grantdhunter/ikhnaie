# ikhnaie
Location sharing app to help coordinate with friend. Hopefully to be a Matrix/Element widget.

# Config Example
config.toml
```
host = "localhost"
port = "8080"
mapbox_api_key = "<api key here>"
```

# generating develoment certs
```
openssl genrsa -out server.key 2048
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
