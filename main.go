package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Proxy is part of a config
type Proxy struct {
	WriteCors   bool   `json:"writeCors"`
	Proxyhost   string `json:"proxyhost"`
	Host        string `json:"host"`
	RewritePath string `json:"rewritePath"`
	Port        int    `json:"port"`
}

// Config thats a config
type Config struct {
	Version            int  `json:"version"`
	Port               int  `json:"port"`
	InsecureSkipVerify bool `json:"insecureSkipVerify"`
	Certs              []struct {
		Key  string `json:"key"`
		Cert string `json:"cert"`
	} `json:"certs"`
	Proxy []Proxy `json:"proxy"`
}

var (
	config     *Config
	configLock = new(sync.RWMutex)
)

func loadConfig(fail bool) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("open config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
		if fail {
			os.Exit(1)
		}
	}
	configLock.Lock()
	config = temp
	configLock.Unlock()
}

func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

type transport struct{}

func getRewritePath(path string) string {
	for _, proxy := range config.Proxy {
		if strings.HasPrefix(path, proxy.RewritePath) {
			return proxy.RewritePath
		}
	}
	return ""
}

func getProxyConfig(host string) (*Proxy, error) {
	for _, proxy := range config.Proxy {
		if proxy.Host == host {
			return &proxy, nil
		}
	}
	return nil, errors.New("nothing found")
}

func rewrite(path string) string {
	rewritePath := getRewritePath(path)
	return strings.TrimPrefix(path, rewritePath)
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Path = rewrite(req.URL.Path)
	req.Host = req.URL.Host
	return http.DefaultTransport.RoundTrip(req)
}

func modifyResponse(res *http.Response) error {
	res.Header.Set("Access-Control-Allow-Origin", "*")
	if res.StatusCode > 299 && res.StatusCode < 400 {
		url, err := res.Location()
		if err != nil {
			return nil
		}
		host := fmt.Sprintf("%v://%v", url.Scheme, url.Hostname())
		proxyConf, err := getProxyConfig(host)
		if err != nil {
			return err
		}
		rewriteLocation := fmt.Sprintf("https://%v%v", proxyConf.Proxyhost, url.Path)
		url, err = url.Parse(rewriteLocation)
		if err != nil {
			return err
		}

		rewriteLocation = fmt.Sprintf("%v://%v:%v%v", url.Scheme, url.Hostname(), config.Port, url.Path)
		res.Header.Set("Location", rewriteLocation)
	}

	return nil
}

func main() {
	loadConfig(true)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify}

	cfg := &tls.Config{}

	for _, certPair := range config.Certs {
		cert, err := tls.LoadX509KeyPair(certPair.Cert, certPair.Key)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Certificates = append(cfg.Certificates, cert)
	}
	server := http.Server{
		Addr:      ":" + strconv.Itoa(config.Port),
		TLSConfig: cfg,
	}
	var reverseProxies []*httputil.ReverseProxy
	for _, p := range config.Proxy {
		u, err := url.Parse(p.Host)
		if err != nil {
			log.Fatal(err)
		}

		reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: u.Scheme,
			Host:   u.Host + ":" + strconv.Itoa(p.Port),
			Path:   u.Path,
		})

		reverseProxy.Transport = &transport{}
		if p.WriteCors {
			reverseProxy.ModifyResponse = modifyResponse
		}
		reverseProxies = append(reverseProxies, reverseProxy)
	}

	for idx, p := range config.Proxy {
		http.Handle(p.Proxyhost+"/", reverseProxies[idx])
		fmt.Printf("Proxy: https://%v:%v --> %v:%v \n", p.Proxyhost, config.Port, p.Host, p.Port)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Unknown")
	})

	log.Fatal(server.ListenAndServeTLS("", ""))
}
