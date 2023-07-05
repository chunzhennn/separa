package pool

import (
	"separa/common/log"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	//母版函数
	Function func(interface{})
	//Pool输入队列
	in chan interface{}
	//size用来表明池的大小，不能超发。
	threads int
	//启动协程等待时间
	Interval time.Duration
	//正在工作的协程数量
	active int32
	//用于阻塞
	wg *sync.WaitGroup
	//提前结束标识符
	Done bool
}

// instance of the pool
func New(threads int) *Pool {
	return &Pool{
		threads:  threads,
		wg:       &sync.WaitGroup{},
		in:       make(chan interface{}),
		Function: nil,
		Done:     true,
		Interval: time.Duration(0),
	}
}

// Push element to the pool
func (p *Pool) Push(i interface{}) {
	if p.Done {
		return
	}
	p.in <- i
}

// Stop recieving new incoming elements
func (p *Pool) Stop() {
	if !p.Done {
		close(p.in)
	}
	p.Done = true
}

// Run the pool
func (p *Pool) Run() {
	for i := 0; i < p.threads; i++ {
		p.wg.Add(1)
		time.Sleep(p.Interval)
		go p.work()
		if p.Done {
			break
		}
	}
	p.wg.Wait()
}

// pool worker
func (p *Pool) work() {
	p.Done = false
	var param interface{}

	defer func() {
		defer func() {
			if e := recover(); e != nil {
				log.Log.Printf("[pool] %s: %s", param, e)
			}
		}()
		p.wg.Done()
	}()
	for param = range p.in {
		if p.Done {
			return
		}
		atomic.AddInt32(&p.active, 1)
		go p.Function(param)
		atomic.AddInt32(&p.active, -1)
	}
}

// 获取线程数
func (p *Pool) Threads() int {
	return p.threads
}

func (p *Pool) RunningThreads() int {
	return int(p.active)
}
