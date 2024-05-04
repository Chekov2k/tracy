package main

import (
	"runtime"
	"unsafe"
)

// #cgo CFLAGS: -DTRACY_ENABLE -DTRACY_NAME_BUFFER -DTRACY_FIBERS
// #cgo LDFLAGS: -lTracyClient
// #include <stdlib.h>
// #include <tracy.h>
import "C"

type Fiber struct {
	ptr *C.char
	id  C.uint16_t
}

type Zone struct {
	fiber    *Fiber
	name     *C.char
	file     *C.char
	function *C.char
	zone     C.TracyCZoneCtx
}

func NewFiber(name string) *Fiber {
	if C.TracyEnabled() == C.uint8_t(0) {
		return nil
	}

	var cName *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cId = C.uint16_t(0)
	ptr := C.FiberStart(cName, &cId)
	if ptr == nil {
		return nil
	}

	return &Fiber{ptr: ptr, id: cId}
}

func (f *Fiber) Close() {
	if f == nil {
		return
	}
	C.FiberLeave()
}

func (f *Fiber) NewZone(name string, color Color, depth int) *Zone {
	if f == nil {
		return nil
	}

	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return nil
	}

	function := runtime.FuncForPC(pc)
	if function == nil {
		return nil
	}

	cName := C.CString(name)
	cFile := C.CString(file)
	cFunction := C.CString(function.Name())

	return &Zone{
		fiber:    f,
		name:     cName,
		file:     cFile,
		function: cFunction,
		zone: C.ZoneStart(
			f.id,
			C.uint32_t(line),
			cFile,
			C.size_t(len(file)),
			cFunction,
			C.size_t(len(function.Name())),
			cName,
			C.size_t(len(name)),
			C.uint32_t(color),
			C.int(depth),
		),
	}
}

func (z *Zone) Close() {
	if z == nil {
		return
	}

	C.ZoneEnd(z.fiber.id, z.zone)

	C.free(unsafe.Pointer(z.name))
	C.free(unsafe.Pointer(z.file))
	C.free(unsafe.Pointer(z.function))
}

func (z *Zone) Text(text string) {
	if z == nil {
		return
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.___tracy_emit_zone_text(z.zone, cText, C.size_t(len(text)))
}

func (z *Zone) Name(name string) {
	if z == nil {
		return
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.___tracy_emit_zone_name(z.zone, cName, C.size_t(len(name)))
}

func (z *Zone) Color(color Color) {
	if z == nil {
		return
	}
	C.___tracy_emit_zone_color(z.zone, C.uint32_t(color))
}

type Frame struct {
	id C.uint16_t
}

func NewFrame(name string, doMark bool) *Frame {
	if C.TracyEnabled() == C.uint8_t(0) {
		return nil
	}

	var cName *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cMark = C.uint16_t(0)
	if doMark {
		cMark = C.uint16_t(1)
	}

	var cId = C.uint16_t(0)
	C.FrameCreate(cName, &cId, cMark)

	return &Frame{id: cId}
}

func (f *Frame) Mark() {
	if f == nil {
		return
	}
	C.FrameMark(f.id)
}

func (f *Frame) Start() {
	if f == nil {
		return
	}
	C.FrameStart(f.id)
}

func (f *Frame) End() {
	if f == nil {
		return
	}
	C.FrameEnd(f.id)
}
