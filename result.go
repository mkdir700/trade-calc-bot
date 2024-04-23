package main

import (
	"errors"
	"math"

	"strconv"
)

type OpenPositionResult struct {
	OpenPositionParams
	Leverage               int     // 杠杆
	Margin                 float64 // 保证金
	PositionSize           float64 // 仓位大小
	MaxLoss                float64 // 基于本金的最大亏损
	RemainCapitalAfterLoss float64 // 止损后剩余本金
}

func NewOpenPositionResult(op OpenPositionParams) *OpenPositionResult {
	return &OpenPositionResult{
		OpenPositionParams: op,
	}
}

func (op *OpenPositionResult) Calculate() error {
	if op.Capital == 0 {
		return errors.New("本金不能为0")
	}
	if op.CapitalLossRatio == 0 {
		return errors.New("本金亏损比例不能为0")
	}
	if op.LossRatio == 0 {
		return errors.New("止损比例不能为0")
	}
	op.PositionSize = math.Ceil(op.Capital * op.CapitalLossRatio / op.LossRatio)
	op.Leverage = int(math.Ceil(op.PositionSize / op.Capital))
	op.MaxLoss = math.Round(op.Capital * op.CapitalLossRatio)
	op.RemainCapitalAfterLoss = op.Capital - op.MaxLoss
	op.Margin = math.Round(op.PositionSize / float64(op.Leverage))
	return nil
}

func (op *OpenPositionResult) BuildText() string {
	var message string
	message += "杠杆: " + strconv.Itoa(op.Leverage) + "\n"
	message += "保证金: " + strconv.FormatFloat(op.Margin, 'f', -1, 64) + "\n"
	message += "仓位大小: " + strconv.FormatFloat(op.PositionSize, 'f', -1, 64) + "\n"
	message += "本金的最大亏损: " + strconv.FormatFloat(op.MaxLoss, 'f', -1, 64) + "\n"
	message += "止损后剩余本金: " + strconv.FormatFloat(op.RemainCapitalAfterLoss, 'f', -1, 64) + "\n"
	return message
}
