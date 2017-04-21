package main

import (
	"bytes"
	"io"
)

type delayedReadSeeker struct {
	rc  io.ReadCloser
	buf *bytes.Reader
}

func newDelayedReadSeeker(rc io.ReadCloser) io.ReadSeeker {

	drs := delayedReadSeeker{}
	drs.rc = rc

	return &drs
}

func (drs *delayedReadSeeker) assert() {

	if drs.buf == nil {
		var buf bytes.Buffer
		io.Copy(&buf, drs.rc)
		drs.buf = bytes.NewReader((&buf).Bytes())
	}
}

func (drs *delayedReadSeeker) Read(p []byte) (n int, err error) {
	drs.assert()
	return drs.buf.Read(p)
}

func (drs *delayedReadSeeker) Seek(offset int64, whence int) (int64, error) {
	drs.assert()
	return drs.buf.Seek(offset, whence)
}
