# Local testing reverse proxy

## Usage
Use `rpgo -target=https://jurekbarth.de` and open `https://localhost:9001`

## Enable Cors
With an optional flag the proxy will add a headers to enable cors.
`rpgo -target=https://jurekbarth.de -cors`

## Installation
`curl -sf https://raw.githubusercontent.com/jurekbarth/rpgo/master/install-rpgo.sh | sh`
