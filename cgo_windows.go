package systray

/*
#cgo windows CFLAGS: -DWIN32

#include "systray.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"
)

func nativeLoop() {
	C.nativeLoop()
}
var prevIconPath string
func quit() {
	os.Remove(prevIconPath)
	C.quit()
}



func SetIcon(iconBytes []byte) {
	f, err := ioutil.TempFile("", "systray_temp_icon")
	if err != nil {
		log.Errorf("Unable to create temp icon: %v", err)
		return
	}
	defer f.Close()
	_, err = f.Write(iconBytes)
	if err != nil {
		log.Errorf("Unable to write icon to temp file %v: %v", f.Name(), f)
		return
	}
	iconPath := f.Name()
	// Need to close file before we load it to make sure contents is flushed.
	f.Close()
	name, err := syscall.UTF16PtrFromString(f.Name())
	if err != nil {
		log.Errorf("Unable to convert name to string pointer: %v", err)
		return
	}
	C.setIcon((*C.wchar_t)(unsafe.Pointer(name)), 0)
	// remove old icon file
	os.Remove(prevIconPath)
	prevIconPath = iconPath
}

// SetTitle sets the systray title, only available on Mac.
func SetTitle(title string) {
	t, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		panic(err)
	}

	C.setTitle((*C.wchar_t)(unsafe.Pointer(t)))
}

// SetTooltip sets the systray tooltip to display on mouse hover of the tray icon,
// only available on Mac and Windows.
func SetTooltip(tooltip string) {
	t, err := syscall.UTF16PtrFromString(tooltip)
	if err != nil {
		panic(err)
	}

	C.setTooltip((*C.wchar_t)(unsafe.Pointer(t)))

}

func addOrUpdateMenuItem(item *MenuItem) {
	var disabled C.short = 0
	if item.disabled {
		disabled = 1
	}
	var checked C.short = 0
	if item.checked {
		checked = 1
	}

	title, err := syscall.UTF16PtrFromString(item.title)
	if err != nil {
		panic(err)
	}

	tooltip, err := syscall.UTF16PtrFromString(item.tooltip)
	if err != nil {
		panic(err)
	}

	C.add_or_update_menu_item(
		C.int(item.id),
		(*C.wchar_t)(unsafe.Pointer(title)),
		(*C.wchar_t)(unsafe.Pointer(tooltip)),
		disabled,
		checked,
	)
}

//export systray_ready
func systray_ready() {
	fmt.Println("systray ready")
	systrayReady()
}

//export systray_menu_item_selected
func systray_menu_item_selected(cId C.int) {
	systrayMenuItemSelected(int32(cId))
}
