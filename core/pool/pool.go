package pool

import (
	"math/rand"
	"separa/common/log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Worker struct {
	f func(interface{})
}

func generateWorker(f func(interface{})) *Worker {
	return &Worker{
		f: func(in interface{}) {
			f(in)
		},
	}
}

func (t *Worker) run(in interface{}) {
	t.f(in)
}

type Pool struct {
	//母版函数
	Function func(interface{})
	//Pool输入队列
	in chan interface{}
	//size用来表明池的大小，不能超发。
	threads int
	//正在执行的任务清单
	JobsList *sync.Map
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
		JobsList: &sync.Map{},
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
	p.Done = false
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
	var Tick string
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
		Tick = p.generateTick()
		p.JobsList.Store(Tick, param)
		f := generateWorker(p.Function)
		f.run(param)
		p.JobsList.Delete(Tick)
		atomic.AddInt32(&p.active, -1)
	}
}

func (p *Pool) generateTick() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(rand.Int63(), 10)
}

func (p *Pool) Threads() int {
	return p.threads
}

func (p *Pool) RunningThreads() int {
	return int(p.active)
}
