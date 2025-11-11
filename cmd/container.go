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
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/fleeks-inc/fleeks-cli/internal/client"
	"github.com/fleeks-inc/fleeks-cli/internal/config"
)

// containerCmd represents the container command
var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "🐳 Container orchestration and management",
	Long: `
🐳 Revolutionary Cloud Container Orchestration

Unique container management capabilities no competitor has:

✅ Ready Container Pool System:
   • Sub-second container assignment (<100ms)
   • Pre-warmed containers for instant workspace access
   • OverlayFS for instant workspace mounting
   • Template-based container management

✅ Real-time Container Monitoring:
   • Live resource usage streaming
   • Container health monitoring
   • Performance metrics and alerting
   • Multi-language environment support

✅ Advanced Container Operations:
   • Execute commands in cloud containers
   • Stream logs with filtering
   • Scale container resources
   • Container lifecycle management

Examples:
  # Get container information
  fleeks container info my-api
  
  # Monitor real-time stats
  fleeks container stats my-api --watch
  
  # Execute commands in container
  fleeks container exec my-api "npm install"
  
  # Stream container logs
  fleeks container logs my-api --follow --tail 100
  
  # Scale container resources
  fleeks container scale my-api --cpu 2 --memory 4G
`,
}

var containerInfoCmd = &cobra.Command{
	Use:   "info [project-id]",
	Short: "Get container information",
	Long: `Get detailed information about a workspace container including:
- Container status and health
- Resource allocations and usage
- Template and language support
- Network configuration
- Mount points and storage`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return getContainerInfo(args[0], cmd)
	},
}

var containerStatsCmd = &cobra.Command{
	Use:   "stats [project-id]",
	Short: "Show container resource statistics",
	Long: `Display real-time resource usage statistics for a container.

Shows:
- CPU usage and limits
- Memory usage and limits  
- Disk I/O and usage
- Network I/O
- Process count`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return getContainerStats(args[0], cmd)
	},
}

var containerLogsCmd = &cobra.Command{
	Use:   "logs [project-id]",
	Short: "Show container logs",
	Long: `Show logs from the workspace container.

Supports:
- Real-time log streaming
- Historical log retrieval
- Log filtering and search
- Multiple output formats`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return getContainerLogs(args[0], cmd)
	},
}

var containerExecCmd = &cobra.Command{
	Use:   "exec [project-id] [command]",
	Short: "Execute command in container",
	Long: `Execute a command inside the workspace container.

Features:
- Interactive and non-interactive execution
- Environment variable support
- Working directory specification
- Output streaming`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := args[0]
		command := strings.Join(args[1:], " ")
		return execInContainer(projectID, command, cmd)
	},
}

var containerScaleCmd = &cobra.Command{
	Use:   "scale [project-id]",
	Short: "Scale container resources",
	Long: `Scale container CPU and memory resources.

This allows dynamic resource allocation based on workload requirements.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return scaleContainer(args[0], cmd)
	},
}

func init() {
	// Add subcommands
	containerCmd.AddCommand(containerInfoCmd)
	containerCmd.AddCommand(containerStatsCmd)
	containerCmd.AddCommand(containerLogsCmd)
	containerCmd.AddCommand(containerExecCmd)
	containerCmd.AddCommand(containerScaleCmd)

	// Stats command flags
	containerStatsCmd.Flags().BoolP("watch", "w", false, "Watch stats in real-time")
	containerStatsCmd.Flags().IntP("interval", "i", 5, "Update interval in seconds")

	// Logs command flags
	containerLogsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
	containerLogsCmd.Flags().IntP("tail", "t", 50, "Number of lines to show from the end")
	containerLogsCmd.Flags().StringP("since", "s", "", "Show logs since timestamp (e.g. 2023-01-01T00:00:00Z)")
	containerLogsCmd.Flags().StringP("filter", "", "", "Filter logs by pattern")

	// Exec command flags
	containerExecCmd.Flags().BoolP("interactive", "i", false, "Interactive mode")
	containerExecCmd.Flags().BoolP("tty", "t", false, "Allocate a pseudo-TTY")
	containerExecCmd.Flags().StringP("workdir", "w", "", "Working directory")
	containerExecCmd.Flags().StringSliceP("env", "e", []string{}, "Environment variables")

	// Scale command flags
	containerScaleCmd.Flags().StringP("cpu", "", "", "CPU allocation (e.g. 1, 2, 0.5)")
	containerScaleCmd.Flags().StringP("memory", "", "", "Memory allocation (e.g. 1G, 512M, 2048M)")
}

// ContainerInfo represents container information
type ContainerInfo struct {
	ContainerID string            `json:"container_id"`
	ProjectID   string            `json:"project_id"`
	Status      string            `json:"status"`
	Template    string            `json:"template"`
	Languages   []string          `json:"languages"`
	Created     time.Time         `json:"created"`
	Started     time.Time         `json:"started"`
	Image       string            `json:"image"`
	Platform    string            `json:"platform"`
	Resources   ResourceInfo      `json:"resources"`
	Network     NetworkInfo       `json:"network"`
	Mounts      []MountInfo       `json:"mounts"`
	Environment map[string]string `json:"environment"`
	Health      HealthInfo        `json:"health"`
}

// ResourceInfo represents container resource information
type ResourceInfo struct {
	CPU       string `json:"cpu"`
	Memory    string `json:"memory"`
	Disk      string `json:"disk"`
	CPULimit  string `json:"cpu_limit"`
	MemLimit  string `json:"mem_limit"`
	DiskLimit string `json:"disk_limit"`
}

// NetworkInfo represents container network information
type NetworkInfo struct {
	IPAddress string            `json:"ip_address"`
	Ports     map[string]string `json:"ports"`
	Network   string            `json:"network"`
}

// MountInfo represents container mount information
type MountInfo struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Type        string `json:"type"`
	ReadOnly    bool   `json:"read_only"`
}

// HealthInfo represents container health information
type HealthInfo struct {
	Status      string    `json:"status"`
	LastCheck   time.Time `json:"last_check"`
	FailCount   int       `json:"fail_count"`
	Description string    `json:"description"`
}

// ContainerStats represents real-time container statistics
type ContainerStats struct {
	ContainerID   string    `json:"container_id"`
	ProjectID     string    `json:"project_id"`
	Timestamp     time.Time `json:"timestamp"`
	CPU           float64   `json:"cpu_percent"`
	Memory        int64     `json:"memory_bytes"`
	MemoryPercent float64   `json:"memory_percent"`
	DiskRead      int64     `json:"disk_read_bytes"`
	DiskWrite     int64     `json:"disk_write_bytes"`
	NetRx         int64     `json:"network_rx_bytes"`
	NetTx         int64     `json:"network_tx_bytes"`
	Processes     int       `json:"process_count"`
}

// ExecRequest represents command execution request
type ExecRequest struct {
	Command     string            `json:"command"`
	Interactive bool              `json:"interactive"`
	TTY         bool              `json:"tty"`
	WorkDir     string            `json:"workdir,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

// ExecResponse represents command execution response
type ExecResponse struct {
	ExecID   string `json:"exec_id"`
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
}

func getContainerInfo(projectID string, cmd *cobra.Command) error {
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

	// Get container info
	var container ContainerInfo
	endpoint := fmt.Sprintf("/api/v1/sdk/containers/%s", projectID)
	if err := apiClient.GET(endpoint, &container); err != nil {
		return fmt.Errorf("failed to get container info: %w", err)
	}

	// Display container information
	fmt.Printf("\n%s %s\n\n",
		color.New(color.Bold).Sprint("🐳 Container Information:"),
		color.CyanString(projectID))

	fmt.Printf("%-15s %s\n", "Container ID:", color.BlueString(container.ContainerID))
	fmt.Printf("%-15s %s\n", "Project ID:", color.CyanString(container.ProjectID))
	fmt.Printf("%-15s %s\n", "Status:", getStatusColor(container.Status))
	fmt.Printf("%-15s %s\n", "Template:", color.YellowString(container.Template))
	fmt.Printf("%-15s %s\n", "Image:", container.Image)
	fmt.Printf("%-15s %s\n", "Platform:", container.Platform)
	fmt.Printf("%-15s %s\n", "Created:", color.MagentaString(container.Created.Format("2006-01-02 15:04:05")))
	fmt.Printf("%-15s %s\n", "Started:", color.MagentaString(container.Started.Format("2006-01-02 15:04:05")))

	// Languages
	if len(container.Languages) > 0 {
		fmt.Printf("%-15s %s\n", "Languages:", color.GreenString(strings.Join(container.Languages, ", ")))
	}

	// Resources
	fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("📊 Resources:"))
	fmt.Printf("%-15s %s (limit: %s)\n", "CPU:", container.Resources.CPU, container.Resources.CPULimit)
	fmt.Printf("%-15s %s (limit: %s)\n", "Memory:", container.Resources.Memory, container.Resources.MemLimit)
	fmt.Printf("%-15s %s (limit: %s)\n", "Disk:", container.Resources.Disk, container.Resources.DiskLimit)

	// Network
	fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("🌐 Network:"))
	fmt.Printf("%-15s %s\n", "IP Address:", container.Network.IPAddress)
	fmt.Printf("%-15s %s\n", "Network:", container.Network.Network)
	if len(container.Network.Ports) > 0 {
		fmt.Printf("%-15s\n", "Ports:")
		for internal, external := range container.Network.Ports {
			fmt.Printf("  %s → %s\n", internal, external)
		}
	}

	// Health
	fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("❤️  Health:"))
	fmt.Printf("%-15s %s\n", "Status:", getHealthColor(container.Health.Status))
	fmt.Printf("%-15s %s\n", "Last Check:", color.MagentaString(container.Health.LastCheck.Format("2006-01-02 15:04:05")))
	if container.Health.FailCount > 0 {
		fmt.Printf("%-15s %s\n", "Fail Count:", color.RedString(fmt.Sprintf("%d", container.Health.FailCount)))
	}
	if container.Health.Description != "" {
		fmt.Printf("%-15s %s\n", "Description:", container.Health.Description)
	}

	// Mounts
	if len(container.Mounts) > 0 {
		fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("💾 Mounts:"))
		for _, mount := range container.Mounts {
			access := "rw"
			if mount.ReadOnly {
				access = "ro"
			}
			fmt.Printf("  %s → %s (%s, %s)\n",
				mount.Source, mount.Destination, mount.Type, access)
		}
	}

	return nil
}

func getContainerStats(projectID string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	watch, _ := cmd.Flags().GetBool("watch")
	interval, _ := cmd.Flags().GetInt("interval")

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	if !watch {
		// One-time stats
		var stats ContainerStats
		endpoint := fmt.Sprintf("/api/v1/sdk/containers/%s/stats", projectID)
		if err := apiClient.GET(endpoint, &stats); err != nil {
			return fmt.Errorf("failed to get container stats: %w", err)
		}

		displayStats(stats)
		return nil
	}

	// Watch mode - real-time stats
	fmt.Printf("%s Monitoring container %s (Press Ctrl+C to stop)\n\n",
		color.CyanString("📊"), color.YellowString(projectID))

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n%s Stopping stats monitoring...\n",
			color.YellowString("🛑"))
		cancel()
	}()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			var stats ContainerStats
			endpoint := fmt.Sprintf("/api/v1/sdk/containers/%s/stats", projectID)
			if err := apiClient.GET(endpoint, &stats); err != nil {
				fmt.Printf("Error getting stats: %v\n", err)
				continue
			}

			// Clear screen and display stats
			fmt.Print("\033[2J\033[H")
			fmt.Printf("%s Container Stats - %s\n\n",
				color.New(color.Bold).Sprint("📊"),
				color.CyanString(projectID))
			displayStats(stats)
		}
	}
}

func displayStats(stats ContainerStats) {
	timestamp := stats.Timestamp.Format("15:04:05")

	fmt.Printf("%-15s %s\n", "Timestamp:", color.MagentaString(timestamp))
	fmt.Printf("%-15s %s\n", "CPU Usage:", color.GreenString(fmt.Sprintf("%.1f%%", stats.CPU)))
	fmt.Printf("%-15s %s (%s)\n", "Memory:",
		formatBytes(stats.Memory),
		color.BlueString(fmt.Sprintf("%.1f%%", stats.MemoryPercent)))
	fmt.Printf("%-15s %s\n", "Processes:", color.YellowString(fmt.Sprintf("%d", stats.Processes)))

	fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("💾 Disk I/O:"))
	fmt.Printf("%-15s %s\n", "Read:", formatBytes(stats.DiskRead))
	fmt.Printf("%-15s %s\n", "Write:", formatBytes(stats.DiskWrite))

	fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("🌐 Network I/O:"))
	fmt.Printf("%-15s %s\n", "RX:", formatBytes(stats.NetRx))
	fmt.Printf("%-15s %s\n", "TX:", formatBytes(stats.NetTx))
}

func getContainerLogs(projectID string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	follow, _ := cmd.Flags().GetBool("follow")
	tail, _ := cmd.Flags().GetInt("tail")
	since, _ := cmd.Flags().GetString("since")
	filter, _ := cmd.Flags().GetString("filter")

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Build query parameters
	params := make([]string, 0)
	if tail > 0 {
		params = append(params, fmt.Sprintf("tail=%d", tail))
	}
	if since != "" {
		params = append(params, "since="+since)
	}
	if filter != "" {
		params = append(params, "filter="+filter)
	}

	endpoint := fmt.Sprintf("/api/v1/sdk/containers/%s/logs", projectID)
	if len(params) > 0 {
		endpoint += "?" + strings.Join(params, "&")
	}

	if !follow {
		// One-time logs
		var logs []string
		if err := apiClient.GET(endpoint, &logs); err != nil {
			return fmt.Errorf("failed to get container logs: %w", err)
		}

		for _, line := range logs {
			fmt.Println(line)
		}
		return nil
	}

	// Follow mode - stream logs
	fmt.Printf("%s Following logs for %s (Press Ctrl+C to stop)\n\n",
		color.CyanString("📜"), color.YellowString(projectID))

	// Create stream reader for logs
	streamPath := fmt.Sprintf("/ws/containers/%s/logs", projectID)
	stream, err := apiClient.NewStreamReader(streamPath)
	if err != nil {
		return fmt.Errorf("failed to connect to log stream: %w", err)
	}
	defer stream.Close()

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	// Stream logs
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-stream.Messages():
			if !ok {
				return nil
			}
			fmt.Println(msg.Content)
		case err, ok := <-stream.Errors():
			if !ok {
				return nil
			}
			return fmt.Errorf("stream error: %w", err)
		}
	}
}

func execInContainer(projectID, command string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	// Get flags
	interactive, _ := cmd.Flags().GetBool("interactive")
	tty, _ := cmd.Flags().GetBool("tty")
	workdir, _ := cmd.Flags().GetString("workdir")
	envVars, _ := cmd.Flags().GetStringSlice("env")

	// Parse environment variables
	environment := make(map[string]string)
	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			environment[parts[0]] = parts[1]
		}
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Prepare request
	request := ExecRequest{
		Command:     command,
		Interactive: interactive,
		TTY:         tty,
		WorkDir:     workdir,
		Environment: environment,
	}

	// Start spinner for non-interactive commands
	var s *spinner.Spinner
	if !interactive {
		s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Executing command..."
		s.Start()
		defer s.Stop()
	}

	// Execute command
	var response ExecResponse
	endpoint := fmt.Sprintf("/api/v1/sdk/containers/%s/exec", projectID)
	if err := apiClient.POST(endpoint, request, &response); err != nil {
		if s != nil {
			s.Stop()
		}
		return fmt.Errorf("failed to execute command: %w", err)
	}

	if s != nil {
		s.Stop()
	}

	// Display output
	if response.Output != "" {
		fmt.Print(response.Output)
	}

	if response.Error != "" {
		fmt.Fprintf(os.Stderr, "%s\n", color.RedString(response.Error))
	}

	// Exit with same code as the command
	if response.ExitCode != 0 {
		os.Exit(response.ExitCode)
	}

	return nil
}

func scaleContainer(projectID string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	cpu, _ := cmd.Flags().GetString("cpu")
	memory, _ := cmd.Flags().GetString("memory")

	if cpu == "" && memory == "" {
		return fmt.Errorf("at least one of --cpu or --memory must be specified")
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Prepare scale request
	scaleRequest := make(map[string]string)
	if cpu != "" {
		scaleRequest["cpu"] = cpu
	}
	if memory != "" {
		scaleRequest["memory"] = memory
	}

	// Start spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Scaling container resources..."
	s.Start()
	defer s.Stop()

	// Scale container
	endpoint := fmt.Sprintf("/api/v1/sdk/containers/%s/scale", projectID)
	if err := apiClient.POST(endpoint, scaleRequest, nil); err != nil {
		s.Stop()
		return fmt.Errorf("failed to scale container: %w", err)
	}

	s.Stop()

	fmt.Printf("%s Container %s scaled successfully\n",
		color.GreenString("📈"), color.CyanString(projectID))

	if cpu != "" {
		fmt.Printf("CPU:    %s\n", color.YellowString(cpu))
	}
	if memory != "" {
		fmt.Printf("Memory: %s\n", color.BlueString(memory))
	}

	return nil
}

func getHealthColor(status string) string {
	switch status {
	case "healthy":
		return color.GreenString(status)
	case "unhealthy":
		return color.RedString(status)
	case "starting":
		return color.YellowString(status)
	default:
		return color.WhiteString(status)
	}
}

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
