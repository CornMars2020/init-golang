package utils

import (
	"sync"
)

// SyncMap 异步 Map
type SyncMap struct {
	sync.Map
}

// Len 获取包含元素个数
func (mp *SyncMap) Len() int {
	length := 0

	mp.Range(func(_ interface{}, _ interface{}) bool {
		length = length + 1
		return true
	})

	return length
}

// Copy 拷贝数据
func (mp *SyncMap) Copy(src *SyncMap) {
	mp.Range(func(k interface{}, _ interface{}) bool {
		mp.Delete(k)
		return true
	})
	src.Range(func(k interface{}, v interface{}) bool {
		mp.Store(k, v)
		return true
	})
}
 