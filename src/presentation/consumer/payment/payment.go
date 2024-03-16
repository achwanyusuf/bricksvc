package payment

import (
	"context"
	"strings"

	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/usecase/transfer"
	"github.com/achwanyusuf/bricksvc/utils/kafkalib"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	jsoniter "github.com/json-iterator/go"
)

type Payment struct {
	conf     Conf
	transfer transfer.TransferInterface
}

type Conf struct {
	TransferTopic string `mapstructure:"transfer_topic"`
}

type PaymentInterface interface {
	Transfer(req *kafkalib.Request, resp *kafkalib.Response)
}

func New(conf Conf, transfer transfer.TransferInterface) PaymentInterface {
	return &Payment{
		conf:     conf,
		transfer: transfer,
	}
}

func (p *Payment) Transfer(req *kafkalib.Request, resp *kafkalib.Response) {
	var (
		data model.CreateTransfer
		ctx  = context.Background()
	)
	if err := jsoniter.Unmarshal([]byte(req.Value), &data); err == nil {
		key := strings.Split(req.Key, "pubTransfer:")
		if e := p.transfer.Create(ctx, key[1], data); e != nil {
			logger.Log.Error(e)
		}
	} else {
		logger.Log.Error(err)
	}
	resp.Commit()
}
