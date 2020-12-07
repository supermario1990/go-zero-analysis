package collection

import (
	"fmt"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/timex"
)

type (
	// 滚动窗口选项函数
	RollingWindowOption func(rollingWindow *RollingWindow)

	// 滚动窗口定义
	RollingWindow struct {
		lock          sync.RWMutex	// 读写锁
		size          int			// 滚动窗口大小
		win           *window		// 滚动窗口，一个window有size个Bucket（桶）
		interval      time.Duration	// 间隔时间
		offset        int			// 偏移位置，用来定位当前所在的Bucket（桶）
		ignoreCurrent bool			// 是否忽略当前的Bucket
		lastTime      time.Duration	// 最近时间，通过从lastTime开始到现在流逝的时间和interval，就能确定现在的offset
	}
)

// 创建滚动窗口
func NewRollingWindow(size int, interval time.Duration, opts ...RollingWindowOption) *RollingWindow {
	if size < 1 {
		panic("size must be greater than 0")
	}

	w := &RollingWindow{
		size:     size,
		win:      newWindow(size),
		interval: interval,
		lastTime: timex.Now(),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// 向当前offset指向的bucket放入v
func (rw *RollingWindow) Add(v float64) {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	rw.updateOffset()
	rw.win.add(rw.offset, v)
}

// 将fn分别作用于滚动窗口中没有过期的bucket
func (rw *RollingWindow) Reduce(fn func(b *Bucket)) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	var diff int
	span := rw.span()
	// ignore current bucket, because of partial data
	if span == 0 && rw.ignoreCurrent {
		diff = rw.size - 1
		fmt.Println("----------------------")
	} else {
		diff = rw.size - span
	}
	if diff > 0 {
		offset := (rw.offset + span + 1) % rw.size
		rw.win.reduce(offset, diff, fn)
		fmt.Printf("offset: %v, diff: %v, span: %v\n", offset, diff, span)
	}
}

// 计算跨度
func (rw *RollingWindow) span() int {
	offset := int(timex.Since(rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	} else {
		return rw.size
	}
}

// 更新offset值，重置过期的buckets，更新lastTime
func (rw *RollingWindow) updateOffset() {
	span := rw.span()
	if span <= 0 {
		return
	}

	offset := rw.offset
	start := offset + 1
	steps := start + span
	var remainder int
	if steps > rw.size {
		remainder = steps - rw.size
		steps = rw.size
	}

	// reset expired buckets
	for i := start; i < steps; i++ {
		rw.win.resetBucket(i)
	}
	for i := 0; i < remainder; i++ {
		rw.win.resetBucket(i)
	}

	rw.offset = (offset + span) % rw.size
	rw.lastTime = timex.Now()
}

func (rw *RollingWindow) getOffset() int{
	return rw.offset
}

// 桶
type Bucket struct {
	Sum   float64	// 计算桶里数之和
	Count int64		// 计算桶里有多少个数
}

func (b *Bucket) add(v float64) {
	b.Sum += v
	b.Count++
}

func (b *Bucket) reset() {
	b.Sum = 0
	b.Count = 0
}

type window struct {
	buckets []*Bucket
	size    int
}

// 创建窗口，同时初始化Buckets
func newWindow(size int) *window {
	buckets := make([]*Bucket, size)
	for i := 0; i < size; i++ {
		buckets[i] = new(Bucket)
	}
	return &window{
		buckets: buckets,
		size:    size,
	}
}

func (w *window) add(offset int, v float64) {
	w.buckets[offset%w.size].add(v)
}

func (w *window) reduce(start, count int, fn func(b *Bucket)) {
	for i := 0; i < count; i++ {
		fn(w.buckets[(start+i)%w.size])
	}
}

func (w *window) resetBucket(offset int) {
	w.buckets[offset%w.size].reset()
}

func IgnoreCurrentBucket() RollingWindowOption {
	return func(w *RollingWindow) {
		w.ignoreCurrent = true
	}
}
