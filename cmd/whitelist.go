package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"gohub/internal/dao"
	"gohub/internal/model"
	"gohub/pkg/btcapi"
	"gohub/pkg/console"
	"io"
	"os"
	"strconv"
)

var CmdWhiteList = &cobra.Command{
	Use:   "whitelist",
	Short: "whitelist import",
	Run:   runWhiteList,
	Args:  cobra.NoArgs,
}

type Data struct {
	HSeed   string `json:"hSeed"`
	Address string `json:"address"`
	TempAdd string `json:"tempAdd"`
}

func runWhiteList(cmd *cobra.Command, args []string) {
	// 打开 JSON 文件
	jsonFile, err := os.Open("white_list_cm.json")
	if err != nil {
		console.Exit(err.Error())
	}

	defer func(jsonFile *os.File) {
		errf := jsonFile.Close()
		if errf != nil {
			console.Exit(errf.Error())
		}
	}(jsonFile)

	// 读取文件内容
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		console.Exit(err.Error())
	}

	// 定义一个 map 来存储解码后的数据
	var dataMap map[string]Data

	// 解码 JSON 数据到 map
	err = json.Unmarshal(byteValue, &dataMap)
	if err != nil {
		console.Exit(err.Error())
	}

	console.Success(strconv.Itoa(len(dataMap)))

	count := 1
	whiteListMap := make(map[string]string)
	for k, _ := range dataMap {
		address, err := btcapi.Client.GetSAddressByInscriptionId(k)
		if err != nil {
			console.Exit(err.Error())
		}
		whiteListMap[address] = k
		count++
		console.Success("count:" + strconv.Itoa(count) + " address: " + address + " inscriptionId: " + k)
	}
	whiteListDos := make([]*model.WhiteListDO, 0)
	for k, _ := range whiteListMap {
		whiteListDos = append(whiteListDos, &model.WhiteListDO{
			Address: k,
			Used:    false,
		})
	}

	if err = dao.WhiteList.New().CreateInBatches(whiteListDos, len(whiteListDos)).Error; err != nil {
		console.Exit(err.Error())
	}
}
