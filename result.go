package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
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
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("本金: %0.2f\n", op.Capital))
	builder.WriteString(fmt.Sprintf("本金亏损比例: %0.2f%%\n", op.CapitalLossRatio*100))
	builder.WriteString(fmt.Sprintf("开仓止损比例: %0.2f%%\n", op.LossRatio*100))
	builder.WriteString(fmt.Sprintf("杠杆: %d\n", op.Leverage))
	builder.WriteString(fmt.Sprintf("保证金: %0.2f\n", op.Margin))
	builder.WriteString(fmt.Sprintf("仓位大小: %0.2f\n", op.PositionSize))
	builder.WriteString(fmt.Sprintf("预计亏损: %0.2f\n", op.MaxLoss))
	builder.WriteString(fmt.Sprintf("剩余本金: %0.2f\n", op.RemainCapitalAfterLoss))
	message = builder.String()
	return message
}
