package worker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	defaultServerHost     = "localhost"
	defaultServerPort     = "1025"
	defaultMaxConnections = 20
	defaultMessagesPerSec = 100 // NB
	ConfigFileName        = "worker.conf"
)

type Config struct {
	ServerHost     string
	ServerPort     string
	Address        string
	MaxConnections int
	MessagesPerSec int
}

func NewConfig() Config {
	return Config{
		ServerHost:     defaultServerHost,
		ServerPort:     defaultServerPort,
		Address:        defaultServerHost + ":" + defaultServerPort,
		MaxConnections: defaultMaxConnections,
		MessagesPerSec: defaultMessagesPerSec,
	}
}

func (c *Config) LoadFromFile() error {
	file, err := os.Open(ConfigFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Create default config file...")
			return createDefaultConfigFile()
		}
		return fmt.Errorf("open config file error %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "server_host":
			c.ServerHost = value
			log.Println("Load host from config file: ", c.ServerHost)
		case "server_port":
			c.ServerPort = value
			log.Println("Load port from config file: ", c.ServerPort)
		case "max_connections":
			maxCon, err := strconv.Atoi(value)
			if err == nil {
				c.MaxConnections = maxCon
				log.Println("Load max connections from config file: ", c.MaxConnections)
			}
		case "messages_per_sec":
			msgPerSec, err := strconv.Atoi(value)
			if err == nil {
				c.MessagesPerSec = msgPerSec
				log.Println("Load messages per second from config file: ", c.MessagesPerSec)
			}
		}
	}
	c.Address = c.ServerHost + ":" + c.ServerPort
	log.Println("Connection address: ", c.Address)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read config file error: %w", err)
	}
	return nil
}

func createDefaultConfigFile() error {
	file, err := os.Create(ConfigFileName)
	if err != nil {
		return fmt.Errorf("ошибка создания файла конфигурации: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	configTemplate := `# Конфигурационный файл приложения

# Хост сервера
server_host=%s

# Порт сервера
server_port=%s

# Максимальное количество подключений
max_connections=%s

# Количество сообщений в секунду
messages_per_sec=%s
`
	_, err = fmt.Fprintf(writer, configTemplate,
		defaultServerHost,
		defaultServerPort,
		strconv.Itoa(defaultMaxConnections),
		strconv.Itoa(defaultMessagesPerSec))

	if err != nil {
		return fmt.Errorf("ошибка записи дефолтной конфигурации: %w", err)
	}

	writer.Flush()
	fmt.Printf("Создан новый конфигурационный файл: %s\n", ConfigFileName)

	return nil
}
