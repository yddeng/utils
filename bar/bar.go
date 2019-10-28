/*
 * 进度条
 */
package bar

import (
	"fmt"
	"github.com/yddeng/dutil/dstring"
	"sync"
	"time"
)

type Bar struct {
	mu      sync.Mutex
	name    string
	total   int
	current int

	add      int
	lastTime time.Time
	done     chan struct{}
}

var (
	tickDur  = time.Millisecond * 100
	barWidth = 50 //进度条宽度
	lbarChar string
	rbarChar string
)

func init() {
	for i := 0; i < barWidth; i++ {
		lbarChar += "="
		rbarChar += "-"
	}
}

func NewBar(name string, total int) *Bar {
	if total <= 0 {
		return nil
	}

	now := time.Now()
	bar := &Bar{
		name:     name,
		total:    total,
		lastTime: now,
		done:     make(chan struct{}),
	}

	go bar.run()

	return bar
}

func (b *Bar) Add(count int) {

	b.mu.Lock()
	defer b.mu.Unlock()

	b.current += count
	b.add += count

}

func (b *Bar) run() {
	ticker := time.NewTicker(tickDur)
	for {
		select {
		case now := <-ticker.C:
			b.print(now)
		case <-b.done:
			ticker.Stop()
			return
		}
	}
}

func (b *Bar) print(now time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	tmp := "%s [%s%s]%s %s" //name, lbar, rbar, rate, speed

	rate := b.current * 100 / b.total
	lbw := barWidth * b.current / b.total
	rbw := barWidth - lbw
	speed := int(float64(b.add) / now.Sub(b.lastTime).Seconds())

	var lbar, rbar = "", ""
	lbar = string(lbarChar[:lbw])
	if rbw > 0 {
		rbar = ">" + string(rbarChar[:rbw])
	}

	rateStr := fmt.Sprintf("%2d%%", rate)
	speedStr := fmt.Sprintf("%s/s", dstring.ByteSizeFromat(int64(speed)))
	txt := fmt.Sprintf(tmp, b.name, lbar, rbar, rateStr, speedStr)
	fmt.Printf("%s\r", txt)

	if b.current >= b.total {
		close(b.done)
	}
	b.add = 0
	b.lastTime = now
}
