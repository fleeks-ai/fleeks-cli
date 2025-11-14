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
	"fmt"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/fleeks-inc/fleeks-cli/internal/client"
	"github.com/fleeks-inc/fleeks-cli/internal/config"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview [project-id]",
	Short: "🌐 Get preview URL for workspace",
	Long: `Display the preview URL to access your workspace application in a browser.

Preview URLs provide instant HTTPS access to applications running in your workspace.
No configuration required - just start your application and access it via the URL.

🌟 Features:
  • Instant HTTPS preview URLs
  • WebSocket support for real-time features
  • No configuration or port forwarding needed
  • Open directly in browser
  • Copy to clipboard

Examples:
  # Get preview URL
  fleeks preview my-app

  # Open preview URL in browser
  fleeks preview my-app --open

  # Copy preview URL to clipboard
  fleeks preview my-app --copy

  # Do both
  fleeks preview my-app --open --copy
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return getPreviewURL(args[0], cmd)
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)

	previewCmd.Flags().BoolP("open", "o", false, "Open preview URL in browser")
	previewCmd.Flags().BoolP("copy", "c", false, "Copy preview URL to clipboard")
}

// PreviewURLResponse contains preview URL information
type PreviewURLResponse struct {
	ProjectID    string `json:"project_id"`
	PreviewURL   string `json:"preview_url"`
	WebSocketURL string `json:"websocket_url"`
	Status       string `json:"status"`
	ContainerID  string `json:"container_id"`
}

func getPreviewURL(projectID string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Run 'fleeks auth login' first")
	}

	// Get flags
	openBrowser, _ := cmd.Flags().GetBool("open")
	copyClipboard, _ := cmd.Flags().GetBool("copy")

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Fetch preview URL
	var preview PreviewURLResponse
	endpoint := fmt.Sprintf("/api/v1/sdk/workspaces/%s/preview-url", projectID)
	if err := apiClient.GET(endpoint, &preview); err != nil {
		color.Red("❌ Failed to get preview URL: %v", err)
		fmt.Println()
		color.Yellow("💡 Make sure workspace '%s' exists:", projectID)
		color.Yellow("   fleeks workspace get %s", projectID)
		return nil
	}

	// Display information
	fmt.Println()
	fmt.Printf("🌐 Preview URL: %s\n", color.CyanString(preview.PreviewURL))
	fmt.Printf("🔌 WebSocket URL: %s\n", color.CyanString(preview.WebSocketURL))
	fmt.Println()
	fmt.Printf("📋 Status: %s\n", getStatusColor(preview.Status))
	fmt.Printf("📦 Container: %s\n", color.BlueString(preview.ContainerID))
	fmt.Println()

	// Tips
	fmt.Println(color.YellowString("💡 Tips:"))
	fmt.Println("   • Start a web server in your workspace")
	fmt.Println("   • Access your app via the preview URL")
	fmt.Println("   • WebSocket URL supports real-time features")
	fmt.Println()

	// Open in browser
	if openBrowser {
		fmt.Printf("🌐 Opening %s in your browser...\n", preview.PreviewURL)
		if err := openURL(preview.PreviewURL); err != nil {
			color.Yellow("⚠️  Could not open browser: %v", err)
			color.Yellow("   Please open the URL manually")
		} else {
			fmt.Println(color.GreenString("✅ Browser opened!"))
		}
		fmt.Println()
	}

	// Copy to clipboard
	if copyClipboard {
		if err := copyToClipboard(preview.PreviewURL); err == nil {
			fmt.Println(color.GreenString("✅ Preview URL copied to clipboard!"))
			fmt.Printf("   %s\n", preview.PreviewURL)
			fmt.Println()
		} else {
			color.Yellow("⚠️  Could not copy to clipboard: %v", err)
			fmt.Println()
		}
	}

	return nil
}

// openURL opens a URL in the default browser
func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux, freebsd, openbsd, netbsd
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

// copyToClipboard copies text to the system clipboard
func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("clip")
	case "darwin":
		cmd = exec.Command("pbcopy")
	default: // linux
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}

	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}

	if err := in.Close(); err != nil {
		return err
	}

	return cmd.Wait()
}
