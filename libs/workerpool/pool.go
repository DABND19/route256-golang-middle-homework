package workerpool

import "sync"

// Вдохновлено репозиторием https://github.com/gammazero/workerpool :)
type WorkerPool interface {
	Submit(task func()) // Отправить задачу в пул. Отправка в закрытый пул вызывает панику записи в закрытый канал.
	IsClosed() bool
	Close()     // Закрыть пул и дождаться задачи, которые уже исполняются воркерами.
	WaitClose() // Закрыть пул и дождаться задачи, которые были отправлены.
}

type pool struct {
	amountWorkers int
	workersWg     *sync.WaitGroup // Для ожидания завершения работы всех воркеров после остановки event loop'а.
	workersQueue  chan func()     // Канал отправки задачек воркерам из event loop'а.

	tasksQueue   chan func() // Канал отправки задачек из метода Submit в event loop.
	waitingTasks []func()    // Для сохранения задачек, которые еще не успели выполниться.

	waitOnStop  bool
	closeOnce   *sync.Once
	stopped     chan struct{} // Канал получения сигнала остановки event loop'а.
	isClosed    bool
	isClosedMtx *sync.RWMutex
}

func New(size int) WorkerPool {
	p := &pool{
		amountWorkers: size,
		workersWg:     &sync.WaitGroup{},
		tasksQueue:    make(chan func()),
		workersQueue:  make(chan func()),
		waitingTasks:  make([]func(), 0),
		waitOnStop:    false,
		stopped:       make(chan struct{}),
		closeOnce:     &sync.Once{},
		isClosed:      false,
		isClosedMtx:   &sync.RWMutex{},
	}
	p.bootstrap()
	return p
}

func (p *pool) Submit(task func()) {
	// Избежание паники записи в закрытый канал на совести пользователя,
	// для этого ему предоставлен метод IsClosed.
	go func() {
		p.tasksQueue <- task
	}()
}

func (p *pool) close(wait bool) {
	p.closeOnce.Do(func() {
		// Вначале сообщим, что пул закрыт:
		p.isClosedMtx.Lock()
		p.isClosed = true
		p.isClosedMtx.Unlock()

		// И только потом его остановим:
		p.waitOnStop = wait
		close(p.tasksQueue)
	})
	// Дождемся остановки горутины с event loop'ом:
	<-p.stopped
}

func (p *pool) IsClosed() bool {
	p.isClosedMtx.RLock()
	defer p.isClosedMtx.RUnlock()
	return p.isClosed
}

func (p *pool) Close() {
	p.close(false)
}

func (p *pool) WaitClose() {
	p.close(true)
}

// Для разбора задач, которые в очереди ожидания. Если канал приема задач закрыт, возвращает false.
func (p *pool) processWaitingTasks() bool {
	select {
	// По возможности отправляем задачи из очереди ожидания на выполнение:
	case p.workersQueue <- p.waitingTasks[0]:
		p.waitingTasks = p.waitingTasks[1:len(p.waitingTasks)]

	// Ждем новые задачки и кладем их в очередь:
	case task, ok := <-p.tasksQueue:
		if !ok {
			return false
		}
		p.waitingTasks = append(p.waitingTasks, task)
	}
	return true
}

func (p *pool) startEventLoop() {
	defer close(p.stopped)

eventLoop:
	for {
		// Первым делом следует разобрать задачки из очереди ожидания:
		if len(p.waitingTasks) != 0 {
			if !p.processWaitingTasks() {
				break eventLoop
			}
			continue
		}

		// Если нет ожидающих задач, то ждем пока прилетят из канала:
		select {
		case task, ok := <-p.tasksQueue:
			// Если канал приема задач закрыт, то выходим из event loop'а:
			if !ok {
				break eventLoop
			}

			select {
			// Попробуем отправить задачку на выполнение:
			case p.workersQueue <- task:

			// Если не удалось, то отправляем ее в очередь ожидания:
			default:
				p.waitingTasks = append(p.waitingTasks, task)
			}
		}
	}

	// event loop закончился, если была команда доделать оставшиеся задачи, то доделываем их:
	if p.waitOnStop {
		for len(p.waitingTasks) != 0 {
			p.workersQueue <- p.waitingTasks[0]
			p.waitingTasks = p.waitingTasks[1:len(p.waitingTasks)]
		}
	}

	// Закрываем канал отправки задач на выполнение. Здесь воркеры получают сигнал на выключение.
	close(p.workersQueue)

	// Ждем пока все воркеры выключатся.
	p.workersWg.Wait()
}

func worker(tasks <-chan func()) {
	for task := range tasks {
		task()
	}
}

func (p *pool) bootstrap() {
	go p.startEventLoop()

	for i := 0; i < p.amountWorkers; i++ {
		p.workersWg.Add(1)
		go func() {
			defer p.workersWg.Done()
			worker(p.workersQueue)
		}()
	}
}
