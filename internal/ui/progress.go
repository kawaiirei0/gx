package ui

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// ProgressBar 进度条显示器
type ProgressBar struct {
	writer      io.Writer
	total       int64
	current     int64
	width       int
	prefix      string
	startTime   time.Time
	lastUpdate  time.Time
	updateDelay time.Duration
}

// NewProgressBar 创建新的进度条
func NewProgressBar(writer io.Writer, total int64, prefix string) *ProgressBar {
	return &ProgressBar{
		writer:      writer,
		total:       total,
		current:     0,
		width:       40,
		prefix:      prefix,
		startTime:   time.Now(),
		lastUpdate:  time.Time{},
		updateDelay: 100 * time.Millisecond, // 限制更新频率
	}
}

// Update 更新进度
func (pb *ProgressBar) Update(current int64) {
	pb.current = current

	// 限制更新频率，避免闪烁
	now := time.Now()
	if now.Sub(pb.lastUpdate) < pb.updateDelay && current < pb.total {
		return
	}
	pb.lastUpdate = now

	pb.render()
}

// Finish 完成进度条
func (pb *ProgressBar) Finish() {
	pb.current = pb.total
	pb.render()
	fmt.Fprintln(pb.writer) // 换行
}

// render 渲染进度条
func (pb *ProgressBar) render() {
	if pb.total <= 0 {
		return
	}

	percent := float64(pb.current) / float64(pb.total)
	if percent > 1.0 {
		percent = 1.0
	}

	// 计算进度条填充
	filled := int(float64(pb.width) * percent)
	if filled > pb.width {
		filled = pb.width
	}

	// 构建进度条
	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.width-filled)

	// 计算速度和剩余时间
	elapsed := time.Since(pb.startTime)
	var speed float64
	var eta time.Duration

	if elapsed.Seconds() > 0 {
		speed = float64(pb.current) / elapsed.Seconds()
		if speed > 0 && pb.current < pb.total {
			remaining := pb.total - pb.current
			eta = time.Duration(float64(remaining)/speed) * time.Second
		}
	}

	// 格式化输出
	fmt.Fprintf(pb.writer, "\r%s [%s] %.1f%% %s/%s",
		pb.prefix,
		bar,
		percent*100,
		formatBytes(pb.current),
		formatBytes(pb.total),
	)

	// 显示速度和预计时间
	if speed > 0 {
		fmt.Fprintf(pb.writer, " | %s/s", formatBytes(int64(speed)))
	}
	if eta > 0 && pb.current < pb.total {
		fmt.Fprintf(pb.writer, " | ETA: %s", formatDuration(eta))
	}
}

// formatBytes 格式化字节数
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDuration 格式化时间
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
