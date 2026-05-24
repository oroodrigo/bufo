package cli

import (
	"fmt"

	"github.com/oroodrigo/bufo/internal/config"
	"github.com/oroodrigo/bufo/internal/store"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new project in the proxy",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		socketFile := config.Load().SocketFile

		name := args[0]
		port, _ := cmd.Flags().GetInt("port")

		if port <= 0 || port > 65535 {
			fmt.Println("Erro: --port deve estar entre 1 e 65535")
			return
		}

		if err := ensureDaemon(socketFile); err != nil {
			fmt.Println("Erro ao garantir que o daemon está em execução:", err)
			return
		}

		client := NewApiClient(socketFile)
		err := client.AddRoute(name, store.Route{
			Port: port,
		})
		if err != nil {
			fmt.Println("Erro ao adicionar rota:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().Int("port", 0, "Port number of the project")
}
