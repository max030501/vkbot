package types

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Шаблон для Reply клавиатуры
var MainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить данные"),
		tgbotapi.NewKeyboardButton("Получить данные"),
		tgbotapi.NewKeyboardButton("Удалить данные"),
	))

// Карта, необходимая для обработки команд с Reply клавиатуры
var Comm = map[string]string{"Добавить данные": "set", "Получить данные": "get", "Удалить данные": "del"}

// Режими, необходимые для разделения Inline клавиатур
const (
	DelMode = "del"
	GetMode = "get"
)
