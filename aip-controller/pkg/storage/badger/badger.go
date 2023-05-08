package badger

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/config"
	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type StorageConfig struct {
}

type Storage struct {
	db *badger.DB
}

func NewStorage(
	lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	rsm *config.ResourceManager,
) *Storage {
	path := rsm.GetDataDirectory("ford/forddb/badger")

	options := badger.DefaultOptions(path)

	options.Logger = &ZapEventLogger{SugaredLogger: logger.Named("badger")}

	db, err := badger.Open(options)

	if err != nil {
		log.Fatal(err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return &Storage{
		db: db,
	}
}

func (s *Storage) List(ctx context.Context, typ forddb2.TypeID) ([]forddb2.BasicResource, error) {
	var result []forddb2.BasicResource

	tx := s.db.NewTransaction(false)
	defer tx.Discard()

	it := tx.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(fmt.Sprintf("ford/%s/", typ.Name()))

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		value, err := item.ValueCopy(nil)

		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(value)

		resource, err := forddb2.DeserializeFrom(reader, forddb2.Json)

		if err != nil {
			return nil, err
		}

		result = append(result, resource)
	}

	return result, nil
}

func (s *Storage) Get(ctx context.Context, typ forddb2.TypeID, id forddb2.BasicResourceID) (forddb2.BasicResource, error) {
	tx := s.db.NewTransaction(false)
	defer tx.Discard()

	key := []byte(fmt.Sprintf("ford/%s/%s", typ.Name(), id.String()))

	item, err := tx.Get(key)

	if err == badger.ErrKeyNotFound {
		return nil, forddb2.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	value, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	if len(value) == 0 {
		return nil, forddb2.ErrNotFound
	}

	resource, err := forddb2.Deserialize(value, forddb2.Json)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *Storage) Put(ctx context.Context, resource forddb2.BasicResource) (forddb2.BasicResource, error) {
	for {
		res, _, err := s.doPut(resource)

		if err == badger.ErrConflict {
			continue
		} else if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func (s *Storage) doPut(resource forddb2.BasicResource) (forddb2.BasicResource, forddb2.BasicResource, error) {
	var existing forddb2.BasicResource

	tx := s.db.NewTransaction(true)

	typ := resource.GetResourceTypeID()
	id := resource.GetResourceBasicID()
	key := []byte(fmt.Sprintf("ford/%s/%s", typ.Name(), id.String()))

	item, err := tx.Get(key)

	if err = badger.ErrKeyNotFound; err != nil {
		item = nil
	} else if err != nil {
		return nil, nil, err
	}

	if item != nil {
		value, err := item.ValueCopy(nil)

		if err != nil {
			return nil, nil, err
		}

		existing, err = forddb2.Deserialize(value, forddb2.Json)

		if err != nil {
			return nil, nil, err
		}

		if existing.GetResourceVersion() != resource.GetResourceVersion() {
			return nil, existing, forddb2.ErrVersionMismatch
		}
	}

	buffer := bytes.NewBuffer(nil)

	if err := forddb2.SerializeTo(buffer, forddb2.Json, resource); err != nil {
		return nil, existing, err
	}

	if err := tx.Set(key, buffer.Bytes()); err != nil {
		return nil, existing, err
	}

	if err := tx.Commit(); err != nil {
		return nil, existing, err
	}

	return resource, existing, nil
}

func (s *Storage) Delete(ctx context.Context, resource forddb2.BasicResource) (forddb2.BasicResource, error) {
	for {
		res, err := s.doDelete(resource)

		if err == badger.ErrConflict {
			continue
		} else if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func (s *Storage) doDelete(resource forddb2.BasicResource) (forddb2.BasicResource, error) {
	tx := s.db.NewTransaction(true)

	typ := resource.GetResourceTypeID()
	id := resource.GetResourceBasicID()
	key := []byte(fmt.Sprintf("ford/%s/%s", typ.Name(), id.String()))

	item, err := tx.Get(key)

	if err = badger.ErrKeyNotFound; err != nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	value, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}

	existing, err := forddb2.Deserialize(value, forddb2.Json)

	if err != nil {
		return nil, err
	}

	if existing.GetResourceVersion() != resource.GetResourceVersion() {
		return nil, forddb2.ErrVersionMismatch
	}

	if err := tx.Delete(key); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

type ZapEventLogger struct {
	*zap.SugaredLogger
}

func (el *ZapEventLogger) Warningf(format string, args ...interface{}) {
	el.SugaredLogger.Warnf(format, args...)
}

func (el *ZapEventLogger) Errorf(f string, v ...interface{}) {
	el.SugaredLogger.Errorf(f, v...)
}

func (el *ZapEventLogger) Infof(f string, v ...interface{}) {
}

func (el *ZapEventLogger) Debugf(f string, v ...interface{}) {
}
