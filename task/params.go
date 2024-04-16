package task

import (
	"capital_calculator_tgbot/err"
)

type OpenPosition struct {
	Capital          float64 // 本金
	CapitalLossRatio float64 // 本金亏损比例
	LossRatio        float64 // 止损比例
}

type AddPosition struct {
	Profit           float64 // 当前盈利
	Capital          float64 // 当前本金
	CapitalLossRatio float64 // 本金亏损比例
	LossRatio        float64 // 本单止损比例
}

func NewOpenPosition(capital, capitalLossRatio, lossRatio float64) *OpenPosition {
	return &OpenPosition{
		Capital:          capital,
		CapitalLossRatio: capitalLossRatio,
		LossRatio:        lossRatio,
	}
}

func NewAddPosition(profit, capital, capitalLossRatio, lossRatio float64) *AddPosition {
	return &AddPosition{
		Profit:           profit,
		Capital:          capital,
		CapitalLossRatio: capitalLossRatio,
		LossRatio:        lossRatio,
	}
}

type SetCapital interface {
	SetCapital(float64) error
}

type SetCapitalLossRatio interface {
	SetCapitalLossRatio(float64) error
}

type SetLossRatio interface {
	SetLossRatio(float64) error
}

type SetProfit interface {
	SetProfit(float64) error
}

func (op *OpenPosition) SetCapital(capital float64) error {
	if capital < 0 {
		return &err.InvalidInputError{Msg: "本金不能为负数"}
	}
	op.Capital = capital
	return nil
}

func (op *OpenPosition) SetCapitalLossRatio(capitalLossRatio float64) error {
	if capitalLossRatio < 0 {
		return &err.InvalidInputError{Msg: "本金亏损比例不能为负数"}
	} else if capitalLossRatio > 1 {
		return &err.InvalidInputError{Msg: "本金亏损比例不能大于1"}
	}

	op.CapitalLossRatio = capitalLossRatio
	return nil
}

func (op *OpenPosition) SetLossRatio(lossRatio float64) error {
	if lossRatio < 0 {
		return &err.InvalidInputError{Msg: "本单止损比例不能为负数"}
	} else if lossRatio > 1 {
		return &err.InvalidInputError{Msg: "本单止损比例不能大于1"}
	}

	op.LossRatio = lossRatio
	return nil
}
