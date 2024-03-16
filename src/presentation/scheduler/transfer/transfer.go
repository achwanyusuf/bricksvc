package transfer

import (
	"context"
	"time"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/usecase/transfer"
	"github.com/volatiletech/null/v8"
)

type Transfer struct {
	conf     Conf
	transfer transfer.TransferInterface
}

type Conf struct {
	GetTransferCallback GetTransferCallback `mapstructure:"get_transfer_callback"`
}

type GetTransferCallback struct {
	Name     string        `mapstructure:"name"`
	Interval time.Duration `mapstructure:"interval"`
	Limit    int64         `mapstructure:"limit"`
}

type TransferInterface interface {
	Transfer()
}

func New(conf Conf, transfer transfer.TransferInterface) TransferInterface {
	return &Transfer{
		conf:     conf,
		transfer: transfer,
	}
}

func (t *Transfer) Transfer() {
	tNow := time.Now()
	t.transfer.ProccessGetCallback(context.Background(), &model.GetTransferJobsByParam{
		GetTransferJobByParam: model.GetTransferJobByParam{
			Status:      null.StringFrom(entity.TransferstatusPending.String()),
			CreatedAtGT: null.TimeFrom(time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 0, 0, 0, 0, time.UTC)),
			CreatedAtLT: null.TimeFrom(time.Date(tNow.Year(), tNow.Month(), tNow.Day()+1, 0, 0, 0, 0, time.UTC)),
		},
		Limit: t.conf.GetTransferCallback.Limit,
	})
}
