package scheduler

import (
	"container/heap"
	"sync"
	"time"
)

type Scheduler interface {
	Schedule(after time.Time, task func())
	IsClosed() bool
	Close()
	WaitClose()
}

type WorkerPool interface {
	Submit(task func())
}

type Task struct {
	Priority time.Time
	Timer    <-chan time.Time
	Callback func()
}

type PriorityQueue []Task

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority.Before(pq[j].Priority)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(item any) {
	*pq = append(*pq, item.(Task))
}

func (pq *PriorityQueue) Pop() any {
	extractedItem := (*pq)[len(*pq)-1]
	*pq = (*pq)[0 : len(*pq)-1]
	return extractedItem
}

type scheduler struct {
	workerPool WorkerPool

	tasksQueue     chan Task
	scheduledTasks PriorityQueue

	stopped    chan struct{}
	waitOnStop bool

	closeOnce   *sync.Once
	isClosed    bool
	isClosedMtx *sync.RWMutex
}

func New(workerPool WorkerPool) Scheduler {
	s := &scheduler{
		workerPool:     workerPool,
		tasksQueue:     make(chan Task),
		scheduledTasks: make(PriorityQueue, 0),
		stopped:        make(chan struct{}),
		waitOnStop:     false,
		closeOnce:      new(sync.Once),
		isClosed:       false,
		isClosedMtx:    new(sync.RWMutex),
	}
	s.bootstrap()
	return s
}

func (s *scheduler) IsClosed() bool {
	s.isClosedMtx.RLock()
	defer s.isClosedMtx.RUnlock()
	return s.isClosed
}

func (s *scheduler) close(wait bool) {
	s.closeOnce.Do(func() {
		s.isClosedMtx.Lock()
		s.isClosed = true
		s.isClosedMtx.Unlock()

		s.waitOnStop = wait
		close(s.tasksQueue)
	})
	<-s.stopped
}

func (s *scheduler) Close() {
	s.close(false)
}

func (s *scheduler) WaitClose() {
	s.close(true)
}

func (s *scheduler) Schedule(after time.Time, task func()) {
	s.tasksQueue <- Task{
		Priority: after,
		Timer:    time.After(time.Until(after)),
		Callback: task,
	}
}

func (s *scheduler) processScheduledTasks() bool {
	select {
	case <-s.scheduledTasks[0].Timer:
		task := heap.Pop(&s.scheduledTasks).(Task)
		s.workerPool.Submit(task.Callback)

	case task, ok := <-s.tasksQueue:
		if !ok {
			return false
		}
		heap.Push(&s.scheduledTasks, task)
	}

	return true
}

func (s *scheduler) runEventLoop() {
	defer close(s.stopped)

eventLoop:
	for {
		if len(s.scheduledTasks) != 0 {
			if !s.processScheduledTasks() {
				break eventLoop
			}
			continue eventLoop
		}

		select {
		case task, ok := <-s.tasksQueue:
			if !ok {
				break eventLoop
			}
			heap.Push(&s.scheduledTasks, task)
		}
	}

	if s.waitOnStop {
		for len(s.scheduledTasks) != 0 {
			select {
			case <-s.scheduledTasks[0].Timer:
				task := heap.Pop(&s.scheduledTasks).(Task)
				s.workerPool.Submit(task.Callback)
			}
		}
	}
}

func (s *scheduler) bootstrap() {
	go s.runEventLoop()
}
