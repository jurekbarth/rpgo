package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var rsaKeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAoHPLSLplMH1YKQNmp/kxiC8LX4gNTz2RRHtx/tRGJKddwtua
+QTLpQ965rD2winaNOv7zcAtjMsQmepG02vxXoy+cvy3rXgHZKh1N4+A+99TRj5P
cY1eOMEzEdUVnWNKiL5x17fPUOxOc/iJQI1zth7+LI5IbfEGcxgxhONYu3g55hAS
iNVkCTx7aLGDGcuSipH4fVB0M9MkO8CO0EiMeHUYd3ieqrP2FYDsoZtmANrogrUl
ysRdP3I623jr9wF0RshT0qANdnF90aGx6glIzMRpUTCcD3tNoGPc/Ow+Nr2SpPbG
8zrEMriSzwov8edU7HeGur1gssEH5ArwQuYPvQIDAQABAoIBAQCUMP5SyJzGwS3Y
i2SXxUbjIZgefnjUc+ekWWM62fGCzvWBD/S9A5nWdEqtoEn3oFIByOaC7HjlbXOC
xGbvw+Vkzxbi+tfmJlKlvBSu4SJe/q9Z1BjppoicYIv7b1OMTnU7gLGCbCjU87ut
zqFtdnelgFB+9Fae/BpZ2MF7m8KLOZOWE0Am2lgUx5nyLtx9bdYEiIog5pOgATJ1
Yc9HtVoHIfnqZzIrLAoJRndWIqnU2ZeISEK2MyY+9NOd643tIJSa8VL6awexkO/8
VugLWtuk2+zyZBbT5ufZueZGYp1mYRJ48Amh07WY4azk1aLOtNXnPJq6WEQkJGtg
mZa8oYlhAoGBANJWZxEbyCoOGF5ddkwXYIRtU2L1afLxooP1UZJI25bMCrHkzXiS
OFUw+GN/uPPeqCPS5toJDF5eqeyzmOdfpPkeaR4Nyk8eASqpm6DjhEvpr6JhbOKL
+FUVL80+LyxP19QkLm/5vIZBKWZs8J+X+vAO0XQguxn6flZhiEGhAkrZAoGBAMNJ
ACr4gT0UpC4ffafSmGEl5k9Z1unH610jx/csI+pKSjGaQ5gkuvZozhS6mD3VzSBU
VW/v5UpVfNI4GjtTMbkkKOWOO8+Fzfx9scX8mzJBKDMABJTMbEgD7G20JyOvvo4s
hdN+BOy5byzC/h725OA2CU1utWNqE/6g0X7EgHWFAoGAUa1dnoYkRzhr/BDdBBU7
1JDDhbT47G8qhYV4pI6IPtmC+at4om5dU6+NdM2/G2wF7MtT+6zx0Z9+6ryfDpHU
dSx680G1ot1q5I8yMNrIn9Xh7vNYHezuhNOSWWfhV5q1m9pk8fSPYa7iDbUWB1M0
DY4jha3EGgVsk8yR5bJJOpkCgYEAna1vyUJld6AXAHbEyqCsEKS9VQzBDnoxfD7L
0rN9PEtHpM1eDpZ5r0PoQax4CFV9DsGJSpx0kpR7+HD8HTKLT2X274LsoB71twz2
YVoZJXaesq8tA8gbFfq1B88SWyonvjwMwjtaVplTPt0iunW3T6HR2Qeuxdp80nef
L7AR2NECgYB0afYyJ3XDjKeUnf6EhLRFSqswNIj+l6u3k9sRYoyC3VBXSKAT5T+W
r4s0r8xvAg3air34lqnfVxonbsat0TFAO4SrJYwnwnSoIdvesNkrKmr7NzwUtwpx
SiFNZrKeYsaGp978l0vp0f4s7YZqHisVlkV4rxIYUaTByUAyWMuuHw==
-----END RSA PRIVATE KEY-----`)

var rsaCertPEM = []byte(`-----BEGIN CERTIFICATE-----
MIIDiDCCAnACCQC9tV/d9CmTijANBgkqhkiG9w0BAQsFADCBhTELMAkGA1UEBhMC
REUxDzANBgNVBAgMBkJheWVybjEPMA0GA1UEBwwGTXVuaWNoMQswCQYDVQQKDAJW
STEMMAoGA1UECwwDREVWMRYwFAYDVQQDDA1qdXJla2JhcnRoLmRlMSEwHwYJKoZI
hvcNAQkBFhJwb3N0QGp1cmVrYmFydGguZGUwHhcNMTgwNTA0MDU1MTM3WhcNMjgw
NTAxMDU1MTM3WjCBhTELMAkGA1UEBhMCREUxDzANBgNVBAgMBkJheWVybjEPMA0G
A1UEBwwGTXVuaWNoMQswCQYDVQQKDAJWSTEMMAoGA1UECwwDREVWMRYwFAYDVQQD
DA1qdXJla2JhcnRoLmRlMSEwHwYJKoZIhvcNAQkBFhJwb3N0QGp1cmVrYmFydGgu
ZGUwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCgc8tIumUwfVgpA2an
+TGILwtfiA1PPZFEe3H+1EYkp13C25r5BMulD3rmsPbCKdo06/vNwC2MyxCZ6kbT
a/FejL5y/LeteAdkqHU3j4D731NGPk9xjV44wTMR1RWdY0qIvnHXt89Q7E5z+IlA
jXO2Hv4sjkht8QZzGDGE41i7eDnmEBKI1WQJPHtosYMZy5KKkfh9UHQz0yQ7wI7Q
SIx4dRh3eJ6qs/YVgOyhm2YA2uiCtSXKxF0/cjrbeOv3AXRGyFPSoA12cX3RobHq
CUjMxGlRMJwPe02gY9z87D42vZKk9sbzOsQyuJLPCi/x51Tsd4a6vWCywQfkCvBC
5g+9AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJnEK26Yu1qLQld9knhCa1fWjBBk
NtZWRNxfykkLU+aeA5yQzr+rMRpIazIP5KcJ80eCqXue0h7N9PYarY33WSkvLEBC
8Tc3Hm69vfMguqKWo/oqQlsMSG1o3HrwU7Sw5d/smFpj0SHet6/aIVMQUEaqez/u
3DywGlYIKe64gvtHqCMgXkAFaxm/Er2l85hyPdWAxiR0ejOGd1+psHeEH2rqCMoT
XvUg+Qw5Eep6XDyq43MaNFywBqcZYai1YZnacJ2Cc6fmraKDWPtVMvwh4Jj0LBGb
F5Eyba7Xn8syaOD8U1dhOa8A4Q3rMe0hA3LWI34O6goGbUzpBeXjWfBbnhU=
-----END CERTIFICATE-----`)

type transport struct{}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = req.URL.Host
	return http.DefaultTransport.RoundTrip(req)
}

func modifyResponse(res *http.Response) error {
	res.Header.Add("Access-Control-Allow-Origin", "*")
	return nil
}

func main() {
	proxyURL := flag.String("target", "https://jurekbarth.de", "url you want to proxy")
	portFlag := flag.String("port", "7777", "port the proxy should listen on")
	cors := flag.Bool("cors", false, "enable cors")
	flag.Parse()
	u, err := url.Parse(*proxyURL)
	if err != nil {
		log.Fatal(err)
	}

	if u.Scheme == "" {
		log.Fatal("missing scheme, in proxyUrl: ", *proxyURL)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	})
	reverseProxy.Transport = &transport{}

	var cert tls.Certificate
	var certErr error

	cert, certErr = tls.X509KeyPair(rsaCertPEM, rsaKeyPEM)

	if certErr != nil {
		log.Fatal(certErr)
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	server := http.Server{
		// Other options
		Addr:      ":" + *portFlag,
		TLSConfig: tlsConfig,
	}
	if *cors {
		reverseProxy.ModifyResponse = modifyResponse
	}

	http.Handle("/", reverseProxy)
	fmt.Printf("Proxy %v on https://localhost:%v\n", *proxyURL, *portFlag)
	log.Fatal(server.ListenAndServeTLS("", ""))
}
