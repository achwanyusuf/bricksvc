package kafkalib

import (
	"context"

	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerInterface interface {
	Publish(ctx context.Context, topic string, key []byte, msg []byte) (int, error)
}

type producer struct {
	producer *kafka.Producer
	signal   chan struct{}
}

func NewProducer(conf *Conf) (ProducerInterface, error) {
	kafkaConf := &defaultConf
	if conf != nil {
		kafkaConf = &kafka.ConfigMap{
			"bootstrap.servers":                       conf.BootstrapServers,
			"group.id":                                conf.GroupID,
			"auto.offset.reset":                       conf.AutoOffsetReset,
			"security.protocol":                       conf.SecurityProtocol,
			"broker.address.family":                   conf.BrokerAddressFamily,
			"broker.address.ttl":                      conf.BrokerAddressTtl,
			"fetch.max.bytes":                         conf.FetchMaxBytes,
			"check.crcs":                              conf.CheckCRCS,
			"session.timeout.ms":                      conf.SessionTimeout,
			"heartbeat.interval.ms":                   conf.HeartbeatInterval,
			"enable.auto.commit":                      conf.EnableAutoCommit,
			"fetch.error.backoff.ms":                  conf.FetchErrorBackoff,
			"auto.commit.interval.ms":                 conf.AutoCommitInterval,
			"enable.partition.eof":                    conf.EnablePartitionEof,
			"enable.auto.offset.store":                conf.EnableAutoOffsetStore,
			"max.poll.interval.ms":                    conf.MaxPollInterval,
			"message.max.bytes":                       conf.MessageMaxBytes,
			"message.copy.max.bytes":                  conf.MessageCopyMaxBytes,
			"receive.message.max.bytes":               conf.ReceiveMessageMaxBytes,
			"max.in.flight.requests.per.connection":   conf.MaxInFlightRequestsPerConnection,
			"topic.metadata.refresh.interval.ms":      conf.TopicMetadataRefreshInterval,
			"metadata.max.age.ms":                     conf.MetadataMaxAge,
			"topic.metadata.refresh.fast.interval.ms": conf.TopicMetadataRefreshFastInterval,
			"topic.metadata.refresh.sparse":           conf.TopicMetadataRefreshSparse,
			"socket.timeout.ms":                       conf.SocketTimeout,
			"socket.keepalive.enable":                 conf.SocketKeepaliveEnable,
			"partition.assignment.strategy":           conf.PartitionAssignmentStrategy,
			"coordinator.query.interval.ms":           conf.CoordinatorQueryInterval,
		}

		if conf.SASLOption.Enabled {
			err := kafkaConf.SetKey("sasl.mechanisms", conf.SASLOption.Mechanism)
			if err != nil {
				logger.Log.Error(err)
			}
			err = kafkaConf.SetKey("sasl.username", conf.SASLOption.Username)
			if err != nil {
				logger.Log.Error(err)
			}
			err = kafkaConf.SetKey("sasl.password", conf.SASLOption.Password)
			if err != nil {
				logger.Log.Error(err)
			}
		}
	}
	p, err := kafka.NewProducer(kafkaConf)
	if err != nil {
		return &producer{}, errormsg.WrapErr(errormsg.Error500, err)
	}
	return &producer{
		producer: p,
		signal:   make(chan struct{}, 1),
	}, nil
}

func (p *producer) Publish(ctx context.Context, topic string, key []byte, msg []byte) (int, error) {
	deliveryChan := make(chan kafka.Event)
	err := p.producer.Produce(&kafka.Message{
		Key:   key,
		Value: msg,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}, deliveryChan)
	if err != nil {
		err = errormsg.WrapErr(errormsg.Error400, nil, "error kafka producer")
		return -1, err
	}

	e := <-deliveryChan
	close(deliveryChan)
	m, ok := e.(*kafka.Message)
	if ok && m != nil {
		if m.TopicPartition.Error != nil {
			err = errormsg.WrapErr(errormsg.Error400, nil, "error kafka producer")
			return -1, err
		}
		return int(m.TopicPartition.Partition), nil
	}

	// return error on nil kafka.Message
	return -1, errormsg.WrapErr(errormsg.Error400, nil, "error kafka producer")
}
