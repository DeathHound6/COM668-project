#!/bin/bash

# Create the cert
openssl req -x509 -sha256 -nodes -newkey rsa:2048 -days 365 -keyout localhost.key -out localhost.crt

# Add the cert to keychain
open localhost.crt