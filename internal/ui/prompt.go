package ui

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Prompter 交互式提示器
type Prompter struct {
	reader io.Reader
	writer io.Writer
}

// NewPrompter 创建新的提示器
func NewPrompter(reader io.Reader, writer io.Writer) *Prompter {
	return &Prompter{
		reader: reader,
		writer: writer,
	}
}

// Confirm 确认提示（是/否）
func (p *Prompter) Confirm(message string, defaultYes bool) (bool, error) {
	prompt := message
	if defaultYes {
		prompt += " [Y/n]: "
	} else {
		prompt += " [y/N]: "
	}

	fmt.Fprint(p.writer, prompt)

	scanner := bufio.NewScanner(p.reader)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}
		// EOF 或无输入，使用默认值
		return defaultYes, nil
	}

	input := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if input == "" {
		return defaultYes, nil
	}

	return input == "y" || input == "yes", nil
}

// Select 选择提示（从列表中选择）
func (p *Prompter) Select(message string, options []string) (int, error) {
	fmt.Fprintln(p.writer, message)
	fmt.Fprintln(p.writer)

	for i, option := range options {
		fmt.Fprintf(p.writer, "  %d) %s\n", i+1, option)
	}

	fmt.Fprintln(p.writer)
	fmt.Fprint(p.writer, "Enter your choice (1-", len(options), "): ")

	scanner := bufio.NewScanner(p.reader)
	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return -1, err
			}
			return -1, fmt.Errorf("no input received")
		}

		input := strings.TrimSpace(scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Fprintf(p.writer, "Invalid choice. Please enter a number between 1 and %d: ", len(options))
			continue
		}

		return choice - 1, nil
	}
}

// Input 文本输入提示
func (p *Prompter) Input(message string, defaultValue string) (string, error) {
	prompt := message
	if defaultValue != "" {
		prompt += fmt.Sprintf(" [%s]", defaultValue)
	}
	prompt += ": "

	fmt.Fprint(p.writer, prompt)

	scanner := bufio.NewScanner(p.reader)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return defaultValue, nil
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}

// SelectVersion 版本选择提示（带搜索和分页）
func (p *Prompter) SelectVersion(versions []string, pageSize int) (string, error) {
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions available")
	}

	// 如果版本数量较少，直接显示所有版本
	if len(versions) <= pageSize {
		fmt.Fprintln(p.writer, "Available versions:")
		fmt.Fprintln(p.writer)

		for i, version := range versions {
			fmt.Fprintf(p.writer, "  %d) %s\n", i+1, version)
		}

		fmt.Fprintln(p.writer)
		fmt.Fprint(p.writer, "Enter your choice (1-", len(versions), ") or version number: ")

		scanner := bufio.NewScanner(p.reader)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", err
			}
			return "", fmt.Errorf("no input received")
		}

		input := strings.TrimSpace(scanner.Text())

		// 尝试作为索引解析
		if choice, err := strconv.Atoi(input); err == nil {
			if choice >= 1 && choice <= len(versions) {
				return versions[choice-1], nil
			}
		}

		// 尝试作为版本号解析
		for _, version := range versions {
			if strings.Contains(version, input) || strings.Contains(strings.TrimPrefix(version, "go"), input) {
				return version, nil
			}
		}

		return "", fmt.Errorf("invalid selection: %s", input)
	}

	// 版本数量较多，使用分页显示
	page := 0
	totalPages := (len(versions) + pageSize - 1) / pageSize

	for {
		start := page * pageSize
		end := start + pageSize
		if end > len(versions) {
			end = len(versions)
		}

		fmt.Fprintf(p.writer, "\nAvailable versions (page %d/%d):\n\n", page+1, totalPages)

		for i := start; i < end; i++ {
			fmt.Fprintf(p.writer, "  %d) %s\n", i+1, versions[i])
		}

		fmt.Fprintln(p.writer)
		fmt.Fprint(p.writer, "Enter choice, version number, 'n' for next page, 'p' for previous page, or 'q' to quit: ")

		scanner := bufio.NewScanner(p.reader)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", err
			}
			return "", fmt.Errorf("no input received")
		}

		input := strings.TrimSpace(strings.ToLower(scanner.Text()))

		// 处理导航命令
		switch input {
		case "n", "next":
			if page < totalPages-1 {
				page++
			} else {
				fmt.Fprintln(p.writer, "Already on last page")
			}
			continue
		case "p", "prev", "previous":
			if page > 0 {
				page--
			} else {
				fmt.Fprintln(p.writer, "Already on first page")
			}
			continue
		case "q", "quit":
			return "", fmt.Errorf("selection cancelled")
		}

		// 尝试作为索引解析
		if choice, err := strconv.Atoi(input); err == nil {
			if choice >= 1 && choice <= len(versions) {
				return versions[choice-1], nil
			}
		}

		// 尝试作为版本号解析
		for _, version := range versions {
			if strings.Contains(version, input) || strings.Contains(strings.TrimPrefix(version, "go"), input) {
				return version, nil
			}
		}

		fmt.Fprintf(p.writer, "Invalid selection: %s\n", input)
	}
}
