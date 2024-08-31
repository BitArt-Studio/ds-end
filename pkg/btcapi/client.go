package btcapi

import (
	"gohub/pkg/app"
	"gohub/pkg/config"
	"io"
)

type ApiClient struct {
	baseURL     string
	unisatURL   string
	bearerToken string
}

func NewClient() *ApiClient {
	baseURL := ""
	unisatURL := ""
	if app.IsProduction() {
		baseURL = "https://mempool.space/api"
		unisatURL = "https://open-api.unisat.io/v1/indexer"
	} else {
		baseURL = "https://mempool-testnet.fractalbitcoin.io/api"
		unisatURL = "https://open-api-fractal-testnet.unisat.io/v1/indexer"
	}
	return &ApiClient{
		baseURL:     baseURL,
		unisatURL:   unisatURL,
		bearerToken: config.Get("unisat_api_key"),
	}
}

func (c *ApiClient) mempoolRequest(method, subPath string, requestBody io.Reader) ([]byte, error) {
	return Request(method, c.baseURL, subPath, requestBody, "")
}

func (c *ApiClient) mempoolBaseRequest(method, basePath string, requestBody io.Reader) ([]byte, error) {
	return Request(method, basePath, "", requestBody, "")
}

func (c *ApiClient) unisatRequest(method, subPath string, requestBody io.Reader) ([]byte, error) {
	return Request(method, c.unisatURL, subPath, requestBody, c.bearerToken)
}
