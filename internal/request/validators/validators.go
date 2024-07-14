package validators

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/thedevsaddam/govalidator"
	"gohub/internal/request"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

type ValidatorI interface {
	Validator() map[string][]string
}

func Validate(c *gin.Context, valObj ValidatorI) bool {

	// 1. 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := c.ShouldBind(valObj); err != nil {
		logger.Errorv(err)
		response.Error405(c, errors.New("请求解析错误，请确认请求格式是否正确。参数请使用 JSON 格式。"))
		return false
	}

	// 2. 表单验证
	errs := valObj.Validator()

	jsonData, err := json.Marshal(errs)
	if err != nil {
		response.Error(c, err)
		return false
	}

	// 3. 判断验证是否通过
	if len(errs) > 0 {
		response.Error405(c, errors.New(fmt.Sprintf("参数校验错误: %s", string(jsonData))))
		return false
	}

	return true
}

func ValidateData(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	// 配置选项
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid", // 模型中的 Struct 标签标识符
		Messages:      messages,
	}
	// 开始验证
	return govalidator.New(opts).ValidateStruct()
}

// ValidatePage 自定义规则，验证分页参数
func ValidatePage(pageReq request.PageReq, errs map[string][]string) map[string][]string {

	maxPageSize := config.GetInt("page.max_page_size")
	if pageReq.PageSize > maxPageSize {
		errs["pageNo"] = append(errs["pageNo"], fmt.Sprintf("页码超出最大限制%d", maxPageSize))
	}

	fields := pageReq.Fields
	orders := pageReq.Orders

	if len(fields) != len(orders) {
		errs["fields"] = append(errs["fields"], "fields 和 orders 长度不一致")
	}

	for _, field := range fields {
		if field == "" {
			errs["fields"] = append(errs["fields"], "fields 不能为空")
			break
		}
	}

	for _, order := range orders {
		if order != "asc" && order != "desc" {
			errs["orders"] = append(errs["orders"], "order 只能是 asc 或 desc")
			break
		}
	}
	return errs
}
