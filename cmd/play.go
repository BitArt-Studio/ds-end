package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gohub/pkg/btcapi"
	"gohub/pkg/logger"
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
	price, err := btcapi.Client.LastBlockHeight()
	if err != nil {
		logger.Errorf("%+v", err)
	}
	fmt.Printf("price: %d\n", price)
}
