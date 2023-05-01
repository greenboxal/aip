package forddb

import (
	"io"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

type logStream struct {
	LogIterator

	id LogStreamID

	fa   *os.File
	fb   *os.File
	next *os.File
}

func NewLogStream(ls LogStore, id LogStreamID, fa, fb *os.File) (LogStream, error) {
	l := &logStream{
		LogIterator: ls.Iterator(WithBlockingIterator()),

		id: id,
		fa: fa,
		fb: fb,
	}

	a, err := l.readCheckpoint(fa)

	if err != nil {
		return nil, err
	}

	b, _ := l.readCheckpoint(fa)

	if err != nil {
		return nil, err
	}

	if a.Clock > b.Clock {
		l.Reset(a)
	} else {
		l.Reset(b)
	}

	return l, nil
}

func (l *logStream) StreamID() LogStreamID {
	return l.id
}

func (l *logStream) SaveCheckpoint() error {
	f := l.nextFile()

	if err := f.Truncate(0); err != nil {
		return err
	}

	serialized, err := bson.Marshal(l.CurrentLsn())

	if err != nil {
		return err
	}

	if _, err := f.Write(serialized); err != nil {
		return err
	}

	return f.Sync()
}

func (l *logStream) Close() error {
	if err := l.fa.Close(); err != nil {
		return err
	}

	if err := l.fb.Close(); err != nil {
		return err
	}

	return l.LogIterator.Close()
}

func (l *logStream) readCheckpoint(f *os.File) (lsn LSN, err error) {
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
		err = bson.Unmarshal(data, &lsn)
	}

	return
}

func (l *logStream) nextFile() *os.File {
	if l.next == l.fa {
		l.next = l.fb
	} else {
		l.next = l.fa
	}

	return l.next
}
