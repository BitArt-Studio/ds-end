package btcapi

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"gohub/pkg/logger"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
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

type Brc20Detail struct {
	Address                string `json:"address"`
	OverallBalance         string `json:"overallBalance"`
	TransferableBalance    string `json:"transferableBalance"`
	AvailableBalance       string `json:"availableBalance"`
	AvailableBalanceSafe   string `json:"availableBalanceSafe"`
	AvailableBalanceUnSafe string `json:"availableBalanceUnSafe"`
}

type Brc20PageResponse struct {
	Height int           `json:"height"`
	Total  int           `json:"total"`
	Start  int           `json:"start"`
	Detail []Brc20Detail `json:"detail"`
}

func (c *ApiClient) ListUnspent(address btcutil.Address) ([]*UnspentOutput, error) {

	type Data struct {
		Cursor                int    `json:"cursor"`
		Total                 int    `json:"total"`
		TotalConfirmed        int    `json:"totalConfirmed"`
		TotalUnconfirmed      int    `json:"totalUnconfirmed"`
		TotalUnconfirmedSpend int    `json:"totalUnconfirmedSpend"`
		Utxo                  []Utxo `json:"utxo"`
	}

	res, err := c.unisatRequest(http.MethodGet, fmt.Sprintf("/address/%s/utxo-data?cursor=%d&size=%d", address.EncodeAddress(), 0, 16), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resData Response
	err = json.Unmarshal(res, &resData)
	if err != nil {
		logger.Error(string(res))
		return nil, errors.WithStack(err)
	}

	unspentOutputs := make([]*UnspentOutput, 0)
	data, ok := resData.Data.(Data)
	if !ok {
		return nil, errors.New("failed to parse data")
	}
	for _, utxo := range data.Utxo {
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

func (c *ApiClient) GetSAddressByInscriptionId(inscriptionId string) (string, error) {
	res, err := c.unisatRequest(http.MethodGet, fmt.Sprintf("/inscription/info/%s", inscriptionId), nil)
	if err != nil {
		return "", errors.WithStack(err)
	}

	var resData Response
	err = json.Unmarshal(res, &resData)
	if err != nil {
		logger.Error(string(res))
		return "", errors.WithStack(err)
	}

	dataMap, ok := resData.Data.(map[string]interface{})
	if !ok {
		return "", errors.New("failed to parse data")
	}

	return dataMap["address"].(string), nil
}

func (c *ApiClient) GetBrc20Page(ticker string, start int, limit int) (*Brc20PageResponse, error) {
	res, err := c.unisatRequest(http.MethodGet, fmt.Sprintf("/brc20/%s/holders?start=%d&limit=%d", ticker, start, limit), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resData Response
	err = json.Unmarshal(res, &resData)
	if err != nil {
		logger.Error(string(res))
		return nil, errors.WithStack(err)
	}

	var page Brc20PageResponse
	dataBytes, err := json.Marshal(resData.Data)
	if err != nil {
		logger.Error(resData.Data)
		return nil, errors.New("failed to marshal data")
	}

	err = json.Unmarshal(dataBytes, &page)
	if err != nil {
		logger.Error(string(dataBytes))
		return nil, errors.New("failed to parse data")
	}

	return &page, nil
}
