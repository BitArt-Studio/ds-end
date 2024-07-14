package btcapi

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}

type Data struct {
	Cursor                int    `json:"cursor"`
	Total                 int    `json:"total"`
	TotalConfirmed        int    `json:"totalConfirmed"`
	TotalUnconfirmed      int    `json:"totalUnconfirmed"`
	TotalUnconfirmedSpend int    `json:"totalUnconfirmedSpend"`
	Utxo                  []Utxo `json:"utxo"`
}

type Utxo struct {
	Address      string        `json:"address"`
	CodeType     int           `json:"codeType"`
	Height       int           `json:"height"`
	Idx          int           `json:"idx"`
	Inscriptions []Inscription `json:"inscriptions"`
	IsOpInRBF    bool          `json:"isOpInRBF"`
	Satoshi      int64         `json:"satoshi"`
	ScriptPk     string        `json:"scriptPk"`
	ScriptType   string        `json:"scriptType"`
	Txid         string        `json:"txid"`
	Vout         int           `json:"vout"`
}

type Inscription struct {
	InscriptionId     string `json:"inscriptionId"`
	InscriptionNumber int    `json:"inscriptionNumber"`
	IsBRC20           bool   `json:"isBRC20"`
	Moved             bool   `json:"moved"`
	Offset            int    `json:"offset"`
}

// UTXOs is a slice of UTXO
type UTXOs []Response

func (c *ApiClient) ListUnspent(address btcutil.Address) ([]*UnspentOutput, error) {
	res, err := c.unisatRequest(http.MethodGet, fmt.Sprintf("/address/%s/utxo-data?cursor=%d&size=%d", address.EncodeAddress(), 0, 16), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resData Response
	err = json.Unmarshal(res, &resData)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	unspentOutputs := make([]*UnspentOutput, 0)
	for _, utxo := range resData.Data.Utxo {
		txHash, err := chainhash.NewHashFromStr(utxo.Txid)
		if err != nil {
			return nil, err
		}
		scriptPk, err := hex.DecodeString(utxo.ScriptPk)
		if err != nil {
			return nil, err
		}

		unspentOutputs = append(unspentOutputs, &UnspentOutput{
			Outpoint: wire.NewOutPoint(txHash, uint32(utxo.Vout)),
			Output:   wire.NewTxOut(utxo.Satoshi, scriptPk),
		})
	}
	return unspentOutputs, nil
}
