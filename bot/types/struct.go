package types

// Переменная состояния для реализации конечного автомата
type State int

// Структура для хранения промежуточной информации о создаваемой записи о сервисе
type UserData struct {
	Service  string
	Login    string
	Password string
	State    State
}

// Список состояний
const (
	InitState = iota
	ServiceSet
	LoginSet
	PassSet
)

// Структура для хранения информации о сообщениях, которые необходимо удалить
type DelMessage struct {
	ChatID    int64
	MessageID int
}
