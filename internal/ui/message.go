package ui

import (
	"fmt"
	"io"
	"strings"
)

// MessageType 消息类型
type MessageType int

const (
	// MessageTypeInfo 信息消息
	MessageTypeInfo MessageType = iota
	// MessageTypeSuccess 成功消息
	MessageTypeSuccess
	// MessageTypeWarning 警告消息
	MessageTypeWarning
	// MessageTypeError 错误消息
	MessageTypeError
)

// Messenger 消息显示器
type Messenger struct {
	writer io.Writer
}

// NewMessenger 创建新的消息显示器
func NewMessenger(writer io.Writer) *Messenger {
	return &Messenger{
		writer: writer,
	}
}

// Info 显示信息消息
func (m *Messenger) Info(message string) {
	m.print(MessageTypeInfo, message)
}

// Success 显示成功消息
func (m *Messenger) Success(message string) {
	m.print(MessageTypeSuccess, message)
}

// Warning 显示警告消息
func (m *Messenger) Warning(message string) {
	m.print(MessageTypeWarning, message)
}

// Error 显示错误消息
func (m *Messenger) Error(message string) {
	m.print(MessageTypeError, message)
}

// print 打印消息
func (m *Messenger) print(msgType MessageType, message string) {
	var prefix string
	switch msgType {
	case MessageTypeInfo:
		prefix = "ℹ"
	case MessageTypeSuccess:
		prefix = "✓"
	case MessageTypeWarning:
		prefix = "⚠"
	case MessageTypeError:
		prefix = "✗"
	}

	fmt.Fprintf(m.writer, "%s %s\n", prefix, message)
}

// Section 显示分节标题
func (m *Messenger) Section(title string) {
	fmt.Fprintln(m.writer)
	fmt.Fprintln(m.writer, title)
	fmt.Fprintln(m.writer, strings.Repeat("─", len(title)))
}

// List 显示列表
func (m *Messenger) List(items []string, marker string) {
	if marker == "" {
		marker = "•"
	}

	for _, item := range items {
		fmt.Fprintf(m.writer, "  %s %s\n", marker, item)
	}
}

// Table 显示简单表格
func (m *Messenger) Table(headers []string, rows [][]string) {
	if len(headers) == 0 || len(rows) == 0 {
		return
	}

	// 计算每列的最大宽度
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// 打印表头
	for i, header := range headers {
		fmt.Fprintf(m.writer, "%-*s  ", colWidths[i], header)
	}
	fmt.Fprintln(m.writer)

	// 打印分隔线
	for _, width := range colWidths {
		fmt.Fprint(m.writer, strings.Repeat("─", width), "  ")
	}
	fmt.Fprintln(m.writer)

	// 打印行
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				fmt.Fprintf(m.writer, "%-*s  ", colWidths[i], cell)
			}
		}
		fmt.Fprintln(m.writer)
	}
}

// Spinner 简单的加载指示器
type Spinner struct {
	writer  io.Writer
	message string
	frames  []string
	current int
}

// NewSpinner 创建新的加载指示器
func NewSpinner(writer io.Writer, message string) *Spinner {
	return &Spinner{
		writer:  writer,
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		current: 0,
	}
}

// Tick 更新加载指示器
func (s *Spinner) Tick() {
	fmt.Fprintf(s.writer, "\r%s %s", s.frames[s.current], s.message)
	s.current = (s.current + 1) % len(s.frames)
}

// Stop 停止加载指示器
func (s *Spinner) Stop(finalMessage string) {
	fmt.Fprintf(s.writer, "\r%s\n", finalMessage)
}

// Clear 清除当前行
func (s *Spinner) Clear() {
	fmt.Fprint(s.writer, "\r", strings.Repeat(" ", len(s.message)+10), "\r")
}
