package main

import (
	"fmt"
	lb "lbotomy/algo"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MainAddress   string
	MockAddresses []string
	Algorithm     string
}

var cfg Config
var requestCount uint64

func main() {
	initCfg()
	openMockServers(cfg.MockAddresses)
	openMainServer(cfg.MainAddress)
}

func initCfg() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	cfg = Config{
		MainAddress:   os.Getenv("LB_URL"),
		MockAddresses: strings.Split(os.Getenv("MOCK_SERVER_URLS"), ","),
		Algorithm:     strings.ToLower(os.Getenv("ALGO")),
	}

	if cfg.MainAddress == "" || len(cfg.MockAddresses) == 0 {
		panic("Configuration is missing! Check your .env file.")
	}
}

func openMainServer(address string) *http.ServeMux {
	fmt.Printf("Load Balancer on %s", address)
	addressUrl, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		distributionAlgorithm().ServeHTTP(w, r)
	})

	err = http.ListenAndServe(":"+addressUrl.Port(), mux)
	if err != nil {
		panic(err)
	}

	return mux
}

func distributionAlgorithm() *httputil.ReverseProxy {
	algo := strings.ToLower(string(os.Getenv("ALGO")))
	addresses := cfg.MockAddresses
	switch algo {
	case "round robin":
		return lb.RoundRobinInit(addresses, &requestCount)
	default:
		panic("Invalid algo")
	}
}

func openMockServers(addresses []string) {
	for _, address := range addresses {
		go func() {
			url, err := url.Parse(address)
			if err != nil {
				panic(err)
			}

			mux := http.NewServeMux()
			// fmt.Printf("\nListening on %s", address)
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				currentTime := time.Now().Format("2006-01-02 15:04:05")
				fmt.Printf("\n[%s]-[ %s >>> %s]-[%s] -- %s \n", currentTime, cfg.MainAddress, url.String(), r.Method, r.URL.Path)
			})

			err = http.ListenAndServe(":"+url.Port(), mux)
			if err != nil {
				panic(err)
			}
		}()
	}
	//printout
	fmt.Printf("\nListening on : \n")
	for _, address := range addresses {
		fmt.Printf("\n%s\n", address)
	}

}
