package controllers

import (
	"learn-hub-backend/database"
	"learn-hub-backend/models"
	"learn-hub-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RuleController struct{}

func NewRuleController() *RuleController {
	return &RuleController{}
}

// GetRuleList 获取规则列表
func (ctrl *RuleController) GetRuleList(c *gin.Context) {
	// 获取分页参数
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if current < 1 {
		current = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	// 计算偏移量
	offset := (current - 1) * pageSize

	// 查询总数
	var total int64
	database.DB.Model(&models.Rule{}).Count(&total)

	// 查询列表
	var rules []models.Rule
	database.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&rules)

	// 转换为前端格式
	var list []map[string]interface{}
	for _, rule := range rules {
		list = append(list, map[string]interface{}{
			"key":       rule.ID,
			"name":      rule.Name,
			"desc":      rule.Desc,
			"status":    rule.Status,
			"callNo":    rule.CallNo,
			"owner":     rule.Owner,
			"avatar":    rule.Avatar,
			"href":      rule.Href,
			"disabled":  rule.Disabled,
			"progress":  rule.Progress,
			"updatedAt": rule.UpdatedAt.Format("2006-01-02 15:04:05"),
			"createdAt": rule.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 前端期望的格式：RuleList = { data, total, success }
	// 注意：不使用 utils.Success，因为前端期望的格式不同
	c.JSON(200, gin.H{
		"data":    list,
		"total":   total,
		"success": true,
	})
}

// CreateOrUpdateRule 创建或更新规则
func (ctrl *RuleController) CreateOrUpdateRule(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取 method
	method, _ := req["method"].(string)
	if method == "" {
		method = "post" // 默认创建
	}

	switch method {
	case "post":
		// 创建规则
		var rule models.Rule
		
		// 从 data 字段或直接获取数据
		if data, ok := req["data"].(map[string]interface{}); ok {
			if name, ok := data["name"].(string); ok {
				rule.Name = name
			}
			if desc, ok := data["desc"].(string); ok {
				rule.Desc = desc
			}
			if status, ok := data["status"].(float64); ok {
				rule.Status = int(status)
			}
		} else {
			// 直接绑定
			if name, ok := req["name"].(string); ok {
				rule.Name = name
			}
			if desc, ok := req["desc"].(string); ok {
				rule.Desc = desc
			}
		}

		if rule.Name == "" {
			utils.BadRequest(c, "规则名称不能为空")
			return
		}
		if rule.Status == 0 {
			rule.Status = 1
		}

		if err := database.DB.Create(&rule).Error; err != nil {
			utils.InternalError(c, "创建失败: "+err.Error())
			return
		}

		utils.Success(c, map[string]interface{}{
			"key":       rule.ID,
			"name":      rule.Name,
			"desc":      rule.Desc,
			"status":    rule.Status,
			"callNo":    rule.CallNo,
			"updatedAt": rule.UpdatedAt.Format("2006-01-02 15:04:05"),
			"createdAt": rule.CreatedAt.Format("2006-01-02 15:04:05"),
		})

	case "update":
		// 更新规则
		var rule models.Rule
		var key uint

		// 获取 key
		if k, ok := req["key"].(float64); ok {
			key = uint(k)
		} else if data, ok := req["data"].(map[string]interface{}); ok {
			if k, ok := data["key"].(float64); ok {
				key = uint(k)
			}
		}

		if key == 0 {
			utils.BadRequest(c, "缺少规则ID")
			return
		}

		if err := database.DB.First(&rule, key).Error; err != nil {
			utils.Error(c, 200, "RULE_NOT_FOUND", "规则不存在")
			return
		}

		// 更新字段
		if data, ok := req["data"].(map[string]interface{}); ok {
			if name, ok := data["name"].(string); ok && name != "" {
				rule.Name = name
			}
			if desc, ok := data["desc"].(string); ok {
				rule.Desc = desc
			}
			if status, ok := data["status"].(float64); ok {
				rule.Status = int(status)
			}
		}

		if err := database.DB.Save(&rule).Error; err != nil {
			utils.InternalError(c, "更新失败: "+err.Error())
			return
		}

		utils.Success(c, map[string]interface{}{
			"key":       rule.ID,
			"name":      rule.Name,
			"desc":      rule.Desc,
			"status":    rule.Status,
			"callNo":    rule.CallNo,
			"updatedAt": rule.UpdatedAt.Format("2006-01-02 15:04:05"),
			"createdAt": rule.CreatedAt.Format("2006-01-02 15:04:05"),
		})

	case "delete":
		// 删除规则
		var key uint

		// 获取 key
		if k, ok := req["key"].(float64); ok {
			key = uint(k)
		} else if data, ok := req["data"].(map[string]interface{}); ok {
			if k, ok := data["key"].(float64); ok {
				key = uint(k)
			}
		}

		if key == 0 {
			utils.BadRequest(c, "缺少规则ID")
			return
		}

		if err := database.DB.Delete(&models.Rule{}, key).Error; err != nil {
			utils.InternalError(c, "删除失败: "+err.Error())
			return
		}

		utils.Success(c, gin.H{})

	default:
		utils.BadRequest(c, "不支持的 method: "+method)
	}
}

