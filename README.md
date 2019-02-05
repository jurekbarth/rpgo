# Local testing reverse proxy

## Usage
First install the package using npm (you'll need at least Node.js LTS)
```
npm install -g rpgo
```
or use curl to install it globally without node and npm
```
curl -sf https://raw.githubusercontent.com/jurekbarth/rpgo/master/install-rpgo.sh | sh
```

Once that's done, you can create a `config.json` file and run `rpgo`
```
{
  "version": 1,
  "port": 1234,
  "insecureSkipVerify": true,
  "certs": [
    {
      "key": "mycert.key",
      "cert": "mycert.crt"
    }
  ],
  "proxy": [
    {
      "writeCors": true,
      "proxyhost": "mydomain.local/api",
      "host": "http://api.domain.com",
      "rewritePath": "/api",
      "port": 80
    },
    {
      "writeCors": true,
      "proxyhost": "mydomain.local/frontend",
      "host": "https://frontend-domain.com/random/sub/root",
      "rewritePath": "/frontend",
      "port": 443
    },
    {
      "writeCors": true,
      "proxyhost": "mydomain.local",
      "host": "http://localhost",
      "port": 8080
    }
  ]
}
```

Finally, run this command to list all available options
```
rpgo --help
```

## Using RPGO with Docker
`docker run -v ${PWD}/config.json:/root/config.json -v ${PWD}/certs:/root/certs jurekbarth/rpgo:latest`

## Contributing
0. You may need some experience in go to get up and running
1. Fork this repository to your own GitHub account and then clone it to your local device
2. Make your branch
3. Add your PR
