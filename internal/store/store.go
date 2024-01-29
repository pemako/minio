// Copyright (c) 2015-2023 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	xioutil "github.com/minio/minio/internal/ioutil"
	xnet "github.com/minio/pkg/v2/net"
)

const (
	retryInterval = 3 * time.Second
)

type logger = func(ctx context.Context, err error, id string, errKind ...interface{})

// ErrNotConnected - indicates that the target connection is not active.
var ErrNotConnected = errors.New("not connected to target server/service")

// Target - store target interface
type Target interface {
	Name() string
	SendFromStore(key Key) error
}

// Store - Used to persist items.
type Store[I any] interface {
	Put(item I) error
	Get(key string) (I, error)
	Len() int
	List() ([]string, error)
	Del(key string) error
	DelList(key []string) error
	Open() error
	Extension() string
}

// Key denotes the key present in the store.
type Key struct {
	Name   string
	IsLast bool
}

// replayItems - Reads the items from the store and replays.
func replayItems[I any](store Store[I], doneCh <-chan struct{}, log logger, id string) <-chan Key {
	keyCh := make(chan Key)

	go func() {
		defer xioutil.SafeClose(keyCh)

		retryTicker := time.NewTicker(retryInterval)
		defer retryTicker.Stop()

		for {
			names, err := store.List()
			if err != nil {
				log(context.Background(), fmt.Errorf("store.List() failed with: %w", err), id)
			} else {
				keyCount := len(names)
				for i, name := range names {
					select {
					case keyCh <- Key{strings.TrimSuffix(name, store.Extension()), keyCount == i+1}:
					// Get next key.
					case <-doneCh:
						return
					}
				}
			}

			select {
			case <-retryTicker.C:
			case <-doneCh:
				return
			}
		}
	}()

	return keyCh
}

// sendItems - Reads items from the store and re-plays.
func sendItems(target Target, keyCh <-chan Key, doneCh <-chan struct{}, logger logger) {
	retryTicker := time.NewTicker(retryInterval)
	defer retryTicker.Stop()

	send := func(key Key) bool {
		for {
			err := target.SendFromStore(key)
			if err == nil {
				break
			}

			if err != ErrNotConnected && !xnet.IsConnResetErr(err) {
				logger(context.Background(),
					fmt.Errorf("target.SendFromStore() failed with '%w'", err),
					target.Name())
			}

			// Retrying after 3secs back-off

			select {
			case <-retryTicker.C:
			case <-doneCh:
				return false
			}
		}
		return true
	}

	for {
		select {
		case key, ok := <-keyCh:
			if !ok {
				// closed channel.
				return
			}

			if !send(key) {
				return
			}
		case <-doneCh:
			return
		}
	}
}

// StreamItems reads the keys from the store and replays the corresponding item to the target.
func StreamItems[I any](store Store[I], target Target, doneCh <-chan struct{}, logger logger) {
	go func() {
		// Replays the items from the store.
		keyCh := replayItems(store, doneCh, logger, target.Name())
		// Send items from the store.
		sendItems(target, keyCh, doneCh, logger)
	}()
}
