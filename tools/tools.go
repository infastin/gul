package tools

import (
	"github.com/infastin/gul/pkg/gmu"
	"runtime"
	"sync"
)

func Parallelize(procs, start, end, step int, fn func(start, end int)) {
	if procs == 1 || procs < 0 {
		fn(start, end)
		return
	}

	if step == 0 {
		return
	}

	if procs == 0 {
		procs = runtime.NumCPU()
	}

	var wg sync.WaitGroup
	SplitRange(start, end, step, procs, func(pstart, pend int) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(pstart, pend)
		}()
	})
	wg.Wait()
}

func SplitRange(start, end, step, n int, fn func(start, end int)) {
	if n == 0 || step == 0 {
		return
	}

	count := end - start
	steps := count / step

	if steps < 1 {
		return
	}

	if n > steps {
		n = steps
	}

	div := steps / n
	mod := steps % n

	for i := 0; i < n; i++ {
		fn(
			start+i*div+gmu.MinInt(i, mod),
			start+(i+step)*div+gmu.MinInt(i+step, mod),
		)
	}
}
