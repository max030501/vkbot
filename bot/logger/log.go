package logger

import (
	"log"
	"os"
)

var (
	out, _ = os.OpenFile("/var/log/VKbot/VKbot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	Logs   = log.New(out, "", log.Ldate|log.Ltime)
)

// Функция логирования ошибок
func ForError(e error) {
	if e != nil {
		Logs.Fatalln(e)
	}
}
