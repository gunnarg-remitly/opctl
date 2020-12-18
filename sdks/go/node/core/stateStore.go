package core

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/dgraph-io/badger/v2"
	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o internal/fakes/stateStore.go . stateStore
// stateStore allows efficiently querying the current state of opctl.
//
// State is materialized by applying events in the order in which they are/were received.
//
// efficient startup:
// A lastAppliedEventTimestamp is maintained and used at startup to pickup applying events
// from where we left off.
type stateStore interface {
	// TryGetCreds returns creds for a ref if any exist
	TryGetAuth(resource string) *model.Auth

	AddAuth(req model.AuthAdded) error
}

func newStateStore(
	ctx context.Context,
	db *badger.DB,
) stateStore {
	return &_stateStore{
		authsByResourcesKeyPrefix:    "authsByResources_",
		callsByID:                    make(map[string]*model.Call),
		db:                           db,
		lastAppliedEventTimestampKey: "lastAppliedEventTimestamp",
	}
}

type _stateStore struct {
	lastAppliedEventTimestampKey string
	authsByResourcesKeyPrefix    string
	callsByID                    map[string]*model.Call
	db                           *badger.DB
	// synchronize access via mutex
	mux sync.RWMutex
}

func (ss *_stateStore) AddAuth(authAdded model.AuthAdded) error {
	return ss.db.Update(func(txn *badger.Txn) error {
		auth := authAdded.Auth
		encodedAuth, err := json.Marshal(auth)
		if nil != err {
			return err
		}

		return txn.Set(
			[]byte(ss.authsByResourcesKeyPrefix+strings.ToLower(auth.Resources)),
			encodedAuth,
		)
	})
}

func (ss *_stateStore) TryGetAuth(
	ref string,
) *model.Auth {
	ref = strings.ToLower(ref)
	var auth *model.Auth
	ss.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefixBytes := []byte(ss.authsByResourcesKeyPrefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			key := string(item.Key())
			prefix := strings.TrimPrefix(key, ss.authsByResourcesKeyPrefix)

			if strings.HasPrefix(ref, prefix) {
				item.Value(func(value []byte) error {
					auth = &model.Auth{}
					return json.Unmarshal(value, auth)
				})
			}
		}
		return nil
	})

	return auth
}
