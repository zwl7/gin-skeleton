package utils

/**
 * @desc		基础工具包：异步任务管理器（叠加器）
 * ---------------------------------------------------------------------
 * @author		unphp <unphp@qq.com>
 * @date		2019-08-16
 * @copyright	PPOSUtils 0.1
 * ---------------------------------------------------------------------
 */

import (
	"log"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

func NewTimerCounter() *TimerCounter {
	return &TimerCounter{
		counter: sync.Map{},
	}
}

type TimerCounter struct {
	counter sync.Map
}

func (that *TimerCounter) GetTime(key string, maxLimit int64) int64 {
	_v, _ok := that.counter.Load(key)
	_cleanTime := int64(0)
	if !_ok {
		that.counter.Store(key, ToString(time.Now().Unix(), true)+"|1")
		_v, _ok = that.counter.Load(key)
	}
	_slice := strings.Split(ToString(_v), "|")
	_i := ToInt64(_slice[1]) + 1
	if _i > maxLimit {
		_cleanTime = ToInt64(_slice[0])
		//重新归零计算
		_slice[0] = ToString(time.Now().Unix(), true)
		_slice[1] = "1"
		that.counter.Store(key, strings.Join(_slice, "|"))
	} else {
		_slice[1] = ToString(_i)
		that.counter.Store(key, strings.Join(_slice, "|"))
	}
	return _cleanTime
}

// NewStackChan ...
func NewStackChan() StackChan {
	return make(chan []interface{}, 1)
}

// StackChan 无限叠加器：实现多协程之间无阻塞写入
type StackChan chan []interface{}

// Stack 实现无阻塞的信道“叠加器”方法
func (that StackChan) Stack(value ...interface{}) {
	newdata := make([]interface{}, 0)
	newdata = append(newdata, value...)
	for {
		select {
		case that <- newdata:
			return
		case old := <-that:
			old = append(old, newdata...)
			newdata = old
		}
	}
}

func (that StackChan) Get() []interface{} {
	//_data := []interface{}{}
	//_ticker := time.NewTicker(1 * time.Second)
	//defer _ticker.Stop()
	//select {
	//case _data = <-that:
	//default:
	//}
	//return _data
	return <-that
}

type StackChanInterface interface {
	Stack(value ...interface{})
	Get() []interface{}
}

type StackChanTask interface {
	TaskHandler()
}

// taskData 任务结构
type taskData struct {
	index    int         //任务队列的“倒排索引值”
	runTime  int         //任务执行的时间
	taskData interface{} //任务

}

// TaskerBuilder 任务管理器接口方法
type TaskerBuilder interface {
	AddTask(runTime int, task interface{}) error //
}

// NewTasker 创建一个"异步任务管理器"
func NewTasker(poolMax ...int) TaskerBuilder {
	_tasker := new(Tasker)
	go _tasker.Run(poolMax...)
	return _tasker
}

// Tasker 异步任务管理器
type Tasker struct {
	//taskTransferStack      StackChan //叠加器
	taskTransferStackSlice []StackChan
	taskTransferPool       chan int //
	taskTransferPoolMaxNum int      //

	taskWorkStack   StackChan //主程序 任务调度器
	taskWorkPool    chan int  //
	taskWorkPoolMax int       //

	lock       sync.Mutex //锁
	sortSlice  []int      //倒排索引
	sortLength int
	sortSign   chan int
	taskMap    map[int]interface{} //任务池

	runTime int64 //
}

// AddTask 添加任务
func (that *Tasker) AddTask(runTime int, task interface{}) error {
	if runTime < 86400*365 {
		runTime = ToInt(time.Now().Unix()) + runTime
	}
	that._stack(runTime, task)
	return nil
}

// Clean 清除叠加器与队列积压的任务，重新初始化
func (that *Tasker) Clean() {
	_i := 0
	_j := len(that.taskTransferStackSlice) + 5
	for {
		if _i > _j {
			break
		}
		for _, _stack := range that.taskTransferStackSlice {
			select {
			case <-_stack:
			default:
			}
		}
		select {
		case <-that.taskWorkStack:
		default:
			that.sortSlice = []int{}
			that.taskMap = map[int]interface{}{}
			_i++
		}
	}
}

// Run 任务主进程
func (that *Tasker) Run(poolMax ...int) {
	_taskTransferLen := 2
	that.taskWorkPoolMax = 4
	that.taskTransferPoolMaxNum = 2
	_paramsNum := len(poolMax)
	switch _paramsNum {
	case 1:
		that.taskWorkPoolMax = poolMax[0]
	case 2:
		that.taskWorkPoolMax = poolMax[0]
		that.taskTransferPoolMaxNum = poolMax[1]
		_taskTransferLen = poolMax[1]
	case 3:
		that.taskWorkPoolMax = poolMax[0]
		that.taskTransferPoolMaxNum = poolMax[1]
		_taskTransferLen = int(poolMax[2])
	}

	if that.taskWorkPoolMax > 20 {
		that.taskTransferPoolMaxNum = 20
	}
	if that.taskTransferPoolMaxNum > 5 {
		that.taskTransferPoolMaxNum = 5
	}
	if _taskTransferLen > 5 {
		_taskTransferLen = 5
	}

	if _paramsNum > 2 {
		log.Println("task run params error!")
	}
	that.sortSlice = make([]int, 0)
	that.taskMap = make(map[int]interface{})
	//
	that.taskTransferStackSlice = []StackChan{}
	for _i := 0; _i < _taskTransferLen; _i++ {
		that.taskTransferStackSlice = append(that.taskTransferStackSlice, NewStackChan())
	}
	//
	that.taskWorkStack = NewStackChan()
	//
	that.taskTransferPool = make(chan int, that.taskTransferPoolMaxNum)
	that.taskWorkPool = make(chan int, that.taskWorkPoolMax)
	that.sortSign = make(chan int, 1)

	//合并压入的任务缓存，进入索引
	go func() {
		for {
			select {
			case that.taskTransferPool <- 0:
				//异步
				_slice := that._rangeTransferStackSlice(0)
				if len(_slice) > 0 {
					//fmt.Println("======= _transfer =======")
					go func(transferSlice []interface{}, stackStatus chan int) {
						defer func() {
							<-stackStatus
						}()
						that._transfer(transferSlice)
					}(_slice, that.taskTransferPool)
				} else {
					//fmt.Println("======= _transfer 001 =======")
					<-that.taskTransferPool
					runtime.Gosched()
					time.Sleep(time.Duration(500) * time.Millisecond)
				}
			default:
				//fmt.Println("======= _transfer 002 =======")
				runtime.Gosched()
				time.Sleep(time.Duration(500) * time.Millisecond)
			}
		}
	}()

	//
	go func() {
		for {
			_nowTime := time.Now().Unix()
			if that.runTime > _nowTime {
				_ticker := time.NewTicker(time.Duration(that.runTime-_nowTime) * time.Second)
				select {
				case <-that.sortSign:
					if !that._sort() {
						//fmt.Println("======= _sort 001 =======")
						runtime.Gosched()
						time.Sleep(time.Duration(500) * time.Millisecond)
					}
				case <-_ticker.C:
					if !that._sort() {
						//fmt.Println("======= _sort 002 =======")
						runtime.Gosched()
						time.Sleep(time.Duration(500) * time.Millisecond)
					}
				}
				_ticker.Stop()
				_ticker = nil
			} else {
				if !that._sort() {
					//fmt.Println("======= _sort 003 =======")
					runtime.Gosched()
					time.Sleep(time.Duration(500) * time.Millisecond)
				}
			}
		}
	}()

	//主进程：执行任务
	for {
		select {
		case that.taskWorkPool <- 0:
			select {
			case _workSlice := <-that.taskWorkStack:
				//任务采用异步方式执行
				//fmt.Println("======= _work  =======")
				go func(taskSlice []interface{}, stackStatus chan int) {
					defer func() {
						<-stackStatus
					}()
					for _, _t := range taskSlice {
						that._work(_t)
					}
				}(_workSlice, that.taskWorkPool)
			default:
				<-that.taskWorkPool
				//fmt.Println("======= _work 001 =======")
				runtime.Gosched()
				time.Sleep(time.Duration(500) * time.Millisecond)
			}
		default:
			//暂停
			//fmt.Println("======= _work 002 =======")
			runtime.Gosched()
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}
}

// _sort "倒排索引队列"中找出到期需要执行的"任务"
func (that *Tasker) _sort() bool {
	defer that.lock.Unlock()
	that.lock.Lock()
	_length := len(that.sortSlice)
	if _length == 0 {
		that.runTime = 0
		return false
	}
	//排序
	if _length != that.sortLength {
		sort.Sort(sort.IntSlice(that.sortSlice))
	}
	//
	_index := that.sortSlice[0]
	_data := that.taskMap[_index].(*taskData)
	_runTime := ToInt64(_data.runTime)
	//
	if that.runTime == 0 {
		that.runTime = _runTime
	}
	if that.runTime > _runTime {
		that.runTime = _runTime
	}
	//到期执行任务
	_status := false
	if time.Now().Unix() >= _runTime {
		that.sortSlice = that.sortSlice[1:len(that.sortSlice)]
		delete(that.taskMap, _index)
		that.taskWorkStack.Stack(_data)
		_status = true
	}
	that.sortLength = len(that.sortSlice)
	return _status
}

func (that *Tasker) _stack(runTime int, task interface{}) {
	_taskData := &taskData{
		runTime:  runTime,
		taskData: task,
	}
	_len := len(that.taskTransferStackSlice)
	if _len == 1 {
		that.taskTransferStackSlice[0].Stack(_taskData)
		return
	}
	_i := RandInt(_len)
	switch len(that.taskTransferStackSlice) {
	case 5:
		select {
		case that.taskTransferStackSlice[_len-1] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-2] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-3] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-4] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-5] <- []interface{}{_taskData}:
		default:
			that.taskTransferStackSlice[_i].Stack(_taskData)
		}
	case 4:
		select {
		case that.taskTransferStackSlice[_len-1] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-2] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-3] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-4] <- []interface{}{_taskData}:
		default:
			that.taskTransferStackSlice[_i].Stack(_taskData)
		}
	case 3:
		select {
		case that.taskTransferStackSlice[_len-1] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-2] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-3] <- []interface{}{_taskData}:
		default:
			that.taskTransferStackSlice[_i].Stack(_taskData)
		}
	case 2:
		select {
		case that.taskTransferStackSlice[_len-1] <- []interface{}{_taskData}:
		case that.taskTransferStackSlice[_len-2] <- []interface{}{_taskData}:
		default:
			that.taskTransferStackSlice[_i].Stack(_taskData)
		}
	}
}

// _workRange 执行任务
func (that *Tasker) _work(taskDataSlice ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("task doing fail! ", err)
		}
	}()
	for _, _taskInterface := range taskDataSlice {
		_task := _taskInterface.(*taskData)
		switch _t := _task.taskData.(type) {
		case func(): //函数类型
			_t()
		}
	}
}

// _rangeTransferStackSlice 合并缓存
func (that *Tasker) _rangeTransferStackSlice(i int, slice ...interface{}) []interface{} {
	_len := len(that.taskTransferStackSlice)
	if i >= _len {
		return slice
	}
	//
	select {
	case _transferSlice := <-that.taskTransferStackSlice[i]:
		slice = append(slice, _transferSlice...)
	default:
	}
	//
	i++
	return that._rangeTransferStackSlice(i, slice...)
}

// _transferRange 从"生产者叠加器" 转运 "任务" 到 "倒排索引队列"
func (that *Tasker) _transfer(taskDataSlice []interface{}) {
	for _, _taskInterface := range taskDataSlice {
		_taskData := _taskInterface.(*taskData)
		that._transferToSort(_taskData)
	}
}

// _transferAddTask 将"任务"添加到"倒排索引队列"
func (that *Tasker) _transferToSort(taskData *taskData) {
	defer that.lock.Unlock()
	that.lock.Lock()
	_index := that._getIndex(taskData.runTime)
	taskData.index = _index
	leng := len(that.sortSlice)
	if leng < 1 {
		that.sortSlice = append(that.sortSlice, _index)
		that.taskMap[_index] = taskData
	} else {
		that.sortSlice = append(that.sortSlice, _index)
		that.taskMap[_index] = taskData
	}
	//比较时间
	_runTime := int64(taskData.runTime)
	//
	if that.runTime == 0 {
		that.runTime = _runTime
	}
	if that.runTime > _runTime {
		that.runTime = _runTime
	}
	select {
	case that.sortSign <- 0:
	default:

	}
}

// _getIndex 获取计划任务的唯一键名索引
func (that *Tasker) _getIndex(index int) int {
	return that._getOnlyIndex(index * 1000)
}

// _getOnlyIndex
func (that *Tasker) _getOnlyIndex(index int) int {
	if _, found := that.taskMap[index]; found {
		return that._getOnlyIndex(index + 1)
	}
	return index

}

func NewStackChanPool(length ...int) *StackChanPool {
	_len := 4
	if len(length) > 0 {
		_len = length[0]
		if _len > 10 {
			_len = 10
		}
		if _len < 4 {
			_len = 4
		}
	}
	_stackChanSlice := []StackChan{}
	for _i := 0; _i < _len; _i++ {
		_stackChanSlice = append(_stackChanSlice, NewStackChan())
	}
	//fmt.Println("================== NewStackChanPool ===================", len(_stackChanSlice))
	return &StackChanPool{
		stackLen:       _len,
		stackChanSlice: _stackChanSlice,
	}
}

type StackChanPool struct {
	stackLen       int
	stackChanSlice []StackChan
}

func (that StackChanPool) Stack(value ...interface{}) {
	_i := RandInt(that.stackLen)
	that.stackChanSlice[_i].Stack(value...)
}

func (that StackChanPool) Get() []interface{} {
	switch that.stackLen {
	case 1:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		}
	case 2:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		case _t := <-that.stackChanSlice[1]:
			return _t
		}
	case 3:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		case _t := <-that.stackChanSlice[1]:
			return _t
		case _t := <-that.stackChanSlice[2]:
			return _t
		}
	case 4:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		case _t := <-that.stackChanSlice[1]:
			return _t
		case _t := <-that.stackChanSlice[2]:
			return _t
		case _t := <-that.stackChanSlice[3]:
			return _t
		}
	case 5:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		case _t := <-that.stackChanSlice[1]:
			return _t
		case _t := <-that.stackChanSlice[2]:
			return _t
		case _t := <-that.stackChanSlice[3]:
			return _t
		case _t := <-that.stackChanSlice[4]:
			return _t
		}
	case 6:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		case _t := <-that.stackChanSlice[1]:
			return _t
		case _t := <-that.stackChanSlice[2]:
			return _t
		case _t := <-that.stackChanSlice[3]:
			return _t
		case _t := <-that.stackChanSlice[4]:
			return _t
		case _t := <-that.stackChanSlice[5]:
			return _t
		}
	case 7:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		case _t := <-that.stackChanSlice[1]:
			return _t
		case _t := <-that.stackChanSlice[2]:
			return _t
		case _t := <-that.stackChanSlice[3]:
			return _t
		case _t := <-that.stackChanSlice[4]:
			return _t
		case _t := <-that.stackChanSlice[5]:
			return _t
		case _t := <-that.stackChanSlice[6]:
			return _t
		}
	case 8:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		case _t := <-that.stackChanSlice[1]:
			return _t
		case _t := <-that.stackChanSlice[2]:
			return _t
		case _t := <-that.stackChanSlice[3]:
			return _t
		case _t := <-that.stackChanSlice[4]:
			return _t
		case _t := <-that.stackChanSlice[5]:
			return _t
		case _t := <-that.stackChanSlice[6]:
			return _t
		case _t := <-that.stackChanSlice[7]:
			return _t
		}
	case 9:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		case _t := <-that.stackChanSlice[1]:
			return _t
		case _t := <-that.stackChanSlice[2]:
			return _t
		case _t := <-that.stackChanSlice[3]:
			return _t
		case _t := <-that.stackChanSlice[4]:
			return _t
		case _t := <-that.stackChanSlice[5]:
			return _t
		case _t := <-that.stackChanSlice[6]:
			return _t
		case _t := <-that.stackChanSlice[7]:
			return _t
		case _t := <-that.stackChanSlice[8]:
			return _t
		}
	case 10:
		select {
		case _t := <-that.stackChanSlice[0]:
			return _t
		case _t := <-that.stackChanSlice[1]:
			return _t
		case _t := <-that.stackChanSlice[2]:
			return _t
		case _t := <-that.stackChanSlice[3]:
			return _t
		case _t := <-that.stackChanSlice[4]:
			return _t
		case _t := <-that.stackChanSlice[5]:
			return _t
		case _t := <-that.stackChanSlice[6]:
			return _t
		case _t := <-that.stackChanSlice[7]:
			return _t
		case _t := <-that.stackChanSlice[8]:
			return _t
		case _t := <-that.stackChanSlice[9]:
			return _t
		}
	}
	return nil
}
