package main

import (
	"container/list"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/madflojo/tasks"
	"log"
	"strconv"
	"strings"
	"time"
	"vkbot/db"
	"vkbot/logger"
	"vkbot/types"
)

// Карта для хранения списка сервисов каждого пользователя локально для обработки Inline клавиатур
var services map[int64][]string

// Список идентификаторов чатов и сообщений для последующего удаления спустя время
var queueDelMessages *list.List

// Планировщик для удаления сообщений с паролями
var scheduler *tasks.Scheduler

// Карта для хранения промежуточной информации при добавлении данных нового сервиса
var usersData map[int64]*types.UserData

// Объект бота для отправки сообщений
var bot *tgbotapi.BotAPI
var err error

// Объект конфигурации
var cfg *types.Config

func main() {
	cfgPath, err := types.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err = types.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	queueDelMessages = list.New()
	scheduler = tasks.New()
	defer scheduler.Stop()
	usersData = make(map[int64]*types.UserData)
	services = make(map[int64][]string)
	bot, err = tgbotapi.NewBotAPI(cfg.Bot.BotToken)
	logger.ForError(err)
	logger.Logs.Printf(types.AuthorizationInfo, bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	bot.Debug = true
	u.Timeout = cfg.Bot.UpdateTimeout
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			chatId := update.Message.Chat.ID
			message := update.Message.Text
			messageId := update.Message.MessageID
			msg := tgbotapi.NewMessage(chatId, "")
			if el, ok := types.Comm[message]; update.Message.IsCommand() || ok {
				usersData[chatId] = &types.UserData{
					State: types.InitState,
				}
				var com string
				if !update.Message.IsCommand() {
					com = el
				} else {
					com = update.Message.Command()
				}
				switch com {
				case "start":
					msg.ReplyMarkup = types.MainMenu
					msg.Text = types.HelpCom
				case "set":
					usersData[chatId].State = types.ServiceSet
					msg.Text = types.ServiceNameSet
				case "get":
					services[chatId] = db.GetServices(chatId)
					SendServices(chatId, 0, cfg.Bot.InlineKeyboard.CountPerPage, nil, types.GetMode)
					usersData[chatId].State = types.InitState
					continue
				case "del":
					services[chatId] = db.GetServices(chatId)
					SendServices(chatId, 0, cfg.Bot.InlineKeyboard.CountPerPage, nil, types.DelMode)
					usersData[chatId].State = types.InitState
					continue
				case "help":
					msg.Text = types.HelpCom
				default:
					msg.Text = types.ErrCom
				}
			} else {
				if user, ok := usersData[chatId]; ok {
					switch user.State {
					case types.ServiceSet:
						if len(message) > cfg.Bot.SetService.LenService {
							msg.Text = types.ErrLength
						} else {
							_, _, err = db.GetServiceByName(chatId, message)
							if err == nil {
								msg.Text = types.ServiceExists
							} else {
								usersData[chatId].Service = message
								usersData[chatId].State = types.LoginSet
								msg.Text = types.LoginNameSet
							}
						}
					case types.LoginSet:
						if len(message) > cfg.Bot.SetService.LenLogin {
							msg.Text = types.ErrLength
						} else {
							usersData[chatId].Login = message
							usersData[chatId].State = types.PassSet
							msg.Text = types.PasswordNameSet
						}
					case types.PassSet:
						if len(message) > cfg.Bot.SetService.LenPassword {
							msg.Text = types.ErrLength
						} else {
							usersData[chatId].Password = message
							usersData[chatId].State = types.InitState
							db.InsertServiceInfo(chatId,
								usersData[chatId].Service,
								usersData[chatId].Login,
								usersData[chatId].Password)
							msg.Text = fmt.Sprintf(types.AddService, usersData[chatId].Service)
							deleteTimeoutMessage(chatId, messageId, time.Duration(cfg.Bot.DelMessage.TimeoutAfterSet)*time.Second)
						}
					default:
						msg.Text = types.ErrCom
					}
				} else {
					msg.Text = types.ErrCom
				}
			}
			if _, err = bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if update.CallbackQuery != nil {
			CallbackQueryHandler(update.CallbackQuery)
			continue
		}
	}
}

// Функция формирования разметки для Inline клавиатуры
func ServiceMarkup(currentPage, count int, chatId int64, mode string) (markup tgbotapi.InlineKeyboardMarkup) {

	var s []string
	if currentPage*count+count >= len(services[chatId][currentPage*count:]) {
		s = services[chatId][currentPage*count:]
	} else {
		s = services[chatId][currentPage*count : currentPage*count+count]
	}
	a := cfg.Bot.InlineKeyboard.CountPerRow
	strs := len(s) / a
	if len(s)%cfg.Bot.InlineKeyboard.CountPerRow != 0 {
		strs++
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, strs)
	for i, j := 0, -1; i < len(s); i++ {
		if i%cfg.Bot.InlineKeyboard.CountPerRow == 0 {
			j++
			rows[j] = make([]tgbotapi.InlineKeyboardButton, 0)
		}
		rows[j] = append(rows[j], tgbotapi.NewInlineKeyboardButtonData(s[i], fmt.Sprintf("%s:%s", mode, s[i])))

	}
	maxPages := len(services[chatId]) / count
	if len(services[chatId])%cfg.Bot.InlineKeyboard.CountPerPage != 0 {
		maxPages++
	}
	if currentPage > 0 || currentPage < maxPages-1 {
		rows = append(rows, make([]tgbotapi.InlineKeyboardButton, 0))
	}

	if currentPage > 0 {
		rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData("Предыдущая", fmt.Sprintf("pager:prev:%d:%d:%s", currentPage, count, mode)))
	}
	if currentPage < maxPages-1 {
		rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData("Следующая", fmt.Sprintf("pager:next:%d:%d:%s", currentPage, count, mode)))
	}

	markup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
	return
}

// Функция создания/изменения Inline клавиатуры
func SendServices(chatId int64, currentPage, count int, messageId *int, mode string) {
	keyboard := ServiceMarkup(currentPage, count, chatId, mode)
	var cfg tgbotapi.Chattable
	var text string
	if len(keyboard.InlineKeyboard) == 0 {
		text = types.NoServices
	} else {
		text = types.Services
	}
	if messageId == nil {
		msg := tgbotapi.NewMessage(chatId, text)
		msg.ReplyMarkup = keyboard
		cfg = msg
	} else {
		msg := tgbotapi.NewEditMessageText(chatId, *messageId, text)
		msg.ReplyMarkup = &keyboard
		cfg = msg
	}

	bot.Send(cfg)
}

// Функция обработки Callback запроса
func CallbackQueryHandler(query *tgbotapi.CallbackQuery) {
	split := strings.Split(query.Data, ":")
	switch split[0] {
	case "pager":
		HandleNavigationCallbackQuery(query.Message.MessageID, query.Message.Chat.ID, split[1:]...)
		return
	case "get":
		HandleGetCallbackQuery(query.Message.MessageID, query.Message.Chat.ID, split[1])
		return
	case "del":
		HandleDelCallbackQuery(query.Message.MessageID, query.Message.Chat.ID, split[1])
		return
	}

}

// Функци обработки Callback запроса на получение сервиса
func HandleGetCallbackQuery(messageId int, chatId int64, service string) {
	login, password, err := db.GetServiceByName(chatId, service)
	if err != nil {
		services[chatId] = db.GetServices(chatId)
		SendServices(chatId, 0, cfg.Bot.InlineKeyboard.CountPerPage, &messageId, "get")
	} else {
		msg := tgbotapi.NewMessage(chatId, fmt.Sprintf(types.ServiceInfo, service, login, password))
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		m, err := bot.Send(msg)
		if err != nil {
			logger.ForError(err)
		}
		deleteTimeoutMessage(m.Chat.ID, m.MessageID, time.Duration(cfg.Bot.DelMessage.TimeoutAfterGet)*time.Second)

	}
}

// Функци обработки Callback запроса на удаление сервиса
func HandleDelCallbackQuery(messageId int, chatId int64, service string) {
	db.DelServiceByName(chatId, service)
	services[chatId] = db.GetServices(chatId)
	SendServices(chatId, 0, cfg.Bot.InlineKeyboard.CountPerPage, &messageId, "del")

}

// Функци обработки Callback запроса изменение страницы Inline клавиатуры
func HandleNavigationCallbackQuery(messageId int, chatId int64, data ...string) {
	pagerType := data[0]
	currentPage, _ := strconv.Atoi(data[1])
	itemsPerPage, _ := strconv.Atoi(data[2])
	mode := data[3]
	maxPages := len(services[chatId]) / cfg.Bot.InlineKeyboard.CountPerPage
	if len(services[chatId])%cfg.Bot.InlineKeyboard.CountPerPage != 0 {
		maxPages++
	}
	if pagerType == "next" {
		nextPage := currentPage + 1
		if nextPage < maxPages {
			SendServices(chatId, nextPage, itemsPerPage, &messageId, mode)
		}
	}
	if pagerType == "prev" {
		previousPage := currentPage - 1
		if previousPage >= 0 {
			SendServices(chatId, previousPage, itemsPerPage, &messageId, mode)
		}
	}
}

// Функция для формирования задач для планировщика по удалению сообщений
func deleteTimeoutMessage(chatId int64, messageId int, timeout time.Duration) {
	queueDelMessages.PushBack(types.DelMessage{
		ChatID:    chatId,
		MessageID: messageId,
	})
	_, err = scheduler.Add(&tasks.Task{
		Interval: timeout,
		RunOnce:  true,
		TaskFunc: func() error {
			bot.Send(tgbotapi.NewDeleteMessage(queueDelMessages.Front().Value.(types.DelMessage).ChatID, queueDelMessages.Front().Value.(types.DelMessage).MessageID))
			queueDelMessages.Remove(queueDelMessages.Front())
			return nil
		},
	})
	logger.ForError(err)
}
