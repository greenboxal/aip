package logstore

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

const FileSegmentHeaderMagic = 0x46534547 // "FSEG"
const FileSegmentHeaderSize = 64          // "FSEG"
const FileSegmentBaseSeekToHead = 0xFFFFFFFFFFFFFFFF

type OldFileLogStore struct {
	m    sync.RWMutex
	cond *sync.Cond

	logDir string

	currentSegment *fileStoreSegment
}

func (fls *OldFileLogStore) OpenStream(id forddb.LogStreamID) (forddb.LogStream, error) {
	faName := path.Join(fls.logDir, fmt.Sprintf("%s-1", id))
	fbName := path.Join(fls.logDir, fmt.Sprintf("%s-2", id))

	fa, err := os.OpenFile(faName, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	fb, err := os.OpenFile(fbName, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	checkPoint, err := forddb.NewFileCheckpoint(id, fa, fb)

	if err != nil {
		return nil, err
	}

	return forddb.NewLogStream(fls, id, checkPoint)
}

func NewOldFileLogStore(baseDir string) (*OldFileLogStore, error) {
	fls := &OldFileLogStore{
		logDir: baseDir,
	}

	head := forddb.MakeLSN(FileSegmentBaseSeekToHead, time.Now())
	s, err := fls.openSegment(head, false, false)

	if err != nil {
		return nil, err
	}

	fls.cond = sync.NewCond(&fls.m)
	fls.currentSegment = s

	return fls, nil
}

func (fls *OldFileLogStore) Append(ctx context.Context, log forddb.LogEntry) (forddb.LogEntryRecord, error) {
	var record forddb.LogEntryRecord

	record.LogEntry = log

	if err := fls.currentSegment.Append(&record); err != nil {
		return record, nil
	}

	fls.cond.Broadcast()

	return record, nil
}

func (fls *OldFileLogStore) Iterator(options ...forddb.LogIteratorOption) forddb.LogIterator {
	return newOldFileLogIterator(fls, forddb.NewLogIteratorOptions(options...))
}

func (fls *OldFileLogStore) Close() error {
	fls.m.Lock()

	if fls.currentSegment != nil {
		if err := fls.currentSegment.Close(); err != nil {
			return nil
		}
	}

	return nil
}

func (fls *OldFileLogStore) openSegment(
	base forddb.LSN,
	readOnly bool,
	createExclusive bool,
) (*fileStoreSegment, error) {
	var flags int

	name := fls.getSegmentFilePath(0)

	if readOnly {
		flags |= os.O_RDONLY
	} else {
		flags |= os.O_RDWR | os.O_CREATE

		if createExclusive {
			flags |= os.O_EXCL
		}

		flags |= os.O_SYNC
	}

	f, err := os.OpenFile(name, flags, 0644)

	if err != nil {
		return nil, err
	}

	return newFileStoreSegment(f, base, readOnly)
}

func (fls *OldFileLogStore) getSegmentFilePath(seq int) string {
	return path.Join(fls.logDir, fmt.Sprintf("log-%08d.dat", seq))
}

type fileStoreSegment struct {
	sync.RWMutex

	file     *os.File
	readOnly bool

	head   forddb.LSN
	tail   forddb.LSN
	offset uint64

	header [FileSegmentHeaderSize]byte
}

func newFileStoreSegment(f *os.File, base forddb.LSN, readOnly bool) (*fileStoreSegment, error) {
	fss := &fileStoreSegment{
		head:     base,
		file:     f,
		readOnly: readOnly,
	}

	size, err := f.Seek(0, io.SeekEnd)

	if err != nil {
		return nil, err
	}

	if size > 0 {
		n, err := f.ReadAt(fss.header[:], 0)

		if err != nil {
			return nil, err
		}

		if n == 0 {
			return nil, errors.New("empty segment file")
		}

		magic := binary.LittleEndian.Uint32(fss.header[0:])
		baseClock := binary.LittleEndian.Uint64(fss.header[8:])
		baseTs := binary.LittleEndian.Uint64(fss.header[16:])
		tailClock := binary.LittleEndian.Uint64(fss.header[24:])
		tailTs := binary.LittleEndian.Uint64(fss.header[32:])

		if magic != FileSegmentHeaderMagic {
			return nil, fmt.Errorf("invalid segment header magic")
		}

		fss.head.Clock = baseClock
		fss.head.TS = time.Unix(0, int64(baseTs))
		fss.tail.Clock = tailClock
		fss.tail.TS = time.Unix(0, int64(tailTs))
	} else if !readOnly {
		if err := fss.updateHeader(uint64(len(fss.header)), base); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("empty segment file")
	}

	return fss, nil
}

func (fss *fileStoreSegment) Seek(lsn forddb.LSN) error {
	return nil
}

func (fss *fileStoreSegment) Read(current *forddb.LogEntryRecord) (*forddb.LogEntryRecord, error) {
	return nil, nil
}

func (fss *fileStoreSegment) Append(record *forddb.LogEntryRecord) error {
	if fss.readOnly {
		return errors.New("segment is read only")
	}

	fss.Lock()
	defer fss.Unlock()

	record.LSN.Clock = fss.tail.Clock + 1
	record.LSN.TS = time.Now()

	offset, err := fss.file.Seek(0, io.SeekEnd)

	if err != nil {
		return err
	}

	buffer := make([]byte, 4)
	buffer, err = bson.MarshalAppend(buffer, record)

	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(buffer[0:4], uint32(len(buffer)))

	if err != nil {
		return err
	}

	if _, err := fss.file.Write(buffer); err != nil {
		return err
	}

	if err := fss.updateHeader(uint64(offset), record.LSN); err != nil {
		return err
	}

	return nil
}

func (fss *fileStoreSegment) Close() error {
	fss.Lock()
	defer fss.Unlock()

	if !fss.readOnly {
		// Flush header
		if err := fss.updateHeader(fss.offset, fss.tail); err != nil {
			return err
		}
	}

	if err := fss.file.Close(); err != nil {
		return err
	}

	return nil
}

func (fss *fileStoreSegment) updateHeader(lastOffset uint64, tail forddb.LSN) error {
	if fss.readOnly {
		return errors.New("segment is read only")
	}

	binary.LittleEndian.PutUint32(fss.header[0:], uint32(FileSegmentHeaderMagic))
	binary.LittleEndian.PutUint64(fss.header[8:], fss.head.Clock)
	binary.LittleEndian.PutUint64(fss.header[16:], uint64(fss.head.TS.UnixMilli()))
	binary.LittleEndian.PutUint64(fss.header[24:], tail.Clock)
	binary.LittleEndian.PutUint64(fss.header[32:], uint64(tail.TS.UnixMilli()))
	binary.LittleEndian.PutUint64(fss.header[40:], lastOffset)

	_, err := fss.file.WriteAt(fss.header[:], 0)

	if err != nil {
		return err
	}

	if err := fss.file.Sync(); err != nil {
		return err
	}

	fss.tail = tail
	fss.offset = lastOffset

	return nil
}
