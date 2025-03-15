package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "pile",
	Short: "Docker Compose Wrangler",
}

type PileConfig struct {
	APPS []string `yaml:"APPS"`
	DBS  []string `yaml:"DBS"`
}

func readPileConfig() (*PileConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	pilePath := filepath.Join(homeDir, "pile", "pile.config.yaml")
	file, err := os.Open(pilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config PileConfig

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func constructFileFlags(items []string) string {
	var files []string
	for _, item := range items {
		if strings.ContainsAny(item, "&|;><$") {
			continue // Prevent command injection
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err.Error()
		}
		pileDir := filepath.Join(homeDir, "pile")
		files = append(files, fmt.Sprintf("-f %s/%s/compose.yaml", pileDir, item))
	}
	return strings.Join(files, " ")
}

// APPS:
//   - open-webui

func writePileNetworkConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	pileNetworkPath := filepath.Join(homeDir, "pile", "pile.network.yaml")

	content := `networks:
  pile:
    name: pile  # Explicitly name the network
    driver: bridge # Use bridge network (default)
  `

	if err := os.MkdirAll(filepath.Dir(pileNetworkPath), 0755); err != nil {
		return fmt.Errorf("error creating pile directory: %w", err)
	}

	if err := os.WriteFile(pileNetworkPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing pile.network.yaml: %w", err)
	}

	fmt.Println("✅ Successfully wrote pile.network.yaml to", pileNetworkPath)
	return nil
}

func writePileGroupsConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	pileNetworkPath := filepath.Join(homeDir, "pile", "pile.config.yaml")

	content := `APPS:
  - open-webui
`
	if err := os.MkdirAll(filepath.Dir(pileNetworkPath), 0755); err != nil {
		return fmt.Errorf("error creating pile directory: %w", err)
	}

	if err := os.WriteFile(pileNetworkPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing pile.config.yaml: %w", err)
	}

	fmt.Println("✅ Successfully wrote pile.config.yaml to", pileNetworkPath)
	return nil
}

func runCommand(cmdString string) error {
	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pileUp(cmd *cobra.Command, args []string) {
	config, _ := readPileConfig()
	fmt.Println("Config APPS:", config.APPS)
	files := constructFileFlags(config.APPS)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	pileNetwork := filepath.Join(homeDir, "pile", "pile.network.yaml")

	cmdString := fmt.Sprintf("docker compose -f %s %s up -d --remove-orphans", pileNetwork, files)
	_ = runCommand(cmdString)
}

func pileDown(cmd *cobra.Command, args []string) {
	config, _ := readPileConfig()
	files := constructFileFlags(config.APPS)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	pileNetwork := filepath.Join(homeDir, "pile", "pile.network.yaml")

	cmdString := fmt.Sprintf("docker compose -f %s %s down", pileNetwork, files)
	_ = runCommand(cmdString)
}

func initCmd(cmd *cobra.Command, args []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	pilePath := filepath.Join(homeDir, "pile")
	if err := os.MkdirAll(pilePath, 0755); err != nil {
		fmt.Println("Error creating pile directory:", err)
	} else {
		fmt.Println("Created directory:", pilePath)
	}
	writePileNetworkConfig()
	writePileGroupsConfig()
}

func logs(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please specify an app name")
		return
	}
	_ = runCommand(fmt.Sprintf("docker logs -f pile-%s-1", args[0]))
}

func status(cmd *cobra.Command, args []string) {
	_ = runCommand("docker ps --no-trunc --format \"table {{.Names}}\t{{.State}}\"")
}

func images(cmd *cobra.Command, args []string) {
	_ = runCommand("docker ps --no-trunc --format \"table {{.Names}}\t{{.Image}}\"")
}
func ports(cmd *cobra.Command, args []string) {
	_ = runCommand("docker ps --no-trunc --format \"table {{.Names}}\t{{.Ports}}\"")
}
func commands(cmd *cobra.Command, args []string) {
	_ = runCommand("docker ps --no-trunc --format \"table {{.Names}}\t{{.Command}}\"")
}

func install(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please specify a pile to clone")
		return
	}
	clonePile := args[0]
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// Define paths
	pileDir := filepath.Join(homeDir, "pile")
	serviceDestDir := filepath.Join(pileDir, clonePile)
	repoURL := "https://github.com/docker-pile/pile-library.git"
	tempDir := filepath.Join(os.TempDir(), "pile-library")

	// Ensure ~/pile directory exists
	if err := os.MkdirAll(pileDir, 0755); err != nil {
		fmt.Println("Error creating pile directory:", err)
		return
	}

	// Clone repository using go-git
	// fmt.Println("Cloning repository...")
	repo, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:           repoURL,
		Depth:         1, // Shallow clone
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName("main"),
	})
	if err != nil {
		fmt.Println("Error cloning repository:", err)
		return
	}

	// Verify repository
	if repo == nil {
		fmt.Println("Error: repository not cloned properly")
		return
	}

	// Define source test directory from cloned repo
	serviceSrcDir := filepath.Join(tempDir, clonePile)

	// Ensure source directory exists
	if _, err := os.Stat(serviceSrcDir); os.IsNotExist(err) {
		fmt.Println("Pile: ", clonePile, "not found in pile-library")
		// fmt.Println("Cleaning up...")
		os.RemoveAll(tempDir)
		return
	}

	// Copy files from cloned repo to ~/pile/test
	// fmt.Println("Copying test directory to", serviceDestDir)
	if err := copyDir(serviceSrcDir, serviceDestDir); err != nil {
		fmt.Println("Error copying files:", err)
		// fmt.Println("Cleaning up...")
		os.RemoveAll(tempDir)
		return
	}

	// Cleanup: Remove temporary cloned repo
	// fmt.Println("Cleaning up...")
	os.RemoveAll(tempDir)

	fmt.Println("✅ Successfully coppied pile-library to ", serviceDestDir)
}

// copyDir recursively copies a directory and its contents
func copyDir(src string, dest string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// copyFile copies a single file from src to dest
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	return err
}

func configEdit(cmd *cobra.Command, args []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	pileConfig := filepath.Join(homeDir, "pile", "pile.config.yaml")
	configCommand := fmt.Sprintf("vi %s", pileConfig)
	_ = runCommand(configCommand)
}

func envEdit(cmd *cobra.Command, args []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	pileConfig := filepath.Join(homeDir, "pile", ".env")
	configCommand := fmt.Sprintf("vi %s", pileConfig)
	_ = runCommand(configCommand)
}

func main() {
	// pile commands
	rootCmd.AddCommand(&cobra.Command{Use: "init", Run: initCmd})
	rootCmd.AddCommand(&cobra.Command{Use: "up", Run: pileUp})
	rootCmd.AddCommand(&cobra.Command{Use: "down", Run: pileDown})
	rootCmd.AddCommand(&cobra.Command{Use: "install", Run: install})
	rootCmd.AddCommand(&cobra.Command{Use: "config", Run: configEdit})
	rootCmd.AddCommand(&cobra.Command{Use: "env", Run: envEdit})
	// docker shortcut commands
	rootCmd.AddCommand(&cobra.Command{Use: "logs", Run: logs})
	rootCmd.AddCommand(&cobra.Command{Use: "status", Run: status})
	rootCmd.AddCommand(&cobra.Command{Use: "state", Run: status})
	rootCmd.AddCommand(&cobra.Command{Use: "ps", Run: status})
	rootCmd.AddCommand(&cobra.Command{Use: "ports", Run: ports})
	rootCmd.AddCommand(&cobra.Command{Use: "images", Run: images})
	rootCmd.AddCommand(&cobra.Command{Use: "commands", Run: commands})
	rootCmd.Execute()
}
