package transfer

import (
	"context"
	"fmt"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
)

func (t *Transfer) insertKafka(ctx context.Context, data *entity.TransferJob) error {
	key := fmt.Sprintf(model.TransferKey, data.JobID)
	dt, err := data.Payload.MarshalJSON()
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}
	_, err = t.Kafka.Publish(ctx, t.Conf.TransferTopic, []byte(key), dt)
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}
	return nil
}
