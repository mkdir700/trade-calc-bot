package main

// import (
// 	"capital_calculator_tgbot/component"
// 	"context"

// 	"github.com/go-telegram/bot"
// 	"github.com/go-telegram/bot/models"
// )

// type State int

// const (
// 	Init State = iota
// 	StartMenu
// 	OpenPositionMenu
// 	AddPositionMenu
// 	InputCapitalState
// 	InputCapitalErrorState
// 	InputCapitalLossRadioState
// 	InputCapitalLossRadioErrorState
// 	InputLossRatioState
// 	InputLossRatioErrorState
// )

// type Event int

// const (
// 	ClickOpenPosition Event = iota + 100
// 	ClickAddPosition
// 	ClickInputCapitalButton
// 	ClickInputCapitalLossRatioButton
// 	SelectCapitalLossRatio
// 	ClickInputLossRatioButton
// 	ClickCancelButton
// 	ClickCalculateButton
// 	InvalidAnswer
// 	BackToOpenPositionMenu
// 	BackToAddPositionMenu
// 	BackToStartMenu
// )

// type StateMachine struct {
// 	State        State
// 	MainMenu     *component.MainMenu
// 	OpenPosition *component.OpenPositionMenu
// }

// func NewStateMachine() *StateMachine {
// 	return &StateMachine{
// 		State:        State(StartMenu),
// 		MainMenu:     component.NewMainMenu(),
// 		OpenPosition: component.NewOpenPositionMenu(),
// 	}
// }

// func (self *StateMachine) HandleEvent(event Event, ctx context.Context, b *bot.Bot, update *models.Update) {
// 	switch self.State {
// 	case Init:
// 		self.State = StartMenu
// 	case StartMenu:
// 		switch event {
// 		case ClickOpenPosition:
// 			self.State = OpenPositionMenu
// 			self.OpenPosition.ReplaceShow(
// 				ctx,
// 				b,
// 				update.CallbackQuery.Message.Message.Chat.ID,
// 				update.CallbackQuery.Message.Message.ID,
// 				update.CallbackQuery.InlineMessageID,
// 			)
// 		case ClickAddPosition:
// 			self.State = AddPositionMenu
// 		}
// 	case OpenPositionMenu:
// 		switch event {
// 		case ClickInputCapitalButton:
// 			self.State = InputCapitalState
// 		case ClickInputCapitalLossRatioButton:
// 			self.State = InputCapitalLossRadioState
// 		case ClickInputLossRatioButton:
// 			self.State = InputLossRatioState
// 		case ClickCancelButton:
// 			self.State = StartMenu
// 		case ClickCalculateButton:
// 			// 开始计算
// 		}
// 	case AddPositionMenu:
// 	case InputCapitalState:
// 		switch event {
// 		case InvalidAnswer:
// 			// 提示错误信息
// 			self.State = InputCapitalState
// 		case BackToOpenPositionMenu:
// 			self.State = OpenPositionMenu
// 		case BackToAddPositionMenu:
// 			self.State = AddPositionMenu
// 		case BackToStartMenu:
// 			// 取消当前计算任务
// 			self.State = StartMenu
// 		}
// 	}

// }
