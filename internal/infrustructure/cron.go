package infrustructure

import (
	"context"
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	*cron.Cron
}

func NewCronJob() CronJob {
	return CronJob{
		Cron: cron.New(),
	}
}

func (cj *CronJob) RunTasks(ctx context.Context, cmd func()) error {
	_, err := cj.AddFunc("@every 1m", cmd)
	if err != nil {
		return err
	}
	cj.Start()
	return nil
}
