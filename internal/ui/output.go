package ui

import (
	"bufio"
	"io"
	"sync"
)

// OutputStreamer 实时输出流处理器
type OutputStreamer struct {
	reader io.Reader
	writer io.Writer
	prefix string
	mu     sync.Mutex
}

// NewOutputStreamer 创建新的输出流处理器
func NewOutputStreamer(reader io.Reader, writer io.Writer, prefix string) *OutputStreamer {
	return &OutputStreamer{
		reader: reader,
		writer: writer,
		prefix: prefix,
	}
}

// Stream 开始流式输出
func (os *OutputStreamer) Stream() error {
	scanner := bufio.NewScanner(os.reader)
	for scanner.Scan() {
		os.mu.Lock()
		if os.prefix != "" {
			io.WriteString(os.writer, os.prefix)
		}
		io.WriteString(os.writer, scanner.Text())
		io.WriteString(os.writer, "\n")
		os.mu.Unlock()
	}

	return scanner.Err()
}

// StreamWithCallback 带回调的流式输出
func (os *OutputStreamer) StreamWithCallback(callback func(line string)) error {
	scanner := bufio.NewScanner(os.reader)
	for scanner.Scan() {
		line := scanner.Text()

		os.mu.Lock()
		if os.prefix != "" {
			io.WriteString(os.writer, os.prefix)
		}
		io.WriteString(os.writer, line)
		io.WriteString(os.writer, "\n")
		os.mu.Unlock()

		if callback != nil {
			callback(line)
		}
	}

	return scanner.Err()
}
