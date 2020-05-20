package ttyrec

import (
	"encoding/binary"
	"os"
	"sync"
	"time"
)

type TTYRecorder struct {
	wr  *os.File
	mux sync.Mutex
}

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

func NewTTYRecorder(fp string) (*TTYRecorder, error) {
	f, err := os.Create(fp)
	if err != nil {
		return nil, err
	}
	tr := &TTYRecorder{
		wr: f,
	}
	return tr, nil
}

func (tr *TTYRecorder) Close() error {
	return tr.wr.Close()
}

func (tr *TTYRecorder) writeBytes(data []byte) (int, error) {
	t := time.Now()
	timeval := NanosToTimeval(t.UnixNano())
	writeLen := len(data)
	err := binary.Write(tr.wr, binary.LittleEndian, timeval.Sec)
	if err != nil {
		return -1, err
	}
	err = binary.Write(tr.wr, binary.LittleEndian, timeval.Usec)
	if err != nil {
		return -1, err
	}
	err = binary.Write(tr.wr, binary.LittleEndian, int32(writeLen))
	if err != nil {
		return -1, err
	}
	bw, err := tr.wr.Write(data)
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
