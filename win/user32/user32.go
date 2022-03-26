package user32

import (
	"syscall"
	"unsafe"
)

var (
	function = []string{
		"GetCursorPos",
		"SetCursorPos",
		"mouse_event",
	}
)

const (
	// MouseEventLeftDown 左クリックダウン
	MouseEventLeftDown = 0x0002

	// MouseEventLeftUp 左クリックアップ
	MouseEventLeftUp = 0x0004

	// MouseEventRightDown 右クリックダウン
	MouseEventRightDown = 0x0008

	// MouseEventRightUp 右クリックアップ
	MouseEventRightUp = 0x0010

	// MouseEventWheel ホイール
	MouseEventWheel = 0x0800
)

// Point 位置
type Point struct {
	X int32
	Y int32
}

// DLL user32.dll
type DLL struct {
	dll   *syscall.DLL
	procm map[string]*syscall.Proc
}

// NewDLL dll
func NewDLL() (*DLL, error) {
	d, err := syscall.LoadDLL("user32.dll")
	if err != nil {
		return nil, err
	}
	procm := map[string]*syscall.Proc{}
	for _, s := range function {
		p, err := d.FindProc(s)
		if err != nil {
			return nil, err
		}
		procm[s] = p
	}

	return &DLL{
		dll:   d,
		procm: procm,
	}, nil
}

// Release dllのリリース
func (t *DLL) Release() {
	t.dll.Release()
}

// GetCursorPos マウスの位置を取得する
func (t *DLL) GetCursorPos() (*Point, bool, error) {
	var ret Point
	pp := unsafe.Pointer(&ret)
	r1, _, _ := t.procm["GetCursorPos"].Call(uintptr(pp))

	err := syscall.GetLastError()
	if err != nil {
		return nil, false, err
	}

	return &ret, r1 == 1, nil
}

// SetCursorPos マウス位置設定
func (t *DLL) SetCursorPos(x, y int) (bool, error) {
	r1, _, _ := t.procm["SetCursorPos"].Call(
		uintptr(x), uintptr(y),
	)
	err := syscall.GetLastError()
	if err != nil {
		return false, err
	}
	return r1 == 1, nil
}

// MouseEvent user32dllのmouse_event
func (t *DLL) MouseEvent(dwFlags, dx, dy, dwData, dwinfo int) error {
	t.procm["mouse_event"].Call(uintptr(dwFlags), uintptr(dx), uintptr(dy), uintptr(dwData), uintptr(0))
	err := syscall.GetLastError()
	if err != nil {
		return err
	}
	return nil
}

// RightClick 右クリックをCurrentの位置で行う
func (t *DLL) RightClick() error {
	pos, _, err := t.GetCursorPos()
	if err != nil {
		return err
	}
	return t.MouseEvent(MouseEventRightDown|MouseEventRightUp, int(pos.X), int(pos.Y), 0, 0)
}

// RightClickAt 右クリックをx,yの位置で行う
func (t *DLL) RightClickAt(x, y int) error {
	return t.MouseEvent(MouseEventRightDown|MouseEventRightUp, x, y, 0, 0)
}

// LeftClick 左クリック
func (t *DLL) LeftClick() error {
	pos, _, err := t.GetCursorPos()
	if err != nil {
		return err
	}
	return t.LeftClickAt(int(pos.X), int(pos.Y))
}

// LeftClickAt 左クリックをx,yの位置で行う
func (t *DLL) LeftClickAt(x, y int) error {
	return t.MouseEvent(MouseEventLeftDown|MouseEventLeftUp, x, y, 0, 0)
}

// MouseWheel マウスホイールを行う
func (t *DLL) MouseWheel(value int) error {
	return t.MouseEvent(MouseEventWheel, 0, 0, value, 0)
}
