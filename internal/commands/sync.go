package commands

import (
	"github.com/spf13/cobra"
)

var syncAutoremove bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Make system match Aptfile (install then cleanup)",
	Long: `Sync makes the system match the Aptfile: install any missing packages and
repositories, then remove packages that apt-bundle previously installed but are
no longer listed in the Aptfile (state-based cleanup only; no --zap).

This follows the same "desired state" paradigm as Alpine's apk world file,
uv sync, and poetry sync: one command to bring the system in line with the
declared Aptfile.`,
	RunE: runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&syncAutoremove, "autoremove", false, "Run apt autoremove after cleanup")
	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	if err := runInstall(cmd, args); err != nil {
		return err
	}

	// Run cleanup with force so removal actually happens (sync is not dry-run)
	origForce := cleanupForce
	origAutoremove := cleanupAutoremove
	cleanupForce = true
	cleanupAutoremove = syncAutoremove
	defer func() {
		cleanupForce = origForce
		cleanupAutoremove = origAutoremove
	}()
	return runCleanup(cmd, args)
}
