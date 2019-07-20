package timer

import "time"

// WaitTimer 一定時間ブロックするタイマー
type WaitTimer struct {
	isDone chan struct{}
}

// NewWaitTimer WaitTimerオブジェクトを返す
func NewWaitTimer() *WaitTimer {
	return &WaitTimer{}
}

// Start タイマーを実行する
func (t *WaitTimer) Start(intervalMs int) {
	t.Wait()
	t.isDone = make(chan struct{})
	go func() {
		time.Sleep(time.Duration(intervalMs) * time.Millisecond)

		close(t.isDone)
	}()
}

// Wait ブロックします
func (t *WaitTimer) Wait() {
	if t.isDone == nil {
		return
	}
	<-t.isDone
}
