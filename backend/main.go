package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"database"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
	} `yaml:"kafka"`
}

type PingResult struct {
	IP            string  `json:"ip"`
	PingTime      float64 `json:"ping_time"`
	LastSuccessAt string  `json:"last_success_at"`
}

var db *sql.DB
var config Config

func loadConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
}

func initDB() {
	psqlInfo := "host=" + config.Database.Host +
		" port=" + strconv.Itoa(config.Database.Port) +
		" user=" + config.Database.User +
		" password=" + config.Database.Password +
		" dbname=" + config.Database.Dbname + " sslmode=disable"
	var err error

	for i := 1; i <= 5; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Attempt %d: Failed to open database: %v", i, err)
			time.Sleep(5 * time.Second)
			continue
		}
		err = db.Ping()
		if err != nil {
			log.Printf("Attempt %d: Database not ready: %v", i, err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("Database connection established on attempt %d", i)
		break
	}
	if err != nil {
		log.Fatal("Could not connect to the database after several attempts:", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS ping_results (ip TEXT PRIMARY KEY, ping_time DOUBLE PRECISION, last_success_at TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}
}

func kafkaConsumer() {
	configSarama := sarama.NewConfig()
	configSarama.Consumer.Return.Errors = true

	var consumer sarama.Consumer
	var err error
	for {
		consumer, err = sarama.NewConsumer(config.Kafka.Brokers, configSarama)
		if err != nil {
			log.Printf("Failed to create Kafka consumer, retrying in 5 seconds: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	defer consumer.Close()

	var partitionConsumer sarama.PartitionConsumer
	for {
		partitionConsumer, err = consumer.ConsumePartition(config.Kafka.Topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Failed to consume partition, retrying in 5 seconds: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		var pr PingResult
		err := json.Unmarshal(msg.Value, &pr)
		if err != nil {
			log.Printf("Failed to unmarshal Kafka message: %v", err)
			continue
		}
		_, err = db.Exec("INSERT INTO ping_results(ip, ping_time, last_success_at) VALUES($1, $2, $3) ON CONFLICT (ip) DO UPDATE SET ping_time = EXCLUDED.ping_time, last_success_at = EXCLUDED.last_success_at", pr.IP, pr.PingTime, pr.LastSuccessAt)
		if err != nil {
			log.Printf("Failed to insert/update database: %v", err)
			continue
		}
	}
}

func getStats(c *gin.Context) {
	rows, err := db.Query("SELECT ip, ping_time, last_success_at FROM ping_results")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var results []PingResult
	for rows.Next() {
		var pr PingResult
		var ts time.Time
		err := rows.Scan(&pr.IP, &pr.PingTime, &ts)
		if err != nil {
			continue
		}
		pr.LastSuccessAt = ts.Format(time.RFC3339)
		results = append(results, pr)
	}
	c.JSON(200, results)
}

func main() {
	loadConfig()
	initDB()
	go kafkaConsumer()

	r := gin.Default()
	r.GET("/stats", getStats)
	err := r.Run(":" + strconv.Itoa(config.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
}
