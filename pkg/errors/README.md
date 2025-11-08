# Error Handling and Recovery System

This package provides a comprehensive error handling and recovery mechanism for the gx application.

## Features

### 1. Custom Error Types

All errors in the system use a custom `Error` type that provides:

- **Error Codes**: Unique identifiers for each error type
- **Error Messages**: Human-readable descriptions
- **Error Chaining**: Support for wrapping underlying errors
- **Context Information**: Ability to attach additional context data
- **Recoverability Flags**: Mark errors as recoverable or not

Example:
```go
err := errors.ErrVersionNotFound.
    WithMessage("version go1.21.0 not found").
    WithContext("requested_version", "go1.21.0").
    AsRecoverable()
```

### 2. Error Categories

The system defines several categories of errors:

- **Version Errors**: `ErrVersionNotFound`, `ErrVersionAlreadyInstalled`, `ErrVersionNotInstalled`
- **Download Errors**: `ErrDownloadFailed`, `ErrNetworkError`, `ErrChecksumMismatch`
- **Installation Errors**: `ErrInstallFailed`, `ErrUninstallFailed`, `ErrArchiveCorrupted`
- **Environment Errors**: `ErrEnvironmentSetupFailed`, `ErrPermissionDenied`
- **Storage Errors**: `ErrStorageFailed`, `ErrConfigCorrupted`
- **Platform Errors**: `ErrPlatformNotSupported`
- **Operation Errors**: `ErrOperationFailed`, `ErrTimeout`, `ErrCancelled`

### 3. Recovery Manager

The `RecoveryManager` provides automatic cleanup and rollback capabilities:

```go
recovery := errors.NewRecoveryManager()
defer recovery.Cleanup()

// Register cleanup functions
errors.EnsureDirectoryCleanup(recovery, "/path/to/temp/dir")
errors.EnsureFileCleanup(recovery, "/path/to/temp/file")

// Register rollback functions
recovery.AddRollback(func() error {
    return restoreBackup()
})

// On error, execute rollback
if err != nil {
    recovery.CleanupAndRollback()
    return err
}
```

### 4. Error Wrapping and Propagation

All components wrap errors with additional context:

```go
if err := someOperation(); err != nil {
    return errors.Wrap(err, "OPERATION_FAILED", "failed to perform operation").
        WithContext("operation", "install").
        WithContext("version", version)
}
```

### 5. Graceful Degradation

The system implements graceful degradation where appropriate:

- **Config Backup Recovery**: If config is corrupted, attempt to restore from backup
- **Checksum Verification**: If file info unavailable, skip verification but log warning
- **Environment Variables**: If persistent setup fails, at least set for current session

### 6. User-Friendly Error Reporting

The `ErrorReporter` provides user-friendly error messages with suggestions:

```go
reporter := errors.NewErrorReporter(verbose)
fmt.Println(reporter.Report(err))
```

Output example:
```
Error: version not found
Code: VERSION_NOT_FOUND
Context:
  requested_version: go1.21.0

Suggestion: Use 'gx list' to see installed versions or 'gx install <version>' to install a new version.
```

## Usage Guidelines

### 1. Creating Errors

Always use predefined error types:

```go
// Good
return errors.ErrVersionNotFound.WithMessage("version go1.21.0 not found")

// Avoid
return fmt.Errorf("version not found")
```

### 2. Adding Context

Add relevant context information to errors:

```go
return errors.ErrDownloadFailed.
    WithCause(err).
    WithMessage("failed to download Go archive").
    WithContext("version", version).
    WithContext("url", downloadURL).
    WithContext("attempt", attemptNumber)
```

### 3. Using Recovery Manager

Use recovery manager for operations that create temporary resources:

```go
func Install(version string) error {
    recovery := errors.NewRecoveryManager()
    defer recovery.Cleanup()
    
    // Create temp directory
    tmpDir, err := errors.CreateTempDir(recovery, "gx-install-")
    if err != nil {
        return err
    }
    
    // Perform installation
    if err := doInstall(tmpDir, version); err != nil {
        recovery.CleanupAndRollback()
        return err
    }
    
    // Success - clear cleanup functions
    recovery.Clear()
    return nil
}
```

### 4. Checking Error Types

Use `errors.Is` to check error types:

```go
if errors.IsType(err, errors.ErrVersionNotFound) {
    // Handle version not found
}
```

### 5. Graceful Degradation

Implement graceful degradation for non-critical failures:

```go
// Try to backup config
backupPath, err := errors.BackupFile(configPath)
if err != nil {
    // Log warning but continue
    logger.Warn("Failed to backup config: %v", err)
} else {
    // Register cleanup
    recovery.AddCleanup(func() error {
        return errors.SafeRemoveFile(backupPath)
    })
}
```

## Error Handling Patterns

### Pattern 1: Install with Cleanup

```go
func (m *manager) Install(version string) error {
    recovery := errors.NewRecoveryManager()
    defer recovery.Cleanup()
    
    // Register cleanup for temp files
    errors.EnsureFileCleanup(recovery, archivePath)
    errors.EnsureDirectoryCleanup(recovery, versionPath)
    
    // Perform installation
    if err := download(); err != nil {
        recovery.CleanupAndRollback()
        return err
    }
    
    if err := extract(); err != nil {
        recovery.CleanupAndRollback()
        return err
    }
    
    // Success
    recovery.Clear()
    return nil
}
```

### Pattern 2: Config Update with Backup

```go
func (s *store) Save(config *Config) error {
    // Backup existing config
    backupPath, _ := errors.BackupFile(s.configPath)
    defer os.Remove(backupPath)
    
    // Write to temp file
    tmpPath := s.configPath + ".tmp"
    if err := writeFile(tmpPath, data); err != nil {
        return err
    }
    
    // Atomic replace
    if err := os.Rename(tmpPath, s.configPath); err != nil {
        os.Remove(tmpPath)
        return err
    }
    
    return nil
}
```

### Pattern 3: Operation with Rollback

```go
func (m *manager) SwitchVersion(version string) error {
    recovery := errors.NewRecoveryManager()
    
    // Backup current state
    oldVersion := getCurrentVersion()
    recovery.AddRollback(func() error {
        return switchTo(oldVersion)
    })
    
    // Perform switch
    if err := switchTo(version); err != nil {
        recovery.Rollback()
        return err
    }
    
    return nil
}
```

## Testing Error Handling

When testing, verify:

1. **Error Types**: Correct error types are returned
2. **Error Context**: Relevant context is attached
3. **Cleanup**: Resources are properly cleaned up on failure
4. **Rollback**: State is restored on failure
5. **Error Messages**: User-friendly messages are generated

Example test:
```go
func TestInstallWithCleanup(t *testing.T) {
    // Setup
    tmpDir := t.TempDir()
    
    // Simulate failure
    err := manager.Install("invalid-version")
    
    // Verify error type
    assert.True(t, errors.IsType(err, errors.ErrVersionNotFound))
    
    // Verify cleanup occurred
    assert.NoDirExists(t, filepath.Join(tmpDir, "go1.21.0"))
}
```

## Best Practices

1. **Always wrap errors** with additional context
2. **Use recovery manager** for operations with side effects
3. **Implement graceful degradation** for non-critical failures
4. **Provide user-friendly messages** with actionable suggestions
5. **Log errors** at appropriate levels (ERROR, WARN, INFO)
6. **Clean up resources** even on success (use defer)
7. **Test error paths** as thoroughly as success paths
8. **Document error conditions** in function comments
