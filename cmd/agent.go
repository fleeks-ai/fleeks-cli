/*
Copyright  2025 Fleeks Inc.

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
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/fleeks-inc/fleeks-cli/internal/client"
	"github.com/fleeks-inc/fleeks-cli/internal/config"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: " Manage your AI software engineer",
	Long: `
 Universal AI Software Engineer

World's first CLI with a single AI agent that handles ALL development tasks:

 Dynamic Expertise Adaptation:
    Automatically detects project type from conversation
    Loads relevant skills (web, mobile, blockchain, games, AI/ML, IoT)
    Switches expertise seamlessly during conversation
    Works on multiple project types simultaneously

 Persistent Project Memory:
    Remembers project decisions across sessions
    Context-aware recommendations
    Historical reasoning retrieval

 Real-Time Streaming:
    Live progress monitoring
    Tool usage visibility
    File changes tracking

Examples:
  # Start AI software engineer
  fleeks agent start --project my-api --task "Build user authentication with JWT"
  fleeks agent start --project my-api --task "Create React Native mobile app"
  
  # Agent automatically detects and adapts:
  # - "authentication" + "JWT"  Web development skills loaded
  # - "React Native"  Mobile development skills loaded
  # - Can work on backend + mobile simultaneously!

  # Monitor agent in real-time
  fleeks agent watch agent-123
  fleeks agent list --project my-api

  # Chat with your software engineer
  fleeks chat my-project
`,
}

var agentStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the AI software engineer",
	Long: `Start your universal AI software engineer.

The agent automatically adapts its expertise based on your task:
   "Build React Native app"  Mobile development expertise
   "Create smart contracts"  Blockchain expertise  
   "Implement ML model"  AI/ML expertise
   "Setup CI/CD"  DevOps expertise

No need to specify roles - the agent figures it out!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startAgent(cmd)
	},
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active agent sessions",
	Long: `List all active agent sessions with their status and current tasks.

Shows agents across all projects or filtered by specific project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listAgents(cmd)
	},
}

var agentWatchCmd = &cobra.Command{
	Use:   "watch [agent-id]",
	Short: "Watch agent execution in real-time",
	Long: `Stream agent execution progress in real-time.

Features:
- Live streaming of agent thoughts and actions
- Tool usage visibility
- File changes monitoring
- Progress tracking
- Dynamic expertise switching

Watch as your AI software engineer adapts to different project types!`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return watchAgent(args[0], cmd)
	},
}

var agentStatusCmd = &cobra.Command{
	Use:   "status [agent-id]",
	Short: "Get agent status",
	Long: `Get detailed status information for a specific agent including:
- Current task and progress
- Detected project types
- Active skills loaded
- Tool usage statistics
- Execution timeline
- Resource usage`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return getAgentStatus(args[0], cmd)
	},
}

var agentStopCmd = &cobra.Command{
	Use:   "stop [agent-id]",
	Short: "Stop an agent",
	Long: `Stop a running agent and clean up resources.

The agent's state and context will be preserved for potential restart.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return stopAgent(args[0], cmd)
	},
}

func init() {
	// Add subcommands
	agentCmd.AddCommand(agentStartCmd)
	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentWatchCmd)
	agentCmd.AddCommand(agentStatusCmd)
	agentCmd.AddCommand(agentStopCmd)

	// Start command flags
	agentStartCmd.Flags().StringP("project", "p", "", "Project ID (required)")
	agentStartCmd.Flags().StringP("task", "t", "", "Initial task for the agent")
	agentStartCmd.Flags().IntP("max-iterations", "m", 0, "Maximum iterations (0 = use default)")
	agentStartCmd.Flags().BoolP("detached", "d", false, "Run agent in detached mode")
	agentStartCmd.Flags().StringSliceP("context", "c", []string{}, "Additional context files")

	// List command flags
	agentListCmd.Flags().StringP("project", "p", "", "Filter by project ID")
	agentListCmd.Flags().StringP("status", "s", "", "Filter by status")

	// Watch command flags
	agentWatchCmd.Flags().BoolP("follow", "f", true, "Follow new messages")
	agentWatchCmd.Flags().IntP("tail", "", 50, "Number of recent messages to show")

	// Mark required flags
	agentStartCmd.MarkFlagRequired("project")
}

// AgentStartRequest represents agent start request
type AgentStartRequest struct {
	ProjectID     string            `json:"project_id"`
	Task          string            `json:"task,omitempty"`
	MaxIterations int               `json:"max_iterations,omitempty"`
	Context       map[string]string `json:"context,omitempty"`
}

// AgentResponse represents agent response
type AgentResponse struct {
	AgentID       string    `json:"agent_id"`
	ProjectID     string    `json:"project_id"`
	Status        string    `json:"status"`
	Task          string    `json:"task"`
	Progress      int       `json:"progress"`
	DetectedTypes []string  `json:"detected_types,omitempty"`
	ActiveSkills  []string  `json:"active_skills,omitempty"`
	StartedAt     time.Time `json:"started_at"`
	Message       string    `json:"message"`
}

// AgentStatus represents detailed agent status
type AgentStatus struct {
	AgentID         string     `json:"agent_id"`
	ProjectID       string     `json:"project_id"`
	Status          string     `json:"status"`
	Task            string     `json:"task"`
	Progress        int        `json:"progress"`
	CurrentStep     string     `json:"current_step,omitempty"`
	DetectedTypes   []string   `json:"detected_types,omitempty"`
	ActiveSkills    []string   `json:"active_skills,omitempty"`
	Iterations      int        `json:"iterations_completed"`
	MaxIterations   int        `json:"max_iterations"`
	StartedAt       time.Time  `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	ExecutionTimeMs *float64   `json:"execution_time_ms,omitempty"`
	ToolsUsed       []string   `json:"tools_used,omitempty"`
	FilesModified   []string   `json:"files_modified,omitempty"`
}

func startAgent(cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	// Get flags
	projectID, _ := cmd.Flags().GetString("project")
	task, _ := cmd.Flags().GetString("task")
	maxIterations, _ := cmd.Flags().GetInt("max-iterations")
	detached, _ := cmd.Flags().GetBool("detached")
	contextFiles, _ := cmd.Flags().GetStringSlice("context")

	// If no task provided, prompt for it
	if task == "" {
		prompt := promptui.Prompt{
			Label: "Task for AI engineer",
			Validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("task cannot be empty")
				}
				return nil
			},
		}
		task, err = prompt.Run()
		if err != nil {
			return fmt.Errorf("task input cancelled")
		}
	}

	// Build context from files
	context := make(map[string]string)
	for _, file := range contextFiles {
		if content, err := os.ReadFile(file); err == nil {
			context[file] = string(content)
		}
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Start spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Starting AI software engineer..."
	s.Start()
	defer s.Stop()

	// Prepare request
	request := AgentStartRequest{
		ProjectID:     projectID,
		Task:          task,
		MaxIterations: maxIterations,
		Context:       context,
	}

	// Start agent
	var response AgentResponse
	if err := apiClient.POST("/api/v1/sdk/agents", request, &response); err != nil {
		s.Stop()
		return fmt.Errorf("failed to start agent: %w", err)
	}

	s.Stop()

	// Success output
	fmt.Printf("\n%s %s\n",
		color.GreenString(" AI Software Engineer started!"),
		color.CyanString(response.AgentID))

	fmt.Printf("Project:      %s\n", color.BlueString(response.ProjectID))
	fmt.Printf("Task:         %s\n", color.WhiteString(response.Task))
	fmt.Printf("Status:       %s\n", getStatusColor(response.Status))

	if len(response.DetectedTypes) > 0 {
		fmt.Printf("Detected:     %s\n", color.MagentaString(strings.Join(response.DetectedTypes, ", ")))
	}

	if len(response.ActiveSkills) > 0 {
		fmt.Printf("Skills:       %s\n", color.YellowString(fmt.Sprintf("%d skills loaded", len(response.ActiveSkills))))
	}

	fmt.Printf("Started:      %s\n", color.MagentaString(response.StartedAt.Format("2006-01-02 15:04:05")))

	if !detached {
		fmt.Printf("\n%s Streaming agent execution...\n", color.CyanString(""))
		return watchAgent(response.AgentID, cmd)
	}

	// Show monitoring commands
	fmt.Printf("\n%s\n", color.New(color.Bold).Sprint(" Monitor agent:"))
	fmt.Printf("  %s\n", color.CyanString("fleeks agent watch "+response.AgentID))
	fmt.Printf("  %s\n", color.CyanString("fleeks agent status "+response.AgentID))

	return nil
}

func listAgents(cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	// Get filters
	projectID, _ := cmd.Flags().GetString("project")
	status, _ := cmd.Flags().GetString("status")

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Build query parameters
	endpoint := "/api/v1/sdk/agents"
	params := make([]string, 0)
	if projectID != "" {
		params = append(params, "project_id="+projectID)
	}
	if status != "" {
		params = append(params, "status="+status)
	}
	if len(params) > 0 {
		endpoint += "?" + strings.Join(params, "&")
	}

	// Get agents
	var agents []AgentStatus
	if err := apiClient.GET(endpoint, &agents); err != nil {
		return fmt.Errorf("failed to list agents: %w", err)
	}

	if len(agents) == 0 {
		fmt.Printf("%s No active agents found.\n", color.YellowString(""))
		fmt.Printf("Start one with: %s\n",
			color.CyanString("fleeks agent start --project my-project --task \"Build user auth\""))
		return nil
	}

	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Agent ID", "Project", "Status", "Progress", "Detected Types", "Task"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.FgHiMagentaColor},
		tablewriter.Colors{tablewriter.FgHiYellowColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
	)

	for _, agent := range agents {
		task := agent.Task
		if len(task) > 40 {
			task = task[:37] + "..."
		}

		detectedTypes := "auto"
		if len(agent.DetectedTypes) > 0 {
			detectedTypes = strings.Join(agent.DetectedTypes, ", ")
			if len(detectedTypes) > 20 {
				detectedTypes = detectedTypes[:17] + "..."
			}
		}

		table.Append([]string{
			agent.AgentID[:8] + "...",
			agent.ProjectID,
			agent.Status,
			fmt.Sprintf("%d%%", agent.Progress),
			detectedTypes,
			task,
		})
	}

	fmt.Printf("\n%s %s\n\n",
		color.New(color.Bold).Sprint(" Active AI Software Engineers:"),
		color.GreenString(fmt.Sprintf("(%d total)", len(agents))))

	table.Render()
	return nil
}

func watchAgent(agentID string, cmd *cobra.Command) error {
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

	// Create stream reader
	streamPath := fmt.Sprintf("/ws/agents/%s/stream", agentID)
	stream, err := apiClient.NewStreamReader(streamPath)
	if err != nil {
		return fmt.Errorf("failed to connect to agent stream: %w", err)
	}
	defer stream.Close()

	fmt.Printf("%s Watching AI engineer %s (Press Ctrl+C to exit)\n\n",
		color.CyanString(""), color.YellowString(agentID[:12]))

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n%s Disconnecting from agent stream...\n",
			color.YellowString(""))
		cancel()
	}()

	// Stream messages
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-stream.Messages():
			if !ok {
				fmt.Printf("\n%s Agent session ended\n", color.GreenString(""))
				return nil
			}

			timestamp := msg.Timestamp.Format("15:04:05")
			switch msg.Type {
			case "thought":
				fmt.Printf("[%s] %s %s\n",
					color.MagentaString(timestamp),
					color.CyanString(""),
					msg.Content)
			case "tool_call":
				tool := msg.Metadata["tool"]
				fmt.Printf("[%s] %s Using: %s\n",
					color.MagentaString(timestamp),
					color.YellowString(""),
					color.GreenString(fmt.Sprintf("%v", tool)))
			case "skill_loaded":
				skill := msg.Metadata["skill"]
				projectType := msg.Metadata["project_type"]
				fmt.Printf("[%s] %s [%s] Loaded skill: %s\n",
					color.MagentaString(timestamp),
					color.MagentaString(""),
					color.YellowString(fmt.Sprintf("%v", projectType)),
					color.GreenString(fmt.Sprintf("%v", skill)))
			case "type_detected":
				projectType := msg.Metadata["project_type"]
				fmt.Printf("[%s] %s Detected project type: %s\n",
					color.MagentaString(timestamp),
					color.CyanString(""),
					color.YellowString(fmt.Sprintf("%v", projectType)))
			case "output":
				fmt.Printf("[%s] %s %s\n",
					color.MagentaString(timestamp),
					color.BlueString(""),
					msg.Content)
			case "progress":
				progress := msg.Metadata["progress"]
				fmt.Printf("[%s] %s Progress: %s\n",
					color.MagentaString(timestamp),
					color.GreenString(""),
					color.CyanString(fmt.Sprintf("%v%%", progress)))
			case "complete":
				fmt.Printf("[%s] %s Task completed!\n",
					color.MagentaString(timestamp),
					color.GreenString(""))
				return nil
			case "error":
				fmt.Printf("[%s] %s Error: %s\n",
					color.MagentaString(timestamp),
					color.RedString(""),
					color.RedString(msg.Content))
			}

		case err, ok := <-stream.Errors():
			if !ok {
				return nil
			}
			return fmt.Errorf("stream error: %w", err)
		}
	}
}

func getAgentStatus(agentID string, cmd *cobra.Command) error {
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

	// Get agent status
	var agent AgentStatus
	endpoint := fmt.Sprintf("/api/v1/sdk/agents/%s", agentID)
	if err := apiClient.GET(endpoint, &agent); err != nil {
		return fmt.Errorf("failed to get agent status: %w", err)
	}

	// Display agent status
	fmt.Printf("\n%s %s\n\n",
		color.New(color.Bold).Sprint(" AI Software Engineer Status:"),
		color.CyanString(agentID))

	fmt.Printf("%-20s %s\n", "Agent ID:", color.CyanString(agent.AgentID))
	fmt.Printf("%-20s %s\n", "Project:", color.BlueString(agent.ProjectID))
	fmt.Printf("%-20s %s\n", "Status:", getStatusColor(agent.Status))
	fmt.Printf("%-20s %s\n", "Task:", agent.Task)
	fmt.Printf("%-20s %s\n", "Progress:", color.GreenString(fmt.Sprintf("%d%%", agent.Progress)))

	if len(agent.DetectedTypes) > 0 {
		fmt.Printf("%-20s %s\n", "Detected Types:", color.YellowString(strings.Join(agent.DetectedTypes, ", ")))
	}

	if len(agent.ActiveSkills) > 0 {
		fmt.Printf("%-20s %s\n", "Active Skills:", color.MagentaString(fmt.Sprintf("%d loaded", len(agent.ActiveSkills))))
	}

	if agent.CurrentStep != "" {
		fmt.Printf("%-20s %s\n", "Current Step:", agent.CurrentStep)
	}
	fmt.Printf("%-20s %s\n", "Iterations:", color.MagentaString(fmt.Sprintf("%d/%d", agent.Iterations, agent.MaxIterations)))
	fmt.Printf("%-20s %s\n", "Started:", color.MagentaString(agent.StartedAt.Format("2006-01-02 15:04:05")))

	if agent.CompletedAt != nil {
		fmt.Printf("%-20s %s\n", "Completed:", color.MagentaString(agent.CompletedAt.Format("2006-01-02 15:04:05")))
	}

	if agent.ExecutionTimeMs != nil {
		duration := time.Duration(*agent.ExecutionTimeMs) * time.Millisecond
		fmt.Printf("%-20s %s\n", "Execution Time:", color.MagentaString(duration.String()))
	}

	// Tools and files
	if len(agent.ToolsUsed) > 0 {
		fmt.Printf("\n%s\n", color.New(color.Bold).Sprint(" Tools Used:"))
		for _, tool := range agent.ToolsUsed {
			fmt.Printf("   %s\n", color.GreenString(tool))
		}
	}

	if len(agent.ActiveSkills) > 0 {
		fmt.Printf("\n%s\n", color.New(color.Bold).Sprint(" Active Skills:"))
		for i, skill := range agent.ActiveSkills {
			if i < 10 { // Show first 10
				fmt.Printf("   %s\n", color.YellowString(skill))
			}
		}
		if len(agent.ActiveSkills) > 10 {
			fmt.Printf("  ... and %d more\n", len(agent.ActiveSkills)-10)
		}
	}

	if len(agent.FilesModified) > 0 {
		fmt.Printf("\n%s\n", color.New(color.Bold).Sprint(" Files Modified:"))
		for _, file := range agent.FilesModified {
			fmt.Printf("   %s\n", color.BlueString(file))
		}
	}

	return nil
}

func stopAgent(agentID string, cmd *cobra.Command) error {
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

	// Stop agent
	endpoint := fmt.Sprintf("/api/v1/sdk/agents/%s/stop", agentID)
	if err := apiClient.POST(endpoint, nil, nil); err != nil {
		return fmt.Errorf("failed to stop agent: %w", err)
	}

	fmt.Printf("%s AI Software Engineer %s stopped successfully\n",
		color.GreenString(""), color.CyanString(agentID))

	return nil
}
