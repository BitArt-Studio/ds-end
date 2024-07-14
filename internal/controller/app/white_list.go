package app

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/request/app"
	"gohub/internal/request/validators"
	"gohub/internal/service"
	"gohub/pkg/response"
)

type WhiteListController struct {
}

var whiteService = service.WhiteList

func (ctl *WhiteListController) Validate(c *gin.Context) {
	req := app.AddressReq{}
	if ok := validators.Validate(c, &req); !ok {
		return
	}

	flag := whiteService.Validate(req.Address)

	response.SuccessData(c, gin.H{
		"isValid": flag,
	})
}
