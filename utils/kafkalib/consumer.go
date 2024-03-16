package kafkalib

import (
	"fmt"
	"sync"
	"time"

	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Conf struct {
	Enabled                          bool          `mapstructure:"enabled"`
	BootstrapServers                 string        `mapstructure:"bootstrap_servers"`
	GroupID                          string        `mapstructure:"group_id"`
	AutoOffsetReset                  string        `mapstructure:"auto_offset_reset"`
	SecurityProtocol                 string        `mapstructure:"security_protocol"`
	BrokerAddressFamily              string        `mapstructure:"broker_address_family"`
	BrokerAddressTtl                 time.Duration `mapstructure:"broker_address_ttl"`
	ReconnectBackoff                 time.Duration `mapstructure:"reconnect_backoff"`
	ReconnectBackoffMax              time.Duration `mapstructure:"reconnect_backoff_max"`
	FetchMaxBytes                    int           `mapstructure:"fetch_max_bytes"`
	CheckCRCS                        bool          `mapstructure:"check_crcs"`
	SessionTimeout                   time.Duration `mapstructure:"session_timeout"`
	HeartbeatInterval                time.Duration `mapstructure:"heartbeat_interval"`
	GoApplicationRebalanceEnable     bool          `mapstructure:"go_application_rebalance_enable"`
	EnableAutoCommit                 bool          `mapstructure:"enable_auto_commit"`
	FetchErrorBackoff                time.Duration `mapstructure:"fetch_error_backoff"`
	AutoCommitInterval               time.Duration `mapstructure:"auto_commit_interval"`
	EnablePartitionEof               bool          `mapstructure:"enable_partition_eof"`
	EnableAutoOffsetStore            bool          `mapstructure:"enable_auto_offset_store"`
	MaxPollInterval                  time.Duration `mapstructure:"max_poll_interval"`
	MessageMaxBytes                  int           `mapstructure:"message_max_bytes"`
	MessageCopyMaxBytes              int           `mapstructure:"message_copy_max_bytes"`
	ReceiveMessageMaxBytes           int           `mapstructure:"receive_message_max_bytes"`
	MaxInFlightRequestsPerConnection int           `mapstructure:"max_in_flight_requests_per_connection"`
	TopicMetadataRefreshInterval     time.Duration `mapstructure:"topic_metadata_refresh_interval"`
	MetadataMaxAge                   time.Duration `mapstructure:"metadata_max_age"`
	TopicMetadataRefreshFastInterval time.Duration `mapstructure:"topic_metadata_refresh_fast_interval"`
	TopicMetadataRefreshSparse       bool          `mapstructure:"topic_metadata_refresh_sparse"`
	SocketTimeout                    time.Duration `mapstructure:"socket_timeout"`
	SocketKeepaliveEnable            bool          `mapstructure:"socket_keepalive_enable"`
	PartitionAssignmentStrategy      string        `mapstructure:"partition_assignment_strategy"`
	CoordinatorQueryInterval         time.Duration `mapstructure:"coordinator_query_interval"`
	SASLOption                       SASLOption    `mapstructure:"sasl_options"`
}

type SASLOption struct {
	Enabled   bool   `mapstructure:"enabled"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Mechanism string `mapstructure:"mechanism"`
}

var defaultConf = kafka.ConfigMap{
	"bootstrap.servers":                       "localhost:9092",
	"group.id":                                "brick_local",
	"auto.offset.reset":                       "earliest",
	"security.protocol":                       "PLAINTEXT",
	"broker.address.family":                   "v4",
	"broker.address.ttl":                      60 * time.Second,
	"reconnect.backoff.ms":                    200 * time.Millisecond,
	"reconnect.backoff.max.ms":                6 * time.Second,
	"fetch.max.bytes":                         "1024000",
	"check.crcs":                              true,
	"session.timeout.ms":                      60 * time.Second,
	"heartbeat.interval.ms":                   10 * time.Second,
	"go.application.rebalance.enable":         true,
	"enable.auto.commit":                      false,
	"fetch.error.backoff.ms":                  200 * time.Millisecond,
	"auto.commit.interval.ms":                 6 * time.Second,
	"enable.partition.eof":                    false,
	"enable.auto.offset.store":                false,
	"max.poll.interval.ms":                    600 * time.Second,
	"message.max.bytes":                       1024000,
	"message.copy.max.bytes":                  1024000,
	"receive.message.max.bytes":               2048000,
	"max.in.flight.requests.per.connection":   10000,
	"topic.metadata.refresh.interval.ms":      60 * time.Second,
	"metadata.max.age.ms":                     6 * time.Second,
	"topic.metadata.refresh.fast.interval.ms": 250 * time.Millisecond,
	"topic.metadata.refresh.sparse":           true,
	"socket.timeout.ms":                       20 * time.Second,
	"socket.keepalive.enable":                 true,
	"partition.assignment.strategy":           "roundrobin",
	"coordinator.query.interval.ms":           2 * time.Second,
}

type ConsumerInterface interface {
	Serve() error
	Handler(topic string, fn HandlerFunc)
	Stop()
}

type consumer struct {
	mu         *sync.Mutex
	consumer   *kafka.Consumer
	signal     chan struct{}
	topics     []string
	handlers   map[string]HandlerFunc
	autoCommit bool
}

type HandlerFunc func(req *Request, resp *Response)

func (f HandlerFunc) Set(req *Request, resp *Response) {
	f(req, resp)
}

func (c *consumer) Handler(topic string, fn HandlerFunc) {
	c.topics = append(c.topics, topic)
	c.handlers[topic] = fn
}

func NewConsumer(conf *Conf) (ConsumerInterface, error) {
	if !conf.Enabled {
		return &consumer{}, nil
	}
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
			"session.timeout.ms":                      int(conf.SessionTimeout / time.Millisecond),
			"heartbeat.interval.ms":                   int(conf.HeartbeatInterval / time.Millisecond),
			"go.application.rebalance.enable":         conf.GoApplicationRebalanceEnable,
			"enable.auto.commit":                      conf.EnableAutoCommit,
			"fetch.error.backoff.ms":                  int(conf.FetchErrorBackoff / time.Millisecond),
			"auto.commit.interval.ms":                 int(conf.AutoCommitInterval / time.Millisecond),
			"enable.partition.eof":                    conf.EnablePartitionEof,
			"enable.auto.offset.store":                conf.EnableAutoOffsetStore,
			"max.poll.interval.ms":                    int(conf.MaxPollInterval / time.Millisecond),
			"message.max.bytes":                       conf.MessageMaxBytes,
			"message.copy.max.bytes":                  conf.MessageCopyMaxBytes,
			"receive.message.max.bytes":               conf.ReceiveMessageMaxBytes,
			"max.in.flight.requests.per.connection":   conf.MaxInFlightRequestsPerConnection,
			"topic.metadata.refresh.interval.ms":      int(conf.TopicMetadataRefreshInterval / time.Millisecond),
			"metadata.max.age.ms":                     int(conf.MetadataMaxAge / time.Millisecond),
			"topic.metadata.refresh.fast.interval.ms": int(conf.TopicMetadataRefreshFastInterval / time.Millisecond),
			"topic.metadata.refresh.sparse":           conf.TopicMetadataRefreshSparse,
			"socket.timeout.ms":                       int(conf.SocketTimeout / time.Millisecond),
			"socket.keepalive.enable":                 conf.SocketKeepaliveEnable,
			"partition.assignment.strategy":           conf.PartitionAssignmentStrategy,
			"coordinator.query.interval.ms":           int(conf.CoordinatorQueryInterval / time.Millisecond),
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
	c, err := kafka.NewConsumer(kafkaConf)
	if err != nil {
		return &consumer{}, errormsg.WrapErr(errormsg.Error500, err)
	}
	return &consumer{
		mu:         new(sync.Mutex),
		consumer:   c,
		signal:     make(chan struct{}, 1),
		handlers:   make(map[string]HandlerFunc),
		autoCommit: conf.EnableAutoCommit,
	}, nil
}

func (c *consumer) Serve() error {
	if len(c.topics) < 1 {
		return errormsg.WrapErr(errormsg.Error500, nil, "no active topics, disable kafka!")
	}

	if err := c.consumer.SubscribeTopics(c.topics, nil); err != nil {
		return errormsg.WrapErr(errormsg.Error500, err)
	}

	for {
		select {
		case <-c.signal:
			return errormsg.WrapErr(errormsg.Error500, nil, "Kafka process is shutting down")
		default:
			if event := c.consumer.Poll(1000); event != nil {
				switch ev := event.(type) {
				case *kafka.Message:
					topic := ev.TopicPartition.Topic
					if topic == nil {
						continue
					}
					handler := c.handlers[*topic]
					done := make(chan struct{}, 1)
					go func(handler HandlerFunc, msg kafka.Message, done chan struct{}) {
						req := &Request{
							Topic:         *msg.TopicPartition.Topic,
							Partition:     int(msg.TopicPartition.Partition),
							Offset:        int(msg.TopicPartition.Offset),
							Key:           string(msg.Key),
							Value:         string(msg.Value),
							Timestamp:     msg.Timestamp,
							TimestampType: msg.TimestampType.String(),
						}
						req.setHeader(msg.Headers)
						resp := &Response{
							commit: false,
						}
						handler.Set(req, resp)

						if resp.commit || c.autoCommit {
							msg.TopicPartition.Offset += 1
							if _, err := c.consumer.CommitOffsets([]kafka.TopicPartition{msg.TopicPartition}); err != nil {
								newErr := errormsg.WrapErr(errormsg.Error400, err)
								logger.Log.Error(errormsg.WriteErr(newErr))
							}
						}
						done <- struct{}{}
					}(handler, *ev, done)
				case kafka.OffsetsCommitted:
					logger.Log.Info("Kafka commited ", ev.Offsets)
				case kafka.Error:
					err := fmt.Errorf(ev.String())
					newErr := errormsg.WrapErr(errormsg.Error400, err)
					logger.Log.Error(errormsg.WriteErr(newErr))
				default:
					continue
				}
			}
		}
	}
}

func (c *consumer) Stop() {
	c.signal <- struct{}{}
	c.consumer.Close()
}
