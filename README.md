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

Once that's done, you can run this command and open `https://localhost:7777`
```
rpgo -target=https://jurekbarth.de`
```

Finally, run this command to list all available options
```
rpgo --help
```

## Configuration
To customize rpgo's behaviour use the commandline flags

### Enable Cors
```
rpgo -target=https://jurekbarth.de -cors
```

### Change default port
```
rpgo -target=https://jurekbarth.de -port=1234
```

### Ignore backend cert
```
rpgo -target=https://jurekbarth.de -ignoressl
```

## Contributing
0. You may need some experience in go to get up and running
1. Fork this repository to your own GitHub account and then clone it to your local device
2. Make your branch
3. Add your PR
