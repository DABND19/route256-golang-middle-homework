package scheduler

import (
	"container/heap"
	"sync"
	"time"
)

type Scheduler interface {
	Schedule(after time.Time, task func()) // Запланировать отправку задачки в worker pool после определенного момента времени.
	IsClosed() bool
	Close()     // Закрытие без отправки запланированных задач.
	WaitClose() // Дождаться отправки всех запланированных задач и зыкрыться.
}

type WorkerPool interface {
	Submit(task func())
}

type Task struct {
	Priority time.Time
	Timer    <-chan time.Time // Сохраняем канал, полученный из time.After, для избежания утечек памяти.
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

	tasksQueue     chan Task     // Канал между Submit и event loop'ом.
	scheduledTasks PriorityQueue // Очередь запланированных отправок в worker pool.

	stopped    chan struct{} // Сигнал остановки event loop'а.
	waitOnStop bool          // Дождаться ли отправки всех запланированных задач.

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
	// Ждем пока можно будет отправить самую раннюю задачку
	case <-s.scheduledTasks[0].Timer:
		task := heap.Pop(&s.scheduledTasks).(Task)
		s.workerPool.Submit(task.Callback)

	// А также ждем новые задачки и кладем их в очередь.
	case task, ok := <-s.tasksQueue:
		// Если канал приема новых задач закрыт, то выходим из evet loop'а.
		if !ok {
			return false
		}
		heap.Push(&s.scheduledTasks, task)
	}

	return true
}

func (s *scheduler) runEventLoop() {
	// Отправка сигнала остановки event loop'а
	defer close(s.stopped)

eventLoop:
	for {
		// Если в очереди что-то есть, то разгребаем задачки.
		if len(s.scheduledTasks) != 0 {
			// Если канал приема новых задач закрыт, то выходим из evet loop'а.
			if !s.processScheduledTasks() {
				break eventLoop
			}
			continue eventLoop
		}

		// Если нет, то ждем, пока придут новые.
		select {
		case task, ok := <-s.tasksQueue:
			// Если канал приема новых задач закрыт, то выходим из evet loop'а.
			if !ok {
				break eventLoop
			}
			heap.Push(&s.scheduledTasks, task)
		}
	}

	// Если надо, то дождемся отправки всех запланированных задач
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
