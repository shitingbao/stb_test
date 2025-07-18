package main

/*
#cgo CFLAGS: -I/Library/Frameworks/Python.framework/Versions/3.11/include/python3.11
#cgo LDFLAGS: -L/Library/Frameworks/Python.framework/Versions/3.11/lib -lpython3.11 -ldl -framework CoreFoundation ./ccode/libcallpy.a
#include "ccode/callpy.h"
*/
import "C"

func main() {
	C.call_python(C.CString("Gu Yu"), 3)
}
