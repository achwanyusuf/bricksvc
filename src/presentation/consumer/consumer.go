package consumer

import (
	"github.com/achwanyusuf/bricksvc/src/presentation/consumer/payment"
	"github.com/achwanyusuf/bricksvc/src/usecase"
	"github.com/achwanyusuf/bricksvc/utils/kafkalib"
)

type ConsumerHandlerInterface struct {
	payment payment.PaymentInterface
}

type Consumer struct {
	Conf     Conf
	Consumer kafkalib.ConsumerInterface
	Usecase  usecase.UsecaseInterface
}

type Conf struct {
	Payment payment.Conf `mapstructure:"payment"`
}

func (c *Consumer) Serve(cHandler ConsumerHandlerInterface) error {
	c.Consumer.Handler(c.Conf.Payment.TransferTopic, cHandler.payment.Transfer)
	return c.Consumer.Serve()
}

func New(consumer *Consumer) error {
	handlers := ConsumerHandlerInterface{
		payment.New(consumer.Conf.Payment, consumer.Usecase.Transfer),
	}
	return consumer.Serve(handlers)
}
