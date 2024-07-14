package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"gohub/internal/dao"
	"gohub/pkg/console"
	"gorm.io/gorm"
	"io"
	"os"
	"strconv"
)

var CmdPlay = &cobra.Command{
	Use:   "play",
	Short: "Likes the Go Playground, but running at our application context",
	Run:   runPlay,
}

type Data struct {
	HSeed   string `json:"hSeed"`
	Address string `json:"address"`
	TempAdd string `json:"tempAdd"`
}

// 调试完成后请记得清除测试代码
func runPlay(cmd *cobra.Command, args []string) {
	// 打开 JSON 文件
	jsonFile, err := os.Open("minted-data.json")
	if err != nil {
		console.Exit(err.Error())
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			console.Exit(err.Error())
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

	err = dao.Transaction(func(tx *gorm.DB) error {
		for _, value := range dataMap {
			//if err := dao.Order.Tx(tx).New().Create(&model.OrderDO{
			//	OrderId:        snowflakeP.Node.Generate().Int64(),
			//	HSeed:          value.HSeed,
			//	Address:        value.Address,
			//	PayAddress:     value.TempAdd,
			//	Status:         enum.OrderStatusComplete.Code,
			//	InscriptionsId: key,
			//}).Error; err != nil {
			//	return err
			//}

			if err := dao.Seed.Tx(tx).Model().Where("hSeed = ?", value.HSeed).Update("address", value.Address).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		console.Exit(err.Error())
	}

}
