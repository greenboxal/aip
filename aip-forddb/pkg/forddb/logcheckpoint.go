package forddb

import (
	"context"
	"encoding/json"
	"io"
	"os"
)

type LogCheckpoint interface {
	Load(ctx context.Context, id LogStreamID) (LSN, error)
	Save(ctx context.Context, id LogStreamID, lsn LSN) error
	Close() error
}

type fileCheckpoint struct {
	id   LogStreamID
	fa   *os.File
	fb   *os.File
	next *os.File
}

func (l *fileCheckpoint) Load(ctx context.Context, id LogStreamID) (LSN, error) {
	return l.readCheckpoint(l.nextFile())
}

func NewFileCheckpoint(id LogStreamID, fa, fb *os.File) (LogCheckpoint, error) {
	l := &fileCheckpoint{
		id: id,
		fa: fa,
		fb: fb,
	}

	a, err := l.readCheckpoint(fa)

	if err != nil {
		return nil, err
	}

	b, _ := l.readCheckpoint(fb)

	if err != nil {
		return nil, err
	}

	if a.Clock > b.Clock {
		l.next = fb
	} else {
		l.next = fa
	}

	return l, nil
}

func (l *fileCheckpoint) Save(ctx context.Context, id LogStreamID, lsn LSN) error {
	f := l.nextFile()

	if err := f.Truncate(0); err != nil {
		return err
	}

	serialized, err := json.Marshal(lsn)

	if err != nil {
		return err
	}

	if _, err := f.Write(serialized); err != nil {
		return err
	}

	return f.Sync()
}

func (l *fileCheckpoint) Close() error {
	if err := l.fa.Close(); err != nil {
		return err
	}

	if err := l.fb.Close(); err != nil {
		return err
	}

	return nil
}

func (l *fileCheckpoint) readCheckpoint(f *os.File) (lsn LSN, err error) {
	var data []byte

	if _, err = f.Stat(); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return
	}

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return
	}

	data, err = io.ReadAll(f)

	if err != nil {
		if err == io.EOF {
			err = nil
		}

		return
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, &lsn)
	}

	return
}

func (l *fileCheckpoint) nextFile() *os.File {
	next := l.next

	if l.next == l.fa {
		l.next = l.fb
	} else {
		l.next = l.fa
	}

	return next
}
