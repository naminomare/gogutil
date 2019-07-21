package timer

import "time"

var (
	// TimeStampFormat フォーマット
	TimeStampFormat = "200601021504"
)

// TimeStamp 現在のタイムスタンプを返す
func TimeStamp() string {
	return time.Now().Format(TimeStampFormat)
}
