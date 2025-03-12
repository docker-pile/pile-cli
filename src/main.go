package main

import (
	// "bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Docker management CLI",
}

func readEnvVar(varName string) ([]string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("if [ -f 'pile.env' ]; then yq '.%s' pile.env; fi", varName))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Fields(strings.TrimSpace(string(output))), nil
}

func constructFileFlags(prefix, folder string, items []string) string {
	var files []string
	for _, item := range items {
		files = append(files, fmt.Sprintf("-f %s/%s/%s.yaml", folder, item, item))
	}
	return strings.Join(files, " ")
}

func runCommand(cmdString string) error {
	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dbUp(cmd *cobra.Command, args []string) {
	dbFiles, _ := readEnvVar("DBS")
	files := constructFileFlags("db", "db", dbFiles)
	cmdString := fmt.Sprintf("docker compose -f db/db.yaml %s up -d --remove-orphans", files)
	_ = runCommand(cmdString)
}

func dbDown(cmd *cobra.Command, args []string) {
	dbFiles, _ := readEnvVar("DBS")
	files := constructFileFlags("db", "db", dbFiles)
	cmdString := fmt.Sprintf("docker compose -f db/db.yaml %s down", files)
	_ = runCommand(cmdString)
}

func appUp(cmd *cobra.Command, args []string) {
	appFiles, _ := readEnvVar("APPS")
	files := constructFileFlags("app", "app", appFiles)
	cmdString := fmt.Sprintf("docker compose -f app/app.yaml %s up -d --remove-orphans", files)
	_ = runCommand(cmdString)
}

func appDown(cmd *cobra.Command, args []string) {
	appFiles, _ := readEnvVar("APPS")
	files := constructFileFlags("app", "app", appFiles)
	cmdString := fmt.Sprintf("docker compose -f app/app.yaml %s down", files)
	_ = runCommand(cmdString)
}

func status(cmd *cobra.Command, args []string) {
	_ = runCommand("docker ps --no-trunc --format \"table {{.Names}}\t{{.State}}\"")
}

func logs(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please specify an app name")
		return
	}
	_ = runCommand(fmt.Sprintf("docker logs -f app-%s-1", args[0]))
}

func logsDb(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please specify a database name")
		return
	}
	_ = runCommand(fmt.Sprintf("docker logs -f db-%s-1", args[0]))
}

func restartApp(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please specify an app name")
		return
	}
	_ = runCommand(fmt.Sprintf("docker restart app-%s-1", args[0]))
}

func envEdit(cmd *cobra.Command, args []string) {
	_ = runCommand("vi pile.env")
}

func up(cmd *cobra.Command, args []string) {
	dbUp(cmd, args)
	_ = runCommand("sleep 3")
	appUp(cmd, args)
	status(cmd, args)
}

func down(cmd *cobra.Command, args []string) {
	appDown(cmd, args)
	dbDown(cmd, args)
	_ = runCommand("docker ps")
}

func initCmd(cmd *cobra.Command, args []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	pilePath := homeDir + "/pile"
	if err := os.MkdirAll(pilePath, 0755); err != nil {
		fmt.Println("Error creating pile directory:", err)
	} else {
		fmt.Println("Created directory:", pilePath)
	}
}

func main() {
	rootCmd.AddCommand(&cobra.Command{Use: "db-up", Run: dbUp})
	rootCmd.AddCommand(&cobra.Command{Use: "db-down", Run: dbDown})
	rootCmd.AddCommand(&cobra.Command{Use: "app-up", Run: appUp})
	rootCmd.AddCommand(&cobra.Command{Use: "app-down", Run: appDown})
	rootCmd.AddCommand(&cobra.Command{Use: "status", Run: status})
	rootCmd.AddCommand(&cobra.Command{Use: "logs", Run: logs})
	rootCmd.AddCommand(&cobra.Command{Use: "logs-db", Run: logsDb})
	rootCmd.AddCommand(&cobra.Command{Use: "restart-app", Run: restartApp})
	rootCmd.AddCommand(&cobra.Command{Use: "env", Run: envEdit})
	rootCmd.AddCommand(&cobra.Command{Use: "up", Run: up})
	rootCmd.AddCommand(&cobra.Command{Use: "down", Run: down})
	rootCmd.AddCommand(&cobra.Command{Use: "init", Run: initCmd})
	rootCmd.Execute()
}
