package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func AcceptTermsKeyboard() interface{} {
	var keys []tgbotapi.InlineKeyboardButton
	keys = append(keys, tgbotapi.NewInlineKeyboardButtonData("Accept", AcceptTermsState))
	keys = append(keys, tgbotapi.NewInlineKeyboardButtonData("Deny", DenyTermsState))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keys)
	return keyboard
}

func AcceptFeesKeyboard() interface{} {
	var keys []tgbotapi.InlineKeyboardButton
	keys = append(keys, tgbotapi.NewInlineKeyboardButtonData("Accept", AcceptFeesState))
	keys = append(keys, tgbotapi.NewInlineKeyboardButtonData("Deny", DenyFeesState))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keys)
	return keyboard
}

func GenerateNewLinkKeyboard(paymentLink string) interface{} {
	var keys []tgbotapi.InlineKeyboardButton
	keys = append(keys, tgbotapi.NewInlineKeyboardButtonURL("Pay Now", paymentLink))
	keys = append(keys, tgbotapi.NewInlineKeyboardButtonData("Generate new payment link", GeneratePaymentLinkState))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keys)
	return keyboard
}
