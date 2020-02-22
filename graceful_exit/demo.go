package graceful_exit

// 要达到优雅退出，那么在退出前，需要让processor感知到，并快速处理完当前手头工作。
// 常见的工作包括：
// 1：存储内部状态
// 2：拒绝接受新数据，并将已经接收的数据处理掉

import (
	"fmt"
	"sync"
	"time"
)

var (
	FullError = fmt.Errorf("full, reject")
)

type Manager struct {
	jobQueue chan string
	resultQueue chan string

	done chan struct{} // used for closing
	wg   sync.WaitGroup
}

func NewManager(queueSize int) (*Manager, error) {
	m := &Manager{
		jobQueue: make(chan string, queueSize),
		done:      make(chan struct{}),
	}
	return m, nil
}

func (m *Manager) Start() {
	for i := 0; i < 2; i++ {
		go m.process(i)
	}
}

func (m *Manager) Insert(data string) error {
	if m.IsClosed() {
		return fmt.Errorf("not running")
	}

	queueLength := len(m.jobQueue)
	queueCapacity := cap(m.jobQueue)
	if queueLength+1 > queueCapacity {
		return FullError
	}

	m.jobQueue <- data
	m.wg.Add(1)

	return nil
}

func (m *Manager) Stop() {
	fmt.Println("[stop] start to exit")

	// close(m.done)
	// m.wg.Wait()

	// save status

	fmt.Println("[stop] finish to exit")
}

func (m *Manager) IsClosed() bool {
	select {
	case <-m.done:
		return true
	default:
	}

	return false
}

func (m *Manager) process(workNumber int) {
	for {
		m.processUtil(workNumber)
	}
}

func (m *Manager) processUtil(workNumber int) {
	select {
	case data := <-m.jobQueue:
		func() {
			defer m.wg.Done()

			fmt.Printf("[process]start to process data[%s] in worker[%d]\n", data, workNumber)
			time.Sleep(time.Microsecond)
		}()

	default:
		fmt.Printf("[process] no data to process, sleep 1 second in worker[%d]\n", workNumber)
		time.Sleep(time.Second)
	}
}
