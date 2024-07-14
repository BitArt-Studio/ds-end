package btcapi

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

func TestListUnspent(t *testing.T) {
	client := NewClient(&chaincfg.TestNet3Params, "c77d9c1f759036aab3514721c6330807d9230485c9d240ef509e72fdbc9c053b")
	address, _ := btcutil.DecodeAddress("tb1pclrhuqrlwr2ykdd2muhnnxg2t8umepdndmar5wc4ua0z5wn86qssstj6kj", NetParams)
	unspentList, err := client.ListUnspent(address)
	if err != nil {
		t.Error(err)
	} else {
		for _, output := range unspentList {
			t.Log(output.Outpoint.Hash)
		}
	}
}
