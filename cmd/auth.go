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
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/fleeks-inc/fleeks-cli/internal/client"
	"github.com/fleeks-inc/fleeks-cli/internal/config"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "ðŸ” Authentication and API key management",
	Long: `
ðŸ” Secure Authentication System

Manage your Fleeks CLI authentication and API keys:

âœ… Secure API Key Management:
   â€¢ Encrypted storage of API keys
   â€¢ Token validation and refresh
   â€¢ Multi-account support
   â€¢ Secure credential handling

âœ… Authentication Methods:
   â€¢ API key authentication
   â€¢ Token-based authentication
   â€¢ Single sign-on (SSO) support
   â€¢ Enterprise authentication

Examples:
  # Login with API key
  fleeks auth login
  
  # Login with specific API key
  fleeks auth login --api-key sk_your_api_key_here
  
  # Check authentication status
  fleeks auth status
  
  # Logout and clear credentials
  fleeks auth logout
  
  # Switch between accounts
  fleeks auth switch
`,
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Fleeks",
	Long: `Login to Fleeks using your API key.

You can obtain your API key from the Fleeks Dashboard:
https://dashboard.fleeks.dev/settings/api-keys

The API key will be securely stored in your local configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return loginUser(cmd)
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Fleeks",
	Long: `Logout from Fleeks and clear stored credentials.

This will remove your API key and other authentication tokens
from the local configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return logoutUser(cmd)
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long: `Show current authentication status and user information.

Displays:
- Authentication status
- Current user information
- API key status
- Available scopes and permissions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showAuthStatus(cmd)
	},
}

var authWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user",
	Long:  `Show information about the currently authenticated user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showCurrentUser(cmd)
	},
}

func init() {
	// Add subcommands
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authWhoamiCmd)

	// Login command flags
	authLoginCmd.Flags().StringP("api-key", "k", "", "API key for authentication")
	authLoginCmd.Flags().StringP("base-url", "u", "", "Custom API base URL")
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
}

// UserInfo represents user information
type UserInfo struct {
	ID           string   `json:"id"`
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Organization string   `json:"organization,omitempty"`
	Plan         string   `json:"plan"`
	Verified     bool     `json:"verified"`
	Scopes       []string `json:"scopes"`
	CreatedAt    string   `json:"created_at"`
	LastLogin    string   `json:"last_login"`
}

func loginUser(cmd *cobra.Command) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get API key from flag or prompt
	apiKey, _ := cmd.Flags().GetString("api-key")
	baseURL, _ := cmd.Flags().GetString("base-url")

	if apiKey == "" {
		// Prompt for API key
		prompt := promptui.Prompt{
			Label: "API Key",
			Validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("API key cannot be empty")
				}
				if !strings.HasPrefix(input, "fleeks_") && !strings.HasPrefix(input, "sk_") {
					return fmt.Errorf("invalid API key format")
				}
				return nil
			},
			Mask: '*',
		}

		apiKey, err = prompt.Run()
		if err != nil {
			return fmt.Errorf("API key input cancelled")
		}
	}

	// Set custom base URL if provided
	if baseURL != "" {
		cfg.API.BaseURL = baseURL
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(apiKey)

	// Validate API key by making a test request
	fmt.Printf("%s Validating API key...\n", color.CyanString("ðŸ”"))

	if err := apiClient.HealthCheck(); err != nil {
		return fmt.Errorf("API key validation failed: %w", err)
	}

	// Get user info to confirm authentication
	var userInfo UserInfo
	if err := apiClient.GET("/api/v1/auth/me", &userInfo); err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Store API key securely
	if err := cfg.SetAPIKey(apiKey); err != nil {
		return fmt.Errorf("failed to store API key: %w", err)
	}

	// Success
	fmt.Printf("\n%s %s\n",
		color.GreenString("âœ… Authentication successful!"),
		color.CyanString("Welcome to Fleeks"))

	fmt.Printf("User:         %s (%s)\n", color.YellowString(userInfo.Name), userInfo.Email)
	fmt.Printf("Organization: %s\n", color.BlueString(userInfo.Organization))
	fmt.Printf("Plan:         %s\n", color.MagentaString(userInfo.Plan))

	if !userInfo.Verified {
		fmt.Printf("\n%s Please verify your email address to access all features.\n",
			color.YellowString("âš ï¸"))
	}

	// Show next steps
	fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("ðŸš€ Next steps:"))
	fmt.Printf("  %s\n", color.CyanString("fleeks workspace create my-project --template python"))
	fmt.Printf("  %s\n", color.CyanString("fleeks agent start --project my-project --task \"Build authentication system\""))

	return nil
}

func logoutUser(cmd *cobra.Command) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		fmt.Printf("%s You are not logged in.\n", color.YellowString("â„¹ï¸"))
		return nil
	}

	// Confirm logout
	prompt := promptui.Prompt{
		Label:     "Are you sure you want to logout",
		IsConfirm: true,
	}

	_, err = prompt.Run()
	if err != nil {
		fmt.Println("Logout cancelled.")
		return nil
	}

	// Clear API key and tokens
	cfg.Auth.APIKey = ""
	cfg.Auth.APIKeyHash = ""
	cfg.Auth.RefreshToken = ""
	cfg.Auth.TokenExpiry = ""

	// Save configuration
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("%s Logged out successfully.\n", color.GreenString("ðŸ‘‹"))
	return nil
}

func showAuthStatus(cmd *cobra.Command) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("\n%s\n\n", color.New(color.Bold).Sprint("ðŸ” Authentication Status"))

	if cfg.GetAPIKey() == "" {
		fmt.Printf("Status:       %s\n", color.RedString("Not authenticated"))
		fmt.Printf("API Key:      %s\n", color.New(color.FgHiBlack).Sprint("Not configured"))
		fmt.Printf("\n%s Run 'fleeks auth login' to authenticate.\n",
			color.YellowString("ðŸ’¡"))
		return nil
	}

	// Create API client and test connection
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Test API connection
	if err := apiClient.HealthCheck(); err != nil {
		fmt.Printf("Status:       %s\n", color.RedString("Authentication failed"))
		fmt.Printf("API Key:      %s\n", color.RedString("Invalid"))
		fmt.Printf("Error:        %s\n", color.RedString(err.Error()))
		fmt.Printf("\n%s Run 'fleeks auth login' to re-authenticate.\n",
			color.YellowString("ðŸ’¡"))
		return nil
	}

	// Get user info
	var userInfo UserInfo
	if err := apiClient.GET("/api/v1/auth/me", &userInfo); err != nil {
		fmt.Printf("Status:       %s\n", color.YellowString("Partial"))
		fmt.Printf("API Key:      %s\n", color.GreenString("Valid"))
		fmt.Printf("User Info:    %s\n", color.RedString("Unavailable"))
		return nil
	}

	// Display full status
	fmt.Printf("Status:       %s\n", color.GreenString("Authenticated"))
	fmt.Printf("API Key:      %s\n", color.GreenString("Valid"))
	fmt.Printf("User:         %s (%s)\n", color.YellowString(userInfo.Name), userInfo.Email)
	fmt.Printf("Organization: %s\n", color.BlueString(userInfo.Organization))
	fmt.Printf("Plan:         %s\n", color.MagentaString(userInfo.Plan))
	fmt.Printf("Verified:     %s\n", getBoolColor(userInfo.Verified))
	fmt.Printf("API URL:      %s\n", color.CyanString(cfg.API.BaseURL))

	// Scopes
	if len(userInfo.Scopes) > 0 {
		fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("ðŸ”‘ Available Scopes:"))
		for _, scope := range userInfo.Scopes {
			fmt.Printf("  â€¢ %s\n", color.GreenString(scope))
		}
	}

	return nil
}

func showCurrentUser(cmd *cobra.Command) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GetAPIKey() == "" {
		return fmt.Errorf("not authenticated. Run 'fleeks auth login' first")
	}

	// Create API client
	apiClient := client.NewAPIClient()
	apiClient.SetAPIKey(cfg.GetAPIKey())

	// Get user info
	var userInfo UserInfo
	if err := apiClient.GET("/api/v1/auth/me", &userInfo); err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Display user information
	fmt.Printf("\n%s\n\n", color.New(color.Bold).Sprint("ðŸ‘¤ Current User"))

	fmt.Printf("%-15s %s\n", "ID:", color.CyanString(userInfo.ID))
	fmt.Printf("%-15s %s\n", "Name:", color.YellowString(userInfo.Name))
	fmt.Printf("%-15s %s\n", "Email:", userInfo.Email)
	if userInfo.Organization != "" {
		fmt.Printf("%-15s %s\n", "Organization:", color.BlueString(userInfo.Organization))
	}
	fmt.Printf("%-15s %s\n", "Plan:", color.MagentaString(userInfo.Plan))
	fmt.Printf("%-15s %s\n", "Verified:", getBoolColor(userInfo.Verified))
	fmt.Printf("%-15s %s\n", "Created:", color.MagentaString(userInfo.CreatedAt))
	fmt.Printf("%-15s %s\n", "Last Login:", color.MagentaString(userInfo.LastLogin))

	return nil
}

func getBoolColor(value bool) string {
	if value {
		return color.GreenString("Yes")
	}
	return color.RedString("No")
}

// Helper function to securely read password from terminal
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(password), nil
}
