# UI Package

用户交互和进度显示组件。

## 组件

### ProgressBar
下载进度条显示器，支持：
- 实时进度更新
- 速度显示
- 预计剩余时间
- 字节格式化

### Prompter
交互式提示器，支持：
- 确认提示（是/否）
- 选择提示（从列表中选择）
- 文本输入
- 版本选择（带分页）

### Messenger
消息显示器，支持：
- 信息、成功、警告、错误消息
- 分节标题
- 列表显示
- 表格显示
- 加载指示器（Spinner）

### ErrorFormatter
错误格式化器，支持：
- 友好的错误消息
- 错误原因显示
- 解决建议
- 帮助信息

### OutputStreamer
实时输出流处理器，支持：
- 命令输出的实时显示
- 带前缀的输出
- 回调处理

## 使用示例

```go
// 进度条
pb := ui.NewProgressBar(os.Stdout, totalSize, "Downloading")
pb.Update(currentSize)
pb.Finish()

// 确认提示
prompter := ui.NewPrompter(os.Stdin, os.Stdout)
confirmed, _ := prompter.Confirm("Continue?", true)

// 版本选择
version, _ := prompter.SelectVersion(versions, 10)

// 消息显示
messenger := ui.NewMessenger(os.Stdout)
messenger.Success("Installation completed")
messenger.Error("Failed to download")

// 错误格式化
formatter := ui.NewErrorFormatter(os.Stderr)
formatter.Format(err)

// 实时输出
streamer := ui.NewOutputStreamer(cmdOutput, os.Stdout, "")
streamer.Stream()
```
