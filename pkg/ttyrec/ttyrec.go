package ttyrec

import (
	"encoding/binary"
	"io"
	"sync"
	"time"
)

// TTYRecorder is a wrapper around a Writer that implements ttyrec formatting
type TTYRecorder struct {
	wr  io.Writer
	mux sync.Mutex
}

// Timeval is a struct compatible with what is normally returned by NsecToTimeval
type Timeval struct {
	Sec  int32
	Usec int32
}

// NanosToTimeval is roughly compatible with the NsecToTimeval syscall without a syscall
func NanosToTimeval(ns int64) *Timeval {
	tMicros := ns / 1000
	return &Timeval{
		Sec:  int32(tMicros / 1E6),
		Usec: int32(tMicros % 1E6),
	}
}

// NewTTYRecorder instantiates a TTYRecorder wrapped around the writer w
func NewTTYRecorder(w io.Writer) *TTYRecorder {
	tr := &TTYRecorder{
		wr: w,
	}
	return tr
}

func (tr *TTYRecorder) writeBytes(data []byte) (int, error) {
	t := time.Now()
	timeval := NanosToTimeval(t.UnixNano())
	writeLen := len(data)
	headerBuff := make([]byte, 12)
	binary.LittleEndian.PutUint32(headerBuff[0:4], uint32(timeval.Sec))
	binary.LittleEndian.PutUint32(headerBuff[4:8], uint32(timeval.Usec))
	binary.LittleEndian.PutUint32(headerBuff[8:12], uint32(writeLen))
	writeData := append(headerBuff, data...)
	bw, err := tr.wr.Write(writeData)
	if err != nil {
		return -1, err
	}
	return bw, nil
}

func (tr *TTYRecorder) Write(data []byte) (int, error) {
	tr.mux.Lock()
	bw, err := tr.writeBytes(data)
	tr.mux.Unlock()
	if err != nil {
		return -1, err
	}
	return bw, nil
}
