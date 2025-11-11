/*
Copyright  2025 Fleeks Inc.

Fleeks CLI - Revolutionary AI-Powered Development Platform

The world's first universal AI software engineer CLI with hybrid local-cloud architecture.
Features include dynamic expertise adaptation, persistent project memory, real-time
streaming collaboration, and full container orchestration.

Usage:
  fleeks [command]

Available Commands:
  workspace   Manage hybrid local-cloud workspaces
  agent       Manage your AI software engineer
  container   Container orchestration and management
  files       File operations with smart sync
  terminal    Terminal execution and management
  auth        Authentication and API key management
  config      Configuration management
  completion  Generate completion script
  help        Help about any command
  version     Show version information

Flags:
  -c, --config string   config file (default is $HOME/.fleeksconfig.yaml)
  -h, --help           help for fleeks
  -v, --verbose        verbose output
      --version        version for fleeks

Use "fleeks [command] --help" for more information about a command.

 Get started: fleeks workspace create my-project --template python
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	colorful "github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	environment string
	verbose     bool
)

// Version information (set via ldflags at build time)
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// gradientLine applies a smooth gradient across a line of text
func gradientLine(text string, startR, startG, startB, endR, endG, endB int) string {
	runes := []rune(text)
	length := len(runes)
	if length == 0 {
		return text
	}

	var result string
	for i, r := range runes {
		// Calculate color for this character
		t := float64(i) / float64(length-1)
		if length == 1 {
			t = 0
		}

		red := int(float64(startR) + t*float64(endR-startR))
		green := int(float64(startG) + t*float64(endG-startG))
		blue := int(float64(startB) + t*float64(endB-startB))

		result += colorful.RGB(uint8(red), uint8(green), uint8(blue)).Sprint(string(r))
	}
	return result
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fleeks",
	Short: "🌟 World's First Universal AI Software Engineer CLI",
	Long: fmt.Sprintf(`
%s
%s
%s
%s
%s
%s
%s

🌟 %s

%s
%s

%s
  %s
  %s
  %s
  %s
  %s

%s
  %s
  %s
  %s

%s

%s`,
		gradientLine(" ███████╗ ██╗      ███████╗ ███████╗ ██╗  ██╗ ███████╗", 220, 100, 255, 222, 98, 252),
		gradientLine(" ██╔════╝ ██║      ██╔════╝ ██╔════╝ ██║ ██╔╝ ██╔════╝", 222, 98, 252, 225, 95, 250),
		gradientLine(" █████╗   ██║      █████╗   █████╗   █████╔╝  ███████╗", 225, 95, 250, 230, 92, 247),
		gradientLine(" ██╔══╝   ██║      ██╔══╝   ██╔══╝   ██╔═██╗  ╚════██║", 230, 92, 247, 235, 88, 243),
		gradientLine(" ██║      ███████╗ ███████╗ ███████╗ ██║  ██╗ ███████║", 235, 88, 243, 210, 98, 240),
		gradientLine(" ╚═╝      ╚══════╝ ╚══════╝ ╚══════╝ ╚═╝  ╚═╝ ╚══════╝", 210, 98, 240, 170, 115, 245),
		"",
		color.New(color.Bold, color.FgHiCyan).Sprint("World's First Universal AI Software Engineer CLI"),
		color.WhiteString("Fleeks brings the power of a universal AI software engineer to your terminal."),
		color.HiBlackString("One intelligent agent that adapts to ANY project type - no role selection needed!"),
		color.New(color.Bold).Sprint("Features:"),
		color.New(color.FgHiGreen).Sprint("✅ Universal AI software engineer (adapts to web, mobile, blockchain, games, AI/ML, IoT, etc.)"),
		color.New(color.FgGreen).Sprint("✅ Automatic project type detection from conversation"),
		color.New(color.FgHiYellow).Sprint("✅ Real-time streaming with tool execution visibility"),
		color.New(color.FgYellow).Sprint("✅ Secure workspace isolation with Docker containers"),
		color.New(color.FgHiGreen).Sprint("✅ Built-in terminal, file operations, and Git integration"),
		"",
		color.New(color.Bold).Sprint("Getting Started:"),
		color.New(color.FgHiCyan).Sprint("  1. Authenticate:    fleeks auth login"),
		color.New(color.FgCyan).Sprint("  2. Create project:  fleeks workspace create my-project"),
		color.New(color.FgHiBlue).Sprint("  3. Start building:  fleeks agent start --project my-project --task \"Build authentication system\""),
		"",
		color.HiBlackString("The agent automatically detects what you're building and adapts its expertise!"),
		"",
		color.New(color.FgBlue).Sprint("📚 Learn more: https://docs.fleeks.dev")),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Persistent flags (available to all subcommands)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fleeksconfig.yaml)")
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "", "environment to use (development, staging, production)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Register all subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(workspaceCmd)
	rootCmd.AddCommand(agentCmd)
	rootCmd.AddCommand(containerCmd)
	rootCmd.AddCommand(filesCmd)
	rootCmd.AddCommand(terminalCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🚀 Fleeks CLI\n")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Built:      %s\n", BuildTime)
		fmt.Printf("Platform:   Universal Multi-Agent Development\n")
		fmt.Printf("\n🌟 Revolutionary Features: Multi-agent workflows, Hybrid local-cloud, Real-time streaming\n")
	},
}

// initializeConfig reads in config file and ENV variables if set
func initializeConfig() error {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		// Search config in home directory with name .fleeksconfig (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".fleeksconfig")
		viper.SetConfigType("yaml")
	}

	// Read in environment variables that match
	viper.SetEnvPrefix("FLEEKS")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	return nil
}

// GetEnvironment returns the current environment setting
func GetEnvironment() string {
	if environment != "" {
		return environment
	}
	return viper.GetString("environment")
}

// IsVerbose returns whether verbose mode is enabled
func IsVerbose() bool {
	return verbose
}
