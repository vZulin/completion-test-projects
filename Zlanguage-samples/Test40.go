package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Static error definitions for err113 compliance
var (
	ErrToolboxDirCreation     = errors.New("failed to create Toolbox CLI directory")
	ErrToolboxFileNotFound    = errors.New("failed to find downloaded Toolbox CLI file")
	ErrToolboxExtraction      = errors.New("failed to extract Toolbox CLI")
	ErrInstallationMarker     = errors.New("failed to create installation marker")
	ErrHomeDirectory          = errors.New("failed to get home directory")
	ErrEmptyHomeDirectory     = errors.New("empty home directory returned")
	ErrNoRemoteDevEntries     = errors.New("no RemoteDev registry entries found")
	ErrNoRemoteDevConfig      = errors.New("no RemoteDev configuration files found")
	ErrEmptyURL               = errors.New("empty URL")
	ErrBaseURLExtraction      = errors.New("unable to extract base URL")
	ErrConfigFetch            = errors.New("failed to fetch config")
	ErrConfigNotFound         = errors.New("config not found")
	ErrConfigRead             = errors.New("failed to read config")
	ErrConfigParse            = errors.New("failed to parse config")
	ErrUnsupportedOS          = errors.New("unsupported OS")
	ErrUnsupportedProductType = errors.New("unsupported product type")
	ErrDirectoryCreation      = errors.New("failed to create directory")
	ErrFileWrite              = errors.New("failed to write file")
	ErrFileRead               = errors.New("failed to read file")
	ErrDownloadFailure        = errors.New("failed to download file")
	ErrInstallationFailure    = errors.New("installation failed")
	ErrIDENotFound            = errors.New("IDE not found")
	ErrProductsJSON           = errors.New("products JSON file not found")
	ErrBuildNumberNotFound    = errors.New("build number not found")
	ErrSSHConnection          = errors.New("SSH connection failed")
	ErrPermissionDenied       = errors.New("permission denied")
	ErrFileNotExists          = errors.New("file does not exist")
	ErrDirectoryNotExists     = errors.New("directory does not exist")
	ErrTimeout                = errors.New("operation timed out")
	ErrProcessNotFound        = errors.New("process not found")
	ErrInvalidFileType        = errors.New("unsupported file type")
)

// handleCommandOutputAsync runs a command asynchronously and logs its output.
func handleCommandOutputAsync(cmd *exec.Cmd, isIDE bool) {
	go func() {
		output, err := cmd.CombinedOutput()
		if err != nil {
			if isIDE {
				logger.Warning(fmt.Sprintf("IDE may have encountered an issue: %v", err))
			} else {
				logger.Warning(fmt.Sprintf("TBCLi may have encountered an issue: %v", err))
			}
		}
		if len(output) > 0 {
			outputStr := strings.TrimSpace(string(output))
			if isIDE {
				if strings.Contains(outputStr, "not found") || strings.Contains(outputStr, "cannot find") {
					logger.Warning(fmt.Sprintf("remote-dev-server executable not found: %s", outputStr))
				} else if len(outputStr) > 0 {
					maxLen := 200
					if len(outputStr) < maxLen {
						maxLen = len(outputStr)
					}
					logger.Info(fmt.Sprintf("IDE output: %s", outputStr[:maxLen]))
				}
			} else {
				if strings.Contains(outputStr, "~RESULT:") {
					logger.Success("TBCLi started successfully on Windows")
				} else if len(outputStr) > 0 {
					maxLen := 200
					if len(outputStr) < maxLen {
						maxLen = len(outputStr)
					}
					logger.Info(fmt.Sprintf("TBCLi output: %s", outputStr[:maxLen]))
				}
			}
		}
	}()
}

// handleIDEInstallationLogic handles the common IDE installation logic.
func handleIDEInstallationLogic(tbcliDir, tbcliURL string, forceReinstall bool) (bool, error) {
	// Handle force reinstallation - remove existing directory
	if forceReinstall {
		logger.Info("Force reinstallation enabled, reinstalling Toolbox CLI")
		if _, err := os.Stat(tbcliDir); err == nil {
			logger.Info(fmt.Sprintf("Removing existing Toolbox CLI directory: %s", tbcliDir))
			if err := os.RemoveAll(tbcliDir); err != nil {
				logger.Warning(fmt.Sprintf("Failed to remove existing Toolbox CLI directory: %v", err))
			}
		}
	}

	err := os.MkdirAll(tbcliDir, 0o750)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrToolboxDirCreation, err)
	}

	tbcliFile, err := findDownloadedFile(filepath.Base(tbcliURL))
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrToolboxFileNotFound, err)
	}

	err = extractTarGz(tbcliFile, tbcliDir)
	if err != nil {
		logger.Warning(fmt.Sprintf("Extraction failed, cleaning up directory: %s", tbcliDir))
		_ = os.RemoveAll(tbcliDir) // #nosec G104 -- Best effort cleanup
		return false, fmt.Errorf("%w: %v", ErrToolboxExtraction, err)
	}

	// Create .extracted marker file - critical for installation state tracking
	extractedFile := filepath.Join(tbcliDir, ".extracted")
	err = os.WriteFile(extractedFile, []byte(""), 0o600)
	if err != nil {
		// If marker creation fails, the installation is incomplete
		logger.Warning(fmt.Sprintf("Failed to create .extracted marker, cleaning up: %v", err))
		_ = os.RemoveAll(tbcliDir) // #nosec G104 -- Best effort cleanup
		return false, fmt.Errorf("%w: %v", ErrInstallationMarker, err)
	}
	return true, nil // Mark as newly installed
}

// handleRemoteVMOptionsConfiguration handles VM options configuration for remote installation
func handleRemoteVMOptionsConfiguration(config Config, opts InstallOptions, remoteInfo RemoteSystemInfo) error {
	if len(config.VMOptions) > 0 {
		var homeDir string
		if remoteInfo.OS == PlatformWindows {
			cmd := exec.Command("ssh", opts.SSHTarget, `powershell.exe -Command "$env:USERPROFILE"`) // #nosec G204 -- SSH target is validated
			output, err := cmd.Output()
			if err == nil {
				homeDir = strings.TrimSpace(string(output))
			}
		} else {
			var err error
			homeDir, err = getRemoteHomeDirectory(opts.SSHTarget, remoteInfo.OS)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to get remote home directory: %v", err))
			}
		}

		if homeDir != "" {
			configDir, vmoptionsFile := getIDEConfigPaths(remoteInfo.OS, homeDir, opts.IdeType, opts.IdeVersion)
			err := createVMOptionsFileRemote(opts.SSHTarget, configDir, vmoptionsFile, config.VMOptions, remoteInfo.OS)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to configure VM options: %v", err))
			}
		}
	}
	return nil
}

// ConfigIdeSpec represents an IDE specification in config
type ConfigIdeSpec struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

// Configuration structure
type Config struct {
	// String fields
	DefaultTbcliVersion   string `json:"defaultTbcliVersion"`
	DefaultToolboxVersion string `json:"defaultToolboxVersion"`
	DefaultJbrVersion     string `json:"defaultJbrVersion"`
	DefaultJbrBuild       string `json:"defaultJbrBuild"`
	DefaultIdeType        string `json:"defaultIdeType"`
	DefaultIdeVersion     string `json:"defaultIdeVersion"`
	IdeType               string `json:"ide_type"`
	IdeVersion            string `json:"ide_version"`
	IdeInstallPath        string `json:"ide_install_path"`
	TbcliInstallPath      string `json:"tbcli_install_path"`
	JbrInstallPath        string `json:"jbr_install_path"`
	DownloadOS            string `json:"download_os"`
	DownloadArch          string `json:"download_arch"`
	LaunchIdeArgs         string `json:"launch_ide_args"`
	JetbrainsBaseURL      string `json:"jetbrainsBaseURL"`
	JbrBaseURL            string `json:"jbrBaseURL"`
	ProductsJSONURL       string `json:"productsJsonURL"`
	ToolboxEnvironment    string `json:"toolboxEnvironment"`
	// Slice fields
	IdeSpecs    []ConfigIdeSpec `json:"ideSpecs,omitempty"`
	IdeTypes    []string        `json:"ideTypes,omitempty"`
	IdeVersions []string        `json:"ideVersions,omitempty"`
	VMOptions   []string        `json:"vmOptions,omitempty"`
	// Bool fields
	InstallIde     bool `json:"install_ide"`
	InstallClient  bool `json:"install_client"`
	NoClient       bool `json:"no_client"`
	DownloadOnly   bool `json:"download_only"`
	LocalInstall   bool `json:"local_install"`
	NoProgress     bool `json:"no_progress"`
	Force          bool `json:"force"`
	InfoLogging    bool `json:"info_logging"`
	LaunchTbcli    bool `json:"launch_tbcli"`
	OpenToolboxURL bool `json:"open_toolbox_url"`
}

// ToolboxEnvironmentConfig represents the structure of environment.json
type ToolboxEnvironmentConfig struct {
	Tools               Tools `json:"tools"`
	AllowPortForwarding bool  `json:"allowPortForwarding"`
}

// Tools represents the tools section in environment.json
type Tools struct {
	Location            []Location `json:"location"`
	AllowInstallation   bool       `json:"allowInstallation"`
	AllowUpdate         bool       `json:"allowUpdate"`
	AllowUninstallation bool       `json:"allowUninstallation"`
}

// Location represents a location entry in environment.json
type Location struct {
	Path   string `json:"path"`
	Levels int    `json:"levels,omitempty"`
}

// IDE installation specification for multi-IDE support
type IdeSpec struct {
	Type    string
	Version string
}

// SSH connection information
type SSHConnectionInfo struct {
	Hostname string
	Username string
	Port     string
}

// Installation options
type InstallOptions struct {
	// String fields
	TbcliVersion     string
	ToolboxVersion   string
	JbrVersion       string
	JbrBuild         string
	IdeType          string // Legacy single IDE type
	IdeVersion       string // Legacy single IDE version
	IdeInstallPath   string
	TbcliInstallPath string
	JbrInstallPath   string
	DownloadOS       string
	DownloadArch     string
	LaunchIdeArgs    string
	ConfigPath       string
	SSHTarget        string
	// Slice fields
	IdeTypes    []string  // Multiple IDE types
	IdeVersions []string  // Multiple IDE versions
	IdeSpecs    []IdeSpec // Parsed IDE specifications
	// Bool fields
	InstallIde            bool
	InstallClient         bool
	NoClient              bool
	DownloadOnly          bool
	LocalInstall          bool
	InfoLogging           bool
	NoProgress            bool
	Force                 bool
	LaunchTbcli           bool
	OpenToolboxURL        bool
	UserSetTbcliVersion   bool
	UserSetToolboxVersion bool
	UserSetInstallIde     bool
	UserSetInstallClient  bool
}

// Product info structure for existing installation checks
type ProductInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	BuildNumber string `json:"buildNumber"`
	ProductCode string `json:"productCode"`
	Launch      []struct {
		OS   string `json:"os"`
		Arch string `json:"arch"`
	} `json:"launch"`
}

// InstallationResult holds the result of installation check/install
type InstallationResult struct {
	AppPath string
	Exists  bool
}

// System info
type SystemInfo struct {
	OS   string
	Arch string
	Name string
}

// Progress writer for download progress bar
type ProgressWriter struct {
	StartTime  time.Time
	LastUpdate time.Time
	Filename   string
	Total      int64
	Downloaded int64
	mu         sync.Mutex
}

func NewProgressWriter(total int64, filename string) *ProgressWriter {
	return &ProgressWriter{
		Total:     total,
		StartTime: time.Now(),
		Filename:  filename,
	}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	n := len(p)
	pw.Downloaded += int64(n)

	now := time.Now()
	if now.Sub(pw.LastUpdate) >= 100*time.Millisecond || pw.Downloaded == pw.Total {
		pw.LastUpdate = now
		pw.displayProgress()
	}

	return n, nil
}

func (pw *ProgressWriter) displayProgress() {
	if pw.Total <= 0 {
		return
	}

	percent := float64(pw.Downloaded) / float64(pw.Total) * 100
	elapsed := time.Since(pw.StartTime)

	// Calculate speed
	var speed float64
	var speedUnit string
	if elapsed.Seconds() > 0 {
		bytesPerSec := float64(pw.Downloaded) / elapsed.Seconds()
		if bytesPerSec >= 1024*1024 {
			speed = bytesPerSec / (1024 * 1024)
			speedUnit = "MB/s"
		} else if bytesPerSec >= 1024 {
			speed = bytesPerSec / 1024
			speedUnit = "KB/s"
		} else {
			speed = bytesPerSec
			speedUnit = "B/s"
		}
	}

	// Format file size
	var downloadedStr, totalStr string
	if pw.Total >= 1024*1024*1024 {
		downloadedStr = fmt.Sprintf("%.1f GB", float64(pw.Downloaded)/(1024*1024*1024))
		totalStr = fmt.Sprintf("%.1f GB", float64(pw.Total)/(1024*1024*1024))
	} else if pw.Total >= 1024*1024 {
		downloadedStr = fmt.Sprintf("%.1f MB", float64(pw.Downloaded)/(1024*1024))
		totalStr = fmt.Sprintf("%.1f MB", float64(pw.Total)/(1024*1024))
	} else if pw.Total >= 1024 {
		downloadedStr = fmt.Sprintf("%.1f KB", float64(pw.Downloaded)/1024)
		totalStr = fmt.Sprintf("%.1f KB", float64(pw.Total)/1024)
	} else {
		downloadedStr = fmt.Sprintf("%d B", pw.Downloaded)
		totalStr = fmt.Sprintf("%d B", pw.Total)
	}

	// Create progress bar
	barWidth := 40
	filled := int(percent / 100 * float64(barWidth))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	// Clear line and print progress
	fmt.Printf("\r%s[%s] %.1f%% (%s/%s) %.1f %s",
		colorBlue+"[DOWNLOADING]"+colorReset,
		bar,
		percent,
		downloadedStr,
		totalStr,
		speed,
		speedUnit)

	if pw.Downloaded == pw.Total {
		fmt.Println()
	}
}

// Global logger
var logger *Logger

// Global progress settings
var disableProgress bool

// Global force settings
var forceReinstall bool

// Logger with colored output
type Logger struct {
	infoEnabled bool
}

// expandEnvPath expands environment variables in path strings cross-platform
func expandEnvPath(path string) string {
	if path == "" {
		return path
	}

	// Support both $VAR and ${VAR} syntax on all platforms
	expanded := os.ExpandEnv(path)

	// For Windows, also support %VAR% syntax
	if runtime.GOOS == PlatformWindows {
		// Match %VARNAME% pattern
		re := regexp.MustCompile(`%([^%]+)%`)
		expanded = re.ReplaceAllStringFunc(expanded, func(match string) string {
			// Extract variable name (remove % characters)
			varName := match[1 : len(match)-1]
			if value := os.Getenv(varName); value != "" {
				return value
			}
			return match // Return original if variable not found
		})
	}

	return expanded
}

// normalizePathForOS normalizes path separators for the specified OS
func normalizePathForOS(path, osType string) string {
	if path == "" {
		return path
	}

	switch osType {
	case PlatformWindows:
		normalized := strings.ReplaceAll(path, "/", "\\")

		for strings.Contains(normalized, "\\\\") {
			normalized = strings.ReplaceAll(normalized, "\\\\", "\\")
		}

		if len(normalized) > 3 && strings.HasSuffix(normalized, "\\") {
			normalized = strings.TrimSuffix(normalized, "\\")
		}

		return normalized
	default:
		normalized := strings.ReplaceAll(path, "\\", "/")

		for strings.Contains(normalized, "//") {
			normalized = strings.ReplaceAll(normalized, "//", "/")
		}

		if len(normalized) > 1 && strings.HasSuffix(normalized, "/") {
			normalized = strings.TrimSuffix(normalized, "/")
		}

		return normalized
	}
}

// getRemoteHomeDirectory gets the home directory on remote machine via SSH
func getRemoteHomeDirectory(sshTarget, osType string) (string, error) {
	var cmd *exec.Cmd

	switch osType {
	case PlatformWindows:
		cmd = exec.Command("ssh", sshTarget, `powershell.exe -Command "Write-Output $env:USERPROFILE"`) // #nosec G204 -- SSH with validated arguments
	default:
		cmd = exec.Command("ssh", sshTarget, "echo $HOME") // #nosec G204 -- SSH with validated arguments
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%w on %s: %v", ErrHomeDirectory, osType, err)
	}

	homeDir := strings.TrimSpace(string(output))
	if homeDir == "" {
		return "", fmt.Errorf("%w from %s", ErrEmptyHomeDirectory, osType)
	}

	logger.Info(fmt.Sprintf("Remote home directory (%s): %s", osType, homeDir))
	return homeDir, nil
}

// expandEnvPathRemote expands environment variables on remote machine via SSH
func expandEnvPathRemote(sshTarget, path string, remoteOS string) (string, error) {
	if path == "" {
		return path, nil
	}

	// Check if path needs expansion (contains $, %, or ~ symbols)
	if !strings.Contains(path, "$") && !strings.Contains(path, "%") && !strings.HasPrefix(path, "~") {
		return path, nil
	}

	var cmd *exec.Cmd

	if remoteOS == PlatformWindows {
		windowsPath := path
		if strings.Contains(path, "$HOME") {
			// Replace $HOME with %USERPROFILE% for Windows
			windowsPath = strings.ReplaceAll(path, "$HOME", "%USERPROFILE%")
			windowsPath = strings.ReplaceAll(windowsPath, "${HOME}", "%USERPROFILE%")
		}

		escapedPath := strings.ReplaceAll(windowsPath, "'", "''")
		psCommand := fmt.Sprintf(`[System.Environment]::ExpandEnvironmentVariables('%s')`, escapedPath)

		cmd = exec.Command("ssh", sshTarget, fmt.Sprintf(`powershell.exe -Command "%s"`, psCommand)) // #nosec G204 -- SSH with validated arguments
	} else {
		// For Unix systems, use bash to expand both environment variables and tilde
		escapedPath := strings.ReplaceAll(path, `"`, `\"`)
		cmd = exec.Command("ssh", sshTarget, fmt.Sprintf(`bash -c "echo %s"`, escapedPath)) // #nosec G204 -- SSH with validated arguments
	}

	output, err := cmd.Output()
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to expand environment variables remotely on host %s: %v", sshTarget, err))
		return path, err // Return original path if expansion fails
	}

	expanded := strings.TrimSpace(string(output))

	normalizedPath := normalizePathForOS(expanded, remoteOS)

	logger.Info(fmt.Sprintf("Remote path expansion: '%s' -> '%s'", path, normalizedPath))
	return normalizedPath, nil
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"

	PlatformWindows = "windows"
	PlatformDarwin  = "darwin"
	PlatformLinux   = "linux"
	PlatformMac     = "mac"
	PlatformOSX     = "osx"

	ArchX64     = "x64"
	ArchAarch64 = "aarch64"
	ArchAll     = "all"

	ExtGz    = "gz"
	ExtTarGz = "tar.gz"
	ExtDmg   = "dmg"
	ExtExe   = "exe"
	ExtSit   = "sit"
	ExtZip   = "zip"

	TbcliName = "tbcli"

	IdeIdea      = "idea"
	IdePyCharm   = "pycharm"
	IdeGoland    = "goland"
	IdeClion     = "clion"
	IdeRider     = "rider"
	IdeRuby      = "ruby"
	IdeRubymine  = "rubymine"
	IdeRustRover = "rustrover"
	IdePhpStorm  = "phpstorm"
	IdeWebStorm  = "webstorm"

	DownloadsDir        = "downloads"
	RemoteDevServerName = "remote-dev-server"
)

func NewLogger(infoEnabled bool) *Logger {
	return &Logger{infoEnabled: infoEnabled}
}

func (l *Logger) Info(msg string) {
	if l.infoEnabled {
		fmt.Printf("%s[INFO]%s %s\n", colorBlue, colorReset, msg)
	}
}

func (l *Logger) Success(msg string) {
	fmt.Printf("%s[SUCCESS]%s %s\n", colorGreen, colorReset, msg)
}

func (l *Logger) Warning(msg string) {
	fmt.Printf("%s[WARNING]%s %s\n", colorYellow, colorReset, msg)
}

func (l *Logger) Error(msg string) {
	fmt.Printf("%s[ERROR]%s %s\n", colorRed, colorReset, msg)
}

// Default configuration
func getDefaultConfig() Config {
	return Config{
		DefaultTbcliVersion:   "2.8.1.52155",
		DefaultToolboxVersion: "2.8.1.52155",
		DefaultJbrVersion:     "21.0.3",
		DefaultJbrBuild:       "b509.11",
		DefaultIdeType:        "idea",
		DefaultIdeVersion:     "2025.1.3",
		TbcliInstallPath:      "",
		JbrInstallPath:        "",
		InstallIde:            false,
		InstallClient:         false,
		InfoLogging:           false,
		LaunchTbcli:           true,
		// Multi-IDE examples (commented out by default)
		// IdeSpecs: []ConfigIdeSpec{
		//     {Type: "idea", Version: "2025.1.3"},
		//     {Type: "pycharm", Version: "2025.1.3"},
		//     {Type: "webstorm", Version: "2025.1.3"},
		// },
		// Alternative: separate arrays
		// IdeTypes: []string{"idea", "pycharm", "webstorm"},
		// IdeVersions: []string{"2025.1.3", "2025.1.3", "2025.1.3"},
		JetbrainsBaseURL: "https://download-cdn.jetbrains.com",
		JbrBaseURL:       "https://cache-redirector.jetbrains.com",
		ProductsJSONURL:  "https://data.services.jetbrains.com/products",
		// VMOptions: []string{
		// 	"-Xmx4096m",
		// 	"-Xms512m",
		// },
		ToolboxEnvironment: `{
  "allowPortForwarding": true,
  "tools": {
    "allowInstallation": true,
    "allowUpdate": true,
    "allowUninstallation": true,
    "location": [
      {
        "path": "{{CUSTOM_IDE_PATH}}",
        "levels": 2
      }
    ]
  }
}`,
	}
}

type LegacyConfigUrls struct {
	ProductsInfoURL   string
	ClientDownloadURL string
}

// getRemoteDevPaths returns platform-specific paths for RemoteDev configuration
func getRemoteDevPaths() []string {
	var paths []string

	switch runtime.GOOS {
	case PlatformDarwin:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Warning("Failed to get home directory, using default paths")
			homeDir = ""
		}
		paths = []string{
			filepath.Join(homeDir, "Library/Application Support/JetBrains/RemoteDev"),
			"/Library/Application Support/JetBrains/RemoteDev",
		}
	case PlatformLinux:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Warning("Failed to get home directory, using default paths")
			homeDir = ""
		}
		paths = []string{
			filepath.Join(homeDir, ".config/JetBrains/RemoteDev"),
			"/etc/xdg/JetBrains/RemoteDev",
		}
	case PlatformWindows:
		// On Windows, we'll use registry instead of file paths
		paths = []string{}
	}

	return paths
}

// readFileContent safely reads content from a file
func readFileContent(filePath string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", err
	}

	content, err := os.ReadFile(filePath) // #nosec G304 -- File path is controlled by application
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// readWindowsRegistry reads values from Windows registry
func readWindowsRegistry() (*LegacyConfigUrls, error) {
	var urls LegacyConfigUrls

	registryPaths := []string{
		`HKEY_CURRENT_USER\SOFTWARE\JetBrains\RemoteDev`,
		`HKEY_LOCAL_MACHINE\SOFTWARE\JetBrains\RemoteDev`,
	}

	for _, regPath := range registryPaths {
		if productsURL := readRegistryValue(regPath, "productsInfoUrl"); productsURL != "" {
			urls.ProductsInfoURL = productsURL
		}

		if clientURL := readRegistryValue(regPath, "clientDownloadUrl"); clientURL != "" {
			urls.ClientDownloadURL = clientURL
		}

		if urls.ProductsInfoURL != "" && urls.ClientDownloadURL != "" {
			break
		}
	}

	if urls.ProductsInfoURL == "" && urls.ClientDownloadURL == "" {
		return nil, ErrNoRemoteDevEntries
	}

	return &urls, nil
}

// readRegistryValue reads a value from Windows registry using reg.exe
func readRegistryValue(keyPath, valueName string) string {
	cmd := exec.Command("reg", "query", keyPath, "/v", valueName)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, valueName) {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				for i, part := range parts {
					if part == "REG_SZ" || part == "REG_EXPAND_SZ" || part == "REG_MULTI_SZ" {
						if i+1 < len(parts) {
							return strings.Join(parts[i+1:], " ")
						}
					}
				}
			}
		}
	}

	return ""
}

// searchLegacyConfig searches for legacy RemoteDev configuration
func searchLegacyConfig() (*LegacyConfigUrls, error) {
	var urls LegacyConfigUrls

	if runtime.GOOS == PlatformWindows {
		return readWindowsRegistry()
	}

	paths := getRemoteDevPaths()

	for _, basePath := range paths {
		productsInfoPath := filepath.Join(basePath, "productsInfoUrl")
		if content, err := readFileContent(productsInfoPath); err == nil && content != "" {
			urls.ProductsInfoURL = content
		}

		clientDownloadPath := filepath.Join(basePath, "clientDownloadUrl")
		if content, err := readFileContent(clientDownloadPath); err == nil && content != "" {
			urls.ClientDownloadURL = content
		}

		if urls.ProductsInfoURL != "" && urls.ClientDownloadURL != "" {
			break
		}
	}

	if urls.ProductsInfoURL == "" && urls.ClientDownloadURL == "" {
		return nil, ErrNoRemoteDevConfig
	}

	return &urls, nil
}

// extractBaseURL extracts base URL from a RemoteDev URL
func extractBaseURL(url string) (string, error) {
	if url == "" {
		return "", ErrEmptyURL
	}

	// For productsInfoUrl like https://internal.site/backends/<PRODUCT_CODE>/products.json
	// we want to extract https://internal.site/
	re := regexp.MustCompile(`^(https?://[^/]+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("%w from: %s", ErrBaseURLExtraction, url)
	}

	return matches[1], nil
}

// loadLegacyConfig attempts to load configuration from legacy RemoteDev setup
func loadLegacyConfig() (Config, bool, error) {
	legacyUrls, err := searchLegacyConfig()
	if err != nil {
		return Config{}, false, err
	}

	var baseURL string
	if legacyUrls.ProductsInfoURL != "" {
		baseURL, err = extractBaseURL(legacyUrls.ProductsInfoURL)
		if err != nil && legacyUrls.ClientDownloadURL != "" {
			baseURL, err = extractBaseURL(legacyUrls.ClientDownloadURL)
		}
	} else if legacyUrls.ClientDownloadURL != "" {
		baseURL, err = extractBaseURL(legacyUrls.ClientDownloadURL)
	}

	if err != nil {
		return Config{}, false, fmt.Errorf("%w: %v", ErrBaseURLExtraction, err)
	}

	configURL := baseURL + "/config.json"
	resp, err := http.Get(configURL) // #nosec G107 -- URL is from trusted configuration source
	if err != nil {
		return Config{}, false, fmt.Errorf("%w from %s: %v", ErrConfigFetch, configURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Config{}, false, fmt.Errorf("%w at %s (HTTP %d)", ErrConfigNotFound, configURL, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Config{}, false, fmt.Errorf("%w from %s: %v", ErrConfigRead, configURL, err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, false, fmt.Errorf("%w from %s: %v", ErrConfigParse, configURL, err)
	}

	// Note: Don't expand IDE install path from legacy config here as it will be
	// expanded in the appropriate context (local or remote) later

	return config, true, nil
}

func generateConfigTemplate(path string) error {
	config := getDefaultConfig()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}

// Load config from various sources
func loadConfig(configPath string) (Config, error) {
	config := getDefaultConfig()

	var configFile string

	if configPath != "" {
		if strings.HasPrefix(configPath, "http://") || strings.HasPrefix(configPath, "https://") {
			resp, err := http.Get(configPath) // #nosec G107 -- Config path is from trusted source
			if err != nil {
				return config, fmt.Errorf("failed to download config from URL: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return config, fmt.Errorf("failed to download config: HTTP %d from URL: %s", resp.StatusCode, configPath)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return config, fmt.Errorf("failed to read config from URL: %v", err)
			}

			err = json.Unmarshal(data, &config)
			if err != nil {
				return config, fmt.Errorf("failed to parse config from URL: %v", err)
			}

			// Note: Don't expand IDE install path from URL config here as it will be
			// expanded in the appropriate context (local or remote) later

			return config, nil
		} else {
			configFile = configPath
		}
	} else {
		// Check for config.json next to executable
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			configCandidate := filepath.Join(execDir, "config.json")
			if _, err := os.Stat(configCandidate); err == nil {
				configFile = configCandidate
			}
		}

		// If no local config found, try legacy RemoteDev configuration
		if configFile == "" {
			legacyConfig, found, err := loadLegacyConfig()
			if err == nil && found {
				logger.Info("Found legacy RemoteDev configuration, applying remote config")
				return legacyConfig, nil
			} else if err != nil {
				logger.Info(fmt.Sprintf("No legacy RemoteDev configuration found: %v", err))
			}
		}
	}

	// Load local config file if found
	if configFile != "" {
		data, err := os.ReadFile(configFile) // #nosec G304 -- Config file path is from trusted source
		if err != nil {
			return config, fmt.Errorf("failed to read config file %s: %v", configFile, err)
		}

		err = json.Unmarshal(data, &config)
		if err != nil {
			return config, fmt.Errorf("failed to parse config file %s: %v", configFile, err)
		}

		// Note: Don't expand IDE install path from config here as it will be
		// expanded in the appropriate context (local or remote) later

		logger.Info(fmt.Sprintf("Loaded config from: %s", configFile))
	}

	return config, nil
}

// Detect system information
func detectSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	switch info.OS {
	case PlatformDarwin:
		info.OS = PlatformMac
		info.Name = "macOS"
	case PlatformWindows:
		info.Name = PlatformWindows
	case PlatformLinux:
		info.Name = PlatformLinux
	default:
		info.Name = info.OS
	}

	switch info.Arch {
	case "amd64":
		info.Arch = ArchX64
	case "arm64":
		info.Arch = ArchAarch64
	}

	return info
}

// Get platform-specific paths
func getPlatformPaths(osType string, homeDir string) (toolboxPath, idePath, clientPath string) {
	switch osType {
	case PlatformLinux:
		toolboxPath = joinRemotePathForOS(osType, homeDir, ".cache", "JetBrains", "Toolbox-CLI-dist")
		idePath = joinRemotePathForOS(osType, homeDir, ".local", "share", "JetBrains", "Toolbox", "apps")
		clientPath = joinRemotePathForOS(osType, homeDir, ".local", "share", "JetBrains", "Toolbox", "internal-tools")
	case PlatformMac:
		toolboxPath = joinRemotePathForOS(osType, homeDir, "Library", "Caches", "JetBrains", "Toolbox-CLI-dist")
		idePath = joinRemotePathForOS(osType, homeDir, "Applications")
		clientPath = joinRemotePathForOS(osType, homeDir, "Library", "Application Support", "JetBrains", "Toolbox", "internal-tools")
	case PlatformWindows:
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = joinRemotePathForOS(osType, homeDir, "AppData", "Local")
		}
		toolboxPath = joinRemotePathForOS(osType, localAppData, "JetBrains", "Toolbox-CLI-dist")
		idePath = joinRemotePathForOS(osType, localAppData, "Programs")
		clientPath = joinRemotePathForOS(osType, localAppData, "JetBrains", "Toolbox", "internal-tools")
	}
	return
}

// Generate IDE product name from type and version (e.g., "idea" + "2025.1.3" -> "IntelliJIdea2025.1")
func generateIDEProductName(ideType, ideVersion string) string {
	// Extract major.minor version (e.g., "2025.1.3" -> "2025.1")
	versionParts := strings.Split(ideVersion, ".")
	majorMinor := ideVersion
	if len(versionParts) >= 2 {
		majorMinor = versionParts[0] + "." + versionParts[1]
	}

	// Map IDE types to their product names
	productNameMap := map[string]string{
		"idea":      "IntelliJIdea",
		"pycharm":   "PyCharm",
		"goland":    "GoLand",
		"webstorm":  "WebStorm",
		"phpstorm":  "PhpStorm",
		"rider":     "Rider",
		"clion":     "CLion",
		"rubymine":  "RubyMine",
		"appcode":   "AppCode",
		"datagrip":  "DataGrip",
		"dataspell": "DataSpell",
	}

	productName, exists := productNameMap[ideType]
	if !exists {
		productName = strings.ToUpper(ideType[:1]) + ideType[1:]
	}

	return productName + majorMinor
}

// Get IDE-specific configuration directory paths for .vmoptions files
func getIDEConfigPaths(osType, homeDir, ideType, ideVersion string) (configDir, vmoptionsFile string) {
	productName := generateIDEProductName(ideType, ideVersion)

	switch osType {
	case PlatformLinux:
		configDir = joinRemotePathForOS(osType, homeDir, ".config", "JetBrains", productName)
		vmoptionsFile = ideType + "64.vmoptions"
	case PlatformMac:
		configDir = joinRemotePathForOS(osType, homeDir, "Library", "Application Support", "JetBrains", productName)
		vmoptionsFile = ideType + ".vmoptions"
	case PlatformWindows:
		// For Windows, construct path with proper backslashes
		// Don't use filepath.Join as it may use forward slashes on non-Windows hosts
		configDir = joinRemotePathForOS(PlatformWindows, homeDir, "AppData", "Roaming", "JetBrains", productName)
		vmoptionsFile = ideType + "64.exe.vmoptions"
	}

	return configDir, vmoptionsFile
}

// Create or update .vmoptions file locally
func createVMOptionsFile(configDir, vmoptionsFile string, vmOptions []string, osType string) error {
	if len(vmOptions) == 0 {
		return nil // No VM options to set
	}

	logger.Info("Configuring VM options for IDE...")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		logger.Info(fmt.Sprintf("Creating IDE config directory: %s", configDir))
		err = os.MkdirAll(configDir, 0o750)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}

		// Create migrate.config file
		migrateConfigPath := filepath.Join(configDir, "migrate.config")
		err = os.WriteFile(migrateConfigPath, []byte("merge-configs"), 0o600)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to create migrate.config: %v", err))
		} else {
			logger.Info("Created migrate.config with merge-configs directive")
		}
	}

	vmoptionsPath := filepath.Join(configDir, vmoptionsFile)

	var existingContent []string
	if _, err := os.Stat(vmoptionsPath); err == nil {
		data, err := os.ReadFile(vmoptionsPath)
		if err != nil {
			return fmt.Errorf("failed to read existing .vmoptions file: %v", err)
		}

		content := strings.TrimSpace(string(data))
		if content != "" {
			existingContent = strings.Split(content, "\n")
		}
		logger.Info(fmt.Sprintf("Found existing .vmoptions file %s with %d entries", vmoptionsPath, len(existingContent)))
	} else {
		logger.Info(fmt.Sprintf("Creating new .vmoptions file: %s", vmoptionsPath))
	}

	var newOptions []string
	for _, newOption := range vmOptions {
		newOption = strings.TrimSpace(newOption)
		if newOption == "" {
			continue
		}

		found := false
		for _, existing := range existingContent {
			if strings.TrimSpace(existing) == newOption {
				found = true
				break
			}
		}

		if !found {
			newOptions = append(newOptions, newOption)
		} else {
			logger.Info(fmt.Sprintf("VM option already exists, skipping: %s", newOption))
		}
	}

	if len(newOptions) == 0 {
		logger.Info("All VM options already exist in .vmoptions file")
		return nil
	}

	allOptions := append(existingContent, newOptions...)

	content := strings.Join(allOptions, "\n") + "\n"

	err := os.WriteFile(vmoptionsPath, []byte(content), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write .vmoptions file: %v", err)
	}

	logger.Success(fmt.Sprintf("VM options configured in %s", vmoptionsPath))
	for _, option := range newOptions {
		logger.Info(fmt.Sprintf("Added VM option: %s", option))
	}

	return nil
}

// Create or update .vmoptions file remotely via SSH
func createVMOptionsFileRemote(sshTarget, configDir, vmoptionsFile string, vmOptions []string, osType string) error {
	if len(vmOptions) == 0 {
		return nil // No VM options to set
	}

	logger.Info("Configuring VM options for IDE on remote host...")

	normalizedConfigDir := joinRemotePathForOS(osType, configDir)

	var checkDirCmd string
	if osType == PlatformWindows {
		checkDirCmd = fmt.Sprintf(`powershell.exe -Command "Test-Path '%s'"`, normalizedConfigDir)
	} else {
		checkDirCmd = fmt.Sprintf("[ -d %s ]", escapeUnixArgument(normalizedConfigDir))
	}

	cmd := exec.Command("ssh", sshTarget, checkDirCmd) // #nosec G204 -- SSH target is validated
	dirOutput, dirErr := cmd.CombinedOutput()

	// For Windows PowerShell, check the output content, not exit code
	var dirExists bool
	if osType == PlatformWindows {
		dirExists = strings.TrimSpace(string(dirOutput)) == "True"
	} else {
		dirExists = (dirErr == nil)
	}

	vmoptionsPath := joinRemotePathForOS(osType, normalizedConfigDir, vmoptionsFile)

	var checkFileExistsCmd string
	var readFileCmd string
	var fileExists bool

	if osType == PlatformWindows {
		checkFileExistsCmd = fmt.Sprintf(`powershell.exe -Command "Test-Path '%s'"`, vmoptionsPath)
		readFileCmd = fmt.Sprintf(`powershell.exe -Command "Get-Content '%s' -Raw"`, vmoptionsPath)
	} else {
		checkFileExistsCmd = fmt.Sprintf("[ -f %s ]", escapeUnixArgument(vmoptionsPath))
		readFileCmd = fmt.Sprintf("cat %s", escapeUnixArgument(vmoptionsPath))
	}

	cmd = exec.Command("ssh", sshTarget, checkFileExistsCmd) // #nosec G204 -- SSH target is validated
	output, err := cmd.CombinedOutput()

	// For Windows PowerShell, check the output content, not exit code
	if osType == PlatformWindows {
		fileExists = strings.TrimSpace(string(output)) == "True"
	} else {
		fileExists = (err == nil)
	}

	var existingContent []string
	if fileExists {
		cmd = exec.Command("ssh", sshTarget, readFileCmd) // #nosec G204 -- SSH target is validated
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			content := strings.TrimSpace(string(output))
			if content != "" {
				// For Windows with -Raw flag, content comes as single string with actual line breaks
				// Need to handle both \r\n (Windows) and \n (Unix) line endings
				if osType == PlatformWindows {
					// Normalize Windows line endings to Unix style
					content = strings.ReplaceAll(content, "\r\n", "\n")
				}
				lines := strings.Split(content, "\n")
				// Filter out empty lines and trim whitespace
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" {
						existingContent = append(existingContent, line)
					}
				}
			}
		}
		logger.Info(fmt.Sprintf("Found existing .vmoptions file %s with %d entries", vmoptionsPath, len(existingContent)))
	} else {
		logger.Info(fmt.Sprintf("Creating new .vmoptions file: %s", vmoptionsPath))

		if !dirExists {
			logger.Info(fmt.Sprintf("Creating IDE config directory on remote host: %s", normalizedConfigDir))

			var createDirCmd string
			if osType == PlatformWindows {
				createDirCmd = fmt.Sprintf(`powershell.exe -Command "New-Item -ItemType Directory -Force -Path '%s'"`, normalizedConfigDir)
			} else {
				createDirCmd = fmt.Sprintf("mkdir -p %s", escapeUnixArgument(normalizedConfigDir))
			}

			cmd = exec.Command("ssh", sshTarget, createDirCmd) // #nosec G204 -- SSH target is validated
			output, err := cmd.CombinedOutput()
			if err != nil {
				logger.Info(fmt.Sprintf("Create directory command: %s", createDirCmd))
				logger.Info(fmt.Sprintf("Create directory output: %s", strings.TrimSpace(string(output))))
				return fmt.Errorf("failed to create config directory: %v, output: %s", err, string(output))
			}
			logger.Info("Created IDE config directory successfully")

			// Create migrate.config file
			migrateConfigPath := joinRemotePathForOS(osType, normalizedConfigDir, "migrate.config")
			var createMigrateCmd string
			if osType == PlatformWindows {
				createMigrateCmd = fmt.Sprintf(`powershell.exe -Command "Set-Content -Path '%s' -Value 'merge-configs'"`, migrateConfigPath)
			} else {
				createMigrateCmd = fmt.Sprintf("echo 'merge-configs' > %s", escapeUnixArgument(migrateConfigPath))
			}

			cmd = exec.Command("ssh", sshTarget, createMigrateCmd) // #nosec G204 -- SSH target is validated
			output, err = cmd.CombinedOutput()
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to create migrate.config: %v, output: %s", err, string(output)))
			} else {
				logger.Info("Created migrate.config with merge-configs directive")
			}
		}
	}

	// Filter out VM options that already exist
	var newOptions []string
	for _, newOption := range vmOptions {
		newOption = strings.TrimSpace(newOption)
		if newOption == "" {
			continue
		}

		found := false
		for _, existing := range existingContent {
			if strings.TrimSpace(existing) == newOption {
				found = true
				break
			}
		}

		if !found {
			newOptions = append(newOptions, newOption)
		} else {
			logger.Info(fmt.Sprintf("VM option already exists, skipping: %s", newOption))
		}
	}

	if len(newOptions) == 0 {
		logger.Info("All VM options already exist in .vmoptions file")
		return nil
	}

	allOptions := append(existingContent, newOptions...)
	content := strings.Join(allOptions, "\n") + "\n"

	if osType == PlatformWindows {
		// For Windows, write each line separately for maximum reliability
		// First, clear/create the file
		clearFileCmd := fmt.Sprintf(`powershell.exe -Command "Set-Content -Path '%s' -Value '' -Encoding UTF8"`, vmoptionsPath)
		cmd = exec.Command("ssh", sshTarget, clearFileCmd) // #nosec G204 -- SSH target is validated
		output, err := cmd.CombinedOutput()
		if err != nil {
			logger.Info(fmt.Sprintf("Clear file command: %s", clearFileCmd))
			logger.Info(fmt.Sprintf("Clear file output: %s", strings.TrimSpace(string(output))))
			return fmt.Errorf("failed to create .vmoptions file: %v, output: %s", err, string(output))
		}

		lines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				escapedLine := strings.ReplaceAll(line, "'", "''")
				addLineCmd := fmt.Sprintf(`powershell.exe -Command "Add-Content -Path '%s' -Value '%s' -Encoding UTF8"`, vmoptionsPath, escapedLine)
				cmd = exec.Command("ssh", sshTarget, addLineCmd) // #nosec G204 -- SSH target is validated
				output, err := cmd.CombinedOutput()
				if err != nil {
					logger.Info(fmt.Sprintf("Add line command: %s", addLineCmd))
					logger.Info(fmt.Sprintf("Add line output: %s", strings.TrimSpace(string(output))))
					return fmt.Errorf("failed to add line to .vmoptions file: %v, output: %s", err, string(output))
				}
			}
		}
	} else {
		// Use cat with here-document for better handling of special characters
		writeFileCmd := fmt.Sprintf("cat > %s << 'EOF'\n%sEOF", escapeUnixArgument(vmoptionsPath), content)
		cmd = exec.Command("ssh", sshTarget, writeFileCmd) // #nosec G204 -- SSH target is validated
		output, err := cmd.CombinedOutput()
		if err != nil {
			logger.Info(fmt.Sprintf("Write command: %s", writeFileCmd))
			logger.Info(fmt.Sprintf("Write command output: %s", strings.TrimSpace(string(output))))
			return fmt.Errorf("failed to write .vmoptions file: %v, output: %s", err, string(output))
		}
	}

	logger.Success(fmt.Sprintf("VM options configured in %s", vmoptionsPath))
	for _, option := range newOptions {
		logger.Info(fmt.Sprintf("Added VM option: %s", option))
	}

	return nil
}

// Sync versions between tbcli and toolbox
func syncVersions(opts *InstallOptions) {
	if opts.UserSetTbcliVersion && !opts.UserSetToolboxVersion {
		opts.ToolboxVersion = opts.TbcliVersion
		logger.Info(fmt.Sprintf("Synchronized Toolbox version to TBCLI version: %s", opts.ToolboxVersion))
	} else if !opts.UserSetTbcliVersion && opts.UserSetToolboxVersion {
		opts.TbcliVersion = opts.ToolboxVersion
		logger.Info(fmt.Sprintf("Synchronized TBCLI version to Toolbox version: %s", opts.TbcliVersion))
	}
}

// Normalize OS names from user input
func normalizeOS(osName string) string {
	switch strings.ToLower(osName) {
	case "win", "win32", "win64", PlatformWindows:
		return PlatformWindows
	case PlatformMac, "macos", PlatformDarwin:
		return PlatformMac
	case PlatformLinux, "ubuntu", "debian", "centos", "rhel":
		return PlatformLinux
	case ArchAll:
		return ArchAll
	default:
		return osName
	}
}

// generateJBRURL generates JBR download URL
func generateJBRURL(config Config, version, os, arch, build string) (string, error) {
	var jbrOS string
	switch os {
	case PlatformLinux:
		jbrOS = PlatformLinux
	case PlatformMac:
		jbrOS = PlatformOSX
	case PlatformWindows:
		jbrOS = PlatformWindows
	default:
		return "", fmt.Errorf("unsupported OS for JBR: %s", os)
	}

	return fmt.Sprintf("%s/intellij-jbr/jbr-%s-%s-%s-%s.%s",
		config.JbrBaseURL, version, jbrOS, arch, build, ExtTarGz), nil
}

// generateTBCLIURL generates Toolbox CLI download URL
func generateTBCLIURL(config Config, version, os, arch string) (string, error) {
	tbcliOS := os
	if os == PlatformWindows {
		tbcliOS = "win32"
	}

	return fmt.Sprintf("%s/toolbox/agent/%s-%s-%s-%s-nojre.%s",
		config.JetbrainsBaseURL, TbcliName, tbcliOS, arch, version, ExtTarGz), nil
}

// generateToolboxURL generates Toolbox application download URL
func generateToolboxAppURL(config Config, version, os, arch string) (string, error) {
	var extension string
	switch os {
	case PlatformLinux:
		extension = ExtTarGz
	case PlatformMac:
		extension = ExtDmg
	case PlatformWindows:
		extension = ExtExe
	default:
		return "", fmt.Errorf("unsupported OS for Toolbox: %s", os)
	}

	productName := "jetbrains-toolbox"
	if arch == ArchAarch64 {
		return fmt.Sprintf("%s/toolbox/%s-%s-arm64.%s",
			config.JetbrainsBaseURL, productName, version, extension), nil
	}
	return fmt.Sprintf("%s/toolbox/%s-%s.%s",
		config.JetbrainsBaseURL, productName, version, extension), nil
}

// generateIDEURL generates IDE download URL
func generateIDEURL(config Config, productType, version, os, arch string) (string, error) {
	var extension string
	switch os {
	case PlatformLinux:
		extension = ExtTarGz
	case PlatformMac:
		extension = ExtDmg
	case PlatformWindows:
		extension = ExtExe
	default:
		return "", fmt.Errorf("unsupported OS for IDE %s: %s", productType, os)
	}

	var pathPrefix, productName string
	switch productType {
	case IdeIdea:
		pathPrefix = "idea"
		productName = "ideaIU"
	case IdePyCharm:
		pathPrefix = "python"
		productName = "pycharm"
	case IdeGoland:
		pathPrefix = "go"
		productName = "goland"
	case IdeClion:
		pathPrefix = "cpp"
		productName = "CLion"
	case IdeRider:
		pathPrefix = "rider"
		productName = "JetBrains.Rider"
	case IdeRuby, IdeRubymine:
		pathPrefix = "ruby"
		productName = "RubyMine"
	case IdeRustRover:
		pathPrefix = "rustrover"
		productName = "RustRover"
	default:
		return "", fmt.Errorf("unsupported IDE type: %s", productType)
	}

	archSuffix := ""
	if arch == ArchAarch64 {
		archSuffix = "-aarch64"
	}

	return fmt.Sprintf("%s/%s/%s-%s%s.%s",
		config.JetbrainsBaseURL, pathPrefix, productName, version, archSuffix, extension), nil
}

// Generate JetBrains download URL
func generateJetBrainsURL(config Config, productType, version, os, arch, build string) (string, error) {
	switch productType {
	case "jbr":
		return generateJBRURL(config, version, os, arch, build)
	case TbcliName:
		return generateTBCLIURL(config, version, os, arch)
	case "toolbox":
		return generateToolboxAppURL(config, version, os, arch)

	case IdeIdea, IdePyCharm, IdeGoland, IdeClion, IdeRider, IdeRuby, IdeRubymine, IdeRustRover:
		return generateIDEURL(config, productType, version, os, arch)

	default:
		return "", fmt.Errorf("unsupported product type: %s", productType)
	}
}

// Check if product is already installed by examining product-info.json
func checkExistingInstallation(installPath, ideType, ideVersion string, sysInfo SystemInfo) (InstallationResult, error) {
	result := InstallationResult{Exists: false, AppPath: ""}

	if forceReinstall {
		logger.Info("Force reinstallation enabled, skipping existence check")
		return result, nil
	}

	if sysInfo.OS == PlatformMac {
		entries, err := os.ReadDir(installPath)
		if err != nil {
			return result, err
		}

		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".app") {
				appPath := filepath.Join(installPath, entry.Name())
				productInfoPath := filepath.Join(appPath, "Contents", "Resources", "product-info.json")

				if _, err := os.Stat(productInfoPath); os.IsNotExist(err) {
					continue
				}

				data, err := os.ReadFile(productInfoPath)
				if err != nil {
					continue
				}

				var productInfo ProductInfo
				err = json.Unmarshal(data, &productInfo)
				if err != nil {
					continue
				}

				if isMatchingIDE(productInfo, ideType, ideVersion) {
					logger.Info(fmt.Sprintf("Found existing installation: %s %s (build %s) at %s",
						productInfo.Name, productInfo.Version, productInfo.BuildNumber, appPath))
					result.Exists = true
					result.AppPath = appPath
					return result, nil
				}
			}
		}
	} else {
		ideDir := fmt.Sprintf("%s-%s", ideType, ideVersion)
		productInfoPath := filepath.Join(installPath, ideDir, "product-info.json")

		if _, err := os.Stat(productInfoPath); os.IsNotExist(err) {
			return result, nil
		}

		data, err := os.ReadFile(productInfoPath)
		if err != nil {
			return result, err
		}

		var productInfo ProductInfo
		err = json.Unmarshal(data, &productInfo)
		if err != nil {
			return result, err
		}

		if productInfo.Version == ideVersion {
			logger.Info(fmt.Sprintf("Found existing installation: %s %s (build %s)",
				productInfo.Name, productInfo.Version, productInfo.BuildNumber))
			result.Exists = true
			return result, nil
		}
	}

	return result, nil
}

func isMatchingIDE(productInfo ProductInfo, ideType, ideVersion string) bool {
	// Check version match
	if productInfo.Version != ideVersion {
		return false
	}

	switch ideType {
	case IdeIdea:
		return productInfo.ProductCode == "IU" || productInfo.ProductCode == "IC" ||
			strings.Contains(strings.ToLower(productInfo.Name), "intellij idea")
	case IdePyCharm:
		return productInfo.ProductCode == "PY" || productInfo.ProductCode == "PC" ||
			strings.Contains(strings.ToLower(productInfo.Name), "pycharm")
	case IdeGoland:
		return productInfo.ProductCode == "GO" ||
			strings.Contains(strings.ToLower(productInfo.Name), "goland")
	case IdeClion:
		return productInfo.ProductCode == "CL" ||
			strings.Contains(strings.ToLower(productInfo.Name), "clion")
	case IdeRider:
		return productInfo.ProductCode == "RD" ||
			strings.Contains(strings.ToLower(productInfo.Name), "rider")
	case IdeRuby, IdeRubymine:
		return productInfo.ProductCode == "RM" ||
			strings.Contains(strings.ToLower(productInfo.Name), IdeRubymine)
	case IdeRustRover:
		return productInfo.ProductCode == "RR" ||
			strings.Contains(strings.ToLower(productInfo.Name), "rustrover")
	case IdeWebStorm:
		return productInfo.ProductCode == "WS" ||
			strings.Contains(strings.ToLower(productInfo.Name), "webstorm")
	case IdePhpStorm:
		return productInfo.ProductCode == "PS" ||
			strings.Contains(strings.ToLower(productInfo.Name), "phpstorm")
	default:
		return strings.Contains(strings.ToLower(productInfo.Name), strings.ToLower(ideType))
	}
}

// generateUniqueAppName creates a unique app name by checking existing directories
func generateUniqueAppName(installPath, appName string, isRemote bool, sshTarget string) (string, error) {
	if isRemote {
		return generateUniqueAppNameRemote(installPath, appName, sshTarget)
	}

	baseName := strings.TrimSuffix(appName, ".app")
	originalPath := filepath.Join(installPath, appName)

	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		return appName, nil
	}

	for i := 2; i <= 999; i++ {
		uniqueName := fmt.Sprintf("%s %d.app", baseName, i)
		uniquePath := filepath.Join(installPath, uniqueName)

		if _, err := os.Stat(uniquePath); os.IsNotExist(err) {
			return uniqueName, nil
		}
	}

	return "", fmt.Errorf("unable to generate unique app name after 999 attempts")
}

// generateUniqueAppNameRemote generates unique app name for remote installation
func generateUniqueAppNameRemote(installPath, appName, sshTarget string) (string, error) {
	baseName := strings.TrimSuffix(appName, ".app")
	originalPath := joinRemotePathForOS(PlatformMac, installPath, appName)

	checkCmd := exec.Command("ssh", sshTarget, fmt.Sprintf("[ ! -e '%s' ]", originalPath)) // #nosec G204 -- SSH target is validated
	if checkCmd.Run() == nil {
		return appName, nil // Original name is available
	}

	for i := 2; i <= 999; i++ {
		uniqueName := fmt.Sprintf("%s %d.app", baseName, i)
		uniquePath := joinRemotePathForOS(PlatformMac, installPath, uniqueName)

		checkCmd := exec.Command("ssh", sshTarget, fmt.Sprintf("[ ! -e '%s' ]", uniquePath)) // #nosec G204 -- SSH target is validated
		if checkCmd.Run() == nil {
			return uniqueName, nil
		}
	}

	return "", fmt.Errorf("unable to generate unique app name after 999 attempts")
}

// findIDEAppPath finds the app path for macOS by scanning directory via SSH or locally
func findIDEAppPath(idePath, ideType, ideVersion string, isRemote bool, sshTarget string) (string, error) {
	if isRemote {
		// For remote, scan via SSH
		cmd := exec.Command("ssh", sshTarget, fmt.Sprintf("find '%s' -name '*.app' -type d -maxdepth 1", idePath)) // #nosec G204 -- SSH target is validated
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}

		apps := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, appPath := range apps {
			if appPath == "" {
				continue
			}

			productInfoPath := joinRemotePathForOS(PlatformMac, appPath, "Contents", "Resources", "product-info.json")
			checkCmd := exec.Command("ssh", sshTarget, fmt.Sprintf("cat %s", escapeUnixArgument(productInfoPath))) // #nosec G204 -- SSH target is validated
			data, err := checkCmd.Output()
			if err != nil {
				continue
			}

			var productInfo ProductInfo
			if json.Unmarshal(data, &productInfo) == nil && isMatchingIDE(productInfo, ideType, ideVersion) {
				return appPath, nil
			}
		}
		return "", fmt.Errorf("no matching .app found")
	} else {
		// For local, use existing logic
		entries, err := os.ReadDir(idePath)
		if err != nil {
			return "", err
		}

		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".app") {
				appPath := filepath.Join(idePath, entry.Name())
				productInfoPath := filepath.Join(appPath, "Contents", "Resources", "product-info.json")

				data, err := os.ReadFile(productInfoPath)
				if err != nil {
					continue
				}

				var productInfo ProductInfo
				if json.Unmarshal(data, &productInfo) == nil && isMatchingIDE(productInfo, ideType, ideVersion) {
					return appPath, nil
				}
			}
		}
		return "", fmt.Errorf("no matching .app found")
	}
}

// Download file with progress (like bash version)
func downloadFile(url, filename string) error {
	logger.Info(fmt.Sprintf("Downloading %s...", filename))

	var urlPath string

	if idx := strings.Index(url, "://"); idx != -1 {
		if remaining := url[idx+3:]; remaining != "" {
			if slashIdx := strings.Index(remaining, "/"); slashIdx != -1 {
				fullPath := remaining[slashIdx+1:]
				if lastSlash := strings.LastIndex(fullPath, "/"); lastSlash != -1 {
					urlPath = fullPath[:lastSlash]
				}
			}
		}
	}

	var localDir string
	if urlPath != "" {
		localDir = filepath.Join("./downloads", urlPath)
	} else {
		localDir = "./downloads"
	}

	destPath := filepath.Join(localDir, filename)

	err := os.MkdirAll(localDir, 0o750)
	if err != nil {
		return err
	}

	if _, statErr := os.Stat(destPath); statErr == nil {
		logger.Info(fmt.Sprintf("File %s already exists, skipping download", filename))
		return nil
	}

	resp, err := http.Get(url) // #nosec G107 -- URL is validated before calling this function
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: HTTP %d from URL: %s", resp.StatusCode, url)
	}

	out, err := os.Create(destPath) // #nosec G304 -- File path is controlled by application
	if err != nil {
		return err
	}
	defer out.Close()

	if resp.ContentLength > 0 && !disableProgress {
		pw := NewProgressWriter(resp.ContentLength, filename)
		_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
	} else {
		if disableProgress {
			logger.Info(fmt.Sprintf("Downloading %s...", filename))
		} else {
			logger.Info("Downloading (size unknown)...")
		}
		_, err = io.Copy(out, resp.Body)
	}

	if err != nil {
		return err
	}

	logger.Success(fmt.Sprintf("Downloaded %s", filename))
	if logger.infoEnabled {
		logger.Info(fmt.Sprintf("Saved to: %s", destPath))
	}
	return nil
}

// Extract tar.gz archive
func extractTarGz(srcPath, destDir string) error {
	logger.Info(fmt.Sprintf("Extracting %s...", filepath.Base(srcPath)))

	file, err := os.Open(srcPath) // #nosec G304 -- File path is controlled by application
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		pathParts := strings.Split(header.Name, "/")
		if len(pathParts) <= 1 {
			continue
		}
		destPath := filepath.Join(destDir, filepath.Join(pathParts[1:]...))

		switch header.Typeflag {
		case tar.TypeDir:
			mode := os.FileMode(header.Mode) & 0o777
			if mode&0o700 == 0 {
				mode |= 0o750
			}
			err = os.MkdirAll(destPath, mode)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to create directory %s: %v", destPath, err))
				err = os.MkdirAll(destPath, 0o750)
				if err != nil {
					return err
				}
			}
		case tar.TypeReg:
			err = os.MkdirAll(filepath.Dir(destPath), 0o750)
			if err != nil {
				return err
			}

			outFile, err := os.Create(destPath) // #nosec G304 -- File path is controlled by application
			if err != nil {
				return err
			}

			// #nosec G110 -- Limited by tar header size field
			_, err = io.Copy(outFile, tr)
			_ = outFile.Close() // #nosec G104 -- File close after successful operation
			if err != nil {
				return err
			}

			mode := os.FileMode(header.Mode) & 0o777
			if mode&0o400 == 0 {
				mode |= 0o644
			}
			err = os.Chmod(destPath, mode)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to set permissions for %s: %v", destPath, err))
			}
		}
	}

	logger.Success(fmt.Sprintf("Extracted %s", filepath.Base(srcPath)))
	return nil
}

// setupLocalInstallation initializes paths and system info for local installation
func setupLocalInstallation(opts InstallOptions) (SystemInfo, string, string, string, error) {
	sysInfo := detectSystemInfo()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return SystemInfo{}, "", "", "", fmt.Errorf("failed to get home directory: %v", err)
	}

	toolboxPath, idePath, _ := getPlatformPaths(sysInfo.OS, homeDir)

	// Apply custom paths if specified
	if opts.TbcliInstallPath != "" {
		toolboxPath = expandEnvPath(opts.TbcliInstallPath)
	}
	if opts.IdeInstallPath != "" {
		idePath = expandEnvPath(opts.IdeInstallPath)
	}

	logger.Info(fmt.Sprintf("Installing on local system: %s (%s)", sysInfo.Name, sysInfo.Arch))
	logger.Info(fmt.Sprintf("Toolbox CLI path: %s", toolboxPath))
	logger.Info(fmt.Sprintf("IDE path: %s", idePath))

	downloadDir := DownloadsDir
	err = os.MkdirAll(downloadDir, 0o750)
	if err != nil {
		return SystemInfo{}, "", "", "", fmt.Errorf("failed to create download directory: %v", err)
	}

	return sysInfo, homeDir, toolboxPath, idePath, nil
}

// installTBCLILocally handles TBCLI installation locally
func installTBCLILocally(config Config, opts InstallOptions, sysInfo SystemInfo, toolboxPath string) (bool, error) {
	tbcliDir := filepath.Join(toolboxPath, fmt.Sprintf("tbcli-%s", opts.TbcliVersion))
	extractedFile := filepath.Join(tbcliDir, ".extracted")

	// Check if TBCLI is already installed
	if !forceReinstall {
		if _, err := os.Stat(extractedFile); err == nil {
			logger.Info(fmt.Sprintf("Toolbox CLI already installed at: %s", tbcliDir))
			return false, nil
		}
	}

	tbcliURL, err := generateJetBrainsURL(config, TbcliName, opts.TbcliVersion, sysInfo.OS, sysInfo.Arch, "")
	if err != nil {
		return false, fmt.Errorf("failed to generate Toolbox CLI URL: %v", err)
	}

	logger.Info(fmt.Sprintf("Toolbox CLI URL: %s", tbcliURL))
	err = downloadFile(tbcliURL, filepath.Base(tbcliURL))
	if err != nil {
		return false, fmt.Errorf("failed to download Toolbox CLI: %v", err)
	}

	return handleIDEInstallationLogic(tbcliDir, tbcliURL, forceReinstall)
}

// installJBRLocally handles JBR installation locally
func installJBRLocally(config Config, opts InstallOptions, sysInfo SystemInfo, toolboxPath string) (bool, error) {
	var jbrOS string
	switch sysInfo.OS {
	case PlatformLinux:
		jbrOS = PlatformLinux
	case PlatformMac:
		jbrOS = PlatformOSX
	case PlatformWindows:
		jbrOS = PlatformWindows
	}

	// Use custom JBR path if specified, otherwise use toolbox path
	jbrBasePath := toolboxPath
	if opts.JbrInstallPath != "" {
		jbrBasePath = expandEnvPath(opts.JbrInstallPath)
	}

	jbrDir := filepath.Join(jbrBasePath, fmt.Sprintf("jbr-%s-%s-%s-%s-1", opts.JbrVersion, jbrOS, sysInfo.Arch, opts.JbrBuild))
	extractedFile := filepath.Join(jbrDir, ".extracted")

	// Check if JBR is already installed
	if !forceReinstall {
		if _, err := os.Stat(extractedFile); err == nil {
			logger.Info(fmt.Sprintf("JBR already installed at: %s", jbrDir))
			return false, nil
		}
	}

	jbrURL, err := generateJetBrainsURL(config, "jbr", opts.JbrVersion, sysInfo.OS, sysInfo.Arch, opts.JbrBuild)
	if err != nil {
		return false, fmt.Errorf("failed to generate JBR URL: %v", err)
	}

	logger.Info(fmt.Sprintf("JBR URL: %s", jbrURL))
	err = downloadFile(jbrURL, filepath.Base(jbrURL))
	if err != nil {
		return false, fmt.Errorf("failed to download JBR: %v", err)
	}

	return installJBRInDirectory(jbrURL, jbrDir)
}

// installJBRInDirectory installs JBR in the specified directory
func installJBRInDirectory(jbrURL, jbrDir string) (bool, error) {
	// Handle force reinstallation - remove existing directory
	if forceReinstall {
		logger.Info("Force reinstallation enabled, reinstalling JBR")
		if _, statErr := os.Stat(jbrDir); statErr == nil {
			logger.Info(fmt.Sprintf("Removing existing JBR directory: %s", jbrDir))
			err := os.RemoveAll(jbrDir)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to remove existing JBR directory: %v", err))
			}
		}
	}

	return extractAndMarkJBR(jbrURL, jbrDir)
}

// extractAndMarkJBR extracts JBR and creates installation marker
func extractAndMarkJBR(jbrURL, jbrDir string) (bool, error) {
	err := os.MkdirAll(jbrDir, 0o750)
	if err != nil {
		return false, fmt.Errorf("failed to create JBR directory: %v", err)
	}

	jbrFile, findErr := findDownloadedFile(filepath.Base(jbrURL))
	if findErr != nil {
		return false, fmt.Errorf("failed to find downloaded JBR file: %v", findErr)
	}

	err = extractTarGz(jbrFile, jbrDir)
	if err != nil {
		logger.Warning(fmt.Sprintf("Extraction failed, cleaning up directory: %s", jbrDir))
		_ = os.RemoveAll(jbrDir) // #nosec G104 -- Best effort cleanup
		return false, fmt.Errorf("failed to extract JBR: %v", err)
	}

	// Create .extracted marker file - critical for installation state tracking
	extractedFile := filepath.Join(jbrDir, ".extracted")
	err = os.WriteFile(extractedFile, []byte(""), 0o600)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to create .extracted marker, cleaning up: %v", err))
		_ = os.RemoveAll(jbrDir) // #nosec G104 -- Best effort cleanup
		return false, fmt.Errorf("failed to create installation marker: %v", err)
	}
	return true, nil // Mark as newly installed
}

// Install locally (new mode)
func installLocally(config Config, opts InstallOptions) error {
	var tbcliInstalled, jbrInstalled, ideInstalled bool

	sysInfo, homeDir, toolboxPath, idePath, err := setupLocalInstallation(opts)
	if err != nil {
		return err
	}

	tbcliInstalled, err = installTBCLILocally(config, opts, sysInfo, toolboxPath)
	if err != nil {
		return err
	}

	jbrInstalled, err = installJBRLocally(config, opts, sysInfo, toolboxPath)
	if err != nil {
		return err
	}

	// Create environment files for custom paths
	// Calculate actual JBR installation directory
	var jbrPath string

	// Determine JBR base path (custom or default)
	jbrBasePath := toolboxPath
	if opts.JbrInstallPath != "" {
		jbrBasePath = expandEnvPath(opts.JbrInstallPath)
	}

	// Calculate OS string for JBR directory name
	var jbrOS string
	switch sysInfo.OS {
	case PlatformLinux:
		jbrOS = PlatformLinux
	case PlatformMac:
		jbrOS = PlatformOSX
	case PlatformWindows:
		jbrOS = PlatformWindows
	}

	// Always construct the full JBR directory path
	jbrPath = filepath.Join(jbrBasePath, fmt.Sprintf("jbr-%s-%s-%s-%s-1", opts.JbrVersion, jbrOS, sysInfo.Arch, opts.JbrBuild))

	// Calculate full TBCLI directory path
	var tbcliBasePath string
	if opts.TbcliInstallPath != "" {
		tbcliBasePath = expandEnvPath(opts.TbcliInstallPath)
	} else {
		tbcliBasePath = toolboxPath
	}
	tbcliFullPath := filepath.Join(tbcliBasePath, fmt.Sprintf("tbcli-%s", opts.TbcliVersion))

	err = createEnvFiles(opts, sysInfo.OS, sysInfo.Arch, tbcliFullPath, jbrPath, false, "")
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to create environment files: %v", err))
	}

	var installedAppPath string

	if opts.InstallIde {
		installResult, installErr := checkExistingInstallation(idePath, opts.IdeType, opts.IdeVersion, sysInfo)
		if installErr != nil {
			logger.Warning(fmt.Sprintf("Failed to check existing installation: %v", installErr))
		} else if installResult.Exists {
			logger.Info(fmt.Sprintf("IDE %s %s is already installed", opts.IdeType, opts.IdeVersion))
			installedAppPath = installResult.AppPath
		} else {
			ideURL, ideURLErr := generateJetBrainsURL(config, opts.IdeType, opts.IdeVersion, sysInfo.OS, sysInfo.Arch, "")
			if ideURLErr != nil {
				return fmt.Errorf("failed to generate IDE URL: %v", ideURLErr)
			}

			logger.Info(fmt.Sprintf("IDE URL: %s", ideURL))
			err = downloadFile(ideURL, filepath.Base(ideURL))
			if err != nil {
				return fmt.Errorf("failed to download IDE: %v", err)
			}

			// Find downloaded file
			ideFile, findErr := findDownloadedFile(filepath.Base(ideURL))
			if findErr != nil {
				return fmt.Errorf("failed to find downloaded IDE file: %v", findErr)
			}

			appPath, installErr := installIDE(ideFile, idePath, opts.IdeType, opts.IdeVersion, sysInfo)
			if installErr != nil {
				return fmt.Errorf("failed to install IDE: %v", installErr)
			}
			installedAppPath = appPath
			ideInstalled = true // Mark as newly installed
		}

		if opts.IdeInstallPath != "" {
			expandedPath := expandEnvPath(opts.IdeInstallPath)
			err = createLocalToolboxEnvironmentConfig(homeDir, expandedPath, config, sysInfo.OS)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to create Toolbox environment configuration: %v", err))
			}
		}

		if len(config.VMOptions) > 0 {
			configDir, vmoptionsFile := getIDEConfigPaths(sysInfo.OS, homeDir, opts.IdeType, opts.IdeVersion)
			err = createVMOptionsFile(configDir, vmoptionsFile, config.VMOptions, sysInfo.OS)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to configure VM options: %v", err))
			}
		}

		if opts.LaunchIdeArgs != "" {
			err = launchIDELocally(idePath, opts.IdeType, opts.IdeVersion, sysInfo.OS, opts.LaunchIdeArgs, installedAppPath)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to launch IDE: %v", err))
			}
		}

		if opts.InstallClient {
			remoteInfo := RemoteSystemInfo(sysInfo)
			err = installClientLocally(config, opts, remoteInfo)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to install JetBrains Client locally: %v", err))
			}
		}
	}

	if opts.InstallClient && !opts.InstallIde && !opts.NoClient {
		logger.Info("Installing JetBrains Client separately (without IDE)...")
		remoteInfo := RemoteSystemInfo(sysInfo)
		err = installClientLocally(config, opts, remoteInfo)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to install JetBrains Client: %v", err))
		}
	}

	if opts.LaunchIdeArgs != "" && !opts.InstallIde && opts.IdeType != "" {
		logger.Info("Checking for existing IDE to launch...")
		installResult, checkErr := checkExistingInstallation(idePath, opts.IdeType, opts.IdeVersion, sysInfo)
		if checkErr != nil {
			logger.Warning(fmt.Sprintf("Failed to check existing installation: %v", checkErr))
		} else if installResult.Exists {
			logger.Info(fmt.Sprintf("Found existing IDE %s %s, launching...", opts.IdeType, opts.IdeVersion))

			if len(config.VMOptions) > 0 {
				configDir, vmoptionsFile := getIDEConfigPaths(sysInfo.OS, homeDir, opts.IdeType, opts.IdeVersion)
				err = createVMOptionsFile(configDir, vmoptionsFile, config.VMOptions, sysInfo.OS)
				if err != nil {
					logger.Warning(fmt.Sprintf("Failed to configure VM options: %v", err))
				}
			}

			err = launchIDELocally(idePath, opts.IdeType, opts.IdeVersion, sysInfo.OS, opts.LaunchIdeArgs, installResult.AppPath)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to launch IDE: %v", err))
			}
		} else {
			logger.Warning(fmt.Sprintf("IDE %s %s not found for launch. Use --install-ide to install first.", opts.IdeType, opts.IdeVersion))
		}
	}

	if opts.LaunchTbcli && (tbcliInstalled || jbrInstalled || ideInstalled) {
		logger.Info("Components were installed, restarting TBCLi...")
		err = launchTbcli(config, opts, sysInfo.OS, sysInfo.Arch, toolboxPath, false, "")
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to launch TBCLi: %v", err))
		}
	} else if opts.LaunchTbcli {
		logger.Info("No components were installed, TBCLi restart not needed")
	}

	logger.Success("Local installation completed successfully!")
	return nil
}

// Install IDE based on platform and file type
func installIDE(ideFile, installPath, ideType, ideVersion string, sysInfo SystemInfo) (string, error) {
	switch filepath.Ext(ideFile) {
	case ExtGz:
		ideDir := filepath.Join(installPath, fmt.Sprintf("%s-%s", ideType, ideVersion))
		err := os.MkdirAll(ideDir, 0o750)
		if err != nil {
			return "", err
		}
		err = extractTarGz(ideFile, ideDir)
		return "", err

	case ".dmg":
		return installDMG(ideFile, installPath)

	case ".exe":
		err := installEXE(ideFile, installPath, ideType, ideVersion)
		return "", err

	default:
		return "", fmt.Errorf("unsupported IDE file type: %s", filepath.Ext(ideFile))
	}
}

// Install DMG on macOS
func installDMG(dmgFile, installPath string) (string, error) {
	mountPoint := joinRemotePathForOS(PlatformMac, "/tmp", "ide_mount_"+strconv.FormatInt(time.Now().Unix(), 10))

	err := os.MkdirAll(mountPoint, 0o750)
	if err != nil {
		return "", fmt.Errorf("failed to create mount point: %v", err)
	}

	cmd := exec.Command("hdiutil", "attach", dmgFile, "-mountpoint", mountPoint, "-nobrowse", "-quiet") // #nosec G204 -- macOS utility with validated arguments
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to mount DMG: %v, output: %s", err, string(output))
	}
	defer func() {
		detachCmd := exec.Command("hdiutil", "detach", mountPoint, "-quiet") // #nosec G204 -- macOS utility with validated arguments
		if detachErr := detachCmd.Run(); detachErr != nil {
			logger.Warning(fmt.Sprintf("Failed to run detach command: %v", detachErr))
		}
	}()

	time.Sleep(2 * time.Second)

	entries, err := os.ReadDir(mountPoint)
	if err != nil {
		return "", fmt.Errorf("failed to read mount point: %v", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			src := filepath.Join(mountPoint, entry.Name())

			uniqueAppName, err := generateUniqueAppName(installPath, entry.Name(), false, "")
			if err != nil {
				return "", fmt.Errorf("failed to generate unique app name: %v", err)
			}

			dst := filepath.Join(installPath, uniqueAppName)

			if uniqueAppName != entry.Name() {
				logger.Info(fmt.Sprintf("App %s already exists, installing as %s", entry.Name(), uniqueAppName))
			}

			cmd := exec.Command("ditto", src, dst) // #nosec G204 -- macOS utility with validated paths
			output, err := cmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("failed to copy app bundle: %v, output: %s", err, string(output))
			}

			logger.Success(fmt.Sprintf("Installed %s", uniqueAppName))
			return dst, nil
		}
	}

	return "", errors.New("no .app bundle found in DMG")
}

// Install EXE on Windows
func installEXE(exeFile, installPath, ideType, ideVersion string) error {
	ideDir := filepath.Join(installPath, fmt.Sprintf("%s-%s", ideType, ideVersion))

	cmd := exec.Command(exeFile, "/S", "/D="+ideDir) // #nosec G204 -- Installer executable path is validated
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run installer: %v", err)
	}

	logger.Success(fmt.Sprintf("Installed %s %s", ideType, ideVersion))
	return nil
}

// installTBCLIRemotely installs Toolbox CLI on remote host
func installTBCLIRemotely(config Config, opts InstallOptions, remoteInfo RemoteSystemInfo, remotePaths RemotePaths) (bool, error) {
	tbcliExists, _, err := checkSSHExistingInstallation(opts.SSHTarget, remotePaths.ToolboxPath, opts, remoteInfo)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to check existing TBCLI installation: %v", err))
	}

	if tbcliExists {
		logger.Info("Toolbox CLI already installed")
		return false, nil
	}

	tbcliURL, err := generateJetBrainsURL(config, TbcliName, opts.TbcliVersion, remoteInfo.OS, remoteInfo.Arch, "")
	if err != nil {
		return false, fmt.Errorf("failed to generate Toolbox CLI URL: %v", err)
	}

	err = transferAndInstallSSHWithURL(opts.SSHTarget, "", tbcliURL, remotePaths.ToolboxPath, TbcliName, opts.TbcliVersion, remoteInfo.OS)
	if err != nil {
		return false, fmt.Errorf("failed to install Toolbox CLI: %v", err)
	}

	return true, nil
}

// installJBRRemotely installs JBR on remote host
func installJBRRemotely(config Config, opts InstallOptions, remoteInfo RemoteSystemInfo, remotePaths RemotePaths) (bool, error) {
	_, jbrExists, err := checkSSHExistingInstallation(opts.SSHTarget, remotePaths.ToolboxPath, opts, remoteInfo)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to check existing JBR installation: %v", err))
	}

	if jbrExists {
		logger.Info("JBR already installed")
		return false, nil
	}

	jbrURL, err := generateJetBrainsURL(config, "jbr", opts.JbrVersion, remoteInfo.OS, remoteInfo.Arch, opts.JbrBuild)
	if err != nil {
		return false, fmt.Errorf("failed to generate JBR URL: %v", err)
	}

	var jbrOS string
	switch remoteInfo.OS {
	case PlatformLinux:
		jbrOS = PlatformLinux
	case PlatformMac:
		jbrOS = PlatformOSX
	case PlatformWindows:
		jbrOS = PlatformWindows
	}

	// Determine JBR installation path (custom or default)
	var jbrInstallPath string
	if opts.JbrInstallPath != "" {
		expandedPath, expandErr := expandEnvPathRemote(opts.SSHTarget, opts.JbrInstallPath, remoteInfo.OS)
		if expandErr != nil {
			logger.Warning(fmt.Sprintf("Failed to expand JBR install path remotely for installation: %v", expandErr))
			jbrInstallPath = opts.JbrInstallPath
		} else {
			jbrInstallPath = expandedPath
		}
	} else {
		jbrInstallPath = remotePaths.ToolboxPath
	}

	jbrDirName := fmt.Sprintf("jbr-%s-%s-%s-%s-1", opts.JbrVersion, jbrOS, remoteInfo.Arch, opts.JbrBuild)
	err = transferAndInstallSSHWithURL(opts.SSHTarget, "", jbrURL, jbrInstallPath, "jbr", jbrDirName, remoteInfo.OS)
	if err != nil {
		return false, fmt.Errorf("failed to install JBR: %v", err)
	}

	return true, nil
}

// installIDEsRemotely installs IDEs on remote host
func installIDEsRemotely(config Config, opts InstallOptions, remoteInfo RemoteSystemInfo, remotePaths RemotePaths) (bool, error) {
	if !opts.InstallIde {
		return false, nil
	}

	ideSpecs := opts.IdeSpecs
	if len(ideSpecs) == 0 {
		logger.Warning("No IDE specifications found for installation")
		return false, nil
	}

	logger.Info(fmt.Sprintf("Installing %d IDE(s)...", len(ideSpecs)))
	ideInstalled := false

	for i, spec := range ideSpecs {
		logger.Info(fmt.Sprintf("[%d/%d] Processing IDE: %s %s", i+1, len(ideSpecs), spec.Type, spec.Version))

		idePath := remotePaths.IDEPath
		if opts.IdeInstallPath != "" {
			expandedPath, err := expandEnvPathRemote(opts.SSHTarget, opts.IdeInstallPath, remoteInfo.OS)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to expand IDE install path remotely on host %s, using original: %v", opts.SSHTarget, err))
				idePath = opts.IdeInstallPath
			} else {
				idePath = expandedPath
			}
		}

		ideExists, err := checkSSHExistingIDE(opts.SSHTarget, idePath, spec.Type, spec.Version, remoteInfo)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to check existing IDE installation for %s %s: %v", spec.Type, spec.Version, err))
		}

		if !ideExists {
			ideOpts := opts
			ideOpts.IdeType = spec.Type
			ideOpts.IdeVersion = spec.Version

			err = installIDEOnRemote(config, ideOpts, opts.SSHTarget, remotePaths, remoteInfo)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to install IDE %s %s: %v", spec.Type, spec.Version, err))
				continue
			}
			ideInstalled = true
			logger.Success(fmt.Sprintf("Successfully installed IDE: %s %s", spec.Type, spec.Version))
		} else {
			logger.Info(fmt.Sprintf("IDE %s %s is already installed, skipping installation", spec.Type, spec.Version))
		}
	}

	return ideInstalled, nil
}

// handlePostInstallationTasks handles configuration and launch tasks after installation
func handlePostInstallationTasks(config Config, opts InstallOptions, remoteInfo RemoteSystemInfo, remotePaths RemotePaths) error {
	if opts.IdeInstallPath != "" {
		expandedPath, err := expandEnvPathRemote(opts.SSHTarget, opts.IdeInstallPath, remoteInfo.OS)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to expand IDE install path remotely for toolbox config on host %s, using original: %v", opts.SSHTarget, err))
			expandedPath = opts.IdeInstallPath
		}
		err = createToolboxEnvironmentConfig(opts.SSHTarget, expandedPath, config, remoteInfo)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to create Toolbox environment configuration: %v", err))
		}
	}

	if err := handleRemoteVMOptionsConfiguration(config, opts, remoteInfo); err != nil {
		logger.Warning(fmt.Sprintf("Failed to configure VM options: %v", err))
	}

	if opts.LaunchIdeArgs != "" {
		idePath := remotePaths.IDEPath
		if opts.IdeInstallPath != "" {
			expandedPath, err := expandEnvPathRemote(opts.SSHTarget, opts.IdeInstallPath, remoteInfo.OS)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to expand IDE install path remotely for launch on host %s, using original: %v", opts.SSHTarget, err))
				idePath = opts.IdeInstallPath
			} else {
				idePath = expandedPath
			}
		}

		err := launchIDERemotely(opts.SSHTarget, idePath, opts.IdeType, opts.IdeVersion, remoteInfo.OS, opts.LaunchIdeArgs)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to launch IDE remotely: %v", err))
		}
	}

	return nil
}

// Install via SSH
func installViaSSH(config Config, opts InstallOptions) error {
	logger.Success(fmt.Sprintf("Connecting to remote host: %s", opts.SSHTarget))

	var tbcliInstalled, jbrInstalled, ideInstalled bool

	remoteInfo, err := detectRemoteSystem(opts.SSHTarget)
	if err != nil {
		return fmt.Errorf("failed to detect remote system: %v", err)
	}

	logger.Info(fmt.Sprintf("Remote system: %s (%s)", remoteInfo.Name, remoteInfo.Arch))

	remotePaths, err := getRemotePathsWithOptions(opts.SSHTarget, remoteInfo.OS, opts)
	if err != nil {
		return fmt.Errorf("failed to get remote paths: %v", err)
	}

	downloadDir := DownloadsDir
	err = os.MkdirAll(downloadDir, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create download directory: %v", err)
	}

	// Install TBCLI
	tbcliInstalled, err = installTBCLIRemotely(config, opts, remoteInfo, remotePaths)
	if err != nil {
		return err
	}

	// Install JBR
	jbrInstalled, err = installJBRRemotely(config, opts, remoteInfo, remotePaths)
	if err != nil {
		return err
	}

	// Create environment files for custom paths
	// Calculate actual JBR installation directory for remote installation
	var jbrPath string

	// Determine JBR base path (custom or default)
	var jbrBasePath string
	if opts.JbrInstallPath != "" {
		expandedPath, expandErr := expandEnvPathRemote(opts.SSHTarget, opts.JbrInstallPath, remoteInfo.OS)
		if expandErr != nil {
			logger.Warning(fmt.Sprintf("Failed to expand JBR install path remotely for env file creation: %v", expandErr))
			jbrBasePath = opts.JbrInstallPath
		} else {
			jbrBasePath = expandedPath
		}
	} else {
		jbrBasePath = remotePaths.ToolboxPath
	}

	// Calculate OS string for JBR directory name
	var jbrOS string
	switch remoteInfo.OS {
	case PlatformLinux:
		jbrOS = PlatformLinux
	case PlatformMac:
		jbrOS = PlatformOSX
	case PlatformWindows:
		jbrOS = PlatformWindows
	}

	// Always construct the full JBR directory path
	jbrPath = joinRemotePathForOS(remoteInfo.OS, jbrBasePath, fmt.Sprintf("jbr-%s-%s-%s-%s-1", opts.JbrVersion, jbrOS, remoteInfo.Arch, opts.JbrBuild))

	// Calculate full TBCLI directory path for remote
	var tbcliRemoteBasePath string
	if opts.TbcliInstallPath != "" {
		expandedPath, expandErr := expandEnvPathRemote(opts.SSHTarget, opts.TbcliInstallPath, remoteInfo.OS)
		if expandErr != nil {
			logger.Warning(fmt.Sprintf("Failed to expand TBCLI install path remotely for env file creation: %v", expandErr))
			tbcliRemoteBasePath = opts.TbcliInstallPath
		} else {
			tbcliRemoteBasePath = expandedPath
		}
	} else {
		tbcliRemoteBasePath = remotePaths.ToolboxPath
	}
	tbcliRemoteFullPath := joinRemotePathForOS(remoteInfo.OS, tbcliRemoteBasePath, fmt.Sprintf("tbcli-%s", opts.TbcliVersion))

	err = createEnvFiles(opts, remoteInfo.OS, remoteInfo.Arch, tbcliRemoteFullPath, jbrPath, true, opts.SSHTarget)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to create remote environment files: %v", err))
	}

	// Install IDEs
	ideInstalled, err = installIDEsRemotely(config, opts, remoteInfo, remotePaths)
	if err != nil {
		return err
	}

	if opts.InstallClient {
		err = installClientLocally(config, opts, remoteInfo)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to install JetBrains Client locally: %v", err))
		}
	}

	// Handle post-installation tasks
	if opts.InstallIde {
		err = handlePostInstallationTasks(config, opts, remoteInfo, remotePaths)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to handle post-installation tasks: %v", err))
		}
	}

	// Install JetBrains Client separately if requested but IDE installation is not enabled
	if opts.InstallClient && !opts.InstallIde && !opts.NoClient {
		logger.Info("Installing JetBrains Client separately (without IDE) via SSH...")
		err = installClientLocally(config, opts, remoteInfo)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to install JetBrains Client locally: %v", err))
		}
	}

	// Launch IDE if arguments provided but not installing (check if already installed)
	if opts.LaunchIdeArgs != "" && !opts.InstallIde && opts.IdeType != "" {
		logger.Info("Checking for existing IDE to launch...")
		idePath := remotePaths.IDEPath
		if opts.IdeInstallPath != "" {
			expandedPath, expandErr := expandEnvPathRemote(opts.SSHTarget, opts.IdeInstallPath, remoteInfo.OS)
			if expandErr != nil {
				logger.Warning(fmt.Sprintf("Failed to expand IDE install path remotely for existing check on host %s, using original: %v", opts.SSHTarget, expandErr))
				idePath = opts.IdeInstallPath
			} else {
				idePath = expandedPath
			}
		}

		exists, checkErr := checkSSHExistingIDE(opts.SSHTarget, idePath, opts.IdeType, opts.IdeVersion, remoteInfo)
		if checkErr != nil {
			logger.Warning(fmt.Sprintf("Failed to check existing IDE installation: %v", checkErr))
		} else if exists {
			logger.Info(fmt.Sprintf("Found existing IDE %s %s, launching...", opts.IdeType, opts.IdeVersion))

			if vmErr := handleRemoteVMOptionsConfiguration(config, opts, remoteInfo); vmErr != nil {
				logger.Warning(fmt.Sprintf("Failed to configure VM options: %v", vmErr))
			}

			err = launchIDERemotely(opts.SSHTarget, idePath, opts.IdeType, opts.IdeVersion, remoteInfo.OS, opts.LaunchIdeArgs)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to launch IDE remotely: %v", err))
			}
		} else {
			logger.Warning(fmt.Sprintf("IDE %s %s not found for launch. Use --install-ide to install first.", opts.IdeType, opts.IdeVersion))
		}
	}

	// Launch TBCLi if requested and something was actually installed
	if opts.LaunchTbcli && (tbcliInstalled || jbrInstalled || ideInstalled) {
		logger.Info("Components were installed, restarting TBCLi...")
		err = launchTbcli(config, opts, remoteInfo.OS, remoteInfo.Arch, remotePaths.ToolboxPath, true, opts.SSHTarget)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to launch TBCLi via SSH: %v", err))
		}
	} else if opts.LaunchTbcli {
		logger.Info("No components were installed, TBCLi restart not needed")
	}

	logger.Success("SSH installation completed successfully!")

	// Generate and display JetBrains Toolbox SSH URL if SSH installation
	sshInfo := parseSSHConnectionString(opts.SSHTarget)

	includeIDELaunch := opts.LaunchIdeArgs != ""
	var buildNumber string
	var primaryIdeType string
	var primaryIdeVersion string

	if len(opts.IdeSpecs) > 0 {
		primaryIdeType = opts.IdeSpecs[0].Type
		primaryIdeVersion = opts.IdeSpecs[0].Version
	} else if opts.IdeType != "" {
		primaryIdeType = opts.IdeType
		primaryIdeVersion = opts.IdeVersion
	}

	if includeIDELaunch && primaryIdeType != "" && primaryIdeVersion != "" {
		ideURL, err := generateJetBrainsURL(config, primaryIdeType, primaryIdeVersion, remoteInfo.OS, remoteInfo.Arch, "")
		if err == nil {
			ideFilename := filepath.Base(ideURL)
			if bn, err := getIDEBuildNumber(ideFilename); err == nil {
				buildNumber = bn
			}
		}
	}

	// Generate Toolbox URL
	toolboxURL := generateToolboxURL(sshInfo, buildNumber, primaryIdeType, includeIDELaunch)

	// Open Toolbox URL automatically if requested
	if opts.OpenToolboxURL {
		logger.Info("Opening JetBrains Toolbox SSH URL in default application...")
		err := openURLInBrowser(toolboxURL)
		if err != nil {
			fmt.Printf("\nJetBrains Toolbox SSH URL:\n%s\n", toolboxURL)
			logger.Warning(fmt.Sprintf("Failed to open Toolbox URL automatically: %v", err))
			logger.Info("You can manually copy and open the URL above in JetBrains Toolbox.")
		} else {
			logger.Success("Toolbox URL opened successfully!")
			fmt.Printf("\nJetBrains Toolbox SSH URL:\n%s\n", toolboxURL)
		}
	} else {
		fmt.Printf("\nJetBrains Toolbox SSH URL:\n%s\n", toolboxURL)
	}

	return nil
}

// Install JetBrains Client locally for SSH remote IDE access
func installClientLocally(config Config, opts InstallOptions, remoteInfo RemoteSystemInfo) error {
	logger.Info("Installing JetBrains Client locally for remote IDE access...")

	// Download products JSON if not already downloaded
	if !downloadProductsJSON(config) {
		return fmt.Errorf("failed to download products metadata for Client installation")
	}

	localSysInfo := detectSystemInfo()
	logger.Info(fmt.Sprintf("Installing Client on local system: %s (%s)", localSysInfo.Name, localSysInfo.Arch))

	ideURL, err := generateJetBrainsURL(config, opts.IdeType, opts.IdeVersion, remoteInfo.OS, remoteInfo.Arch, "")
	if err != nil {
		return fmt.Errorf("failed to generate IDE URL for build number: %v", err)
	}

	ideFilename := filepath.Base(ideURL)
	buildNumber, err := getIDEBuildNumber(ideFilename)
	if err != nil || buildNumber == "" {
		return fmt.Errorf("failed to get IDE build number for Client: %v", err)
	}

	logger.Info(fmt.Sprintf("Using IDE build number: %s", buildNumber))

	clientURL, err := generateClientURL(config, buildNumber, localSysInfo.OS, localSysInfo.Arch, opts.IdeType)
	if err != nil {
		return fmt.Errorf("failed to generate Client URL: %v", err)
	}

	clientFilename := filepath.Base(clientURL)
	logger.Info(fmt.Sprintf("Client URL: %s", clientURL))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	_, _, clientPath := getPlatformPaths(localSysInfo.OS, homeDir)

	var clientDirName string
	switch filepath.Ext(clientFilename) {
	case ExtGz:
		clientDirName = strings.TrimSuffix(clientFilename, ".tar.gz")
	case ExtSit:
		clientDirName = strings.TrimSuffix(clientFilename, ExtSit)
	case ExtZip:
		if strings.HasSuffix(clientFilename, ".jbr.win.zip") {
			clientDirName = strings.TrimSuffix(clientFilename, ".jbr.win.zip")
		} else {
			clientDirName = strings.TrimSuffix(clientFilename, ExtZip)
		}
	default:
		clientDirName = strings.TrimSuffix(clientFilename, filepath.Ext(clientFilename))
	}

	clientInstallDir := filepath.Join(clientPath, clientDirName)
	if isClientInstalled(clientInstallDir, localSysInfo.OS) {
		logger.Info(fmt.Sprintf("JetBrains Client already installed at: %s", clientInstallDir))
		logger.Success("JetBrains Client installation check completed")
		logger.Info(fmt.Sprintf("You can use JetBrains Client from: %s", clientPath))
		return nil
	}

	err = downloadFile(clientURL, clientFilename)
	if err != nil {
		return fmt.Errorf("failed to download JetBrains Client: %v", err)
	}

	clientFile, err := findDownloadedFile(clientFilename)
	if err != nil {
		return fmt.Errorf("failed to find downloaded Client file: %v", err)
	}

	logger.Info(fmt.Sprintf("Client installation path: %s", clientPath))

	err = installClientFile(clientFile, clientPath, clientFilename, localSysInfo)
	if err != nil {
		return fmt.Errorf("failed to install JetBrains Client: %v", err)
	}

	logger.Success("JetBrains Client installed locally")
	logger.Info(fmt.Sprintf("You can now connect to remote IDE using JetBrains Client from: %s", clientPath))
	return nil
}

// Install client file locally based on OS and file type
func installClientFile(clientFile, installPath, filename string, sysInfo SystemInfo) error {
	logger.Info(fmt.Sprintf("Installing Client: %s", filename))

	err := os.MkdirAll(installPath, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create Client installation directory: %v", err)
	}

	var clientDirName string
	switch filepath.Ext(filename) {
	case ExtGz:
		// Handle .tar.gz
		clientDirName = strings.TrimSuffix(filename, ".tar.gz")
	case ExtSit:
		clientDirName = strings.TrimSuffix(filename, ExtSit)
	case ExtZip:
		if strings.HasSuffix(filename, ".jbr.win.zip") {
			clientDirName = strings.TrimSuffix(filename, ".jbr.win.zip")
		} else {
			clientDirName = strings.TrimSuffix(filename, ExtZip)
		}
	default:
		clientDirName = strings.TrimSuffix(filename, filepath.Ext(filename))
	}

	clientInstallDir := filepath.Join(installPath, clientDirName)

	if isClientInstalled(clientInstallDir, sysInfo.OS) {
		logger.Info(fmt.Sprintf("JetBrains Client already installed at: %s", clientInstallDir))
		return nil
	}

	switch filepath.Ext(clientFile) {
	case ExtGz:
		err = os.MkdirAll(clientInstallDir, 0o750)
		if err != nil {
			return err
		}
		return extractTarGz(clientFile, clientInstallDir)

	case ExtSit:
		err = os.MkdirAll(clientInstallDir, 0o750)
		if err != nil {
			return err
		}
		return extractSitFile(clientFile, clientInstallDir)

	case ExtZip:
		return extractZipFile(clientFile, clientInstallDir)

	default:
		return fmt.Errorf("unsupported Client file type: %s", filepath.Ext(clientFile))
	}
}

// Extract .sit file on macOS
func extractSitFile(sitFile, targetDir string) error {
	logger.Info(fmt.Sprintf("Extracting .sit archive: %s", filepath.Base(sitFile)))

	cmd := exec.Command("unzip", "-q", sitFile, "-d", targetDir)
	err := cmd.Run()
	if err == nil {
		logger.Success("Extracted .sit archive using unzip")
		return nil
	}

	cmd = exec.Command("ditto", "-x", "-k", sitFile, targetDir) // #nosec G204 -- macOS utility with validated paths
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to extract .sit archive: %v", err)
	}

	logger.Success("Extracted .sit archive using ditto")
	return nil
}

// Extract ZIP file using native Go archive/zip
func extractZipFile(zipFile, installDir string) error {
	logger.Info(fmt.Sprintf("Extracting ZIP archive: %s", filepath.Base(zipFile)))

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %v", err)
	}
	defer reader.Close()

	err = os.MkdirAll(installDir, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create installation directory: %v", err)
	}

	for _, file := range reader.File {
		if file.FileInfo().IsDir() || strings.HasPrefix(filepath.Base(file.Name), ".") {
			continue
		}

		targetPath := filepath.Join(installDir, filepath.FromSlash(file.Name))

		if !strings.HasPrefix(targetPath, filepath.Clean(installDir)+string(os.PathSeparator)) {
			logger.Warning(fmt.Sprintf("Skipping file outside target directory: %s", file.Name))
			continue
		}

		err = os.MkdirAll(filepath.Dir(targetPath), 0o750)
		if err != nil {
			return fmt.Errorf("failed to create directory structure for %s: %v", file.Name, err)
		}

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s in archive: %v", file.Name, err)
		}

		fileMode := file.FileInfo().Mode()
		if fileMode&0o400 == 0 {
			fileMode |= 0o644
		}

		outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileMode) // #nosec G304 -- File path is controlled by application
		if err != nil {
			_ = rc.Close() // #nosec G104 -- Resource cleanup
			return fmt.Errorf("failed to create target file %s: %v", targetPath, err)
		}

		// #nosec G110 -- Size limited by ZIP metadata
		_, err = io.Copy(outFile, rc)
		_ = rc.Close()      // #nosec G104 -- Resource cleanup
		_ = outFile.Close() // #nosec G104 -- File close after successful operation

		if err != nil {
			return fmt.Errorf("failed to extract file %s: %v", file.Name, err)
		}

		if logger.infoEnabled {
			logger.Info(fmt.Sprintf("Extracted: %s", file.Name))
		}
	}

	logger.Success("Extracted ZIP archive")
	return nil
}

// Check if JetBrains Client is properly installed by verifying product-info.json
func isClientInstalled(clientInstallDir, osType string) bool {
	if forceReinstall {
		logger.Info("Force reinstallation enabled, skipping Client existence check")
		return false
	}

	if _, err := os.Stat(clientInstallDir); os.IsNotExist(err) {
		return false
	}

	var productInfoPaths []string

	switch osType {
	case PlatformMac:
		entries, err := os.ReadDir(clientInstallDir)
		if err != nil {
			return false
		}

		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".app") {
				appPath := filepath.Join(clientInstallDir, entry.Name())
				productInfoPaths = append(productInfoPaths, filepath.Join(appPath, "Contents", "Resources", "product-info.json"))
			}
		}
	case PlatformLinux:
		productInfoPaths = []string{
			filepath.Join(clientInstallDir, "product-info.json"),
			filepath.Join(clientInstallDir, "bin", "product-info.json"),
		}
	case PlatformWindows:
		productInfoPaths = []string{
			filepath.Join(clientInstallDir, "product-info.json"),
			filepath.Join(clientInstallDir, "bin", "product-info.json"),
		}

		windowsExecutables := []string{
			filepath.Join(clientInstallDir, "bin", "jetbrains_client.exe"),
		}

		var hasExecutable bool
		for _, exePath := range windowsExecutables {
			if _, err := os.Stat(exePath); err == nil {
				hasExecutable = true
				if logger.infoEnabled {
					logger.Info(fmt.Sprintf("Found JetBrains Client executable: %s", exePath))
				}
				break
			}
		}

		// For Windows, if we find an executable but no product-info.json,
		// still consider it a valid installation
		if hasExecutable {
			// Check if any product-info.json exists
			for _, path := range productInfoPaths {
				if _, err := os.Stat(path); err == nil {
					return true
				}
			}
			// If executable exists but no product-info.json found, still consider valid
			logger.Info("Windows client executable found, considering installation valid")
			return true
		}
	}

	// Check if any product-info.json exists and is readable
	for _, productInfoPath := range productInfoPaths {
		if _, err := os.Stat(productInfoPath); err == nil {
			if logger.infoEnabled {
				logger.Info(fmt.Sprintf("Found product-info.json: %s", productInfoPath))
			}
			return true
		}
	}

	return false
}

// Remote system info
type RemoteSystemInfo struct {
	OS   string
	Arch string
	Name string
}

// Remote paths
type RemotePaths struct {
	ToolboxPath string
	IDEPath     string
	ClientPath  string
}

// Detect remote system via SSH
func detectRemoteSystem(sshTarget string) (RemoteSystemInfo, error) {
	var info RemoteSystemInfo

	cmd := exec.Command("ssh", sshTarget, "uname -s") // #nosec G204 -- SSH target is validated
	output, err := cmd.Output()

	if err != nil {
		logger.Info("Unix commands failed, trying Windows detection...")

		cmd = exec.Command("ssh", sshTarget, `powershell.exe -Command "$env:OS"`) // #nosec G204 -- SSH target is validated
		output, err = cmd.Output()
		if err != nil {
			return info, fmt.Errorf("failed to detect OS (tried both Unix and Windows): %v", err)
		}

		osName := strings.TrimSpace(string(output))
		if osName == "Windows_NT" {
			info.OS = PlatformWindows
			info.Name = "Windows"
			logger.Info("Detected Windows system")
		} else {
			return info, fmt.Errorf("unknown Windows OS type: %s", osName)
		}
	} else {
		// Unix commands worked
		osName := strings.TrimSpace(string(output))
		switch osName {
		case "Linux":
			info.OS = PlatformLinux
			info.Name = "Linux"
		case PlatformDarwin:
			info.OS = PlatformMac
			info.Name = "macOS"
		default:
			return info, fmt.Errorf("unsupported Unix OS: %s", osName)
		}
		logger.Info(fmt.Sprintf("Detected Unix system: %s", info.Name))
	}

	if info.OS == PlatformWindows {
		cmd = exec.Command("ssh", sshTarget, `powershell.exe -Command "$env:PROCESSOR_ARCHITECTURE"`) // #nosec G204 -- SSH target is validated
		output, err = cmd.Output()
		if err != nil {
			return info, fmt.Errorf("failed to detect Windows architecture: %v", err)
		}

		archName := strings.TrimSpace(string(output))
		switch archName {
		case "AMD64", "x86_64":
			info.Arch = ArchX64
		case "ARM64", ArchAarch64:
			info.Arch = ArchAarch64
		case "x86":
			info.Arch = "x86"
		default:
			logger.Warning(fmt.Sprintf("Unknown Windows architecture: %s, defaulting to x64", archName))
			info.Arch = ArchX64
		}
	} else {
		cmd = exec.Command("ssh", sshTarget, "uname -m") // #nosec G204 -- SSH target is validated
		output, err = cmd.Output()
		if err != nil {
			return info, fmt.Errorf("failed to detect Unix architecture: %v", err)
		}

		archName := strings.TrimSpace(string(output))
		switch archName {
		case "x86_64", "amd64":
			info.Arch = ArchX64
		case ArchAarch64, "arm64":
			info.Arch = ArchAarch64
		default:
			logger.Warning(fmt.Sprintf("Unknown Unix architecture: %s, defaulting to x64", archName))
			info.Arch = ArchX64
		}
	}

	logger.Info(fmt.Sprintf("Detected system: %s (%s)", info.Name, info.Arch))
	return info, nil
}

// Get remote paths via SSH
func getRemotePaths(sshTarget, osType string) (RemotePaths, error) {
	return getRemotePathsWithOptions(sshTarget, osType, InstallOptions{})
}

// Get remote paths via SSH with custom options
func getRemotePathsWithOptions(sshTarget, osType string, opts InstallOptions) (RemotePaths, error) {
	var paths RemotePaths
	var homeDir string

	if osType == PlatformWindows {
		cmd := exec.Command("ssh", sshTarget, `powershell.exe -Command "$env:USERPROFILE"`) // #nosec G204 -- SSH target is validated
		output, err := cmd.Output()
		if err != nil {
			return paths, fmt.Errorf("failed to get Windows home directory: %v", err)
		}
		homeDir = strings.TrimSpace(string(output))
		logger.Info(fmt.Sprintf("Windows home directory: %s", homeDir))
	} else {
		var err error
		homeDir, err = getRemoteHomeDirectory(sshTarget, osType)
		if err != nil {
			return paths, fmt.Errorf("failed to get Unix home directory: %v", err)
		}
		logger.Info(fmt.Sprintf("Unix home directory: %s", homeDir))
	}

	switch osType {
	case PlatformLinux:
		paths.ToolboxPath = joinRemotePathForOS(osType, homeDir, ".cache", "JetBrains", "Toolbox-CLI-dist")
		paths.IDEPath = joinRemotePathForOS(osType, homeDir, ".local", "share", "JetBrains", "Toolbox", "apps")
		paths.ClientPath = joinRemotePathForOS(osType, homeDir, ".local", "share", "JetBrains", "Toolbox", "internal-tools")
	case PlatformMac:
		paths.ToolboxPath = joinRemotePathForOS(osType, homeDir, "Library", "Caches", "JetBrains", "Toolbox-CLI-dist")
		paths.IDEPath = joinRemotePathForOS(osType, homeDir, "Applications")
		paths.ClientPath = joinRemotePathForOS(osType, homeDir, "Library", "Application Support", "JetBrains", "Toolbox", "internal-tools")
	case PlatformWindows:
		// Get LOCALAPPDATA for Windows
		cmd := exec.Command("ssh", sshTarget, `powershell.exe -Command "$env:LOCALAPPDATA"`) // #nosec G204 -- SSH target is validated
		output, err := cmd.Output()
		if err != nil {
			return paths, fmt.Errorf("failed to get Windows LOCALAPPDATA: %v", err)
		}
		localAppData := strings.TrimSpace(string(output))
		logger.Info(fmt.Sprintf("Windows LOCALAPPDATA: %s", localAppData))

		paths.ToolboxPath = joinRemotePathForOS(osType, localAppData, "JetBrains", "Toolbox-CLI-dist")
		paths.IDEPath = joinRemotePathForOS(osType, localAppData, "Programs")
		paths.ClientPath = joinRemotePathForOS(osType, localAppData, "JetBrains", "Toolbox", "internal-tools")
	default:
		return paths, fmt.Errorf("unsupported OS: %s", osType)
	}

	// Apply custom paths if specified
	if opts.TbcliInstallPath != "" {
		expandedPath, err := expandEnvPathRemote(sshTarget, opts.TbcliInstallPath, osType)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to expand TBCLI install path remotely, using original: %v", err))
			paths.ToolboxPath = opts.TbcliInstallPath
		} else {
			paths.ToolboxPath = expandedPath
		}
	}
	if opts.IdeInstallPath != "" {
		expandedPath, err := expandEnvPathRemote(sshTarget, opts.IdeInstallPath, osType)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to expand IDE install path remotely, using original: %v", err))
			paths.IDEPath = opts.IdeInstallPath
		} else {
			paths.IDEPath = expandedPath
		}
	}

	logger.Info(fmt.Sprintf("Remote paths - Toolbox: %s, IDE: %s, Client: %s",
		paths.ToolboxPath, paths.IDEPath, paths.ClientPath))
	return paths, nil
}

// createEnvFiles creates environment files for TBCLI and JBR paths
// toolboxPath parameter is now the full path to TBCLI directory (e.g., /path/tbcli-2.8.1.52155)
// jbrPath parameter is the full path to JBR directory (e.g., /path/jbr-21.0.3-osx-aarch64-b509.11-1)
func createEnvFiles(opts InstallOptions, osType, arch string, toolboxPath, jbrPath string, isRemote bool, sshTarget string) error {
	if opts.TbcliInstallPath == "" && opts.JbrInstallPath == "" {
		// No custom paths specified, skip env file creation
		return nil
	}

	var envDir, envFile, content string
	var tbcliPath, jbrPathForEnv string

	// Use the full paths that were passed to us (toolboxPath is now full TBCLI path, jbrPath is full JBR path)
	tbcliPath = normalizePathForOS(toolboxPath, osType)
	jbrPathForEnv = normalizePathForOS(jbrPath, osType)

	// Determine environment file location and content based on OS
	switch osType {
	case PlatformWindows:
		if isRemote {
			envDir = "%LOCALAPPDATA%"
			envFile = joinRemotePathForOS(osType, envDir, "JetBrains", "ToolboxSshDeploy", "env.ps1")
		} else {
			localAppData := os.Getenv("LOCALAPPDATA")
			if localAppData == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %v", err)
				}
				localAppData = joinRemotePathForOS(osType, homeDir, "AppData", "Local")
			}
			envDir = joinRemotePathForOS(osType, localAppData, "JetBrains", "ToolboxSshDeploy")
			envFile = joinRemotePathForOS(osType, envDir, "env.ps1")
		}
		content = fmt.Sprintf("$cli_dir_path=\"%s\"\n$java_path=\"%s\"\n", tbcliPath, jbrPathForEnv)

	case PlatformLinux:
		if isRemote {
			homeDir := "~"
			envFile = joinRemotePathForOS(osType, homeDir, ".config", "JetBrains", "ToolboxSshDeploy", "env.sh")
			envDir = joinRemotePathForOS(osType, homeDir, ".config", "JetBrains", "ToolboxSshDeploy")
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %v", err)
			}
			envDir = joinRemotePathForOS(osType, homeDir, ".config", "JetBrains", "ToolboxSshDeploy")
			envFile = joinRemotePathForOS(osType, envDir, "env.sh")
		}
		content = fmt.Sprintf("cli_dir_path=%s\njava_path=%s\n", tbcliPath, jbrPathForEnv)

	case PlatformMac:
		if isRemote {
			homeDir := "~"
			envFile = joinRemotePathForOS(osType, homeDir, "Library", "Application Support", "JetBrains", "ToolboxSshDeploy", "env.sh")
			envDir = joinRemotePathForOS(osType, homeDir, "Library", "Application Support", "JetBrains", "ToolboxSshDeploy")
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %v", err)
			}
			envDir = joinRemotePathForOS(osType, homeDir, "Library", "Application Support", "JetBrains", "ToolboxSshDeploy")
			envFile = joinRemotePathForOS(osType, envDir, "env.sh")
		}
		content = fmt.Sprintf("cli_dir_path=%s\njava_path=%s\n", tbcliPath, jbrPathForEnv)

	default:
		return fmt.Errorf("unsupported OS for environment file creation: %s", osType)
	}

	if isRemote {
		return createEnvFileRemote(sshTarget, envDir, envFile, content, osType)
	}
	return createEnvFileLocal(envDir, envFile, content)
}

// createEnvFileLocal creates environment file locally
func createEnvFileLocal(envDir, envFile, content string) error {
	err := os.MkdirAll(envDir, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create environment directory: %v", err)
	}

	err = os.WriteFile(envFile, []byte(content), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write environment file: %v", err)
	}

	logger.Info(fmt.Sprintf("Created environment file: %s", envFile))
	return nil
}

// createEnvFileRemote creates environment file on remote host via SSH
func createEnvFileRemote(sshTarget, envDir, envFile, content, osType string) error {
	var mkdirCmd, writeCmd *exec.Cmd

	switch osType {
	case PlatformWindows:
		mkdirCmd = exec.Command("ssh", sshTarget, fmt.Sprintf(`powershell.exe -Command "New-Item -ItemType Directory -Force -Path '%s'"`, envDir)) // #nosec G204 -- SSH target is validated
		writeCmd = exec.Command("ssh", sshTarget, fmt.Sprintf(`powershell.exe -Command "Set-Content -Path '%s' -Value '%s'"`, envFile, content))   // #nosec G204 -- SSH target is validated

	case PlatformLinux, PlatformMac:
		mkdirCmd = exec.Command("ssh", sshTarget, "mkdir", "-p", envDir)                                                                     // #nosec G204 -- SSH target is validated
		writeCmd = exec.Command("ssh", sshTarget, "sh", "-c", fmt.Sprintf("cat > %s << 'EOF'\n%sEOF", escapeUnixArgument(envFile), content)) // #nosec G204 -- SSH target is validated
	}

	err := mkdirCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create remote environment directory: %v", err)
	}

	err = writeCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to write remote environment file: %v", err)
	}

	logger.Info(fmt.Sprintf("Created remote environment file: %s", envFile))
	return nil
}

// Download products JSON for client builds
func downloadProductsJSON(config Config) bool {
	productsJSONFile := "./downloads/products/jetbrains_products.json"

	if _, err := os.Stat(productsJSONFile); err == nil {
		logger.Info("Using cached JetBrains products metadata")
		return true
	}

	logger.Info("Downloading JetBrains products metadata (this may take a moment)...")
	err := os.MkdirAll(filepath.Dir(productsJSONFile), 0o750)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to create products directory: %v", err))
		return false
	}

	resp, err := http.Get(config.ProductsJSONURL) // #nosec G107 -- URL from configuration
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to download products metadata: %v", err))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warning(fmt.Sprintf("Failed to download products metadata: HTTP %d from URL: %s", resp.StatusCode, config.ProductsJSONURL))
		return false
	}

	out, err := os.Create(productsJSONFile) // #nosec G304 -- File path is controlled by application
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to create products file: %v", err))
		return false
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to save products metadata: %v", err))
		return false
	}

	logger.Success("Downloaded products metadata")
	return true
}

// Get IDE build number from products JSON
func getIDEBuildNumber(ideFilename string) (string, error) {
	productsJSONFile := "./downloads/products/jetbrains_products.json"

	if _, err := os.Stat(productsJSONFile); os.IsNotExist(err) {
		return "", fmt.Errorf("products JSON file not found")
	}

	logger.Info(fmt.Sprintf("Searching for build number using filename: %s", ideFilename))

	data, err := os.ReadFile(productsJSONFile)
	if err != nil {
		return "", fmt.Errorf("failed to read products JSON: %v", err)
	}

	var products []map[string]interface{}
	err = json.Unmarshal(data, &products)
	if err != nil {
		return "", fmt.Errorf("failed to parse products JSON: %v", err)
	}

	for _, product := range products {
		if releases, ok := product["releases"].([]interface{}); ok {
			for _, release := range releases {
				if releaseMap, ok := release.(map[string]interface{}); ok {
					if downloads, ok := releaseMap["downloads"].(map[string]interface{}); ok {
						for _, download := range downloads {
							if downloadMap, ok := download.(map[string]interface{}); ok {
								if link, ok := downloadMap["link"].(string); ok && strings.Contains(link, ideFilename) {
									if build, ok := releaseMap["build"].(string); ok {
										logger.Info(fmt.Sprintf("Found build number: %s", build))
										return build, nil
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return "", fmt.Errorf("build number not found for %s in products metadata", ideFilename)
}

// parseSSHConnectionString parses SSH connection string and extracts hostname, username, and port
// Supports format: [username@]hostname [-p port] [additional flags]
// Example: "username@hostname.local -p 22 -i ~/.ssh/id_rsa"
func parseSSHConnectionString(sshTarget string) SSHConnectionInfo {
	info := SSHConnectionInfo{
		// Don't set default port - only add if explicitly specified
	}

	// Split by spaces to handle arguments
	parts := strings.Fields(sshTarget)
	if len(parts) == 0 {
		return info
	}

	// First part contains user@host or just host
	hostPart := parts[0]
	if strings.Contains(hostPart, "@") {
		userHost := strings.Split(hostPart, "@")
		if len(userHost) == 2 {
			info.Username = userHost[0]
			info.Hostname = userHost[1]
		}
	} else {
		info.Hostname = hostPart
	}

	// Parse remaining arguments for port
	for i := 1; i < len(parts); i++ {
		if parts[i] == "-p" && i+1 < len(parts) {
			info.Port = parts[i+1]
			break
		}
	}

	return info
}

// getIDEProductCode maps IDE type to JetBrains product code for Toolbox URL
func getIDEProductCode(ideType string) string {
	switch ideType {
	case IdeIdea:
		return "IU" // IntelliJ IDEA Ultimate
	case IdePyCharm:
		return "PY" // PyCharm Professional
	case IdeGoland:
		return "GO" // GoLand
	case IdeClion:
		return "CL" // CLion
	case IdeRider:
		return "RD" // Rider
	case IdeRuby, IdeRubymine:
		return "RM" // RubyMine
	case IdeRustRover:
		return "RR" // RustRover
	case IdeWebStorm:
		return "WS" // WebStorm
	case IdePhpStorm:
		return "PS" // PhpStorm
	default:
		// Default fallback
		return "IU"
	}
}

// openURLInBrowser opens the provided URL in the default browser/application
func openURLInBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case PlatformDarwin: // macOS
		cmd = exec.Command("open", url)
	case PlatformWindows:
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case PlatformLinux:
		// Try common browsers
		browsers := []string{"xdg-open", "firefox", "chromium-browser", "google-chrome"}
		var lastErr error

		for _, browser := range browsers {
			cmd = exec.Command(browser, url) // #nosec G204 -- Browser path is from system configuration
			err := cmd.Start()
			if err == nil {
				return nil
			}
			lastErr = err
		}
		return fmt.Errorf("failed to open URL with any browser: %v", lastErr)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// generateToolboxURL generates the JetBrains Toolbox SSH URL
func generateToolboxURL(sshInfo SSHConnectionInfo, buildNumber, ideType string, withIDELaunch bool) string {
	baseURL := "jetbrains://gateway/ssh/environment"

	params := url.Values{}

	if sshInfo.Hostname != "" {
		params.Add("h", sshInfo.Hostname)
	}

	if sshInfo.Username != "" {
		params.Add("u", sshInfo.Username)
	}

	if sshInfo.Port != "" && sshInfo.Port != "22" {
		params.Add("p", sshInfo.Port)
	}

	if withIDELaunch && buildNumber != "" && ideType != "" {
		params.Add("launchIde", "true")
		productCode := getIDEProductCode(ideType)
		params.Add("ideHint", productCode+"-"+buildNumber)
	}

	if len(params) > 0 {
		return baseURL + "?" + params.Encode()
	}

	return baseURL
}

func generateClientURL(config Config, buildNumber, os, arch, ideType string) (string, error) {
	var archSuffix string
	var extension string

	if arch == ArchAarch64 {
		archSuffix = "-aarch64"
	}

	pathPrefix := "idea"

	switch os {
	case PlatformLinux:
		extension = ExtTarGz
	case PlatformMac:
		extension = "sit"
	case PlatformWindows:
		extension = "jbr.win.zip"
	default:
		return "", fmt.Errorf("unsupported OS for client: %s", os)
	}

	return fmt.Sprintf("%s/%s/code-with-me/JetBrainsClient-%s%s.%s",
		config.JetbrainsBaseURL, pathPrefix, buildNumber, archSuffix, extension), nil
}

func checkSSHExistingInstallation(sshTarget, installPath string, opts InstallOptions, remoteInfo RemoteSystemInfo) (tbcliExists, jbrExists bool, err error) {
	if forceReinstall {
		logger.Info("Force reinstallation enabled, skipping TBCLI and JBR existence checks")
		return false, false, nil
	}

	logger.Info("Checking for existing installation...")

	// Expand TBCLI installation path (may contain ~ or environment variables)
	tbcliInstallPath, expandErr := expandEnvPathRemote(sshTarget, installPath, remoteInfo.OS)
	if expandErr != nil {
		logger.Warning(fmt.Sprintf("Failed to expand TBCLI install path remotely for check: %v", expandErr))
		tbcliInstallPath = installPath
	}

	var jbrOS string
	switch remoteInfo.OS {
	case PlatformLinux:
		jbrOS = PlatformLinux
	case PlatformMac:
		jbrOS = PlatformOSX
	case PlatformWindows:
		jbrOS = PlatformWindows
	}

	// Determine JBR installation path (custom or default)
	var jbrInstallPath string
	if opts.JbrInstallPath != "" {
		expandedPath, expandErr := expandEnvPathRemote(sshTarget, opts.JbrInstallPath, remoteInfo.OS)
		if expandErr != nil {
			logger.Warning(fmt.Sprintf("Failed to expand JBR install path remotely for check: %v", expandErr))
			jbrInstallPath = opts.JbrInstallPath
		} else {
			jbrInstallPath = expandedPath
		}
	} else {
		jbrInstallPath = tbcliInstallPath
	}

	var checkCmd string
	if remoteInfo.OS == PlatformWindows {
		targetJbrDir := fmt.Sprintf("jbr-%s-%s-%s-%s-1", opts.JbrVersion, jbrOS, remoteInfo.Arch, opts.JbrBuild)

		checkCmd = fmt.Sprintf(`powershell.exe -Command "$result = @(); if (Test-Path '%s') { $result += 'dir_exists' }; if (Test-Path '%s\\tbcli-%s\\.extracted') { $result += 'tbcli_exists' }; if (Test-Path '%s\\%s\\.extracted') { $result += 'jbr_exists' }; $result -join ';'"`,
			tbcliInstallPath, tbcliInstallPath, opts.TbcliVersion, jbrInstallPath, targetJbrDir)
	} else {
		checkCmd = fmt.Sprintf(`
			if [ -d '%s' ]; then
				echo 'dir_exists'
				if [ -f '%s/tbcli-%s/.extracted' ]; then
					echo 'tbcli_exists'
				fi
			fi
			target_dir="jbr-%s-%s-%s-%s-1"
			if [ -f '%s/'$target_dir'/.extracted' ]; then
				echo 'jbr_exists'
			fi
		`, tbcliInstallPath, tbcliInstallPath, opts.TbcliVersion, opts.JbrVersion, jbrOS, remoteInfo.Arch, opts.JbrBuild, jbrInstallPath)
	}

	cmd := exec.Command("ssh", sshTarget, checkCmd) // #nosec G204 -- SSH target is validated
	output, err := cmd.Output()
	if err != nil {
		return false, false, err
	}

	result := strings.TrimSpace(string(output))
	// Log the full output for debugging
	if logger.infoEnabled {
		logger.Info(fmt.Sprintf("Check installation output: %s", result))
	}

	if remoteInfo.OS == PlatformWindows {
		results := strings.Split(result, ";")
		tbcliExists = false
		jbrExists = false
		for _, res := range results {
			if strings.TrimSpace(res) == "tbcli_exists" {
				tbcliExists = true
			}
			if strings.TrimSpace(res) == "jbr_exists" {
				jbrExists = true
			}
		}
	} else {
		tbcliExists = strings.Contains(result, "tbcli_exists")
		jbrExists = strings.Contains(result, "jbr_exists")
	}

	return tbcliExists, jbrExists, nil
}

// Check existing IDE installation via SSH
func checkSSHExistingIDE(sshTarget, idePath, ideType, ideVersion string, remoteInfo RemoteSystemInfo) (bool, error) {
	if forceReinstall {
		logger.Info("Force reinstallation enabled, skipping IDE existence check")
		return false, nil
	}

	logger.Info("Checking for existing IDE installation...")

	var productInfoPath string

	if remoteInfo.OS == PlatformMac {
		appPath, err := findIDEAppPath(idePath, ideType, ideVersion, true, sshTarget)
		if err != nil {
			// Check if it's an expected "not found" error vs real error
			if strings.Contains(err.Error(), "no matching .app found") {
				return false, nil // IDE not found - this is expected
			}
			return false, err // Real error - return it
		}
		productInfoPath = joinRemotePathForOS(remoteInfo.OS, appPath, "Contents", "Resources", "product-info.json")
	} else if remoteInfo.OS == PlatformWindows {
		ideDir := fmt.Sprintf("%s-%s", ideType, ideVersion)
		productInfoPath = joinRemotePathForOS(remoteInfo.OS, idePath, ideDir, "product-info.json")
	} else {
		ideDir := fmt.Sprintf("%s-%s", ideType, ideVersion)
		productInfoPath = joinRemotePathForOS(remoteInfo.OS, idePath, ideDir, "product-info.json")
	}

	var checkCmd string
	if remoteInfo.OS == PlatformWindows {
		checkCmd = fmt.Sprintf(`powershell.exe -Command "if (Test-Path '%s') { Write-Output 'exists' } else { Write-Output 'not_found' }"`, productInfoPath)
	} else {
		checkCmd = fmt.Sprintf("[ -f %s ] && echo 'exists' || echo 'not_found'", escapeUnixArgument(productInfoPath))
	}

	cmd := exec.Command("ssh", sshTarget, checkCmd) // #nosec G204 -- SSH target is validated
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	result := strings.TrimSpace(string(output))
	if result != "exists" {
		return false, nil
	}

	var readCmd string
	if remoteInfo.OS == PlatformWindows {
		readCmd = fmt.Sprintf(`powershell.exe -Command "Get-Content '%s' -Raw"`, productInfoPath)
	} else {
		readCmd = fmt.Sprintf("cat '%s'", productInfoPath)
	}
	cmd = exec.Command("ssh", sshTarget, readCmd) // #nosec G204 -- SSH target is validated
	output, err = cmd.Output()
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to read product-info.json: %v", err))
		return false, nil
	}

	var productInfo ProductInfo
	err = json.Unmarshal(output, &productInfo)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to parse product-info.json: %v", err))
		return false, nil
	}

	if productInfo.Version == ideVersion {
		logger.Info(fmt.Sprintf("Found existing installation: %s %s (build %s)",
			productInfo.Name, productInfo.Version, productInfo.BuildNumber))
		return true, nil
	}

	logger.Info(fmt.Sprintf("Found different version: %s %s (requested: %s)",
		productInfo.Name, productInfo.Version, ideVersion))
	return false, nil
}

// Verify IDE installation success by checking product-info.json
func verifyIDEInstallation(sshTarget, idePath, ideType, ideVersion string, remoteInfo RemoteSystemInfo) (bool, error) {
	var productInfoPath string

	if remoteInfo.OS == PlatformMac {
		appPath, err := findIDEAppPath(idePath, ideType, ideVersion, true, sshTarget)
		if err != nil {
			// Check if it's an expected "not found" error vs real error
			if strings.Contains(err.Error(), "no matching .app found") {
				return false, nil // IDE not found - this is expected for verification
			}
			return false, err // Real error - return it
		}
		productInfoPath = joinRemotePathForOS(remoteInfo.OS, appPath, "Contents", "Resources", "product-info.json")
	} else if remoteInfo.OS == PlatformWindows {
		ideDir := fmt.Sprintf("%s-%s", ideType, ideVersion)
		productInfoPath = joinRemotePathForOS(remoteInfo.OS, idePath, ideDir, "product-info.json")
	} else {
		ideDir := fmt.Sprintf("%s-%s", ideType, ideVersion)
		productInfoPath = joinRemotePathForOS(remoteInfo.OS, idePath, ideDir, "product-info.json")
	}

	var checkCmd string
	if remoteInfo.OS == PlatformWindows {
		checkCmd = fmt.Sprintf(`powershell.exe -Command "if (Test-Path '%s') { Write-Output 'exists' } else { Write-Output 'not_found' }"`, productInfoPath)
	} else {
		checkCmd = fmt.Sprintf("[ -f %s ] && echo 'exists' || echo 'not_found'", escapeUnixArgument(productInfoPath))
	}

	cmd := exec.Command("ssh", sshTarget, checkCmd) // #nosec G204 -- SSH target is validated
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	result := strings.TrimSpace(string(output))
	if result != "exists" {
		return false, nil
	}

	var readCmd string
	if remoteInfo.OS == PlatformWindows {
		readCmd = fmt.Sprintf(`powershell.exe -Command "Get-Content '%s' -Raw"`, productInfoPath)
	} else {
		readCmd = fmt.Sprintf("cat '%s'", productInfoPath)
	}

	cmd = exec.Command("ssh", sshTarget, readCmd) // #nosec G204 -- SSH target is validated
	output, err = cmd.Output()
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to read product-info.json: %v", err))
		return false, nil
	}

	var productInfo ProductInfo
	err = json.Unmarshal(output, &productInfo)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to parse product-info.json: %v", err))
		return false, nil
	}

	if productInfo.Version == ideVersion {
		logger.Info(fmt.Sprintf("Verified installation: %s %s (build %s)",
			productInfo.Name, productInfo.Version, productInfo.BuildNumber))
		return true, nil
	}

	logger.Warning(fmt.Sprintf("Version mismatch: expected %s, found %s", ideVersion, productInfo.Version))
	return false, nil
}

// escapeUnixArgument properly escapes arguments for Unix shell commands
func escapeUnixArgument(arg string) string {
	// If the argument contains spaces, quotes, or special shell characters, quote it
	if strings.ContainsAny(arg, " \t\n\"'\\$`;&|<>(){}[]!?*") {
		// Escape single quotes by ending the quoted string, adding an escaped quote, and starting a new quoted string
		escaped := strings.ReplaceAll(arg, "'", "'\"'\"'")
		return "'" + escaped + "'"
	}
	return arg
}

// parseArgsPreservingQuotes parses a command line string into arguments while preserving quoted strings
func parseArgsPreservingQuotes(cmdLine string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(cmdLine); i++ {
		char := cmdLine[i]

		switch char {
		case '"', '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteByte(char)
			}
		case ' ', '\t':
			if inQuotes {
				current.WriteByte(char)
			} else if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

// joinRemotePathForOS joins path elements with OS-appropriate separators for specific remote OS
func joinRemotePathForOS(osType string, elem ...string) string {
	if len(elem) == 0 {
		return ""
	}

	var separator string
	if osType == PlatformWindows {
		separator = "\\"
	} else {
		separator = "/"
	}

	normalizedElems := make([]string, len(elem))
	for i, el := range elem {
		if osType == PlatformWindows {
			normalizedElems[i] = strings.ReplaceAll(el, "/", "\\")
		} else {
			normalizedElems[i] = strings.ReplaceAll(el, "\\", "/")
		}
	}

	result := normalizedElems[0]
	for i := 1; i < len(normalizedElems); i++ {
		if result != "" && !strings.HasSuffix(result, separator) && !strings.HasPrefix(normalizedElems[i], separator) {
			result += separator
		} else if strings.HasSuffix(result, separator) && strings.HasPrefix(normalizedElems[i], separator) {
			normalizedElems[i] = strings.TrimPrefix(normalizedElems[i], separator)
		}
		result += normalizedElems[i]
	}
	return result
}

// Get remote temporary directory
func getRemoteTempDir(sshTarget string, remoteOS string) (string, error) {
	switch remoteOS {
	case PlatformLinux, PlatformMac:
		return "/tmp", nil
	case PlatformWindows:
		cmd := exec.Command("ssh", sshTarget, `powershell.exe -Command "$env:TEMP"`) // #nosec G204 -- SSH target is validated
		output, err := cmd.Output()
		if err != nil {
			// Fallback to default Windows temp path using PowerShell
			cmd = exec.Command("ssh", sshTarget, `powershell.exe -Command "$env:USERPROFILE"`) // #nosec G204 -- SSH target is validated
			homeOutput, homeErr := cmd.Output()
			if homeErr != nil {
				return "", fmt.Errorf("failed to get temp directory: %v", err)
			}
			homeDir := strings.TrimSpace(string(homeOutput))
			return joinRemotePathForOS(PlatformWindows, homeDir, "AppData", "Local", "Temp"), nil
		}
		return strings.TrimSpace(string(output)), nil
	default:
		return "", fmt.Errorf("unknown OS: %s", remoteOS)
	}
}

// Check directory permissions on remote host
func checkDirectoryPermissions(sshTarget, dirPath, operation, osType string) error {
	logger.Info(fmt.Sprintf("Checking %s permissions for directory: %s", operation, dirPath))

	if osType == PlatformWindows {
		switch operation {
		case "write":
			testFile := filepath.Join(dirPath, ".perm_test_file")
			checkCmd := fmt.Sprintf(`powershell.exe -Command "try { New-Item -ItemType File -Path '%s' -Force | Out-Null; Remove-Item '%s' -Force | Out-Null; Write-Host 'ok' } catch { Write-Host 'fail' }"`, testFile, testFile)

			cmd := exec.Command("ssh", sshTarget, checkCmd) // #nosec G204 -- SSH target is validated
			output, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("failed to check Windows permissions: %v", err)
			}

			result := strings.TrimSpace(string(output))
			if result != "ok" {
				return fmt.Errorf("insufficient %s permissions for Windows directory: %s", operation, dirPath)
			}
		case "read":
			checkCmd := fmt.Sprintf(`powershell.exe -Command "try { Get-ChildItem '%s' | Out-Null; Write-Host 'ok' } catch { Write-Host 'fail' }"`, dirPath)

			cmd := exec.Command("ssh", sshTarget, checkCmd) // #nosec G204 -- SSH target is validated
			output, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("failed to check Windows read permissions: %v", err)
			}

			result := strings.TrimSpace(string(output))
			if result != "ok" {
				return fmt.Errorf("insufficient %s permissions for Windows directory: %s", operation, dirPath)
			}
		default:
			return fmt.Errorf("unsupported Windows permission check operation: %s", operation)
		}
	} else {
		var checkCmd string
		switch operation {
		case "write":
			checkCmd = fmt.Sprintf("[ -w '%s' ] && echo 'ok' || echo 'fail'", dirPath)
		case "read":
			checkCmd = fmt.Sprintf("[ -r '%s' ] && echo 'ok' || echo 'fail'", dirPath)
		case "execute":
			checkCmd = fmt.Sprintf("[ -x '%s' ] && echo 'ok' || echo 'fail'", dirPath)
		default:
			return fmt.Errorf("unknown Unix operation: %s", operation)
		}

		cmd := exec.Command("ssh", sshTarget, checkCmd) // #nosec G204 -- SSH target is validated
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to check Unix permissions: %v", err)
		}

		result := strings.TrimSpace(string(output))
		if result != "ok" {
			return fmt.Errorf("insufficient %s permissions for Unix directory: %s", operation, dirPath)
		}
	}

	logger.Info("Directory permissions check passed")
	return nil
}

// Install IDE on remote host via SSH
func installIDEOnRemote(config Config, opts InstallOptions, sshTarget string, remotePaths RemotePaths, remoteInfo RemoteSystemInfo) error {
	logger.Info(fmt.Sprintf("Installing IDE (%s) on remote host...", opts.IdeType))

	idePath := remotePaths.IDEPath
	if opts.IdeInstallPath != "" {
		expandedPath, err := expandEnvPathRemote(sshTarget, opts.IdeInstallPath, remoteInfo.OS)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to expand IDE install path remotely in installIDEOnRemote on host %s, using original: %v", sshTarget, err))
			idePath = opts.IdeInstallPath
		} else {
			idePath = expandedPath
		}
	}

	exists, err := checkSSHExistingIDE(sshTarget, idePath, opts.IdeType, opts.IdeVersion, remoteInfo)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to check existing IDE installation: %v", err))
	} else if exists {
		logger.Info(fmt.Sprintf("IDE %s %s is already installed, skipping installation", opts.IdeType, opts.IdeVersion))
		return nil
	}

	tempDir, err := getRemoteTempDir(sshTarget, remoteInfo.OS)
	if err != nil {
		return fmt.Errorf("failed to get temporary directory: %v", err)
	}

	logger.Info(fmt.Sprintf("Creating IDE installation directory: %s", idePath))
	var cmd *exec.Cmd
	if remoteInfo.OS == PlatformWindows {
		cmd = exec.Command("ssh", sshTarget, fmt.Sprintf(`powershell.exe -Command "New-Item -ItemType Directory -Force -Path '%s'"`, idePath)) // #nosec G204 -- SSH target is validated
	} else {
		cmd = exec.Command("ssh", sshTarget, "mkdir", "-p", idePath) // #nosec G204 -- SSH target is validated
	}
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create IDE installation directory: %v", err)
	}

	err = checkDirectoryPermissions(sshTarget, idePath, "write", remoteInfo.OS)
	if err != nil {
		return fmt.Errorf("no write permissions for IDE installation directory: %v", err)
	}

	ideURL, err := generateJetBrainsURL(config, opts.IdeType, opts.IdeVersion, remoteInfo.OS, remoteInfo.Arch, "")
	if err != nil {
		return fmt.Errorf("failed to generate IDE URL: %v", err)
	}

	ideFilename := filepath.Base(ideURL)
	logger.Info(fmt.Sprintf("IDE URL: %s", ideURL))

	tempFilePath := joinRemotePathForOS(remoteInfo.OS, tempDir, ideFilename)
	directDownloadSuccess := false

	tool, hasTools := checkRemoteDownloadTools(sshTarget, remoteInfo.OS)
	if hasTools {
		logger.Info("Checking IDE URL accessibility from remote host...")
		if checkRemoteURLAccess(sshTarget, ideURL, remoteInfo.OS) {
			logger.Info(fmt.Sprintf("Downloading IDE directly on remote host using %s...", tool))
			err = downloadFileRemotely(sshTarget, ideURL, tempFilePath, remoteInfo.OS)
			if err == nil {
				logger.Success("IDE downloaded directly on remote host")
				directDownloadSuccess = true
			} else {
				logger.Warning(fmt.Sprintf("Direct IDE download failed: %v. Falling back to local download and transfer.", err))
			}
		} else {
			logger.Warning("IDE URL not accessible from remote host. Using local download and transfer.")
		}
	} else {
		logger.Warning("Download tools not available on remote host. Using local download and transfer.")
	}

	// Fallback: download locally and transfer if direct download failed
	if !directDownloadSuccess {
		logger.Info("Downloading IDE locally...")
		err = downloadFile(ideURL, ideFilename)
		if err != nil {
			return fmt.Errorf("failed to download IDE locally: %v", err)
		}

		downloadedFile, findErr := findDownloadedFile(ideFilename)
		if findErr != nil {
			return fmt.Errorf("failed to find downloaded IDE file: %v", findErr)
		}

		logger.Info("Transferring IDE file to remote host...")
		cmd := exec.Command("scp", downloadedFile, fmt.Sprintf("%s:%s", sshTarget, tempFilePath)) // #nosec G204 -- SCP with validated arguments
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to transfer IDE file: %v", err)
		}
	}

	err = installIDERemotely(sshTarget, tempFilePath, idePath, opts.IdeType, opts.IdeVersion, remoteInfo.OS)
	if err != nil {
		return fmt.Errorf("failed to install IDE: %v", err)
	}

	logger.Success(fmt.Sprintf("IDE (%s) installed successfully", opts.IdeType))
	return nil
}

// Find downloaded file in downloads directory structure
func findDownloadedFile(filename string) (string, error) {
	var foundPath string

	err := filepath.Walk("./downloads", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == filename {
			foundPath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("file %s not found in downloads directory", filename)
	}

	return foundPath, nil
}

// Install IDE remotely based on OS and file type
func installIDERemotely(sshTarget, tempFilePath, idePath, ideType, ideVersion, osType string) error {
	switch filepath.Ext(tempFilePath) {
	case ExtGz:
		// tar.gz file (Linux)
		ideDir := joinRemotePathForOS(PlatformLinux, idePath, fmt.Sprintf("%s-%s", ideType, ideVersion))
		extractCmd := fmt.Sprintf("mkdir -p %s && tar -xzf %s --strip-components=1 -C %s && rm %s",
			escapeUnixArgument(ideDir), escapeUnixArgument(tempFilePath), escapeUnixArgument(ideDir), escapeUnixArgument(tempFilePath))

		cmd := exec.Command("ssh", sshTarget, extractCmd) // #nosec G204 -- SSH target is validated
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to extract IDE archive: %v, output: %s", err, string(output))
		}
		logger.Info("IDE installation completed")

	case ".dmg":
		extractCmd := fmt.Sprintf(`
			mount_point="/tmp/ide_mount_$(date +%%s)"
			mkdir -p "$mount_point"
			hdiutil attach "%s" -mountpoint "$mount_point" -nobrowse -quiet
			sleep 2
			
			# Find .app in mounted DMG
			for app in "$mount_point"/*.app; do
				if [ -d "$app" ]; then
					app_name=$(basename "$app")
					base_name="${app_name%%.app}"
					target_path="%s/$app_name"
					
					# Check if app already exists and find unique name
					counter=1
					while [ -e "$target_path" ]; do
						counter=$((counter + 1))
						unique_name="$base_name $counter.app"
						target_path="%s/$unique_name"
					done
					
					# Copy with unique name
					cp -R "$app" "$target_path"
					
					# Log if using different name
					if [ "$app_name" != "$(basename "$target_path")" ]; then
						echo "App $app_name already exists, installed as $(basename "$target_path")"
					else
						echo "Installed $app_name"
					fi
					
					break
				fi
			done
			
			hdiutil detach "$mount_point" -quiet
			rm "%s"
		`, tempFilePath, idePath, idePath, tempFilePath)

		cmd := exec.Command("ssh", sshTarget, extractCmd) // #nosec G204 -- SSH target is validated
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install DMG: %v, output: %s", err, string(output))
		}
		logger.Info("IDE installation completed")
		logger.Info(fmt.Sprintf("Installation output: %s", string(output)))

	case ".exe":
		ideDir := joinRemotePathForOS("windows", idePath, fmt.Sprintf("%s-%s", ideType, ideVersion))

		createDirCmd := fmt.Sprintf(`powershell.exe -Command "New-Item -ItemType Directory -Force -Path '%s'"`, ideDir)
		cmd := exec.Command("ssh", sshTarget, createDirCmd) // #nosec G204 -- SSH target is validated
		if err := cmd.Run(); err != nil {
			logger.Warning(fmt.Sprintf("Failed to create IDE directory: %v", err))
		}

		installCmd := fmt.Sprintf(`powershell.exe -Command "$timeout = 300; $pinfo = New-Object System.Diagnostics.ProcessStartInfo; $pinfo.FileName = '%s'; $pinfo.Arguments = '/S /D=%s'; $pinfo.UseShellExecute = $false; $pinfo.RedirectStandardOutput = $true; $pinfo.RedirectStandardError = $true; $p = New-Object System.Diagnostics.Process; $p.StartInfo = $pinfo; $p.Start(); if (-not $p.WaitForExit($timeout * 1000)) { $p.Kill(); throw 'Installation timeout after ' + $timeout + ' seconds' }; if ($p.ExitCode -ne 0) { throw 'Installation failed with exit code ' + $p.ExitCode }; Write-Output 'Installation completed successfully'"`,
			tempFilePath, ideDir)

		logger.Info("Starting Windows IDE installation (this may take a few minutes)...")
		cmd = exec.Command("ssh", sshTarget, installCmd) // #nosec G204 -- SSH target is validated
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run installer: %v, output: %s", err, string(output))
		}
		logger.Info(fmt.Sprintf("IDE installation completed. Output: %s", string(output)))

		cleanupCmd := fmt.Sprintf(`powershell.exe -Command "Remove-Item '%s' -Force -ErrorAction SilentlyContinue"`, tempFilePath)
		cleanupCmdExec := exec.Command("ssh", sshTarget, cleanupCmd) // #nosec G204 -- SSH target is validated
		if err := cleanupCmdExec.Run(); err != nil {
			logger.Info(fmt.Sprintf("Cleanup command failed (this is normal): %v", err))
		}

	default:
		return fmt.Errorf("unsupported IDE file type: %s", filepath.Ext(tempFilePath))
	}

	logger.Info("Verifying IDE installation...")
	remoteInfo := RemoteSystemInfo{OS: osType}
	verified, err := verifyIDEInstallation(sshTarget, idePath, ideType, ideVersion, remoteInfo)
	if err != nil {
		return fmt.Errorf("installation verification failed: %v", err)
	}
	if !verified {
		return fmt.Errorf("IDE installation failed - product-info.json not found or invalid")
	}

	logger.Success("IDE installation verified successfully")
	return nil
}

// Create Toolbox environment configuration
func createToolboxEnvironmentConfig(sshTarget, customIDEPath string, config Config, remoteInfo RemoteSystemInfo) error {
	logger.Info("Creating/updating Toolbox environment configuration for custom IDE path...")

	var toolboxDir string
	var environmentFile string

	switch remoteInfo.OS {
	case PlatformWindows:
		cmd := exec.Command("ssh", sshTarget, `powershell.exe -Command "$env:LOCALAPPDATA"`)
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get LOCALAPPDATA: %v", err)
		}
		localAppData := strings.TrimSpace(string(output))
		toolboxDir = filepath.Join(localAppData, "JetBrains", "Toolbox")
	case PlatformLinux:
		homeDir, err := getRemoteHomeDirectory(sshTarget, PlatformLinux)
		if err != nil {
			return fmt.Errorf("failed to get Linux home directory: %v", err)
		}
		toolboxDir = joinRemotePathForOS(PlatformLinux, homeDir, ".local", "share", "JetBrains", "Toolbox")
	case PlatformMac:
		homeDir, err := getRemoteHomeDirectory(sshTarget, PlatformMac)
		if err != nil {
			return fmt.Errorf("failed to get macOS home directory: %v", err)
		}
		toolboxDir = joinRemotePathForOS(PlatformMac, homeDir, "Library", "Application Support", "JetBrains", "Toolbox")
	default:
		return fmt.Errorf("unsupported OS for environment config: %s", remoteInfo.OS)
	}

	environmentFile = joinRemotePathForOS(remoteInfo.OS, toolboxDir, "environment.json")

	var mkdirCmd *exec.Cmd
	switch remoteInfo.OS {
	case PlatformWindows:
		mkdirCmd = exec.Command("ssh", sshTarget, fmt.Sprintf(`powershell.exe -Command "New-Item -ItemType Directory -Force -Path '%s'"`, toolboxDir)) // #nosec G204 -- SSH target is validated
	case PlatformLinux, PlatformMac:
		mkdirCmd = exec.Command("ssh", sshTarget, "mkdir", "-p", toolboxDir) // #nosec G204 -- SSH target is validated
	}

	err := mkdirCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create Toolbox directory: %v", err)
	}

	existingConfig, exists, err := readExistingEnvironmentConfigRemote(sshTarget, environmentFile, remoteInfo)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to read existing environment.json, creating new: %v", err))
		exists = false
	}

	var finalConfig ToolboxEnvironmentConfig
	if exists {
		logger.Info("Found existing environment.json, merging new location...")

		if locationExists(existingConfig.Tools.Location, customIDEPath) {
			logger.Info(fmt.Sprintf("Location '%s' already exists in environment.json", customIDEPath))
			return nil
		}

		finalConfig = mergeEnvironmentConfig(existingConfig, customIDEPath)
		logger.Info(fmt.Sprintf("Added new location '%s' to environment.json", customIDEPath))
	} else {
		logger.Info("Creating new environment.json...")
		finalConfig = createDefaultEnvironmentConfig(customIDEPath)
	}

	jsonData, err := json.MarshalIndent(finalConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal environment configuration: %v", err)
	}

	logger.Info(fmt.Sprintf("Writing environment.json to: %s", environmentFile))

	var writeCmd *exec.Cmd
	var output []byte
	jsonContent := string(jsonData)

	switch remoteInfo.OS {
	case PlatformWindows:
		tmpFile, tmpErr := os.CreateTemp("", "environment-*.json")
		if tmpErr != nil {
			return fmt.Errorf("failed to create temporary file: %v", tmpErr)
		}
		defer os.Remove(tmpFile.Name())

		windowsContent := strings.ReplaceAll(jsonContent, "\n", "\r\n")
		_, err = tmpFile.WriteString(windowsContent)
		if err != nil {
			_ = tmpFile.Close() // #nosec G104 -- Temporary file close
			return fmt.Errorf("failed to write to temporary file: %v", err)
		}
		_ = tmpFile.Close() // #nosec G104 -- Temporary file close

		scpCmd := exec.Command("scp", tmpFile.Name(), fmt.Sprintf("%s:%s", sshTarget, strings.ReplaceAll(environmentFile, "\\", "/"))) // #nosec G204 -- SCP with validated arguments
		output, err = scpCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to copy environment.json to Windows: %v, output: %s", err, string(output))
		}

		if exists {
			logger.Success("Toolbox environment configuration updated with new location")
		} else {
			logger.Success("Toolbox environment configuration created")
		}
		logger.Info("Configuration details: Environment configuration updated successfully")
		return nil
	case PlatformLinux, PlatformMac:
		unixCommand := fmt.Sprintf("cat > '%s' << 'EOF_JSON'\n%s\nEOF_JSON\necho 'Environment configuration updated successfully'", environmentFile, jsonContent)
		writeCmd = exec.Command("ssh", sshTarget, unixCommand) // #nosec G204 -- SSH target is validated
	}

	output, err = writeCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create/update Toolbox environment configuration: %v, output: %s", err, string(output))
	}

	if exists {
		logger.Success("Toolbox environment configuration updated with new location")
	} else {
		logger.Success("Toolbox environment configuration created")
	}
	logger.Info(fmt.Sprintf("Configuration details: %s", string(output)))
	return nil
}

// Create Toolbox environment configuration for local installation
func createLocalToolboxEnvironmentConfig(homeDir, customIDEPath string, config Config, osType string) error {
	logger.Info("Creating/updating Toolbox environment configuration for custom IDE path...")

	var toolboxDir string

	switch osType {
	case PlatformLinux:
		toolboxDir = filepath.Join(homeDir, ".local", "share", "JetBrains", "Toolbox")
	case PlatformMac:
		toolboxDir = filepath.Join(homeDir, "Library", "Application Support", "JetBrains", "Toolbox")
	case PlatformWindows:
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = filepath.Join(homeDir, "AppData", "Local")
		}
		toolboxDir = filepath.Join(localAppData, "JetBrains", "Toolbox")
	default:
		return fmt.Errorf("unsupported OS for environment config: %s", osType)
	}

	environmentFile := filepath.Join(toolboxDir, "environment.json")

	err := os.MkdirAll(toolboxDir, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create Toolbox directory: %v", err)
	}

	existingConfig, exists, err := readExistingEnvironmentConfig(environmentFile)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to read existing environment.json, creating new: %v", err))
		exists = false
	}

	var finalConfig ToolboxEnvironmentConfig
	if exists {
		logger.Info("Found existing environment.json, merging new location...")

		if locationExists(existingConfig.Tools.Location, customIDEPath) {
			logger.Info(fmt.Sprintf("Location '%s' already exists in environment.json", customIDEPath))
			return nil
		}

		finalConfig = mergeEnvironmentConfig(existingConfig, customIDEPath)
		logger.Info(fmt.Sprintf("Added new location '%s' to environment.json", customIDEPath))
	} else {
		logger.Info("Creating new environment.json...")
		finalConfig = createDefaultEnvironmentConfig(customIDEPath)
	}

	jsonData, err := json.MarshalIndent(finalConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal environment configuration: %v", err)
	}

	logger.Info(fmt.Sprintf("Writing environment.json to: %s", environmentFile))
	err = os.WriteFile(environmentFile, jsonData, 0o600)
	if err != nil {
		return fmt.Errorf("failed to create/update Toolbox environment configuration: %v", err)
	}

	if exists {
		logger.Success("Toolbox environment configuration updated with new location")
	} else {
		logger.Success("Toolbox environment configuration created")
	}
	return nil
}

// Check if download tools (curl/wget) are available on remote host
func checkRemoteDownloadTools(sshTarget, osType string) (string, bool) {
	var tools []string
	if osType == PlatformWindows {
		tools = []string{"curl", "powershell"}
	} else {
		tools = []string{"curl", "wget"}
	}

	for _, tool := range tools {
		var cmd *exec.Cmd
		if osType == PlatformWindows {
			if tool == "curl" {
				cmd = exec.Command("ssh", sshTarget, "where", "curl") // #nosec G204 -- SSH target is validated
			} else {
				cmd = exec.Command("ssh", sshTarget, "where", "powershell") // #nosec G204 -- SSH target is validated
			}
		} else {
			cmd = exec.Command("ssh", sshTarget, "which", tool) // #nosec G204 -- SSH target is validated
		}

		err := cmd.Run()
		if err == nil {
			return tool, true
		}
	}
	return "", false
}

// Check if URL is accessible from remote host
func checkRemoteURLAccess(sshTarget, url, osType string) bool {
	var cmd *exec.Cmd

	if osType == PlatformWindows {
		psCmd := fmt.Sprintf(`powershell.exe -Command "try { $response = Invoke-WebRequest -Uri '%s' -Method Head -TimeoutSec 10; exit 0 } catch { exit 1 }"`, url)
		cmd = exec.Command("ssh", sshTarget, psCmd) // #nosec G204 -- SSH target is validated
	} else {
		curlCmd := fmt.Sprintf("curl -s --head --connect-timeout 10 '%s' >/dev/null 2>&1", url)
		cmd = exec.Command("ssh", sshTarget, curlCmd) // #nosec G204 -- SSH target is validated
		err := cmd.Run()
		if err == nil {
			return true
		}

		wgetCmd := fmt.Sprintf("wget --spider --timeout=10 '%s' >/dev/null 2>&1", url)
		cmd = exec.Command("ssh", sshTarget, wgetCmd) // #nosec G204 -- SSH target is validated
	}

	err := cmd.Run()
	return err == nil
}

// Download file directly on remote host
func downloadFileRemotely(sshTarget, url, remotePath, osType string) error {
	var cmd *exec.Cmd

	if osType == PlatformWindows {
		psCmd := fmt.Sprintf(`powershell.exe -Command "Invoke-WebRequest -Uri '%s' -OutFile '%s'"`, url, remotePath)
		cmd = exec.Command("ssh", sshTarget, psCmd) // #nosec G204 -- SSH target is validated
	} else {
		// Try curl first
		curlCmd := fmt.Sprintf("curl -L -o '%s' '%s'", remotePath, url)
		cmd = exec.Command("ssh", sshTarget, curlCmd) // #nosec G204 -- SSH target is validated
		err := cmd.Run()
		if err == nil {
			return nil
		}

		wgetCmd := fmt.Sprintf("wget -O '%s' '%s'", remotePath, url)
		cmd = exec.Command("ssh", sshTarget, wgetCmd) // #nosec G204 -- SSH target is validated
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download file remotely: %v, output: %s", err, string(output))
	}
	return nil
}

// Transfer and install file via SSH with optional direct download
func transferAndInstallSSHWithURL(sshTarget, localFile, downloadURL, remotePath, productType, version, osType string) error {
	// Create remote directory based on OS
	if osType == PlatformWindows {
		cmd := exec.Command("ssh", sshTarget, fmt.Sprintf(`powershell.exe -Command "New-Item -ItemType Directory -Force -Path '%s'"`, remotePath)) // #nosec G204 -- SSH target is validated
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to create Windows remote directory: %v", err)
		}
	} else {
		cmd := exec.Command("ssh", sshTarget, "mkdir", "-p", remotePath) // #nosec G204 -- SSH target is validated
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to create Unix remote directory: %v", err)
		}
	}

	tempDir, err := getRemoteTempDir(sshTarget, osType)
	if err != nil {
		return fmt.Errorf("failed to get remote temp directory: %v", err)
	}

	var remoteTempFile string
	var filename string

	if downloadURL != "" {
		filename = filepath.Base(downloadURL)
		remoteTempFile = joinRemotePathForOS(osType, tempDir, filename)

		tool, hasTools := checkRemoteDownloadTools(sshTarget, osType)
		if hasTools {
			logger.Info(fmt.Sprintf("Checking URL accessibility: %s", downloadURL))
			if checkRemoteURLAccess(sshTarget, downloadURL, osType) {
				logger.Info(fmt.Sprintf("Downloading %s directly on remote host using %s...", filename, tool))
				err = downloadFileRemotely(sshTarget, downloadURL, remoteTempFile, osType)
				if err == nil {
					logger.Success("File downloaded directly on remote host")
				} else {
					logger.Warning(fmt.Sprintf("Direct download failed: %v. Falling back to local download and transfer.", err))
				}
			} else {
				logger.Warning("URL not accessible from remote host. Using local download and transfer.")
			}
		} else {
			logger.Warning("Download tools not available on remote host. Using local download and transfer.")
		}
	}

	// Fallback: transfer local file if direct download didn't work or wasn't attempted
	directDownloadSuccess := false
	if downloadURL != "" && remoteTempFile != "" {
		var checkCmd *exec.Cmd
		if osType == PlatformWindows {
			checkCmd = exec.Command("ssh", sshTarget, fmt.Sprintf(`powershell.exe -Command "Test-Path '%s'"`, remoteTempFile)) // #nosec G204 -- SSH target is validated
		} else {
			checkCmd = exec.Command("ssh", sshTarget, "test", "-f", remoteTempFile) // #nosec G204 -- SSH target is validated
		}
		directDownloadSuccess = checkCmd.Run() == nil
	}

	if !directDownloadSuccess {
		if localFile == "" {
			if downloadURL == "" {
				return fmt.Errorf("no local file or download URL provided")
			}
			filename = filepath.Base(downloadURL)
			logger.Info(fmt.Sprintf("Remote download failed, downloading %s locally...", filename))
			err = downloadFile(downloadURL, filename)
			if err != nil {
				return fmt.Errorf("failed to download file locally as fallback: %v", err)
			}

			localFile, err = findDownloadedFile(filename)
			if err != nil {
				return fmt.Errorf("failed to find downloaded file: %v", err)
			}
		}

		filename = filepath.Base(localFile)
		remoteTempFile = joinRemotePathForOS(osType, tempDir, filename)
		logger.Info(fmt.Sprintf("Transferring %s to remote host...", filename))
		cmd := exec.Command("scp", localFile, fmt.Sprintf("%s:%s", sshTarget, remoteTempFile)) // #nosec G204 -- SCP with validated arguments
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to transfer file: %v", err)
		}
	}

	var targetDir string
	if productType == TbcliName {
		targetDir = joinRemotePathForOS(osType, remotePath, fmt.Sprintf("tbcli-%s", version))
	} else {
		targetDir = joinRemotePathForOS(osType, remotePath, version)
	}

	if osType == PlatformWindows {
		extractCmd := fmt.Sprintf(`powershell.exe -Command "New-Item -ItemType Directory -Force -Path '%s'; tar -xzf '%s' --strip-components=1 -C '%s'; New-Item -ItemType File -Force -Path '%s\.extracted'; Remove-Item '%s'"`,
			targetDir, remoteTempFile, targetDir, targetDir, remoteTempFile)

		cmd := exec.Command("ssh", sshTarget, extractCmd) // #nosec G204 -- SSH target is validated
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to extract on Windows remote host: %v, output: %s", err, string(output))
		}
	} else {
		verifyCmd := fmt.Sprintf("file %s | grep -q 'gzip compressed'", escapeUnixArgument(remoteTempFile))
		cmd := exec.Command("ssh", sshTarget, verifyCmd) // #nosec G204 -- SSH target is validated
		if cmd.Run() != nil {
			debugCmd := fmt.Sprintf("head -c 200 %s", escapeUnixArgument(remoteTempFile))
			debugCmd2 := exec.Command("ssh", sshTarget, debugCmd) // #nosec G204 -- SSH target is validated
			debugOutput, err := debugCmd2.Output()
			if err != nil {
				return fmt.Errorf("downloaded file is not a valid gzip archive. Failed to get file preview: %v", err)
			}

			return fmt.Errorf("downloaded file is not a valid gzip archive. File content preview: %s", string(debugOutput))
		}

		extractCmd := fmt.Sprintf("mkdir -p %s && tar -xzf %s --strip-components=1 -C %s && touch %s/.extracted && rm %s",
			escapeUnixArgument(targetDir), escapeUnixArgument(remoteTempFile), escapeUnixArgument(targetDir), escapeUnixArgument(targetDir), escapeUnixArgument(remoteTempFile))

		cmd = exec.Command("ssh", sshTarget, extractCmd) // #nosec G204 -- SSH target is validated
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to extract on Unix remote host: %v, output: %s", err, string(output))
		}
	}

	logger.Success(fmt.Sprintf("Installed %s on remote host", productType))
	return nil
}

// Download for multiple platforms (like bash version)
func downloadForPlatforms(config Config, opts InstallOptions) error {
	var osTypes []string
	var archTypes []string

	if opts.DownloadOS == ArchAll {
		osTypes = []string{PlatformLinux, PlatformMac, PlatformWindows}
	} else {
		osTypes = []string{opts.DownloadOS}
	}

	if opts.DownloadArch == ArchAll {
		archTypes = []string{ArchX64, ArchAarch64}
	} else {
		archTypes = []string{opts.DownloadArch}
	}

	logger.Info(fmt.Sprintf("Will download for OS: %v", osTypes))
	logger.Info(fmt.Sprintf("Will download for architectures: %v", archTypes))

	for _, osType := range osTypes {
		for _, archType := range archTypes {
			logger.Info(fmt.Sprintf("Processing Toolbox CLI, JBR and Toolbox for %s-%s...", osType, archType))

			tbcliURL, err := generateJetBrainsURL(config, TbcliName, opts.TbcliVersion, osType, archType, "")
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to generate Toolbox CLI URL for %s-%s: %v", osType, archType, err))
			} else {
				logger.Info(fmt.Sprintf("Toolbox CLI URL (%s-%s): %s", osType, archType, tbcliURL))
				err = downloadFile(tbcliURL, filepath.Base(tbcliURL))
				if err != nil {
					logger.Warning(fmt.Sprintf("Failed to download Toolbox CLI for %s-%s: %v", osType, archType, err))
				}
			}

			toolboxURL, err := generateJetBrainsURL(config, "toolbox", opts.ToolboxVersion, osType, archType, "")
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to generate Toolbox URL for %s-%s: %v", osType, archType, err))
			} else {
				logger.Info(fmt.Sprintf("Toolbox URL (%s-%s): %s", osType, archType, toolboxURL))
				err = downloadFile(toolboxURL, filepath.Base(toolboxURL))
				if err != nil {
					logger.Warning(fmt.Sprintf("Failed to download Toolbox for %s-%s: %v", osType, archType, err))
				}
			}

			jbrURL, err := generateJetBrainsURL(config, "jbr", opts.JbrVersion, osType, archType, opts.JbrBuild)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to generate JBR URL for %s-%s: %v", osType, archType, err))
			} else {
				logger.Info(fmt.Sprintf("JBR URL (%s-%s): %s", osType, archType, jbrURL))
				err = downloadFile(jbrURL, filepath.Base(jbrURL))
				if err != nil {
					logger.Warning(fmt.Sprintf("Failed to download JBR for %s-%s: %v", osType, archType, err))
				}
			}

			logger.Success(fmt.Sprintf("Toolbox CLI, JBR and Toolbox download completed for %s-%s", osType, archType))
		}
	}

	if !opts.InstallIde {
		logger.Info("IDE download skipped (--install-ide not specified)")
	} else {
		logger.Info(fmt.Sprintf("Starting IDE download for %s version %s", opts.IdeType, opts.IdeVersion))

		if opts.InstallClient {
			if !downloadProductsJSON(config) {
				logger.Warning("Failed to download products metadata, Client downloads may be skipped")
				opts.InstallClient = false
			}
		}

		for _, osType := range osTypes {
			for _, archType := range archTypes {
				logger.Info(fmt.Sprintf("Processing %s-%s...", osType, archType))

				ideURL, err := generateJetBrainsURL(config, opts.IdeType, opts.IdeVersion, osType, archType, "")
				if err != nil {
					logger.Warning(fmt.Sprintf("Failed to generate IDE URL for %s-%s: %v", osType, archType, err))
					continue
				}

				ideFilename := filepath.Base(ideURL)
				logger.Info(fmt.Sprintf("IDE URL (%s-%s): %s", osType, archType, ideURL))
				logger.Info(fmt.Sprintf("IDE Filename: %s", ideFilename))

				err = downloadFile(ideURL, ideFilename)
				if err != nil {
					logger.Warning(fmt.Sprintf("Failed to download IDE for %s-%s: %v", osType, archType, err))
					continue
				}

				if opts.InstallClient {
					buildNumber, err := getIDEBuildNumber(ideFilename)
					if err != nil || buildNumber == "" {
						logger.Warning(fmt.Sprintf("Failed to get IDE build number for %s, skipping Client download", ideFilename))
						continue
					}

					clientURL, err := generateClientURL(config, buildNumber, osType, archType, opts.IdeType)
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to generate Client URL for %s-%s, skipping Client download", osType, archType))
						continue
					}

					clientFilename := filepath.Base(clientURL)
					logger.Info(fmt.Sprintf("Client URL (%s-%s): %s", osType, archType, clientURL))
					logger.Info(fmt.Sprintf("Client Filename: %s", clientFilename))

					err = downloadFile(clientURL, clientFilename)
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to download JetBrains Client for %s-%s: %v", osType, archType, err))
					}
				}

				logger.Success(fmt.Sprintf("Download completed for %s-%s", osType, archType))
			}
		}
	}

	if opts.InstallClient && !opts.InstallIde && !opts.NoClient {
		logger.Info("Starting Client-only download (without IDE)...")
		logger.Info(fmt.Sprintf("Using default IDE type: %s version: %s for Client build number", opts.IdeType, opts.IdeVersion))

		if !downloadProductsJSON(config) {
			logger.Warning("Failed to download products metadata, Client downloads will be skipped")
		} else {
			for _, osType := range osTypes {
				for _, archType := range archTypes {
					logger.Info(fmt.Sprintf("Processing Client for %s-%s...", osType, archType))

					ideURL, err := generateJetBrainsURL(config, opts.IdeType, opts.IdeVersion, osType, archType, "")
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to generate IDE URL for build number (%s-%s): %v", osType, archType, err))
						continue
					}

					ideFilename := filepath.Base(ideURL)
					buildNumber, err := getIDEBuildNumber(ideFilename)
					if err != nil || buildNumber == "" {
						logger.Warning(fmt.Sprintf("Failed to get IDE build number for %s (%s-%s), skipping Client download", ideFilename, osType, archType))
						continue
					}

					clientURL, err := generateClientURL(config, buildNumber, osType, archType, opts.IdeType)
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to generate Client URL for %s-%s: %v", osType, archType, err))
						continue
					}

					clientFilename := filepath.Base(clientURL)
					logger.Info(fmt.Sprintf("Client URL (%s-%s): %s", osType, archType, clientURL))
					logger.Info(fmt.Sprintf("Client Filename: %s", clientFilename))

					err = downloadFile(clientURL, clientFilename)
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to download JetBrains Client for %s-%s: %v", osType, archType, err))
					} else {
						logger.Success(fmt.Sprintf("Downloaded Client for %s-%s", osType, archType))
					}
				}
			}
		}
	}

	logger.Success("Multi-platform download completed!")
	logger.Info("Files saved to: ./downloads/")
	return nil
}

// Show help
func showHelp() {
	fmt.Println(`Usage: CLIInstaller [OPTIONS] [SSH_TARGET]

Install JetBrains Toolbox CLI on remote host via SSH, locally, or download tools

Arguments:
    SSH_TARGET          SSH target (username@hostname or SSH alias)
                        Required for SSH installation mode

Options:
    -h, --help         Show this help message
    -v, --version      Show version
    --info             Enable detailed INFO logging (default: disabled)
    --config           Path to config file or URL
    --generate-config  Generate config.json template at specified path
    --tbcli-version    Set Toolbox CLI version
    --toolbox-version  Set Toolbox version  
    --jbr-version      Set JBR version
    --jbr-build        Set JBR build
    --install-ide      Install IDE along with Toolbox CLI
    --ide-type         Set IDE type (idea, pycharm, goland, clion, rider, etc.)
    --ide-version      Set IDE version
    --ide-install-path   Set custom IDE installation path
    --tbcli-install-path Set custom Toolbox CLI installation path
    --jbr-install-path   Set custom JBR installation path
    --launch-ide-args    Arguments to pass to IDE remote-dev-server after installation
    --install-client   Install JetBrains Client for IDE
    --no-client        Skip JetBrains Client installation
    --download-tools   Download only mode for all platforms
    --download-os      Download for specific OS (linux, mac, windows, all)
    --download-arch    Download for specific architecture (x64, aarch64, all)
    --local            Install locally instead of via SSH
    --no-progress      Disable progress bar
    --force            Force reinstallation (ignore existing installations)
    --launch-tbcli     Launch TBCLi after installation
    --open-toolbox-url Automatically open JetBrains Toolbox SSH URL in default application

Examples:
    # Install locally
    CLIInstaller --local
    CLIInstaller --local --install-ide --ide-type idea --ide-install-path /opt/jetbrains

    # Install on remote host via SSH
    CLIInstaller user@server.local
    CLIInstaller --install-ide --ide-type pycharm --ide-version 2025.1.3 user@server.local
    CLIInstaller --install-ide --ide-type idea --no-client user@server.local
    CLIInstaller --install-ide --ide-type idea --launch-ide-args status user@server.local
    CLIInstaller --install-ide --ide-type pycharm --open-toolbox-url user@server.local
    
    # Download only (for administrators)
    CLIInstaller --download-tools
    CLIInstaller --download-os linux --download-arch x64
    CLIInstaller --download-os all --download-arch all --install-ide --ide-type idea
    
    # Disable progress bar for scripting
    CLIInstaller --download-tools --no-progress
    CLIInstaller --local --no-progress
    
    # Force reinstallation (ignore existing installations)
    CLIInstaller --local --force
    CLIInstaller --install-ide --ide-type idea --force user@server.local
    CLIInstaller --download-tools --force
    
    # Configuration
    CLIInstaller --generate-config custom-config.json
    CLIInstaller --config https://company.com/jetbrains-config.json --download-tools
    
    # Version synchronization
    CLIInstaller --tbcli-version 2.8.0.48537 user@server  # Toolbox version auto-synced`)
}

// Launch TBCLi with proper environment
func launchTbcli(config Config, opts InstallOptions, osType, arch, toolboxPath string, isSSH bool, sshTarget string) error {
	err := stopExistingTbcli(osType, isSSH, sshTarget)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to stop existing TBCLi processes: %v", err))
	}

	logger.Info("Launching TBCLi...")

	var tbcliDir, jbrPath, tbcliPath string

	// Determine TBCLI base path (custom or default)
	var tbcliBasePath string
	if opts.TbcliInstallPath != "" {
		if isSSH {
			expandedPath, expandErr := expandEnvPathRemote(sshTarget, opts.TbcliInstallPath, osType)
			if expandErr != nil {
				logger.Warning(fmt.Sprintf("Failed to expand TBCLI install path for launch: %v", expandErr))
				tbcliBasePath = opts.TbcliInstallPath
			} else {
				tbcliBasePath = expandedPath
			}
		} else {
			tbcliBasePath = expandEnvPath(opts.TbcliInstallPath)
		}
	} else {
		tbcliBasePath = toolboxPath
	}

	if isSSH {
		tbcliDir = joinRemotePathForOS(osType, tbcliBasePath, fmt.Sprintf("tbcli-%s", opts.TbcliVersion))
	} else {
		tbcliDir = filepath.Join(tbcliBasePath, fmt.Sprintf("tbcli-%s", opts.TbcliVersion))
	}

	var jbrOS string
	switch osType {
	case PlatformLinux:
		jbrOS = PlatformLinux
	case PlatformMac:
		jbrOS = PlatformOSX
	case PlatformWindows:
		jbrOS = PlatformWindows
	}

	// Determine JBR base path (custom or default)
	var jbrBasePath string
	if opts.JbrInstallPath != "" {
		if isSSH {
			expandedPath, expandErr := expandEnvPathRemote(sshTarget, opts.JbrInstallPath, osType)
			if expandErr != nil {
				logger.Warning(fmt.Sprintf("Failed to expand JBR install path for launch: %v", expandErr))
				jbrBasePath = opts.JbrInstallPath
			} else {
				jbrBasePath = expandedPath
			}
		} else {
			jbrBasePath = expandEnvPath(opts.JbrInstallPath)
		}
	} else {
		jbrBasePath = toolboxPath
	}

	jbrDir := fmt.Sprintf("jbr-%s-%s-%s-%s-1", opts.JbrVersion, jbrOS, arch, opts.JbrBuild)

	if isSSH {
		jbrPath = joinRemotePathForOS(osType, jbrBasePath, jbrDir)
		if osType == PlatformMac {
			jbrPath = joinRemotePathForOS(osType, jbrPath, "Contents", "Home")
		}
	} else {
		jbrPath = filepath.Join(jbrBasePath, jbrDir)
		if osType == PlatformMac {
			jbrPath = filepath.Join(jbrPath, "Contents", "Home")
		}
	}

	tbcliExe := TbcliName
	if osType == PlatformWindows {
		tbcliExe = "Test14.bat"
	}

	if isSSH {
		tbcliPath = joinRemotePathForOS(osType, tbcliDir, "bin", tbcliExe)
	} else {
		tbcliPath = filepath.Join(tbcliDir, "bin", tbcliExe)
	}

	args := []string{"--structured-logging", "agent"}

	if isSSH && sshTarget != "" {
		var remoteCommand string
		if osType == PlatformWindows {
			// Windows needs different syntax for environment variables and paths
			// But TBCLi bat expects forward slashes in TB_JAVA_HOME
			windowsTbcliPath := strings.ReplaceAll(tbcliPath, "/", "\\")
			windowsjbrPath := strings.ReplaceAll(jbrPath, "/", "\\")
			// Start TBCLi in background on Windows using PowerShell - simple and reliable approach
			remoteCommand = fmt.Sprintf("powershell -Command \"$env:TB_JAVA_HOME='%s'; & '%s' %s\"",
				windowsjbrPath, // Keep forward slashes for TB_JAVA_HOME
				windowsTbcliPath,
				strings.Join(args, " "))
		} else {
			// Unix-like systems - run in background with nohup
			// Escape paths to handle spaces and special characters
			remoteCommand = fmt.Sprintf("nohup env TB_JAVA_HOME=%s %s %s > /dev/null 2>&1 &",
				escapeUnixArgument(jbrPath),
				escapeUnixArgument(tbcliPath),
				strings.Join(args, " "))
		}

		sshArgs := []string{sshTarget, remoteCommand}

		cmd := exec.Command("ssh", sshArgs...) // #nosec G204 -- SSH with validated arguments
		logger.Info(fmt.Sprintf("Executing SSH command: ssh %s", strings.Join(sshArgs, " ")))

		if osType == PlatformWindows {
			// For Windows, wait briefly to capture initial output then let it run
			handleCommandOutputAsync(cmd, false)

			// Give command time to start
			time.Sleep(1 * time.Second)
			logger.Success("TBCLi launch initiated via SSH - process may take a moment to fully start on Windows")
		} else {
			// For Unix systems, use Start() to not wait for completion
			err := cmd.Start()
			if err != nil {
				return fmt.Errorf("failed to launch TBCLi via SSH: %v", err)
			}

			// Give TBCLi a moment to start
			time.Sleep(1 * time.Second)

			// Verify TBCLi actually started
			if verifyTbcliRunning(osType, isSSH, sshTarget) {
				logger.Success("TBCLi launched successfully via SSH in background")
			} else {
				logger.Warning("TBCLi agent process not found via SSH - it may have exited immediately or an existing Toolbox instance is handling agent duties")
			}
		}
	} else {
		// For local installation - start in background
		var cmd *exec.Cmd

		if osType == PlatformWindows {
			// Windows: set environment and start process directly
			cmd = exec.Command(tbcliPath, args...) // #nosec G204 -- Toolbox CLI path is validated
			cmd.Env = append(os.Environ(), fmt.Sprintf("TB_JAVA_HOME=%s", jbrPath))
		} else {
			// Unix-like systems: use nohup for background execution
			cmd = exec.Command("nohup", tbcliPath, args[0], args[1]) // #nosec G204 -- Background process execution with validated path
			cmd.Env = append(os.Environ(), fmt.Sprintf("TB_JAVA_HOME=%s", jbrPath))
		}

		logger.Info(fmt.Sprintf("Executing local command: TB_JAVA_HOME=%s %s %s (background)",
			jbrPath,
			tbcliPath,
			strings.Join(args, " ")))

		err := cmd.Start() // Use Start() instead of Run() to not wait for completion
		if err != nil {
			return fmt.Errorf("failed to launch TBCLi locally: %v", err)
		}

		// Give TBCLi a moment to start
		time.Sleep(1 * time.Second)

		// Verify TBCLi actually started
		if verifyTbcliRunning(osType, isSSH, sshTarget) {
			logger.Success("TBCLi launched successfully in background")
		} else {
			logger.Warning("TBCLi agent process not found - it may have exited immediately or an existing Toolbox instance is handling agent duties")
		}
	}

	return nil
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// verifyTbcliRunning checks if TBCLi process is actually running
func verifyTbcliRunning(osType string, isSSH bool, sshTarget string) bool {
	hasJavaAgent := checkForJavaAgentProcess(osType, isSSH, sshTarget)

	hasToolboxProcess := checkForToolboxProcess(osType, isSSH, sshTarget)

	return hasJavaAgent || hasToolboxProcess
}

// checkForJavaAgentProcess checks for standalone TBCLi java agent processes
func checkForJavaAgentProcess(osType string, isSSH bool, sshTarget string) bool {
	var findCmd string

	switch osType {
	case PlatformLinux, PlatformMac:
		findCmd = `ps aux | grep "structured-logging agent" | grep -v grep | wc -l`
	case PlatformWindows:
		findCmd = `powershell -Command "([array](tasklist /fo csv /fi 'imagename eq java.exe' | findstr 'structured-logging')).Count"`
	default:
		return false
	}

	return executeCountCommand(findCmd, osType, isSSH, sshTarget)
}

// checkForToolboxProcess checks for JetBrains Toolbox process (which includes TBCLi agent)
func checkForToolboxProcess(osType string, isSSH bool, sshTarget string) bool {
	var findCmd string

	switch osType {
	case PlatformLinux, PlatformMac:
		findCmd = `ps aux | grep "jetbrains-toolbox" | grep -v grep | wc -l`
	case PlatformWindows:
		findCmd = `powershell -Command "([array](tasklist /fo csv /fi 'imagename eq jetbrains-toolbox.exe' | Select-Object -Skip 1)).Count"`
	default:
		return false
	}

	return executeCountCommand(findCmd, osType, isSSH, sshTarget)
}

// executeCountCommand executes a command that returns a count and checks if > 0
func executeCountCommand(findCmd, osType string, isSSH bool, sshTarget string) bool {
	if isSSH && sshTarget != "" {
		sshArgs := []string{sshTarget, findCmd}
		cmd := exec.Command("ssh", sshArgs...) // #nosec G204 -- SSH with validated arguments

		output, err := cmd.Output()
		if err != nil {
			return false
		}

		count := strings.TrimSpace(string(output))
		processCount, err := strconv.Atoi(count)
		return err == nil && processCount > 0
	} else {
		var cmd *exec.Cmd

		if osType == PlatformWindows {
			cmd = exec.Command("cmd", "/C", findCmd) // #nosec G204 -- Windows command with validated arguments
		} else {
			cmd = exec.Command("sh", "-c", findCmd) // #nosec G204 -- Shell command with validated arguments
		}

		output, err := cmd.Output()
		if err != nil {
			return false
		}

		count := strings.TrimSpace(string(output))
		processCount, err := strconv.Atoi(count)
		return err == nil && processCount > 0
	}
}

// Stop existing TBCLi process if running
func stopExistingTbcli(osType string, isSSH bool, sshTarget string) error {
	logger.Info("Checking for existing TBCLi processes...")

	var findCmd, killCmd string

	switch osType {
	case PlatformLinux, PlatformMac:
		findCmd = `ps aux | grep "structured-logging agent" | grep -v grep | awk '{print $2}'`
		killCmd = "kill -TERM"
	case PlatformWindows:
		// Find process with "structured-logging agent" in command line on Windows
		// Use more reliable tasklist command
		findCmd = `tasklist /fo csv /fi "imagename eq java.exe" | findstr "structured-logging"`
		killCmd = "taskkill /PID"
	default:
		return fmt.Errorf("unsupported OS for process management: %s", osType)
	}

	if isSSH && sshTarget != "" {
		sshArgs := []string{sshTarget, findCmd}
		cmd := exec.Command("ssh", sshArgs...) // #nosec G204 -- SSH with validated arguments

		output, err := cmd.Output()
		if err != nil {
			// Check if it's a command failure (no processes) vs connection error
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
				logger.Info("No existing TBCLi processes found via SSH")
				return nil // Expected: command succeeded but found no processes
			}
			return fmt.Errorf("failed to check for processes via SSH: %v", err) // Real error
		}

		pids := strings.TrimSpace(string(output))
		if pids == "" {
			logger.Info("No existing TBCLi processes found via SSH")
			return nil
		}

		for _, pidLine := range strings.Split(pids, "\n") {
			pidLine = strings.TrimSpace(pidLine)
			if pidLine == "" {
				continue
			}

			var pid string
			if osType == PlatformWindows {
				// Parse CSV format from tasklist
				// Format: "java.exe","PID","Session Name","Session#","Mem Usage"
				fields := strings.Split(pidLine, ",")
				if len(fields) >= 2 {
					// Remove quotes from PID field
					pid = strings.Trim(fields[1], `"`)
				}
			} else {
				// For Unix systems, use the line as PID directly
				pid = pidLine
			}

			// Validate PID is numeric
			if pid == "" || !isNumeric(pid) {
				continue
			}

			logger.Info(fmt.Sprintf("Stopping TBCLi process via SSH: PID %s", pid))

			var killCommand string
			if osType == PlatformWindows {
				killCommand = fmt.Sprintf("%s %s /F", killCmd, pid)
			} else {
				killCommand = fmt.Sprintf("%s %s", killCmd, pid)
			}

			killArgs := []string{sshTarget, killCommand}
			killExec := exec.Command("ssh", killArgs...) // #nosec G204 -- SSH with validated arguments

			err := killExec.Run()
			if err != nil {
				// On Windows, taskkill can return exit status 128 if process doesn't exist or can't be killed
				// This is not necessarily an error
				if osType == PlatformWindows {
					logger.Info(fmt.Sprintf("TBCLi process %s may have already stopped or can't be killed via SSH", pid))
				} else {
					logger.Warning(fmt.Sprintf("Failed to stop TBCLi process %s via SSH: %v", pid, err))

					// Try force kill if graceful failed (Unix only)
					forceKillCommand := fmt.Sprintf("kill -KILL %s", pid)
					forceArgs := []string{sshTarget, forceKillCommand}
					forceExec := exec.Command("ssh", forceArgs...) // #nosec G204 -- SSH with validated arguments

					err = forceExec.Run()
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to force stop TBCLi process %s via SSH: %v", pid, err))
					} else {
						logger.Info(fmt.Sprintf("Force stopped TBCLi process via SSH: PID %s", pid))
					}
				}
			} else {
				logger.Success(fmt.Sprintf("Stopped TBCLi process via SSH: PID %s", pid))
			}
		}
	} else {
		var cmd *exec.Cmd

		if osType == PlatformWindows {
			cmd = exec.Command("cmd", "/C", findCmd) // #nosec G204 -- Windows command with validated arguments
		} else {
			cmd = exec.Command("sh", "-c", findCmd) // #nosec G204 -- Shell command with validated arguments
		}

		output, err := cmd.Output()
		if err != nil {
			// Check if it's a command failure (no processes) vs real error
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
				logger.Info("No existing TBCLi processes found locally")
				return nil // Expected: command succeeded but found no processes
			}
			return fmt.Errorf("failed to check for processes locally: %v", err) // Real error
		}

		pids := strings.TrimSpace(string(output))
		if pids == "" {
			logger.Info("No existing TBCLi processes found locally")
			return nil
		}

		// Kill each found PID
		for _, pidLine := range strings.Split(pids, "\n") {
			pidLine = strings.TrimSpace(pidLine)
			if pidLine == "" {
				continue
			}

			// Extract PID from tasklist CSV output or direct output
			var pid string
			if osType == PlatformWindows {
				// Parse CSV format from tasklist
				// Format: "java.exe","PID","Session Name","Session#","Mem Usage"
				fields := strings.Split(pidLine, ",")
				if len(fields) >= 2 {
					// Remove quotes from PID field
					pid = strings.Trim(fields[1], `"`)
				}
			} else {
				// For Unix systems, use the line as PID directly
				pid = pidLine
			}

			if pid == "" || !isNumeric(pid) {
				continue
			}

			logger.Info(fmt.Sprintf("Stopping TBCLi process locally: PID %s", pid))

			var killExec *exec.Cmd
			if osType == PlatformWindows {
				killExec = exec.Command("taskkill", "/PID", pid, "/F") // #nosec G204 -- Process termination with validated PID
			} else {
				killExec = exec.Command("kill", "-TERM", pid) // #nosec G204 -- Process termination with validated PID
			}

			err := killExec.Run()
			if err != nil {
				// On Windows, taskkill can return exit status 128 if process doesn't exist or can't be killed
				// This is not necessarily an error
				if osType == PlatformWindows {
					logger.Info(fmt.Sprintf("TBCLi process %s may have already stopped or can't be killed locally", pid))
				} else {
					logger.Warning(fmt.Sprintf("Failed to stop TBCLi process %s locally: %v", pid, err))

					// Try force kill if graceful failed (Unix only)
					forceExec := exec.Command("kill", "-KILL", pid) // #nosec G204 -- Process termination with validated PID
					err = forceExec.Run()
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to force stop TBCLi process %s locally: %v", pid, err))
					} else {
						logger.Info(fmt.Sprintf("Force stopped TBCLi process locally: PID %s", pid))
					}
				}
			} else {
				logger.Success(fmt.Sprintf("Stopped TBCLi process locally: PID %s", pid))
			}
		}
	}

	// Give processes time to shutdown gracefully
	time.Sleep(2 * time.Second)

	return nil
}

// findRemoteDevServer finds the remote-dev-server executable in the IDE installation directory
func findRemoteDevServer(idePath, ideType, ideVersion, osType, appPath string) (string, error) {
	var binPath string
	var executableName string

	switch osType {
	case PlatformWindows:
		executableName = "remote-dev-server.exe"
	default:
		executableName = RemoteDevServerName
	}

	switch osType {
	case PlatformMac:
		if appPath != "" {
			binPath = filepath.Join(appPath, "Contents", "bin", executableName)
		} else {
			entries, err := os.ReadDir(idePath)
			if err != nil {
				return "", fmt.Errorf("failed to read IDE directory: %v", err)
			}

			for _, entry := range entries {
				if strings.HasSuffix(entry.Name(), ".app") {
					appPath := filepath.Join(idePath, entry.Name())
					productInfoPath := filepath.Join(appPath, "Contents", "Resources", "product-info.json")

					if data, err := os.ReadFile(productInfoPath); err == nil {
						var productInfo ProductInfo
						if json.Unmarshal(data, &productInfo) == nil && isMatchingIDE(productInfo, ideType, ideVersion) {
							binPath = filepath.Join(appPath, "Contents", "bin", executableName)
							break
						}
					}
				}
			}

			if binPath == "" {
				return "", fmt.Errorf("no matching .app found for %s %s", ideType, ideVersion)
			}
		}
	case PlatformWindows:
		ideDir := filepath.Join(idePath, fmt.Sprintf("%s-%s", ideType, ideVersion))
		binPath = filepath.Join(ideDir, "bin", executableName)

		if _, err := os.Stat(binPath); os.IsNotExist(err) {
			if customPath := os.Getenv("IDE_INSTALL_PATH"); customPath != "" {
				customBinPath := filepath.Join(customPath, fmt.Sprintf("%s-%s", ideType, ideVersion), "bin", executableName)
				if _, err := os.Stat(customBinPath); err == nil {
					binPath = customBinPath
				}
			}

			commonPaths := []string{
				filepath.Join(os.Getenv("USERPROFILE"), "Documents", fmt.Sprintf("%s-%s", ideType, ideVersion)),
				filepath.Join(os.Getenv("PROGRAMFILES"), "JetBrains", fmt.Sprintf("%s-%s", ideType, ideVersion)),
				filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "JetBrains", fmt.Sprintf("%s-%s", ideType, ideVersion)),
				filepath.Join(os.Getenv("USERPROFILE"), "Documents", ideType),
				filepath.Join(os.Getenv("PROGRAMFILES"), "JetBrains", ideType),
			}

			for _, commonPath := range commonPaths {
				if commonPath == "" {
					continue
				}
				candidatePath := filepath.Join(commonPath, "bin", executableName)
				if _, err := os.Stat(candidatePath); err == nil {
					binPath = candidatePath
					break
				}
			}
		}
	default:
		ideDir := filepath.Join(idePath, fmt.Sprintf("%s-%s", ideType, ideVersion))
		binPath = filepath.Join(ideDir, "bin", executableName)
	}

	if _, err := os.Stat(binPath); err != nil {
		if os.IsNotExist(err) {
			if osType == PlatformWindows {
				return "", fmt.Errorf("remote-dev-server.exe not found at %s\n"+
					"Suggestions:\n"+
					"1. Verify IDE %s %s is installed\n"+
					"2. Check if IDE is installed in Documents folder: %s\n"+
					"3. Set IDE_INSTALL_PATH environment variable to your IDE installation directory\n"+
					"4. Use --ide-install-path flag to specify custom installation path",
					binPath, ideType, ideVersion,
					filepath.Join(os.Getenv("USERPROFILE"), "Documents", fmt.Sprintf("%s-%s", ideType, ideVersion)))
			}
			return "", fmt.Errorf("remote-dev-server not found at %s", binPath)
		}
		return "", fmt.Errorf("failed to check remote-dev-server: %v", err)
	}

	return binPath, nil
}

// launchIDELocally launches the IDE with specified arguments locally
func launchIDELocally(idePath, ideType, ideVersion, osType, args, appPath string) error {
	if args == "" {
		logger.Info("No IDE launch arguments provided, skipping IDE launch")
		return nil
	}

	logger.Info("Launching IDE (" + ideType + ") with arguments: " + args)

	execPath, err := findRemoteDevServer(idePath, ideType, ideVersion, osType, appPath)
	if err != nil {
		return fmt.Errorf("failed to find remote-dev-server: %v", err)
	}

	logger.Info(fmt.Sprintf("Found remote-dev-server at: %s", execPath))

	escapedArgs := parseAndEscapeArguments(args, osType)

	var cmd *exec.Cmd
	if osType == PlatformWindows {
		// Windows: Execute directly without cmd.exe to avoid path quoting issues
		// Set environment variable and execute the binary directly
		cmd = exec.Command(execPath) // #nosec G204 -- Executable path is validated
		if escapedArgs != "" {
			// Parse arguments properly for Windows
			args := parseArgsPreservingQuotes(escapedArgs)
			cmd.Args = append(cmd.Args, args...)
		}
		cmd.Env = append(os.Environ(), "REMOTE_DEV_NON_INTERACTIVE=1")

		// Set working directory to the IDE bin directory for proper DLL loading
		cmd.Dir = filepath.Dir(execPath)

		logger.Info(fmt.Sprintf("Executing Windows command: %s %s", execPath, escapedArgs))

		// For Windows, launch in goroutine to return control like launchIDERemotely
		handleCommandOutputAsync(cmd, true)

		time.Sleep(1 * time.Second)
		logger.Success("IDE launch initiated - process may take a moment to fully start on Windows")
		return nil
	} else {
		// Unix-like systems: use nohup for background execution
		cmd = exec.Command("nohup", execPath) // #nosec G204 -- Background process execution with validated path
		cmd.Args = append(cmd.Args, parseArgsPreservingQuotes(args)...)
		// Redirect output to avoid hanging
		cmd.Stdout = nil
		cmd.Stderr = nil

		ideDir := filepath.Dir(filepath.Dir(execPath))
		cmd.Dir = ideDir

		// Use Start() instead of Run() or CombinedOutput() to not wait for completion
		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to launch IDE: %v", err)
		}

		// Give IDE a moment to start
		time.Sleep(1 * time.Second)
		logger.Success("IDE launched successfully in background")
	}
	return nil
}

// escapeWindowsArgument properly escapes arguments for Windows cmd.exe
func escapeWindowsArgument(arg string) string {
	if strings.ContainsAny(arg, " \t\"&|<>^") {
		escaped := strings.ReplaceAll(arg, "\"", "\"\"")
		return "\"" + escaped + "\""
	}
	return arg
}

// parseAndEscapeArguments parses argument string and escapes each argument for the target OS
func parseAndEscapeArguments(argsString, osType string) string {
	if argsString == "" {
		return ""
	}

	parsedArgs := parseArgsPreservingQuotes(argsString)

	var escapedArgs []string
	for _, arg := range parsedArgs {
		if osType == PlatformWindows {
			escapedArgs = append(escapedArgs, escapeWindowsArgument(arg))
		} else {
			escapedArgs = append(escapedArgs, escapeUnixArgument(arg))
		}
	}

	return strings.Join(escapedArgs, " ")
}

// launchIDERemotely launches the IDE with specified arguments via SSH
func launchIDERemotely(sshTarget, idePath, ideType, ideVersion, osType, args string) error {
	if args == "" {
		logger.Info("No IDE launch arguments provided, skipping IDE launch")
		return nil
	}

	logger.Info(fmt.Sprintf("Launching IDE (%s) remotely with arguments: %s", ideType, args))

	var ideDir string
	var executableName string
	var execPath string

	switch osType {
	case PlatformWindows:
		executableName = "remote-dev-server.exe"
		ideDir = joinRemotePathForOS(osType, idePath, fmt.Sprintf("%s-%s", ideType, ideVersion))
		execPath = joinRemotePathForOS(osType, ideDir, "bin", executableName)
	case PlatformMac:
		executableName = RemoteDevServerName
		appPath, err := findIDEAppPath(idePath, ideType, ideVersion, true, sshTarget)
		if err != nil {
			return fmt.Errorf("failed to find IDE app: %v", err)
		}
		ideDir = appPath
		execPath = joinRemotePathForOS(osType, ideDir, "Contents", "bin", executableName)
	default:
		executableName = RemoteDevServerName
		ideDir = joinRemotePathForOS(osType, idePath, fmt.Sprintf("%s-%s", ideType, ideVersion))
		execPath = joinRemotePathForOS(osType, ideDir, "bin", executableName)
	}

	escapedArgs := parseAndEscapeArguments(args, osType)

	var runCmd string
	switch osType {
	case PlatformWindows:
		// Windows: use cmd.exe with REMOTE_DEV_NON_INTERACTIVE environment variable
		windowsExecPath := strings.ReplaceAll(execPath, "/", "\\")
		// Use cmd.exe /C with environment variable, similar to working manual command
		runCmd = fmt.Sprintf("cmd.exe /C set REMOTE_DEV_NON_INTERACTIVE=1 && \"%s\" %s",
			windowsExecPath, escapedArgs)
	default:
		// Linux/Mac: use nohup to detach
		runCmd = fmt.Sprintf("cd '%s' && nohup '%s' %s > /dev/null 2>&1 &",
			escapeUnixArgument(ideDir), escapeUnixArgument(execPath), escapedArgs)
	}

	cmd := exec.Command("ssh", sshTarget, runCmd) // #nosec G204 -- SSH target is validated
	logger.Info(fmt.Sprintf("Executing SSH command: ssh %s %s", sshTarget, runCmd))

	if osType == PlatformWindows {
		// For Windows, wait briefly to capture initial output then let it run
		handleCommandOutputAsync(cmd, true)

		time.Sleep(1 * time.Second)
		logger.Success("IDE launch initiated via SSH - process may take a moment to fully start on Windows")
	} else {
		// For Unix systems, use Start() to return control immediately
		// The nohup command ensures the process continues after SSH disconnects
		go func() {
			output, err := cmd.CombinedOutput()
			if err != nil {
				logger.Warning(fmt.Sprintf("IDE may have encountered an issue: %v", err))
			}
			if len(output) > 0 {
				outputStr := strings.TrimSpace(string(output))
				if strings.Contains(outputStr, "not found") {
					logger.Warning(fmt.Sprintf("remote-dev-server executable not found: %s", outputStr))
				} else if len(outputStr) > 0 {
					maxLen := 200
					if len(outputStr) < maxLen {
						maxLen = len(outputStr)
					}
					logger.Info(fmt.Sprintf("IDE output: %s", outputStr[:maxLen]))
				}
			}
		}()

		time.Sleep(1 * time.Second)
		logger.Success("IDE launched successfully on remote host in background")
	}

	return nil
}

// StringSlice implements flag.Value for []string to support multiple values
type StringSlice []string

func (s *StringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// parseIdeSpecs parses IDE types and versions into IdeSpec slice with default values
func parseIdeSpecs(ideTypes, ideVersions []string, defaultType, defaultVersion string) []IdeSpec {
	var specs []IdeSpec

	if len(ideTypes) == 0 {
		if defaultType != "" {
			specs = append(specs, IdeSpec{Type: defaultType, Version: defaultVersion})
		}
		return specs
	}

	if len(ideVersions) == 0 {
		versionToUse := defaultVersion
		if versionToUse == "" {
			versionToUse = ""
		}
		for _, ideType := range ideTypes {
			specs = append(specs, IdeSpec{Type: ideType, Version: versionToUse})
		}
		return specs
	}

	if len(ideTypes) == len(ideVersions) {
		for i, ideType := range ideTypes {
			specs = append(specs, IdeSpec{Type: ideType, Version: ideVersions[i]})
		}
		return specs
	}

	if len(ideVersions) == 1 {
		for _, ideType := range ideTypes {
			specs = append(specs, IdeSpec{Type: ideType, Version: ideVersions[0]})
		}
		return specs
	}

	// If multiple versions for single type, create multiple specs
	if len(ideTypes) == 1 {
		for _, version := range ideVersions {
			specs = append(specs, IdeSpec{Type: ideTypes[0], Version: version})
		}
		return specs
	}

	// Fallback: pair as many as possible, use first version for remaining types
	maxPairs := len(ideTypes)
	if len(ideVersions) < maxPairs {
		maxPairs = len(ideVersions)
	}

	for i := 0; i < maxPairs; i++ {
		specs = append(specs, IdeSpec{Type: ideTypes[i], Version: ideVersions[i]})
	}

	// Add remaining types with first version or default version
	if len(ideTypes) > len(ideVersions) {
		firstVersion := defaultVersion
		if len(ideVersions) > 0 {
			firstVersion = ideVersions[0]
		}
		for i := len(ideVersions); i < len(ideTypes); i++ {
			specs = append(specs, IdeSpec{Type: ideTypes[i], Version: firstVersion})
		}
	}

	return specs
}

// setupCommandLineFlags configures and parses command line arguments
func setupCommandLineFlags() (InstallOptions, string, bool, string, string, StringSlice, StringSlice, *bool) {
	var opts InstallOptions
	var showVersion bool
	var generateConfigPath string
	var originalTbcliVersion, originalToolboxVersion string
	var ideTypes, ideVersions StringSlice

	flag.StringVar(&opts.ConfigPath, "config", "", "Path to config file or URL")
	flag.StringVar(&generateConfigPath, "generate-config", "", "Generate config.json template at specified path")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&opts.InfoLogging, "info", false, "Enable detailed INFO logging")
	flag.StringVar(&originalTbcliVersion, "tbcli-version", "", "Set Toolbox CLI version")
	flag.StringVar(&originalToolboxVersion, "toolbox-version", "", "Set Toolbox version")
	flag.StringVar(&opts.JbrVersion, "jbr-version", "", "Set JBR version")
	flag.StringVar(&opts.JbrBuild, "jbr-build", "", "Set JBR build")
	flag.BoolVar(&opts.InstallIde, "install-ide", false, "Install IDE along with Toolbox CLI")
	flag.Var(&ideTypes, "ide-type", "Set IDE type (can be specified multiple times for multiple IDEs)")
	flag.Var(&ideVersions, "ide-version", "Set IDE version (can be specified multiple times for multiple versions)")
	flag.StringVar(&opts.IdeInstallPath, "ide-install-path", "", "Set custom IDE installation path")
	flag.StringVar(&opts.TbcliInstallPath, "tbcli-install-path", "", "Set custom Toolbox CLI installation path")
	flag.StringVar(&opts.JbrInstallPath, "jbr-install-path", "", "Set custom JBR installation path")
	flag.BoolVar(&opts.InstallClient, "install-client", false, "Install JetBrains Client")
	flag.BoolVar(&opts.NoClient, "no-client", false, "Skip JetBrains Client installation")
	flag.BoolVar(&opts.DownloadOnly, "download-tools", false, "Download only mode")
	flag.StringVar(&opts.DownloadOS, "download-os", "", "Download for specific OS")
	flag.StringVar(&opts.DownloadArch, "download-arch", "", "Download for specific architecture")
	flag.BoolVar(&opts.LocalInstall, "local", false, "Install locally")
	flag.BoolVar(&opts.NoProgress, "no-progress", false, "Disable progress bar")
	flag.BoolVar(&opts.Force, "force", false, "Force reinstallation (ignore existing installations)")
	flag.BoolVar(&opts.LaunchTbcli, "launch-tbcli", false, "Launch TBCLi after installation")
	flag.StringVar(&opts.LaunchIdeArgs, "launch-ide-args", "", "Arguments to pass to IDE remote-dev-server after installation (e.g., 'status')")
	flag.BoolVar(&opts.OpenToolboxURL, "open-toolbox-url", false, "Automatically open JetBrains Toolbox SSH URL in default application")

	helpFlag := flag.Bool("help", false, "Show help")
	flag.BoolVar(helpFlag, "h", false, "Show help")

	flag.Parse()

	return opts, generateConfigPath, showVersion, originalTbcliVersion, originalToolboxVersion, ideTypes, ideVersions, helpFlag
}

// processArguments processes additional command line arguments
func processArguments(opts *InstallOptions, originalTbcliVersion, originalToolboxVersion string) {
	for _, arg := range os.Args {
		if arg == "--install-ide" {
			opts.UserSetInstallIde = true
		}
		if arg == "--install-client" {
			opts.UserSetInstallClient = true
		}
	}

	if originalTbcliVersion != "" {
		opts.TbcliVersion = originalTbcliVersion
		opts.UserSetTbcliVersion = true
	}
	if originalToolboxVersion != "" {
		opts.ToolboxVersion = originalToolboxVersion
		opts.UserSetToolboxVersion = true
	}
}

// handleEarlyExitConditions checks for help, version, or config generation requests
func handleEarlyExitConditions(helpFlag *bool, showVersion bool, generateConfigPath string) bool {
	if *helpFlag {
		showHelp()
		return true
	}

	if showVersion {
		fmt.Println("JetBrains Toolbox CLI Installer v2.0 (Go)")
		return true
	}

	if generateConfigPath != "" {
		err := generateConfigTemplate(generateConfigPath)
		if err != nil {
			log.Fatalf("Failed to generate config template: %v", err)
		}
		fmt.Printf("Generated config template: %s\n", generateConfigPath)
		return true
	}

	return false
}

// initializeGlobalSettings sets up global application state
func initializeGlobalSettings(opts *InstallOptions) {
	logger = NewLogger(opts.InfoLogging)
	disableProgress = opts.NoProgress
	forceReinstall = opts.Force
}

// loadConfigSafely loads configuration with error handling
func loadConfigSafely(configPath string) Config {
	config, err := loadConfig(configPath)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to load config: %v, using defaults", err))
		config = getDefaultConfig()
	}
	return config
}

// applyConfigDefaults applies default values from configuration
func applyConfigDefaults(opts *InstallOptions, config Config) {
	if opts.TbcliVersion == "" {
		opts.TbcliVersion = config.DefaultTbcliVersion
	}
	if opts.ToolboxVersion == "" {
		opts.ToolboxVersion = config.DefaultToolboxVersion
	}
	if opts.JbrVersion == "" {
		opts.JbrVersion = config.DefaultJbrVersion
	}
	if opts.JbrBuild == "" {
		opts.JbrBuild = config.DefaultJbrBuild
	}
	if opts.IdeType == "" {
		if config.IdeType != "" {
			opts.IdeType = config.IdeType
		} else {
			opts.IdeType = config.DefaultIdeType
		}
	}
	if opts.IdeVersion == "" {
		if config.IdeVersion != "" {
			opts.IdeVersion = config.IdeVersion
		} else {
			opts.IdeVersion = config.DefaultIdeVersion
		}
	}
}

// applyConfigFlags applies boolean flags from configuration
func applyConfigFlags(opts *InstallOptions, config Config) {
	if !opts.UserSetInstallIde && config.InstallIde {
		opts.InstallIde = config.InstallIde
		logger.Info("Enabled IDE installation from config")
	}
	if !opts.UserSetInstallClient && !opts.NoClient && config.InstallClient {
		opts.InstallClient = config.InstallClient
		logger.Info("Enabled Client installation from config")
	}

	if opts.IdeInstallPath == "" && config.IdeInstallPath != "" {
		opts.IdeInstallPath = config.IdeInstallPath
		logger.Info(fmt.Sprintf("Using IDE install path from config: %s", opts.IdeInstallPath))
	}
	if opts.TbcliInstallPath == "" && config.TbcliInstallPath != "" {
		opts.TbcliInstallPath = config.TbcliInstallPath
		logger.Info(fmt.Sprintf("Using TBCLI install path from config: %s", opts.TbcliInstallPath))
	}
	if opts.JbrInstallPath == "" && config.JbrInstallPath != "" {
		opts.JbrInstallPath = config.JbrInstallPath
		logger.Info(fmt.Sprintf("Using JBR install path from config: %s", opts.JbrInstallPath))
	}
	if !opts.NoClient && config.NoClient {
		opts.NoClient = config.NoClient
		logger.Info("Disabled Client installation from config")
	}
	if opts.DownloadOS == "" && config.DownloadOS != "" {
		opts.DownloadOS = config.DownloadOS
		logger.Info(fmt.Sprintf("Using download OS from config: %s", opts.DownloadOS))
	}
	if opts.DownloadArch == "" && config.DownloadArch != "" {
		opts.DownloadArch = config.DownloadArch
		logger.Info(fmt.Sprintf("Using download arch from config: %s", opts.DownloadArch))
	}
	if !opts.DownloadOnly && config.DownloadOnly {
		opts.DownloadOnly = config.DownloadOnly
		logger.Info("Enabled download-only mode from config")
	}
	if !opts.LocalInstall && config.LocalInstall {
		opts.LocalInstall = config.LocalInstall
		logger.Info("Enabled local installation from config")
	}
	if !opts.NoProgress && config.NoProgress {
		opts.NoProgress = config.NoProgress
		logger.Info("Disabled progress bar from config")
	}
	if !opts.Force && config.Force {
		opts.Force = config.Force
		logger.Info("Enabled force reinstallation from config")
	}
	if !opts.InfoLogging && config.InfoLogging {
		opts.InfoLogging = config.InfoLogging
		logger = NewLogger(opts.InfoLogging)
		logger.Info("Enabled info logging from config")
	}
	if !opts.LaunchTbcli && config.LaunchTbcli {
		opts.LaunchTbcli = config.LaunchTbcli
		logger.Info("Enabled TBCLi launch from config")
	}
	if opts.LaunchIdeArgs == "" && config.LaunchIdeArgs != "" {
		opts.LaunchIdeArgs = config.LaunchIdeArgs
		logger.Info(fmt.Sprintf("Applied IDE launch arguments from config: %s", config.LaunchIdeArgs))
	}
	if !opts.OpenToolboxURL && config.OpenToolboxURL {
		opts.OpenToolboxURL = config.OpenToolboxURL
		logger.Info("Enabled automatic Toolbox URL opening from config")
	}
}

// configureIDESettings sets up IDE-related configuration
func configureIDESettings(opts *InstallOptions, config Config, ideTypes, ideVersions StringSlice) {
	opts.IdeTypes = []string(ideTypes)
	opts.IdeVersions = []string(ideVersions)

	if len(opts.IdeTypes) == 0 && len(config.IdeSpecs) > 0 {
		for _, spec := range config.IdeSpecs {
			opts.IdeTypes = append(opts.IdeTypes, spec.Type)
			opts.IdeVersions = append(opts.IdeVersions, spec.Version)
		}
		logger.Info(fmt.Sprintf("Using %d IDE specification(s) from config", len(config.IdeSpecs)))
	} else if len(opts.IdeTypes) == 0 && len(config.IdeTypes) > 0 {
		opts.IdeTypes = config.IdeTypes
		if len(config.IdeVersions) > 0 {
			opts.IdeVersions = config.IdeVersions
		}
		logger.Info(fmt.Sprintf("Using IDE types from config: %v", opts.IdeTypes))
		if len(opts.IdeVersions) > 0 {
			logger.Info(fmt.Sprintf("Using IDE versions from config: %v", opts.IdeVersions))
		}
	}

	if len(opts.IdeTypes) == 0 && config.DefaultIdeType != "" {
		opts.IdeTypes = append(opts.IdeTypes, config.DefaultIdeType)
		logger.Info(fmt.Sprintf("Using default IDE type from config: %s", config.DefaultIdeType))
	}

	defaultType := config.DefaultIdeType
	if config.IdeType != "" {
		defaultType = config.IdeType
	}
	defaultVersion := config.DefaultIdeVersion
	if config.IdeVersion != "" {
		defaultVersion = config.IdeVersion
	}
	opts.IdeSpecs = parseIdeSpecs(opts.IdeTypes, opts.IdeVersions, defaultType, defaultVersion)

	if len(opts.IdeSpecs) > 0 && !opts.UserSetInstallIde {
		opts.InstallIde = true
		logger.Info(fmt.Sprintf("Auto-enabling IDE installation for %d IDE(s)", len(opts.IdeSpecs)))
	}

	if opts.InstallIde && !opts.NoClient && !opts.UserSetInstallClient {
		opts.InstallClient = true
	}
}

// finalizeConfiguration finishes configuration setup
func finalizeConfiguration(opts *InstallOptions) {
	args := flag.Args()
	if len(args) > 0 {
		opts.SSHTarget = args[0]
	}

	if opts.DownloadOS != "" {
		opts.DownloadOS = normalizeOS(opts.DownloadOS)
	}

	if opts.DownloadOnly {
		if opts.DownloadOS == "" {
			opts.DownloadOS = ArchAll
		}
		if opts.DownloadArch == "" {
			opts.DownloadArch = ArchAll
		}
	}

	syncVersions(opts)
}

// executeInstallationMode runs the appropriate installation mode
func executeInstallationMode(config Config, opts InstallOptions) {
	var err error

	if opts.DownloadOnly || (opts.DownloadOS != "" && opts.SSHTarget == "") {
		logger.Info("Starting download-only mode")
		err = downloadForPlatforms(config, opts)
		if err != nil {
			logger.Error(fmt.Sprintf("Download failed: %v", err))
			os.Exit(1)
		}
	} else if opts.LocalInstall {
		logger.Info("Starting local installation mode")
		err = installLocally(config, opts)
		if err != nil {
			logger.Error(fmt.Sprintf("Local installation failed: %v", err))
			os.Exit(1)
		}
	} else if opts.SSHTarget != "" {
		if _, lookErr := exec.LookPath("ssh"); lookErr != nil {
			logger.Error("SSH client not found")
			os.Exit(1)
		}

		logger.Info("Starting SSH installation mode")
		err = installViaSSH(config, opts)
		if err != nil {
			logger.Error(fmt.Sprintf("SSH installation failed: %v", err))
			os.Exit(1)
		}
	} else {
		showHelp()
	}
}

func main() {
	opts, generateConfigPath, showVersion, originalTbcliVersion, originalToolboxVersion, ideTypes, ideVersions, helpFlag := setupCommandLineFlags()

	processArguments(&opts, originalTbcliVersion, originalToolboxVersion)

	if handleEarlyExitConditions(helpFlag, showVersion, generateConfigPath) {
		return
	}

	initializeGlobalSettings(&opts)

	config := loadConfigSafely(opts.ConfigPath)

	applyConfigDefaults(&opts, config)

	applyConfigFlags(&opts, config)

	configureIDESettings(&opts, config, ideTypes, ideVersions)

	finalizeConfiguration(&opts)

	executeInstallationMode(config, opts)
}

// createDefaultEnvironmentConfig creates a default environment configuration
func createDefaultEnvironmentConfig(customIDEPath string) ToolboxEnvironmentConfig {
	return ToolboxEnvironmentConfig{
		AllowPortForwarding: true,
		Tools: Tools{
			AllowInstallation:   true,
			AllowUpdate:         true,
			AllowUninstallation: true,
			Location: []Location{
				{
					Path:   customIDEPath,
					Levels: 2,
				},
			},
		},
	}
}

// locationExists checks if a location already exists in the slice
func locationExists(locations []Location, path string) bool {
	for _, loc := range locations {
		if loc.Path == path {
			return true
		}
	}
	return false
}

// mergeEnvironmentConfig merges a new location into existing environment config
func mergeEnvironmentConfig(existing ToolboxEnvironmentConfig, newPath string) ToolboxEnvironmentConfig {
	if !locationExists(existing.Tools.Location, newPath) {
		existing.Tools.Location = append(existing.Tools.Location, Location{
			Path:   newPath,
			Levels: 2,
		})
	}
	return existing
}

// readExistingEnvironmentConfig reads existing environment.json file locally
func readExistingEnvironmentConfig(filePath string) (ToolboxEnvironmentConfig, bool, error) {
	var config ToolboxEnvironmentConfig

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, false, nil // File doesn't exist, not an error
		}
		return config, false, fmt.Errorf("failed to read environment.json: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, false, fmt.Errorf("failed to parse environment.json: %v", err)
	}

	return config, true, nil
}

// readExistingEnvironmentConfigRemote reads existing environment.json file via SSH
func readExistingEnvironmentConfigRemote(sshTarget, filePath string, remoteInfo RemoteSystemInfo) (ToolboxEnvironmentConfig, bool, error) {
	var config ToolboxEnvironmentConfig

	// Check if file exists and read it
	var readCmd *exec.Cmd
	switch remoteInfo.OS {
	case PlatformWindows:
		readCmd = exec.Command("ssh", sshTarget, fmt.Sprintf(`powershell.exe -Command "if (Test-Path '%s') { Get-Content '%s' -Raw } else { Write-Output 'FILE_NOT_EXISTS' }"`, filePath, filePath)) // #nosec G204 -- SSH target is validated
	case PlatformLinux, PlatformMac:
		readCmd = exec.Command("ssh", sshTarget, fmt.Sprintf("if [ -f %s ]; then cat %s; else echo 'FILE_NOT_EXISTS'; fi", escapeUnixArgument(filePath), escapeUnixArgument(filePath))) // #nosec G204 -- SSH target is validated
	}

	output, err := readCmd.Output()
	if err != nil {
		return config, false, fmt.Errorf("failed to check/read environment.json: %v", err)
	}

	content := strings.TrimSpace(string(output))
	if content == "FILE_NOT_EXISTS" || content == "" {
		return config, false, nil // File doesn't exist, not an error
	}

	err = json.Unmarshal([]byte(content), &config)
	if err != nil {
		return config, false, fmt.Errorf("failed to parse environment.json: %v", err)
	}

	return config, true, nil
}
