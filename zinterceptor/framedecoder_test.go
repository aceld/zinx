package zinterceptor

import (
	"bytes"
	"testing"
)

// A peer sending one frame larger than MaxFrameLength must not stop the
// decoder from delivering the frames that follow it.
func TestFrameDecoderDiscardsOversizedFrameThenDecodesNextFrame(t *testing.T) {
	// 2-byte length field at offset 0 holding the payload length only.
	// MaxFrameLength 10 bytes => a payload above 8 bytes is too long.
	decoder := NewFrameDecoderByParams(10, 0, 2, 0, 0)

	// One oversized frame: length 20 -> 22 bytes on the wire.
	oversized := append([]byte{0x00, 0x14}, bytes.Repeat([]byte{0xAA}, 20)...)
	if frames := decoder.Decode(oversized); len(frames) != 0 {
		t.Fatalf("oversized frame delivered as data: %d frames (%x)", len(frames), frames)
	}

	// A well-formed frame that arrives afterwards.
	valid := []byte{0x00, 0x03, 'a', 'b', 'c'}
	frames := decoder.Decode(valid)
	if len(frames) != 1 {
		t.Fatalf("frame after the oversized one was not delivered: got %d frames (%x), want 1", len(frames), frames)
	}
	if !bytes.Equal(frames[0], valid) {
		t.Fatalf("frame after the oversized one decoded wrongly: got %x, want %x", frames[0], valid)
	}
}

// When the tail of an oversized frame is followed by an incomplete frame,
// the discarded bytes must still be released so that the incomplete frame is
// parsed once its remainder arrives.
func TestFrameDecoderDiscardsOversizedFrameThenKeepsHalfFrame(t *testing.T) {
	decoder := NewFrameDecoderByParams(10, 0, 2, 0, 0)

	// Oversized frame claims 20 payload bytes; send all of it except the last
	// byte, then in the second read send that byte plus the first 2 bytes of the
	// next frame (2-byte length + 1 payload byte of 3).
	first := append([]byte{0x00, 0x14}, bytes.Repeat([]byte{0xAA}, 19)...)
	if frames := decoder.Decode(first); len(frames) != 0 {
		t.Fatalf("partial oversized frame delivered as data: %d frames (%x)", len(frames), frames)
	}

	second := append([]byte{0xAA}, []byte{0x00, 0x03, 'a'}...)
	if frames := decoder.Decode(second); len(frames) != 0 {
		t.Fatalf("half frame delivered as data: %d frames (%x)", len(frames), frames)
	}

	// The final 2 bytes complete the frame.
	frames := decoder.Decode([]byte{'b', 'c'})
	if len(frames) != 1 || !bytes.Equal(frames[0], []byte{0x00, 0x03, 'a', 'b', 'c'}) {
		t.Fatalf("frame completed after discard: got %d frames (%x), want [0003616263]", len(frames), frames)
	}
}

// The tail of an oversized frame can arrive in the same read as a complete
// frame. The discarded bytes must be released so the next frame is not
// corrupted by leftovers on the following read.
func TestFrameDecoderDiscardsOversizedTailThenDecodesFrameInSameRead(t *testing.T) {
	decoder := NewFrameDecoderByParams(10, 0, 2, 0, 0)

	// Oversized frame declares 20 payload bytes, send 19 of them.
	first := append([]byte{0x00, 0x14}, bytes.Repeat([]byte{0xAA}, 19)...)
	if frames := decoder.Decode(first); len(frames) != 0 {
		t.Fatalf("partial oversized frame delivered as data: %d frames (%x)", len(frames), frames)
	}

	// Second read: the last oversized byte, then a complete 3-byte-payload frame.
	second := append([]byte{0xAA}, []byte{0x00, 0x03, 'a', 'b', 'c'}...)
	frames := decoder.Decode(second)
	if len(frames) != 1 || !bytes.Equal(frames[0], []byte{0x00, 0x03, 'a', 'b', 'c'}) {
		t.Fatalf("frame behind oversized tail: got %d frames (%x), want [0003616263]", len(frames), frames)
	}

	// Nothing of the discarded frame may remain: the next frame decodes cleanly.
	next := []byte{0x00, 0x02, 'd', 'e'}
	frames = decoder.Decode(next)
	if len(frames) != 1 || !bytes.Equal(frames[0], next) {
		t.Fatalf("leftover bytes corrupted the next frame: got %d frames (%x), want [%x]", len(frames), frames, next)
	}
}

// The same situation when the oversized frame is split across two reads.
func TestFrameDecoderDiscardsSplitOversizedFrameThenDecodesNextFrame(t *testing.T) {
	decoder := NewFrameDecoderByParams(10, 0, 2, 0, 0)

	// First read: length field plus part of the payload.
	if frames := decoder.Decode([]byte{0x00, 0x14, 0xAA, 0xAA, 0xAA}); len(frames) != 0 {
		t.Fatalf("partial oversized frame delivered as data: %d frames (%x)", len(frames), frames)
	}
	// Second read: the rest of the oversized payload.
	if frames := decoder.Decode(bytes.Repeat([]byte{0xAA}, 17)); len(frames) != 0 {
		t.Fatalf("rest of oversized frame delivered as data: %d frames (%x)", len(frames), frames)
	}

	// Third read: a well-formed frame.
	valid := []byte{0x00, 0x03, 'a', 'b', 'c'}
	frames := decoder.Decode(valid)
	if len(frames) != 1 || !bytes.Equal(frames[0], valid) {
		t.Fatalf("frame after the split oversized one: got %d frames (%x), want [%x]", len(frames), frames, valid)
	}
}
