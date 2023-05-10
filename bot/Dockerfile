FROM golang:latest

# Создадим директорию
RUN mkdir /VKBot


# Скопируем всё в директорию VKBot
ADD ./ /VKBot/

# Установим рабочей папкой директорию VKBot
WORKDIR /VKBot

# Получим зависимости, которые использовали в боте
RUN go mod download
RUN go mod tidy


# Соберём приложение
RUN go build -o main .

# Запустим приложение
CMD ["/VKBot/main"]