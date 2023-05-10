package types

const (
	HelpCom = `Я умею:
1. /get - получить данные для сервиса
2. /set - установить данные для сервиса
3. /delete - удалить данные для сервиса`
	ErrCom = `Я не знаю такой команды 😓
Попробуйте /help!`
	ServiceNameSet    = "Введите название сервиса:"
	LoginNameSet      = "Введите логин:"
	PasswordNameSet   = "Введите пароль:"
	ErrLength         = "Слишком длинное значение"
	AddService        = "Данные к сервису %s добавлены"
	NoServices        = "Данных сервисов нет"
	Services          = "Сервисы:"
	ServiceInfo       = "Данные сервиса `%s`:\nЛогин: `%s`\nПароль: `%s`"
	AuthorizationInfo = "Авторизация в аккаунт %s\n"
	ServiceExists     = "Данные сервиса уже существуют\nВведите другое название:"
)
