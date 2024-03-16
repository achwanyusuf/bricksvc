package scheduler

import (
	"github.com/achwanyusuf/bricksvc/src/presentation/scheduler/transfer"
	"github.com/achwanyusuf/bricksvc/src/usecase"
	"github.com/achwanyusuf/bricksvc/utils/schedulerengine"
)

type SchedulerHandlerInterface struct {
	Transfer transfer.TransferInterface
}

type Scheduler struct {
	Conf      Conf
	Scheduler schedulerengine.SchedulerEngineInterface
	Usecase   usecase.UsecaseInterface
}

type Conf struct {
	Transfer transfer.Conf `mapstructure:"transfer"`
}

func (s *Scheduler) Serve(sHandler SchedulerHandlerInterface) {
	s.Scheduler.Schedule(s.Conf.Transfer.GetTransferCallback.Name, s.Conf.Transfer.GetTransferCallback.Interval, sHandler.Transfer.Transfer)
}

func New(scheduler *Scheduler) {
	handlers := SchedulerHandlerInterface{
		Transfer: transfer.New(scheduler.Conf.Transfer, scheduler.Usecase.Transfer),
	}
	scheduler.Serve(handlers)
}
