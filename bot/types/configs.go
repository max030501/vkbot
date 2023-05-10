package types

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// Структура конфига
type Config struct {
	Bot struct {
		BotToken      string `yaml:"botToken"`
		UpdateTimeout int    `yaml:"updateTimeout"`
		DelMessage    struct {
			TimeoutAfterSet int `yaml:"timeoutAfterSet"`
			TimeoutAfterGet int `yaml:"timeoutAfterGet"`
		} `yaml:"delMessage"`
		SetService struct {
			LenService  int `yaml:"lenService"`
			LenLogin    int `yaml:"lenLogin"`
			LenPassword int `yaml:"lenPassword"`
		} `yaml:"setService"`
		InlineKeyboard struct {
			CountPerPage int `yaml:"countPerPage"`
			CountPerRow  int `yaml:"countPerRow"`
		} `yaml:"inlineKeyboard"`
	} `yaml:"bot"`
}

// Создание нового экземпляра конфига
func NewConfig(configPath string) (*Config, error) {

	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// Валидация путик к конфигурационному файлу
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' - директория'", path)
	}
	return nil
}

// Парсинг флагов
func ParseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./types/config.yml", "путь к конфигурационному файлу")
	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}
	return configPath, nil
}
