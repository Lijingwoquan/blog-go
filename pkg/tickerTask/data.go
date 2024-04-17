package ticker

import (
	"blog/logic"
	"go.uber.org/zap"
	"time"
)

func updateDataAboutIndex() error {
	if err := update(); err != nil {
		return err
	}
	ticker := time.NewTicker(time.Hour * 1)
	defer ticker.Stop()
	for range ticker.C {
		//更新DataAboutIndex
		if err := logic.ResponseDataAboutIndex(&GlobalDataAboutIndex); err != nil {
			zap.L().Error("logic.ResponseDataAboutIndex(&DataAboutIndex) failed,err:%v", zap.Error(err))
			return err
		}
	}
	return nil
}

// 更新数据
func update() error {
	mu.Lock()
	defer mu.Unlock()

	if err := logic.ResponseDataAboutIndex(&GlobalDataAboutIndex); err != nil {
		return err
	}
	return nil
}
