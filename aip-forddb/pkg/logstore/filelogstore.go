package logstore

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/tidwall/wal"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type FileLogStore struct {
	m    sync.RWMutex
	cond *sync.Cond

	logDir string

	log     *wal.Log
	lastLsn forddb.LSN
}

func NewFileLogStore(baseDir string) (*FileLogStore, error) {
	fls := &FileLogStore{
		logDir: baseDir,
	}

	options := *wal.DefaultOptions
	walPath := path.Join(baseDir, "wal")
	log, err := wal.Open(walPath, &options)

	if err != nil {
		return nil, err
	}

	last, err := log.LastIndex()

	if err != nil {
		return nil, err
	}

	fls.log = log
	fls.lastLsn.Clock = last
	fls.cond = sync.NewCond(&fls.m)

	return fls, nil
}

func (fls *FileLogStore) OpenStream(id forddb.LogStreamID) (forddb.LogStream, error) {
	idStr := string(id)
	idStr = strings.Replace(idStr, "/", "_", -1)
	idStr = strings.Replace(idStr, ".", "_", -1)
	id = forddb.LogStreamID(idStr)

	faName := path.Join(fls.logDir, fmt.Sprintf("checkpoint-%s-a", id))
	fbName := path.Join(fls.logDir, fmt.Sprintf("checkpoint-%s-b", id))

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

func (fls *FileLogStore) Append(ctx context.Context, log forddb.LogEntry) (forddb.LogEntryRecord, error) {
	var record forddb.LogEntryRecord

	fls.m.Lock()
	defer fls.m.Unlock()

	record.LogEntry = log
	record.LSN.Clock = fls.lastLsn.Clock + 1
	record.LSN.TS = time.Now()

	data, err := ipld.Encode(typesystem.Wrap(record), dagjson.Encode)

	if err != nil {
		return forddb.LogEntryRecord{}, err
	}

	if err := fls.log.Write(record.Clock, data); err != nil {
		return forddb.LogEntryRecord{}, err
	}

	if err := fls.log.Sync(); err != nil {
		return forddb.LogEntryRecord{}, err
	}

	fls.lastLsn.Clock = record.LSN.Clock

	fls.cond.Broadcast()

	return record, nil
}

func (fls *FileLogStore) Iterator(options ...forddb.LogIteratorOption) forddb.LogIterator {
	iterator := newFileLogIterator(fls, forddb.NewLogIteratorOptions(options...))

	if err := iterator.SetLSN(context.Background(), fls.lastLsn); err != nil {
		panic(err)
	}

	return iterator
}

func (fls *FileLogStore) Close() error {
	return fls.log.Close()
}
