package adb

import (
	"bytes"
	"fmt"
	"os/exec"
)

func runAdb(args ...string) (err error) {
	var b bytes.Buffer
	cmd := exec.Command("adb", args...)
	cmd.Stdout = &b
	cmd.Stderr = &b
	err = cmd.Run()
	if err != nil {
		return
	}
	if cmd.Process != nil {
		err = cmd.Process.Kill()
	}
	return
}

//滑动
func slide(x1, y1, x2, y2, slideTime int) {
	runAdb("shell", fmt.Sprintf("input swipe %d %d %d %d %d", x1, y1, x2, y2, int(slideTime)))
}

//长按
func longPress(x, y, pressTime int) {
	runAdb("shell", fmt.Sprintf("input swipe %d %d %d %d %d", x, y, x, y, int(pressTime)))
}

//保存截图
func saveScreenShot(filename string) {
	filePath := "/sdcard/" + filename
	runAdb("shell", "screencap -p "+filePath)
	runAdb("pull", filePath, ".")
}
