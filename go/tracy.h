#pragma once

#include <tracy/TracyC.h>

u_int8_t TracyEnabled() { return TracyCEnabled(); }

const char* FiberStart(const char* name, uint16_t* id) {
  const char* ptr = TracyCNameBufferAdd(name, id);
  if (!ptr) return 0;
  TracyCFiberEnter(ptr);
  return ptr;
}

void FiberEnter(uint16_t id) {
  const char* ptr = TracyCNameBufferGet(id);
  if (!ptr) return;
  TracyCFiberEnter(ptr);
}

void FiberLeave() { TracyCFiberLeave; }

TracyCZoneCtx ZoneStart(uint16_t fibreId, uint32_t line, const char* source,
                        size_t sourceSz, const char* function,
                        size_t functionSz, const char* name, size_t nameSz,
                        uint32_t color, int depth) {
  FiberEnter(fibreId);
  uint64_t srcLoc = ___tracy_alloc_srcloc_name(line, source, sourceSz, function,
                                               functionSz, name, nameSz, color);
  if (depth)
    return ___tracy_emit_zone_begin_alloc_callstack(srcLoc, depth, 1);
  else
    return ___tracy_emit_zone_begin_alloc(srcLoc, 1);
}

void ZoneEnd(uint16_t fibreId, TracyCZoneCtx zone) {
  FiberEnter(fibreId);
  ___tracy_emit_zone_end(zone);
}

void FrameCreate(const char* name, uint16_t* id, uint16_t mark) {
  if (!name) TracyCFrameMark();
  const char* ptr = TracyCNameBufferAdd(name, id);
  if (!ptr || !mark) return;
  ___tracy_emit_frame_mark(ptr);
}

void FrameMark(uint16_t id) {
  const char* ptr = TracyCNameBufferGet(id);
  if (!ptr) return;
  ___tracy_emit_frame_mark(ptr);
}

void FrameStart(uint16_t id) {
  const char* ptr = TracyCNameBufferGet(id);
  if (!ptr) return;
  ___tracy_emit_frame_mark_start(ptr);
}

void FrameEnd(uint16_t id) {
  const char* ptr = TracyCNameBufferGet(id);
  if (!ptr) return;
  ___tracy_emit_frame_mark_end(ptr);
}
