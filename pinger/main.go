package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/IBM/sarama"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Pinger struct {
		Interval int      `yaml:"interval"`
		Targets  []string `yaml:"targets"`
	} `yaml:"pinger"`
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

var config Config

func loadConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}
}

func checksum(data []byte) uint16 {
	var sum int32
	for i := 0; i < len(data)-1; i += 2 {
		sum += int32(binary.BigEndian.Uint16(data[i:]))
	}
	if len(data)%2 == 1 {
		sum += int32(data[len(data)-1]) << 8
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return uint16(^sum)
}

func pingTarget(target string) (float64, error) {
	const timeout = 5 * time.Second
	start := time.Now()

	conn, err := net.DialTimeout("ip4:icmp", target, timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	packet := make([]byte, 8+32)
	packet[0] = 8
	packet[1] = 0

	pid := uint16(os.Getpid() & 0xffff)
	binary.BigEndian.PutUint16(packet[4:6], pid)
	binary.BigEndian.PutUint16(packet[6:8], 1)

	cs := checksum(packet)
	binary.BigEndian.PutUint16(packet[2:4], cs)

	if _, err := conn.Write(packet); err != nil {
		return 0, err
	}

	conn.SetDeadline(time.Now().Add(timeout))
	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		return 0, err
	}

	var icmpReply []byte
	if n >= 28 {
		icmpReply = reply[20:28]
	} else {
		icmpReply = reply[:8]
	}

	if icmpReply[0] != 0 {
		return 0, fmt.Errorf("invalid ICMP response type: %d", icmpReply[0])
	}

	duration := time.Since(start)
	return float64(duration.Milliseconds()), nil
}

func main() {
	loadConfig()

	configSarama := sarama.NewConfig()
	configSarama.Producer.Return.Successes = true

	var producer sarama.SyncProducer
	var err error
	for {
		producer, err = sarama.NewSyncProducer(config.Kafka.Brokers, configSarama)
		if err != nil {
			log.Printf("Failed to create producer, retrying in 5 seconds: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	defer producer.Close()

	ticker := time.NewTicker(time.Duration(config.Pinger.Interval) * time.Second)
	for {
		for _, target := range config.Pinger.Targets {
			latency, err := pingTarget(target)
			if err != nil {
				log.Printf("Ping %s failed: %v", target, err)
				continue
			}
			pr := PingResult{
				IP:            target,
				PingTime:      latency,
				LastSuccessAt: time.Now().Format(time.RFC3339),
			}
			data, err := json.Marshal(pr)
			if err != nil {
				log.Printf("Failed to marshal result: %v", err)
				continue
			}
			msg := &sarama.ProducerMessage{
				Topic: config.Kafka.Topic,
				Value: sarama.ByteEncoder(data),
			}
			_, _, err = producer.SendMessage(msg)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}
		<-ticker.C
	}
}
