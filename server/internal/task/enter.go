package task

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"server/pkg/global"
)

func RegisterScheduledTasks(c *cron.Cron) error {
	if _, err := c.AddFunc("@every 10m", func() {
		if err := UpdateArticleMetricsTask(); err != nil {
			global.Log.Error("Failed to update article metrics:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	if _, err := c.AddFunc("@every 10m", func() {
		if err := GetHotListSyncTask(); err != nil {
			global.Log.Error("Failed to get hot list:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	if _, err := c.AddFunc("@daily", func() {
		if err := GetCalendarSyncTask(); err != nil {
			global.Log.Error("Failed to get calendar:", zap.Error(err))
		}
	}); err != nil {
		return err
	}

	return nil
}
