package timer

import "time"

// StopWatch 時間計測用
type StopWatch struct {
	startTime time.Time

	pauseDuration time.Duration
	pauseStart    time.Time
}

// NewStopWatch StopWatchインスタンスを取得する
func NewStopWatch() *StopWatch {
	return &StopWatch{}
}

// Start ここから時間を計測する
func (t *StopWatch) Start() {
	t.pauseDuration = 0
	t.startTime = time.Now()
}

// Pause 一時中断
func (t *StopWatch) Pause() {
	t.pauseStart = time.Now()
}

// Resume PauseしたところからResumeする
func (t *StopWatch) Resume() {
	pauseEnd := time.Now()
	t.pauseDuration += pauseEnd.Sub(t.pauseStart)
}

// Stop ここまでの時間を取得する
func (t *StopWatch) Stop() time.Duration {
	now := time.Now()
	return now.Sub(t.startTime) - t.pauseDuration
}
