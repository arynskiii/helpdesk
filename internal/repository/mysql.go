package repository

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	createUsersTableSQL = `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		token VARCHAR(255) DEFAULT 'default_token_value'
	);
	`

	createCategoryTableSQL = `
		CREATE TABLE IF NOT EXISTS category (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL,
		user_id INT NOT NULL,
		create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	createTicketStateSQL = `
	CREATE TABLE IF NOT EXISTS tickets_state (
		id INT AUTO_INCREMENT PRIMARY KEY,
		code VARCHAR(255) CHARACTER SET utf8mb4 NOT NULL,
		name_kz VARCHAR(255) CHARACTER SET utf8mb4 NOT NULL,
		name_ru VARCHAR(255) CHARACTER SET utf8mb4 NOT NULL,
		name_en VARCHAR(255) CHARACTER SET utf8mb4 NOT NULL
	);
	`

	requestOnStateSQL = `
	INSERT INTO tickets_state (code, name_kz, name_ru, name_en)
	VALUES 
	('Request accepted', 'Өтіңіш қабылданды', 'Заявка принята', 'Request accepted'),
	('In process', 'Қарастырылуда', 'В процессе', 'In progress'),
	('Waiting for response','Жауап күтүлуде', 'Ожидание ответа',  'Waiting for response'),
	('Resolved', 'Шешілді', 'Решена', 'Resolved'),
	('Closed', 'Жабылды', 'Закрыта', 'Closed'),
	('On hold','Кейінге қалдырылды',  'Отложено', 'On hold');
	`

	createTicketSQL = `
	CREATE TABLE IF NOT EXISTS ticket (
		id INT AUTO_INCREMENT PRIMARY KEY,
		category_id INT NOT NULL,
		user_id INT NOT NULL,
		title VARCHAR(255),
		description TEXT NOT NULL,
		deadline TIMESTAMP DEFAULT NULL,
		create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES category(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	
	);
	`

	historySQL = `
	CREATE TABLE IF NOT EXISTS history_ticket (
		id INT AUTO_INCREMENT PRIMARY KEY,
		ticket_id INT,
		state_id INT,
		create_at TIMESTAMP,
		deadline TIMESTAMP,
		receiver_id INT,
		sender_id INT,
		FOREIGN KEY (ticket_id) REFERENCES ticket(id),
		FOREIGN KEY (state_id) REFERENCES tickets_state(id),
		FOREIGN KEY (sender_id) REFERENCES users(id)
	);
	`
	createDocumentTableSQL = `
CREATE TABLE IF NOT EXISTS ticket_attachments(
    id INT AUTO_INCREMENT PRIMARY KEY,
    ticket_id INT NOT NULL,
    attachment_path VARCHAR(255) NOT NULL,
	
    FOREIGN KEY (ticket_id) REFERENCES ticket(id)
);
`
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewMySqlDB(cfg *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=false", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open MySQL connection: %v", err)
	}

	if err := db.Ping(); err != nil {

		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}
	return db, nil
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func CreateTable(db *sqlx.DB) error {
	tables := []string{
		createUsersTableSQL,
		createCategoryTableSQL,
		createTicketStateSQL,
		createTicketSQL,
		requestOnStateSQL,
		historySQL,
		createDocumentTableSQL,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return fmt.Errorf("failed to execute table creation SQL: %w", err)
		}
	}
	return nil
}
