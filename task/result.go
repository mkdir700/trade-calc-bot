package task

import (
	"math"

	"strconv"
)

type OpenPositionResult struct {
	OpenPosition
	Leverage               int     // 杠杆
	Margin                 float64 // 保证金
	PositionSize           float64 // 仓位大小
	MaxLoss                float64 // 基于本金的最大亏损
	RemainCapitalAfterLoss float64 // 止损后剩余本金
}

func NewOpenPositionResult(op OpenPosition) *OpenPositionResult {
	positionSize := math.Ceil(op.Capital * op.CapitalLossRatio / op.LossRatio)
	leverage := int(math.Ceil(positionSize / op.Capital))
	maxLoss := math.Round(op.Capital * op.CapitalLossRatio)
	remainCapitalAfterLoss := op.Capital - maxLoss
	margin := math.Round(positionSize / float64(leverage))

	return &OpenPositionResult{
		OpenPosition:           op,
		Leverage:               leverage,
		Margin:                 margin,
		PositionSize:           positionSize,
		MaxLoss:                maxLoss,
		RemainCapitalAfterLoss: remainCapitalAfterLoss,
	}
}

func (op *OpenPositionResult) ShowMessage() string {
	var message string
	message += "杠杆: " + strconv.Itoa(op.Leverage) + "\n"
	message += "保证金: " + strconv.FormatFloat(op.Margin, 'f', -1, 64) + "\n"
	message += "仓位大小: " + strconv.FormatFloat(op.PositionSize, 'f', -1, 64) + "\n"
	message += "本金的最大亏损: " + strconv.FormatFloat(op.MaxLoss, 'f', -1, 64) + "\n"
	message += "止损后剩余本金: " + strconv.FormatFloat(op.RemainCapitalAfterLoss, 'f', -1, 64) + "\n"
	return message
}
