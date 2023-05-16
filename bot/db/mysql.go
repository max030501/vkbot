package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"vkbot/logger"
	"vkbot/types"
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
func GetServices(chatId int64) []*types.Service {
	db := mysqlConn()
	rows, err := db.Query("SELECT service,id from services  where chat_id = ? ", chatId)
	logger.ForError(err)
	var services []*types.Service
	for rows.Next() {
		var service types.Service
		err = rows.Scan(&service.Service, &service.Id)
		logger.ForError(err)
		services = append(services, &service)
	}
	defer db.Close()
	return services
}

// Получение сервиса по имени
func GetServiceById(id int) (string, string, error) {
	db := mysqlConn()
	defer db.Close()
	rows, err := db.Query("SELECT login,password from services  where id = ? ", id)
	logger.ForError(err)
	var login, password string
	rows.Next()
	err = rows.Scan(&login, &password)
	if err != nil {
		return "", "", err
	}
	return login, password, nil
}

func GetServiceByName(name string, chatId int64) (string, string, error) {
	db := mysqlConn()
	defer db.Close()
	rows, err := db.Query("SELECT login,password from services  where service = ? and chat_id = ? ", name, chatId)
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
func DelServiceById(id int) error {
	db := mysqlConn()
	defer db.Close()
	_, err := db.Exec("DELETE from services  where id = ? ", id)
	if err != nil {
		return err
	}
	return nil

}
