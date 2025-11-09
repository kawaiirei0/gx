package ui_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kawaiirei0/gx/internal/ui"
	"github.com/kawaiirei0/gx/pkg/errors"
)

// TestProgressBar 测试进度条功能
func TestProgressBar(t *testing.T) {
	var buf bytes.Buffer
	totalSize := int64(1024 * 1024) // 1MB

	pb := ui.NewProgressBar(&buf, totalSize, "Downloading")

	// 模拟下载进度
	for i := int64(0); i <= totalSize; i += totalSize / 10 {
		pb.Update(i)
		time.Sleep(10 * time.Millisecond)
	}

	pb.Finish()

	output := buf.String()
	if !strings.Contains(output, "100.0%") {
		t.Errorf("Expected progress bar to show 100%%, got: %s", output)
	}
}

// TestMessenger 测试消息显示器
func TestMessenger(t *testing.T) {
	var buf bytes.Buffer
	messenger := ui.NewMessenger(&buf)

	messenger.Info("This is an info message")
	messenger.Success("This is a success message")
	messenger.Warning("This is a warning message")
	messenger.Error("This is an error message")

	output := buf.String()

	expectedMessages := []string{
		"This is an info message",
		"This is a success message",
		"This is a warning message",
		"This is an error message",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("Expected output to contain '%s', got: %s", msg, output)
		}
	}
}

// TestMessengerSection 测试分节显示
func TestMessengerSection(t *testing.T) {
	var buf bytes.Buffer
	messenger := ui.NewMessenger(&buf)

	messenger.Section("Test Section")

	output := buf.String()
	if !strings.Contains(output, "Test Section") {
		t.Errorf("Expected section title, got: %s", output)
	}
}

// TestMessengerList 测试列表显示
func TestMessengerList(t *testing.T) {
	var buf bytes.Buffer
	messenger := ui.NewMessenger(&buf)

	items := []string{"Item 1", "Item 2", "Item 3"}
	messenger.List(items, "•")

	output := buf.String()
	for _, item := range items {
		if !strings.Contains(output, item) {
			t.Errorf("Expected list to contain '%s', got: %s", item, output)
		}
	}
}

// TestMessengerTable 测试表格显示
func TestMessengerTable(t *testing.T) {
	var buf bytes.Buffer
	messenger := ui.NewMessenger(&buf)

	headers := []string{"Name", "Version", "Status"}
	rows := [][]string{
		{"Go", "1.21.5", "Active"},
		{"Go", "1.20.0", "Inactive"},
	}

	messenger.Table(headers, rows)

	output := buf.String()

	// 检查表头
	for _, header := range headers {
		if !strings.Contains(output, header) {
			t.Errorf("Expected table to contain header '%s', got: %s", header, output)
		}
	}

	// 检查行数据
	for _, row := range rows {
		for _, cell := range row {
			if !strings.Contains(output, cell) {
				t.Errorf("Expected table to contain cell '%s', got: %s", cell, output)
			}
		}
	}
}

// TestErrorFormatter 测试错误格式化
func TestErrorFormatter(t *testing.T) {
	var buf bytes.Buffer
	formatter := ui.NewErrorFormatter(&buf)

	// 测试自定义错误
	gxErr := errors.ErrVersionNotFound.WithMessage("Go version 1.99.0 not found")
	formatter.Format(gxErr)

	output := buf.String()
	if !strings.Contains(output, "not found") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

// TestSpinner 测试加载指示器
func TestSpinner(t *testing.T) {
	var buf bytes.Buffer
	spinner := ui.NewSpinner(&buf, "Loading...")

	// 模拟几次更新
	for i := 0; i < 5; i++ {
		spinner.Tick()
		time.Sleep(10 * time.Millisecond)
	}

	spinner.Stop("✓ Done")

	output := buf.String()
	if !strings.Contains(output, "Done") {
		t.Errorf("Expected spinner to show completion message, got: %s", output)
	}
}

// ExampleProgressBar 演示进度条使用
func ExampleProgressBar() {
	var buf bytes.Buffer
	totalSize := int64(1000)

	pb := ui.NewProgressBar(&buf, totalSize, "Downloading")

	// 模拟下载
	for i := int64(0); i <= totalSize; i += 100 {
		pb.Update(i)
	}

	pb.Finish()
	fmt.Println("Download complete")
	// Output: Download complete
}

// ExampleMessenger 演示消息显示器使用
func ExampleMessenger() {
	var buf bytes.Buffer
	messenger := ui.NewMessenger(&buf)

	messenger.Info("Starting installation...")
	messenger.Success("Installation completed")

	// 输出会包含消息图标和文本
	fmt.Println("Messages displayed")
	// Output: Messages displayed
}

// ExamplePrompter_Confirm 演示确认提示
func ExamplePrompter_Confirm() {
	// 注意：这个示例需要用户输入，在实际测试中需要模拟输入
	input := strings.NewReader("y\n")
	var output bytes.Buffer

	prompter := ui.NewPrompter(input, &output)
	confirmed, _ := prompter.Confirm("Continue?", true)

	if confirmed {
		fmt.Println("User confirmed")
	}
	// Output: User confirmed
}

// ExamplePrompter_Select 演示选择提示
func ExamplePrompter_Select() {
	input := strings.NewReader("2\n")
	var output bytes.Buffer

	prompter := ui.NewPrompter(input, &output)
	options := []string{"Option 1", "Option 2", "Option 3"}
	selected, _ := prompter.Select("Choose an option:", options)

	fmt.Printf("Selected: %s\n", options[selected])
	// Output: Selected: Option 2
}
