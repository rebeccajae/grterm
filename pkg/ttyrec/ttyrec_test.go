package ttyrec

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"time"
)

func TestNanosToTimeval(t *testing.T) {
	timeval := NanosToTimeval(1590644963489000000)
	expected := &Timeval{
		Sec:  1590644963,
		Usec: 489000,
	}
	if !reflect.DeepEqual(timeval, expected) {
		t.Fatalf("Timestamp Mismatch, got:\n%+v\nexpected:\n%+v\n", timeval, expected)
	}
}
func TestTTYRecWriter(t *testing.T) {
	testBytes := []byte{0x01, 0x02, 0x03, 0x04}
	var bwc bytes.Buffer
	rec := NewTTYRecorder(&bwc)
	ts := time.Now()
	before := NanosToTimeval(ts.UnixNano())
	rec.Write(testBytes)
	ts = time.Now()
	after := NanosToTimeval(ts.UnixNano())
	res := bwc.Bytes()
	sec := int32(binary.LittleEndian.Uint32(res[0:4]))
	usec := int32(binary.LittleEndian.Uint32(res[4:8]))
	length := int32(binary.LittleEndian.Uint32(res[8:12]))

	if int(length) != len(testBytes) {
		t.Fatalf("Mismatched length, got %d, wanted %d\n", int(length), len(testBytes))
	}

	if !(before.Sec <= sec && sec <= after.Sec) {
		t.Fatalf("Seconds Mismatch, got %d, wanted between %d and %d\n", sec, before.Sec, after.Sec)
	}

	if !(before.Usec <= usec && usec <= after.Usec) {
		t.Fatalf("Microseconds Mismatch, got %d, wanted between %d and %d\n", usec, before.Usec, after.Usec)
	}
}
