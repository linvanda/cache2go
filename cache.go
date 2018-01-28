/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"sync"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

// Cache returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
func Cache(table string) *CacheTable {
	// 这里加的是读锁，大部分情况下 table 已经存在，不会走到 if 里面，
	// cache 属于读多写少的场景，需使用 RWM	utex 加锁。
	// (写锁加锁成功的条件：没有锁定或只有读锁定时。如果有写锁或等待中的写锁，则无法加读锁)
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		// 加写锁。此时该锁可能已经被其他协程获得，此处将进入阻塞（有可能被换入等待队列中）。
		// 因而，此处加锁成功后需要再次判断 cache[table] 是否已经存在
		mutex.Lock()
		t, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			t = &CacheTable{
				name:  table,
				// 注意：map 不能使用其零值初始化（零值为 nil）
				items: make(map[interface{}]*CacheItem),
			}
			cache[table] = t
		}
		mutex.Unlock()
	}

	return t
}
