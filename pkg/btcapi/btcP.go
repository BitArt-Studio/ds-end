package btcapi

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"gohub/pkg/config"
)

var (
	NetParams     *chaincfg.Params
	Client        *ApiClient
	ChargeAddress btcutil.Address
)

func InitBtc() {
	NetParams = &chaincfg.MainNetParams
	Client = NewClient()
	addressStr := config.Get("service_fee.receive_address")

	var err error
	ChargeAddress, err = btcutil.DecodeAddress(addressStr, NetParams)
	if err != nil {
		panic(err)
	}
}
