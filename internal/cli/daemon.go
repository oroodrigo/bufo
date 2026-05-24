package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/oroodrigo/bufo/internal/config"
	"github.com/oroodrigo/bufo/internal/daemon"
	"github.com/spf13/cobra"
)

var daemonCommand = &cobra.Command{
	Use:   "daemon [action]",
	Short: "Manage the daemon",
}

var daemonStartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		exe, err := os.Executable()
		if err != nil {
			fmt.Println("Erro ao obter executável:", err)
			return
		}

		proc := exec.Command(exe, "daemon", "serve")
		if err := proc.Start(); err != nil {
			fmt.Println("Erro ao iniciar daemon:", err)
			return
		}

		fmt.Println("Bufo daemon iniciado")
	},
}

var daemonServeCommand = &cobra.Command{
	Use:   "serve",
	Short: "Serve the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Start()
	},
}

var daemonStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Load()

		if err := daemon.Stop(c.PIDFile, c.SocketFile); err != nil {
			fmt.Println("Erro ao parar Bufo daemon:", err)
			return
		}
		fmt.Println("Bufo daemon parado com sucesso")
	},
}

var daemonRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "Restart the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Load()

		if err := daemon.Stop(c.PIDFile, c.SocketFile); err != nil {
			fmt.Println("Erro ao parar Bufo daemon:", err)
			return
		}

		exe, err := os.Executable()
		if err != nil {
			fmt.Println("Erro ao obter executável:", err)
			return
		}

		proc := exec.Command(exe, "daemon", "serve")
		if err := proc.Start(); err != nil {
			fmt.Println("Erro ao iniciar Bufo daemon:", err)
			return
		}

		fmt.Println("Bufo daemon reiniciado com sucesso")
	},
}

var daemonStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Status of the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Load()
		running, pid, err := daemon.Status(c.PIDFile, c.SocketFile)
		if err != nil {
			fmt.Println("Erro ao verificar status do Bufo daemon:", err)
			return
		}

		if running {
			fmt.Printf("Bufo daemon rodando (PID: %d)\n", pid)
		} else {
			fmt.Println("Bufo daemon não está rodando")
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCommand)

	daemonCommand.AddCommand(daemonStartCommand)
	daemonCommand.AddCommand(daemonServeCommand)
	daemonCommand.AddCommand(daemonStopCommand)
	daemonCommand.AddCommand(daemonRestartCommand)
	daemonCommand.AddCommand(daemonStatusCommand)

	daemonServeCommand.Hidden = true
}
