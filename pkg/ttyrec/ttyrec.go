package ttyrec

import (
	"encoding/binary"
	"os"
	"sync"
	"syscall"
	"time"
)

type TTYRecorder struct {
	wr  *os.File
	mux sync.Mutex
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
	timeval := syscall.NsecToTimeval(t.UnixNano())
	writeLen := len(data)
	err := binary.Write(tr.wr, binary.LittleEndian, int32(timeval.Sec))
	if err != nil {
		return -1, err
	}
	err = binary.Write(tr.wr, binary.LittleEndian, int32(timeval.Usec))
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
