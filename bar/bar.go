/*
 * 进度条
 */
package bar

import (
	"fmt"
	"sync"
	"time"
)

var (
	ts = []string{"b", "Kb", "Mb", "Gb"}
)

// 字节长度格式化输出
// 例：2566b -> 2.50Kb
func byteSize(b int) string {
	n := float64(b)
	i := 0
	for n > 1024 {
		n /= 1024
		i++
		if i == len(ts) {
			break
		}
	}
	return fmt.Sprintf("%.2f%s", n, ts[i])
}

type Bar struct {
	mu        sync.Mutex
	name      string
	total     int
	current   int
	startTime time.Time
}

var (
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
		panic("total <= 0")
	}
	return &Bar{
		name:      name,
		total:     total,
		startTime: time.Now(),
	}
}

func (b *Bar) Add(count int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	b.current += count
	rate := b.current * 100 / b.total
	subTime := now.Sub(b.startTime)
	speed := int(float64(b.current)/float64(subTime.Milliseconds())) * 1000
	if speed < 0 {
		speed = 0
	}

	tmp := "%s [%s%s]%s %s " //name, lbar, rbar, rate, speed
	lbw := barWidth * b.current / b.total
	rbw := barWidth - lbw
	var lbar, rbar = "", ""
	lbar = lbarChar[:lbw]
	if rbw > 0 {
		rbar = ">" + rbarChar[:rbw]
	}

	// print
	rateStr := fmt.Sprintf("%2d%%", rate)
	speedStr := fmt.Sprintf("%10s/s", byteSize(speed))
	txt := fmt.Sprintf(tmp, b.name, lbar, rbar, rateStr, speedStr)
	fmt.Printf("%s\r", txt)
	if b.current >= b.total {
		fmt.Println()
	}
}
