package messaging

import (
	"fmt"

	"github.com/IBM/sarama"
)

var (
	syncProducer sarama.SyncProducer
	consumer     sarama.Consumer
)

// Init initialises the Kafka client.  Connection parameters are read from
// PAO-injected environment variables (camelCase secret keys → UPPER_SNAKE):
//
//	MESSAGING_BROKERS        — comma-separated broker list (e.g. broker:9092)
//	MESSAGING_USERNAME       — SASL username (empty → no auth)
//	MESSAGING_PASSWORD       — SASL password
//	MESSAGING_SASL_MECHANISM — PLAIN or SCRAM-SHA-256 or SCRAM-SHA-512 (empty → PLAIN)
func Init(brokers, username, password, saslMechanism string) error {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	if username != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = username
		cfg.Net.SASL.Password = password
		cfg.Net.TLS.Enable = true

		switch saslMechanism {
		case "SCRAM-SHA-256":
			cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
			cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &XDGSCRAMClient{HashGeneratorFcn: SHA256}
			}
		case "SCRAM-SHA-512":
			cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
			cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
			}
		default:
			cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		}
	}

	brokerList := []string{brokers}
	if len(brokers) == 0 {
		brokerList = []string{"localhost:9092"}
	}

	p, err := sarama.NewSyncProducer(brokerList, cfg)
	if err != nil {
		return fmt.Errorf("messaging: producer: %w", err)
	}
	syncProducer = p
	return nil
}

func Close() {
	if syncProducer != nil {
		syncProducer.Close()
	}
	if consumer != nil {
		consumer.Close()
	}
}

// Producer returns the Kafka sync producer.
func Producer() sarama.SyncProducer {
	return syncProducer
}
