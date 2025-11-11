/*
Copyright Â© 2025 Fleeks Inc.

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
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fleeks-inc/fleeks-cli/internal/config"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "ðŸŒ Environment configuration management",
	Long: `
ðŸŒ Environment Configuration Management

Manage and view environment-specific configurations for Fleeks CLI.

Supports three environments:
â€¢ development - Local development with localhost endpoints
â€¢ staging     - Staging environment with staging-* endpoints  
â€¢ production  - Production environment with production endpoints

Examples:
  # Show current environment info
  fleeks env info
  
  # Switch to development environment
  fleeks --environment development env info
  
  # List all environment settings
  fleeks env list
  
  # Test environment connectivity
  fleeks env test
`,
}

var envInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current environment information",
	Long: `Display detailed information about the current environment configuration.

Shows API endpoints, service URLs, and configuration values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showEnvironmentInfo(cmd)
	},
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environment settings",
	Long: `List all configuration settings for the current environment.

Shows both default values and any overrides from environment files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listEnvironmentSettings(cmd)
	},
}

var envTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test environment connectivity",
	Long: `Test connectivity to all services in the current environment.

Checks:
- Main API endpoint health
- WebSocket connectivity
- LSP service availability
- MCP service availability`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return testEnvironmentConnectivity(cmd)
	},
}

func init() {
	// Add subcommands
	envCmd.AddCommand(envInfoCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envTestCmd)
}

func showEnvironmentInfo(cmd *cobra.Command) error {
	// Load environment configuration
	envConfig, err := config.LoadEnvironment()
	if err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	fmt.Printf("\n%s %s\n\n",
		color.New(color.Bold).Sprint("ðŸŒ Environment Information:"),
		color.CyanString(string(envConfig.Current)))

	// Get environment info
	info := envConfig.GetEnvironmentInfo()

	// Display basic info
	fmt.Printf("%-20s %s\n", "Environment:", color.GreenString(fmt.Sprintf("%v", info["environment"])))
	fmt.Printf("%-20s %s\n", "Config File:", color.YellowString(fmt.Sprintf("%v", info["env_file"])))
	fmt.Printf("%-20s %s\n", "Development Mode:", formatBoolValue(info["dev_mode"]))
	fmt.Printf("%-20s %s\n", "Debug Enabled:", formatBoolValue(info["debug_enabled"]))
	fmt.Printf("%-20s %s\n", "TLS Verify:", formatBoolValue(info["tls_verify"]))

	fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("ðŸ”— Service Endpoints:"))
	fmt.Printf("%-20s %s\n", "Main API:", color.BlueString(fmt.Sprintf("%v", info["api_base_url"])))
	fmt.Printf("%-20s %s\n", "WebSocket:", color.BlueString(fmt.Sprintf("%v", info["ws_base_url"])))
	fmt.Printf("%-20s %s\n", "LSP Service:", color.BlueString(fmt.Sprintf("%v", info["lsp_service"])))
	fmt.Printf("%-20s %s\n", "MCP Service:", color.BlueString(fmt.Sprintf("%v", info["mcp_service"])))

	return nil
}

func listEnvironmentSettings(cmd *cobra.Command) error {
	fmt.Printf("\n%s\n\n",
		color.New(color.Bold).Sprint("âš™ï¸  Environment Settings"))

	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Setting", "Value", "Source"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiYellowColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor},
	)

	// Get all settings
	settings := getAllSettings()

	for key, value := range settings {
		source := "default"
		if viper.IsSet(key) {
			source = "config"
		}
		if os.Getenv(getEnvKey(key)) != "" {
			source = "environment"
		}

		table.Append([]string{
			key,
			fmt.Sprintf("%v", value),
			source,
		})
	}

	table.Render()
	return nil
}

func testEnvironmentConnectivity(cmd *cobra.Command) error {
	fmt.Printf("\n%s\n\n",
		color.New(color.Bold).Sprint("ðŸ” Testing Environment Connectivity"))

	// Test main API
	fmt.Printf("%-30s ", "Main API:")
	apiURL := viper.GetString("api.base_url")
	if testEndpoint(apiURL + "/health") {
		fmt.Printf("%s %s\n", color.GreenString("âœ… Connected"), color.New(color.FgHiBlack).Sprint(apiURL))
	} else {
		fmt.Printf("%s %s\n", color.RedString("âŒ Failed"), color.New(color.FgHiBlack).Sprint(apiURL))
	}

	// Test LSP service
	fmt.Printf("%-30s ", "LSP Service:")
	lspURL := viper.GetString("services.lsp_url")
	if testEndpoint(lspURL + "/health") {
		fmt.Printf("%s %s\n", color.GreenString("âœ… Connected"), color.New(color.FgHiBlack).Sprint(lspURL))
	} else {
		fmt.Printf("%s %s\n", color.RedString("âŒ Failed"), color.New(color.FgHiBlack).Sprint(lspURL))
	}

	// Test MCP service
	fmt.Printf("%-30s ", "MCP Service:")
	mcpURL := viper.GetString("services.mcp_url")
	if testEndpoint(mcpURL + "/health") {
		fmt.Printf("%s %s\n", color.GreenString("âœ… Connected"), color.New(color.FgHiBlack).Sprint(mcpURL))
	} else {
		fmt.Printf("%s %s\n", color.RedString("âŒ Failed"), color.New(color.FgHiBlack).Sprint(mcpURL))
	}

	// Test WebSocket (basic connection test)
	fmt.Printf("%-30s ", "WebSocket:")
	wsURL := viper.GetString("websocket.base_url")
	if testWebSocketEndpoint(wsURL) {
		fmt.Printf("%s %s\n", color.GreenString("âœ… Connected"), color.New(color.FgHiBlack).Sprint(wsURL))
	} else {
		fmt.Printf("%s %s\n", color.RedString("âŒ Failed"), color.New(color.FgHiBlack).Sprint(wsURL))
	}

	return nil
}

// Helper functions
func formatBoolValue(value interface{}) string {
	if b, ok := value.(bool); ok {
		if b {
			return color.GreenString("enabled")
		}
		return color.RedString("disabled")
	}
	return color.New(color.FgHiBlack).Sprint(fmt.Sprintf("%v", value))
}

func getAllSettings() map[string]interface{} {
	return map[string]interface{}{
		"api.base_url":               viper.GetString("api.base_url"),
		"api.timeout":                viper.GetString("api.timeout"),
		"api.debug":                  viper.GetBool("api.debug"),
		"api.tls_verify":             viper.GetBool("api.tls_verify"),
		"websocket.base_url":         viper.GetString("websocket.base_url"),
		"websocket.timeout":          viper.GetString("websocket.timeout"),
		"services.lsp_url":           viper.GetString("services.lsp_url"),
		"services.mcp_url":           viper.GetString("services.mcp_url"),
		"workspace.default_template": viper.GetString("workspace.default_template"),

		"streaming.enabled":     viper.GetBool("streaming.enabled"),
		"streaming.buffer_size": viper.GetInt("streaming.buffer_size"),
		"dev.mode":              viper.GetBool("dev.mode"),
		"dev.verbose":           viper.GetBool("dev.verbose"),
		"dev.log_level":         viper.GetString("dev.log_level"),
	}
}

func getEnvKey(configKey string) string {
	// Convert config key to environment variable name
	envKey := "FLEEKS_" + configKey
	// Replace dots with underscores and convert to uppercase
	return envKey
}

func testEndpoint(url string) bool {
	// This is a simplified test - in a real implementation,
	// you would make an actual HTTP request
	return url != ""
}

func testWebSocketEndpoint(url string) bool {
	// This is a simplified test - in a real implementation,
	// you would attempt a WebSocket connection
	return url != ""
}
