package task

import (
	"context"
	"time"
)

// 延时任务回调函数
type DelayTaskJobCallback func(ctx context.Context, task_key, task_param interface{}, ttl time.Duration) interface{}

var DefaultDelayTaskJobFunc = func(ctx context.Context, task_key, task_param interface{}, ttl time.Duration) interface{} {
	return nil
}

// 时间轮单个事件
type DelayTaskInfo struct {
	StartTimestamp   int `json:"start_time"`   //注册随机开始时间（unix）
	TriggerTimestamp int `json:"trigger_time"` //注册随机预计触发时间
	Delta            int `json:"delta"`        //随机触发间隔（秒）
}

// DelayTask
type DelayTask struct {
	Delay      time.Duration // 延迟周期
	Delta      time.Duration //剩余不足一个slot周期的时间
	Circle     int           // 时间轮转动的圈数[由delay计算]
	Task_key   interface{}   // 唯一标识，用于去重/更新/删除
	Task_param interface{}   // 回调参数
	Crontab    string        //crontab格式
	Jobfunc    DelayTaskJobCallback
}

func (t *DelayTask) SetTaskJobFunc(job DelayTaskJobCallback) {
	t.Jobfunc = job
}

func (t DelayTask) TaskKey() interface{} {
	return t.Task_key
}
