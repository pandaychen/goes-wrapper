package timewheel

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pandaychen/goes-wrapper/datastructure/timewheel/task"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type TimingWheel struct {
	ctx context.Context

	sync.RWMutex
	taskMap map[interface{}]int

	//为了timingwheel的通用性，这里就不修改taskMap结构了，再使用多一个map来存储密码随机的时间数据（不做删除处理）
	TaskInfoMap sync.Map

	interval time.Duration //currentSlotPos 每隔interval往前Move一个单元

	//采用双链表存储delaytask链，单链加互斥锁
	slotMutex      []*sync.Mutex
	slots          []*list.List
	currentSlotPos int
	slotNum        int

	//for 异步操作
	addDelayTaskChan chan task.DelayTask
	delDelayTaskChan chan interface{}
	stopChan         chan struct{}

	//内置计时器
	ticker      *time.Ticker
	tickerQueue chan time.Time

	loger *zap.Logger
}

// TimingWheel 创建管理结构
func NewTimingWheel(logger *zap.Logger, conf *TimingwheelConfig) *TimingWheel {
	w := &TimingWheel{
		ctx:              context.Background(),
		interval:         conf.Interval,
		slots:            make([]*list.List, conf.SlotNum),
		slotMutex:        make([]*sync.Mutex, conf.SlotNum),
		taskMap:          make(map[interface{}]int),
		currentSlotPos:   0,
		slotNum:          conf.SlotNum,
		addDelayTaskChan: make(chan task.DelayTask, 1024),
		delDelayTaskChan: make(chan interface{}, 128),
		stopChan:         make(chan struct{}),
		ticker:           time.NewTicker(conf.Interval),
		tickerQueue:      make(chan time.Time, 32),
		loger:            logger,
	}

	//初始化各个slots
	for i := 0; i < w.slotNum; i++ {
		w.slots[i] = list.New()
		w.slotMutex[i] = new(sync.Mutex)
	}

	//启动时间轮
	go w.run(w.ctx)

	return w
}

func (tw *TimingWheel) tickGenerator() {
	if tw.tickerQueue == nil {
		return
	}

	for {
		select {
		case <-tw.ticker.C:
			select {
			case tw.tickerQueue <- time.Now():
			default:
				panic("raise long time blocking")
			}
		}
	}
}

// 根据延时duration：获取定时器在slots中的位置, 时间轮需要转动的圈数
func (tw *TimingWheel) getPosAndCircle(dur time.Duration) (int, int, int) {
	var (
		pos, circle int
		remains     int //second
	)
	delaySeconds := int(dur.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	//计算dur里当前位置要转多少圈
	circle = int(delaySeconds / intervalSeconds / tw.slotNum)

	//根据剩余不足一个slot的时延计算该任务需要放入哪个slot链表
	pos = int(tw.currentSlotPos+delaySeconds/intervalSeconds) % tw.slotNum
	//计算剩余不足一个slot的延迟
	remains = delaySeconds - circle*intervalSeconds*tw.slotNum

	tw.loger.Info("getPositionAndCircle", zap.Int("pos", pos), zap.Int("delaySeconds", delaySeconds), zap.Any("duration", dur), zap.Int("remains", remains))

	return pos, circle, remains
}

// 启动
func (tw *TimingWheel) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			tw.ticker.Stop()
			return
		//case <-tw.ticker.C:
		case <-tw.tickerQueue:
			//处理slot指针移动
			tw.tickerHandler()
		case delaytask := <-tw.addDelayTaskChan:
			tw.handlerDelaytaskAddEvent(&delaytask)
		case task_key := <-tw.delDelayTaskChan:
			tw.handlerDelaytaskDelEvent(task_key)
		case <-tw.stopChan:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimingWheel) GetTaskCount() int {
	tw.RLock()
	defer tw.RUnlock()
	return len(tw.taskMap)
}

//检查当前时间轮是否存在task_key任务
func (tw *TimingWheel) CheckJobExists(task_key interface{}) bool {
	tw.RLock()
	defer tw.RUnlock()
	if _, exists := tw.taskMap[task_key]; exists {
		return true
	} else {
		return false
	}
}

// 获取全部的密码随机事件
func (tw *TimingWheel) GetAllExistsEvents() map[string]task.DelayTaskInfo {
	var (
		taskmap  map[string]task.DelayTaskInfo = make(map[string]task.DelayTaskInfo)
		taskKeyI interface{}
		taskkey  string
		ok       bool
	)
	tw.RLock()
	defer tw.RUnlock()
	for taskKeyI, _ = range tw.taskMap {
		if taskkey, ok = taskKeyI.(string); !ok {
			tw.loger.Info("TimingWheel GetAllExistsEvents transform error", zap.Any("taskKeyI", taskKeyI))
			continue
		}
		if nodeinfo, exists := tw.TaskInfoMap.Load(taskKeyI); exists {
			taskmap[taskkey] = nodeinfo.(task.DelayTaskInfo)
		} else {
			tw.loger.Error("TimingWheel GetAllExistsEvents TaskInfoMap error[no data]", zap.Any("taskKeyI", taskKeyI))
		}
	}
	return taskmap
}

// AddDelayTask ：添加不重复的任务
func (tw *TimingWheel) AddDelayTask(delay time.Duration, crontab string, task_key interface{}, task_param interface{}, callback task.DelayTaskJobCallback) error {
	if delay < 0 {
		return errors.New("illegal param")
	}

	select {
	case tw.addDelayTaskChan <- task.DelayTask{
		Delay:      delay,
		Jobfunc:    callback,
		Task_key:   task_key,
		Crontab:    crontab, //用于判断是否时间更新了
		Task_param: task_param}:
	case <-time.After(time.Millisecond * 200):
		return errors.New("AddDelayTask chan full")
	}
	return nil
}

// GetTaskDelayRemain：获取某个指定任务的执行延迟时间
func (tw *TimingWheel) GetTaskDelayRemain(task_key interface{}) (string, time.Duration, error) {
	var (
		exists   bool
		position int
	)
	tw.RLock()
	if position, exists = tw.taskMap[task_key]; !exists {
		tw.RUnlock()
		return "", time.Duration(0), fmt.Errorf("not found task-key:%v", task_key)
	}
	tw.RUnlock()

	l := tw.slots[position]
	tw.slotMutex[position].Lock()
	defer tw.slotMutex[position].Unlock()
	for node := l.Front(); node != nil; {
		task, ok := node.Value.(*task.DelayTask)
		if !ok {
			tw.loger.Error("handlerDelaytaskDelEvent node error", zap.Any("node", task))
			continue
		}

		if task.Task_key == task_key {
			if position >= tw.currentSlotPos {
				return task.Task_param.(string), time.Duration(position-tw.currentSlotPos)*tw.interval + time.Duration(task.Circle*tw.slotNum)*tw.interval, nil
			} else {
				return task.Task_param.(string), time.Duration((task.Circle)*tw.slotNum)*tw.interval + time.Duration(tw.slotNum-(tw.currentSlotPos-position))*tw.interval, nil
			}
		}

		node = node.Next()
	}

	return "", time.Duration(0), fmt.Errorf("not found task-key[unknown error]:%v", task_key)
}

//
func (tw *TimingWheel) DelDelayTask(task_key interface{}) error {
	if task_key == nil {
		return errors.New("illegal param")
	}

	select {
	case tw.delDelayTaskChan <- task_key:
	case <-time.After(time.Millisecond * 200):
		return errors.New("DelDelayTask chan full")
	}

	return nil
}

//	更新延时任务（todo）
func (tw *TimingWheel) UpdateDelayTask(delay time.Duration, crontab string, task_key interface{}, task_param interface{}, callback task.DelayTaskJobCallback) error {
	if task_key == nil {
		return errors.New("illegal param")
	}

	return nil
}

// 每隔ticker时间扫描对应的任务slot
func (tw *TimingWheel) tickerHandler() {
	tw.scaningDelayTasks(tw.slots[tw.currentSlotPos], tw.currentSlotPos)
	if tw.currentSlotPos == tw.slotNum-1 {
		//需要%tw.slotNum
		tw.currentSlotPos = 0
	} else {
		tw.currentSlotPos++
	}
}

// 新增任务到链表中
func (tw *TimingWheel) handlerDelaytaskAddEvent(delaytask *task.DelayTask) {
	var (
		start, end, delta int
	)
	tw.RLock()
	if _, exists := tw.taskMap[delaytask.Task_key]; exists {
		tw.RUnlock()
		tw.loger.Error("handlerDelaytaskAddEvent add repeated tasks", zap.Any("task", *delaytask))
		return
	}
	tw.RUnlock()

	//计算slot和circle、remains
	pos, circle, remains := tw.getPosAndCircle(delaytask.Delay)
	delaytask.Circle = circle

	delaytask.Delta = time.Duration(remains) * time.Second

	//将task放入指定的slot中
	tw.slots[pos].PushBack(delaytask)

	if delaytask.Task_key != nil {
		tw.Lock()
		tw.taskMap[delaytask.Task_key] = pos
		tw.Unlock()
		//save extra data
		start = int(time.Now().Unix())
		delta = int(int64(delaytask.Delay) / 1e9)
		end = start + delta
		tw.TaskInfoMap.Store(delaytask.Task_key, task.DelayTaskInfo{
			StartTimestamp:   start,
			TriggerTimestamp: end,
			Delta:            delta,
		})
	}
}

func (tw *TimingWheel) handlerDelaytaskDelEvent(task_key interface{}) {
	tw.RLock()
	position, ok := tw.taskMap[task_key]
	tw.RUnlock()
	if !ok {
		tw.loger.Error("handlerDelaytaskDelEvent error,task_key not exists", zap.Any("task_key", task_key))
		return
	}
	// 获取槽指向的链表
	l := tw.slots[position]
	tw.slotMutex[position].Lock()
	defer tw.slotMutex[position].Unlock()
	for node := l.Front(); node != nil; {
		task, ok := node.Value.(*task.DelayTask)
		if !ok {
			tw.loger.Error("handlerDelaytaskDelEvent node error", zap.Any("node", task))
			continue
		}

		//从任务map及延时链表移除task_key
		if task.Task_key == task_key {
			l.Remove(node)
			tw.Lock()
			delete(tw.taskMap, task.Task_key)
			tw.Unlock()
		}

		node = node.Next()
	}
}

// 扫描单个slot，普通时间轮存在效率问题
func (tw *TimingWheel) scaningDelayTasks(task_list *list.List, position int) {
	tw.slotMutex[position].Lock()
	defer tw.slotMutex[position].Unlock()
	for node := task_list.Front(); node != nil; {
		task, ok := node.Value.(*task.DelayTask)
		if !ok {
			tw.loger.Error("scaningDelayTasks node error", zap.Any("node", task))
			continue
		}
		if task.Circle > 0 {
			// 非当前轮次，需要减1
			task.Circle--
			// go to next node
			node = node.Next()
			continue
		}

		// task.circle==0，延时到期

		// 异步调用（需要控制并发量）
		go task.Jobfunc(tw.ctx, task.Task_key, task.Task_param, task.Delay)
		next := node.Next()

		//从slot链表移出延时任务，从任务map移出任务
		task_list.Remove(node)
		if task.Task_key != nil {
			tw.Lock()
			delete(tw.taskMap, task.Task_key)
			tw.Unlock()
		}

		node = next
	}
}

// 停止
func (tw *TimingWheel) Stop() {
	tw.stopChan <- struct{}{}
}
