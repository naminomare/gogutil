package timer

import "time"

var (
	// TimeStampFormat フォーマット
	TimeStampFormat = "200601021504"
	// DateFormat フォーマット
	DateFormat = "20060102"
)

// TimeStamp 現在のタイムスタンプを返す
func TimeStamp() string {
	return time.Now().Format(TimeStampFormat)
}

// Date 現在の日付を返す
func Date() string {
	return time.Now().Format(DateFormat)
}
