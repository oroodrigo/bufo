package cli

import (
	"fmt"

	"github.com/oroodrigo/bufo/internal/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a project from the proxy",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		socketFile := config.Load().SocketFile

		name := args[0]

		if err := ensureDaemon(socketFile); err != nil {
			fmt.Println("Erro ao garantir que o daemon está em execução:", err)
			return
		}

		client := NewApiClient(socketFile)
		err := client.DeleteRoute(name)
		if err != nil {
			fmt.Println("Erro ao remover rota:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
