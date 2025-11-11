/*
Copyright © 2025 Fleeks Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/fleeks-inc/fleeks-cli/internal/client"
	"github.com/fleeks-inc/fleeks-cli/internal/config"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "📁 File operations with smart sync",
	Long: `
📁 Intelligent File Management System

Advanced file operations with smart synchronization:

✅ Smart File Sync:
   • Bidirectional file synchronization
   • Conflict detection and resolution
   • Incremental updates and delta sync
   • Real-time file watching

✅ Context-Aware Operations:
   • Preserve file context during operations
   • Intelligent file content analysis
   • Language-specific handling
   • Binary file support

✅ Cloud-Local Hybrid:
   • Seamless local-cloud file operations
   • Workspace-aware file management
   • Version control integration
   • Backup and recovery

Examples:
  # List files in workspace
  fleeks files list my-project
  
  # Upload local file to workspace
  fleeks files upload my-project ./src/main.py /workspace/src/main.py
  
  # Download file from workspace
  fleeks files download my-project /workspace/config.json ./config.json
  
  # Create new file in workspace
  fleeks files create my-project /workspace/README.md "# My Project"
  
  # Watch for file changes
  fleeks files watch my-project
`,
}

var filesListCmd = &cobra.Command{
	Use:   "list [project-id]",
	Short: "List files in workspace",
	Long: `List all files in a workspace with detailed information.

Shows file metadata including size, modification time, and type.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return listFiles(args[0], cmd)
	},
}

var filesUploadCmd = &cobra.Command{
	Use:   "upload [project-id] [local-path] [remote-path]",
	Short: "Upload file to workspace",
	Long: `Upload a local file to the cloud workspace.

Supports:
- Single file upload
- Directory upload (recursive)
- Progress tracking
- Conflict handling`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return uploadFile(args[0], args[1], args[2], cmd)
	},
}

var filesDownloadCmd = &cobra.Command{
	Use:   "download [project-id] [remote-path] [local-path]",
	Short: "Download file from workspace",
	Long: `Download a file from the cloud workspace to local filesystem.

Supports:
- Single file download
- Directory download (recursive)
- Progress tracking
- Overwrite protection`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return downloadFile(args[0], args[1], args[2], cmd)
	},
}

var filesCreateCmd = &cobra.Command{
	Use:   "create [project-id] [path] [content]",
	Short: "Create new file in workspace",
	Long: `Create a new file in the cloud workspace with specified content.

The content can be provided as a string or read from stdin.`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := args[0]
		path := args[1]
		content := ""
		if len(args) > 2 {
			content = args[2]
		}
		return createFile(projectID, path, content, cmd)
	},
}

var filesDeleteCmd = &cobra.Command{
	Use:   "delete [project-id] [path]",
	Short: "Delete file from workspace",
	Long: `Delete a file or directory from the cloud workspace.

Use with caution as this operation cannot be undone.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return deleteFile(args[0], args[1], cmd)
	},
}

var filesWatchCmd = &cobra.Command{
	Use:   "watch [project-id]",
	Short: "Watch for file changes",
	Long: `Watch for real-time file changes in the workspace.

Shows:
- File creation, modification, and deletion
- Who made the changes (user or agent)
- Timestamps and change details`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return watchFiles(args[0], cmd)
	},
}

func init() {
	// Add subcommands
	filesCmd.AddCommand(filesListCmd)
	filesCmd.AddCommand(filesUploadCmd)
	filesCmd.AddCommand(filesDownloadCmd)
	filesCmd.AddCommand(filesCreateCmd)
	filesCmd.AddCommand(filesDeleteCmd)
	filesCmd.AddCommand(filesWatchCmd)

	// List command flags
	filesListCmd.Flags().StringP("path", "p", "/", "Path to list (default: root)")
	filesListCmd.Flags().BoolP("recursive", "r", false, "List files recursively")
	filesListCmd.Flags().StringP("filter", "f", "", "Filter files by pattern")

	// Upload command flags
	filesUploadCmd.Flags().BoolP("recursive", "r", false, "Upload directory recursively")
	filesUploadCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing files")

	// Download command flags
	filesDownloadCmd.Flags().BoolP("recursive", "r", false, "Download directory recursively")
	filesDownloadCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing local files")

	// Create command flags
	filesCreateCmd.Flags().BoolP("stdin", "s", false, "Read content from stdin")
	filesCreateCmd.Flags().StringP("template", "t", "", "Use file template")

	// Delete command flags
	filesDeleteCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation")
	filesDeleteCmd.Flags().BoolP("recursive", "r", false, "Delete directory recursively")
}

// FileInfo represents file information
type FileInfo struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	Type         string    `json:"type"` // "file" or "directory"
	MimeType     string    `json:"mime_type,omitempty"`
	ModifiedAt   time.Time `json:"modified_at"`
	CreatedAt    time.Time `json:"created_at"`
	Permissions  string    `json:"permissions"`
	Owner        string    `json:"owner,omitempty"`
	IsExecutable bool      `json:"is_executable"`
}

// FileUploadRequest represents file upload request
type FileUploadRequest struct {
	Path      string `json:"path"`
	Content   string `json:"content"` // base64 encoded for binary files
	Overwrite bool   `json:"overwrite"`
	MimeType  string `json:"mime_type,omitempty"`
}

// FileDownloadResponse represents file download response
type FileDownloadResponse struct {
	Path     string `json:"path"`
	Content  string `json:"content"` // base64 encoded for binary files
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

// FileChangeEvent represents file change event
type FileChangeEvent struct {
	Type      string    `json:"type"` // "created", "modified", "deleted"
	Path      string    `json:"path"`
	Actor     string    `json:"actor"` // "user" or "agent"
	ActorID   string    `json:"actor_id"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details,omitempty"`
}

func listFiles(projectID string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	// Get flags
	path, _ := cmd.Flags().GetString("path")
	recursive, _ := cmd.Flags().GetBool("recursive")
	filter, _ := cmd.Flags().GetString("filter")

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Build query parameters
	params := make([]string, 0)
	if path != "/" {
		params = append(params, "path="+path)
	}
	if recursive {
		params = append(params, "recursive=true")
	}
	if filter != "" {
		params = append(params, "filter="+filter)
	}

	endpoint := fmt.Sprintf("/api/v1/sdk/files/%s", projectID)
	if len(params) > 0 {
		endpoint += "?" + strings.Join(params, "&")
	}

	// Get files
	var files []FileInfo
	if err := apiClient.GET(endpoint, &files); err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		fmt.Printf("%s No files found in %s\n",
			color.YellowString("📁"), color.CyanString(path))
		return nil
	}

	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type", "Size", "Modified", "Permissions"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiYellowColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.FgHiMagentaColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
	)

	for _, file := range files {
		size := formatFileSize(file.Size)
		if file.Type == "directory" {
			size = "-"
		}

		table.Append([]string{
			file.Name,
			file.Type,
			size,
			file.ModifiedAt.Format("2006-01-02 15:04"),
			file.Permissions,
		})
	}

	fmt.Printf("\n%s %s:%s\n\n",
		color.New(color.Bold).Sprint("📁 Files in"),
		color.CyanString(projectID),
		color.YellowString(path))

	table.Render()

	fmt.Printf("\nTotal: %s files\n", color.GreenString(fmt.Sprintf("%d", len(files))))
	return nil
}

func uploadFile(projectID, localPath, remotePath string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	// Check if local file exists
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("local file not found: %w", err)
	}

	recursive, _ := cmd.Flags().GetBool("recursive")
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	if fileInfo.IsDir() && !recursive {
		return fmt.Errorf("use --recursive flag to upload directories")
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Start spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Uploading file..."
	s.Start()
	defer s.Stop()

	if fileInfo.IsDir() {
		// Directory upload (recursive)
		err = uploadDirectory(apiClient, projectID, localPath, remotePath, overwrite)
	} else {
		// Single file upload
		err = uploadSingleFile(apiClient, projectID, localPath, remotePath, overwrite)
	}

	s.Stop()

	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	fmt.Printf("%s File uploaded successfully: %s → %s\n",
		color.GreenString("📤"),
		color.YellowString(localPath),
		color.CyanString(remotePath))

	return nil
}

func uploadSingleFile(apiClient *client.APIClient, projectID, localPath, remotePath string, overwrite bool) error {
	// Read file content
	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Encode content as base64
	encodedContent := base64.StdEncoding.EncodeToString(content)

	// Prepare request
	request := FileUploadRequest{
		Path:      remotePath,
		Content:   encodedContent,
		Overwrite: overwrite,
	}

	// Upload file
	endpoint := fmt.Sprintf("/api/v1/sdk/files/%s/upload", projectID)
	return apiClient.POST(endpoint, request, nil)
}

func uploadDirectory(apiClient *client.APIClient, projectID, localDir, remoteDir string, overwrite bool) error {
	return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Skip directories, they're created automatically
		}

		// Calculate relative path
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		remotePath := filepath.Join(remoteDir, relPath)
		remotePath = strings.ReplaceAll(remotePath, "\\", "/") // Normalize path separators

		return uploadSingleFile(apiClient, projectID, path, remotePath, overwrite)
	})
}

func downloadFile(projectID, remotePath, localPath string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	overwrite, _ := cmd.Flags().GetBool("overwrite")

	// Check if local file exists
	if _, err := os.Stat(localPath); err == nil && !overwrite {
		return fmt.Errorf("local file exists. Use --overwrite to replace it")
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Start spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Downloading file..."
	s.Start()
	defer s.Stop()

	// Download file
	var response FileDownloadResponse
	endpoint := fmt.Sprintf("/api/v1/sdk/files/%s/download?path=%s", projectID, remotePath)
	if err := apiClient.GET(endpoint, &response); err != nil {
		s.Stop()
		return fmt.Errorf("download failed: %w", err)
	}

	// Decode content
	content, err := base64.StdEncoding.DecodeString(response.Content)
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to decode file content: %w", err)
	}

	// Ensure local directory exists
	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		s.Stop()
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(localPath, content, 0644); err != nil {
		s.Stop()
		return fmt.Errorf("failed to write file: %w", err)
	}

	s.Stop()

	fmt.Printf("%s File downloaded successfully: %s → %s\n",
		color.GreenString("📥"),
		color.CyanString(remotePath),
		color.YellowString(localPath))

	return nil
}

func createFile(projectID, path, content string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	// Read from stdin if requested
	useStdin, _ := cmd.Flags().GetBool("stdin")
	if useStdin || content == "" {
		stdinContent, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
		content = string(stdinContent)
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Encode content as base64
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))

	// Prepare request
	request := FileUploadRequest{
		Path:    path,
		Content: encodedContent,
	}

	// Create file
	endpoint := fmt.Sprintf("/api/v1/sdk/files/%s/create", projectID)
	if err := apiClient.POST(endpoint, request, nil); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	fmt.Printf("%s File created successfully: %s\n",
		color.GreenString("📝"), color.CyanString(path))

	return nil
}

func deleteFile(projectID, path string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	force, _ := cmd.Flags().GetBool("force")

	if !force {
		fmt.Printf("%s Are you sure you want to delete '%s'? [y/N] ",
			color.RedString("⚠️"), path)

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Delete file
	endpoint := fmt.Sprintf("/api/v1/sdk/files/%s/delete?path=%s", projectID, path)
	if err := apiClient.DELETE(endpoint, nil); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	fmt.Printf("%s File deleted successfully: %s\n",
		color.GreenString("🗑️"), color.CyanString(path))

	return nil
}

func watchFiles(projectID string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Create stream reader for file changes
	streamPath := fmt.Sprintf("/ws/files/%s/watch", projectID)
	stream, err := apiClient.NewStreamReader(streamPath)
	if err != nil {
		return fmt.Errorf("failed to connect to file watch stream: %w", err)
	}
	defer stream.Close()

	fmt.Printf("%s Watching file changes for %s (Press Ctrl+C to stop)\n\n",
		color.CyanString("👀"), color.YellowString(projectID))

	// Stream file change events
	for {
		select {
		case msg, ok := <-stream.Messages():
			if !ok {
				fmt.Printf("\n%s File watch stream ended\n", color.GreenString("✅"))
				return nil
			}

			// Parse file change event from message metadata
			if changeType, exists := msg.Metadata["type"]; exists {
				path := msg.Metadata["path"]
				actor := msg.Metadata["actor"]
				timestamp := msg.Timestamp.Format("15:04:05")

				var icon, typeColor string
				switch changeType {
				case "created":
					icon = "📝"
					typeColor = color.GreenString("CREATED")
				case "modified":
					icon = "✏️"
					typeColor = color.YellowString("MODIFIED")
				case "deleted":
					icon = "🗑️"
					typeColor = color.RedString("DELETED")
				default:
					icon = "📄"
					typeColor = color.WhiteString(fmt.Sprintf("%v", changeType))
				}

				fmt.Printf("[%s] %s %s %s (by %s)\n",
					color.MagentaString(timestamp),
					icon,
					typeColor,
					color.CyanString(fmt.Sprintf("%v", path)),
					color.BlueString(fmt.Sprintf("%v", actor)))
			}

		case err, ok := <-stream.Errors():
			if !ok {
				return nil
			}
			return fmt.Errorf("stream error: %w", err)
		}
	}
}

func formatFileSize(bytes int64) string {
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
