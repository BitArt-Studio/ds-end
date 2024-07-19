package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gohub/internal/errorI"
	"gohub/internal/request/app"
	"gohub/internal/request/validators"
	"gohub/internal/service"
	"gohub/pkg/logger"
	"gohub/pkg/page"
	"gohub/pkg/response"
	"strconv"
)

type OrderController struct {
}

var orderService = service.Order

func (oc *OrderController) Save(c *gin.Context) {
	req := app.OrderCreateReq{}
	if ok := validators.Validate(c, &req); !ok {
		return
	}

	orderDO, err := orderService.Save(req)

	if err != nil {
		logger.Errorv(err)
		if errors.Is(err, errorI.OrderSeedNoFind) {
			response.Error10001(c, err)
		} else {
			response.ErrorStr(c, "订单保存失败")
		}
	} else {
		response.SuccessData(c, gin.H{
			"payAddress":  orderDO.PayAddress,
			"estimateFee": orderDO.EstimateFee,
			"hSeed":       orderDO.HSeed,
			"orderId":     strconv.FormatInt(orderDO.OrderId, 10),
		})
	}
}

func (oc *OrderController) Execute(c *gin.Context) {
	req := app.OrderExecuteReq{}
	if ok := validators.Validate(c, &req); !ok {
		return
	}

	orderId, err := strconv.ParseInt(req.OrderId, 10, 64)
	if err != nil {
		response.ErrorStr(c, "订单号格式错误")
	}

	orderDO, err := orderService.ExecuteOrder(orderId)

	if err != nil {
		logger.Errorv(err)
		if errors.Is(err, errorI.OrderBalanceInsufficientError) {
			response.Error10001(c, err)
		} else if errors.Is(err, errorI.OrderNoExist) {
			response.Error10002(c, err)
		} else {
			response.ErrorStr(c, "订单执行失败")
		}
	} else {
		response.SuccessData(c, gin.H{
			"RevealTxHash":   orderDO.RevealTxHash,
			"InscriptionsId": orderDO.InscriptionsId,
		})
	}
}

func (oc *OrderController) Page(c *gin.Context) {
	req := page.Req{}
	if ok := validators.Validate(c, &req); !ok {
		return
	}

	pageRes, err := orderService.PageOrder(req)

	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "分页失败")
	}

	type PageResp struct {
		PayAddress     string
		Address        string
		EstimateFee    int64
		OrderId        int64
		FeeRate        int64
		HSeed          string
		RevealTxHash   string
		InscriptionsId string
		UsdPrice       float64
		BtcPrice       float64
	}

	r := page.Resp[PageResp]{}
	r.Total = pageRes.Total
	for i := range pageRes.List {
		list := pageRes.List[i]
		r.List = append(r.List, PageResp{
			PayAddress:     list.PayAddress,
			Address:        list.Address,
			EstimateFee:    list.EstimateFee,
			OrderId:        list.OrderId,
			FeeRate:        list.FeeRate,
			HSeed:          list.HSeed,
			RevealTxHash:   list.RevealTxHash,
			InscriptionsId: list.InscriptionsId,
			UsdPrice:       list.UsdPrice,
			BtcPrice:       list.BtcPrice,
		})
	}

	response.SuccessData(c, r)
}
