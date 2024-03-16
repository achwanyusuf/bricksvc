package kafkalib

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Request struct {
	Header        map[string]string
	Topic         string
	Partition     int
	Offset        int
	Key           string
	Value         string
	Timestamp     time.Time
	TimestampType string
}

func (r *Request) setHeader(kh []kafka.Header) map[string]string {
	h := make(map[string]string)
	for _, v := range kh {
		h[v.Key] = string(v.Value)
	}
	return h
}
