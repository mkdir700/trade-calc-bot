package main

import (
	"trade_calc_bot/err"
)

type OpenPositionParams struct {
	Capital          float64 // 本金
	CapitalLossRatio float64 // 本金亏损比例
	LossRatio        float64 // 止损比例
}

type AddPositionParams struct {
	Profit           float64 // 当前盈利
	Capital          float64 // 当前本金
	CapitalLossRatio float64 // 本金亏损比例
	LossRatio        float64 // 本单止损比例
}

func NewOpenPositionParams(capital, capitalLossRatio, lossRatio float64) *OpenPositionParams {
	return &OpenPositionParams{
		Capital:          capital,
		CapitalLossRatio: capitalLossRatio,
		LossRatio:        lossRatio,
	}
}

func NewAddPositionParams(profit, capital, capitalLossRatio, lossRatio float64) *AddPositionParams {
	return &AddPositionParams{
		Profit:           profit,
		Capital:          capital,
		CapitalLossRatio: capitalLossRatio,
		LossRatio:        lossRatio,
	}
}

func (op *OpenPositionParams) SetCapital(capital float64) error {
	if capital < 0 {
		return &err.InvalidInputError{Msg: "本金不能为负数"}
	}
	op.Capital = capital
	return nil
}

func (op *OpenPositionParams) SetCapitalLossRatio(capitalLossRatio float64) error {
	if capitalLossRatio < 0 {
		return &err.InvalidInputError{Msg: "本金亏损比例不能为负数"}
	} else if capitalLossRatio > 1 {
		return &err.InvalidInputError{Msg: "本金亏损比例不能大于100%"}
	}

	op.CapitalLossRatio = capitalLossRatio
	return nil
}

func (op *OpenPositionParams) SetLossRatio(lossRatio float64) error {
	if lossRatio < 0 {
		return &err.InvalidInputError{Msg: "本单止损比例不能为负数"}
	} else if lossRatio > 1 {
		return &err.InvalidInputError{Msg: "本单止损比例不能大于100%"}
	}

	op.LossRatio = lossRatio
	return nil
}
