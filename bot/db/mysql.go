package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"vkbot/logger"
)

// Создание подключения к БД
func mysqlConn() *sql.DB {
	db, err := sql.Open("mysql", "root:12345@tcp(db:3306)/storepass")
	logger.ForError(err)
	return db
}

// Вставка в таблицу информации о данных нового сервиса
func InsertServiceInfo(chatId int64, service, login, password string) {
	db := mysqlConn()
	res, err := db.Prepare("INSERT INTO services (chat_id,service,login,password) VALUES (?,?,?,?)")
	logger.ForError(err)
	_, err = res.Exec(chatId, service, login, password)
	logger.ForError(err)
	defer db.Close()
}

// Получение списка всех сервисов
func GetServices(chatId int64) []string {
	db := mysqlConn()
	rows, err := db.Query("SELECT service from services  where chat_id = ? ", chatId)
	logger.ForError(err)
	services := make([]string, 0)
	for rows.Next() {
		var service string
		err = rows.Scan(&service)
		logger.ForError(err)
		services = append(services, service)
	}
	defer db.Close()
	return services
}

// Получение сервиса по имени
func GetServiceByName(chatId int64, service string) (string, string, error) {
	db := mysqlConn()
	defer db.Close()
	rows, err := db.Query("SELECT login,password from services  where chat_id = ? and service = ? ", chatId, service)
	logger.ForError(err)
	var login, password string
	rows.Next()
	err = rows.Scan(&login, &password)
	if err != nil {
		return "", "", err
	}
	return login, password, nil
}

// Удаление сервиса по имени
func DelServiceByName(chatId int64, service string) error {
	db := mysqlConn()
	defer db.Close()
	_, err := db.Exec("DELETE from services  where chat_id = ? and service = ? ", chatId, service)
	if err != nil {
		return err
	}
	return nil

}
