package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/previder/previder-go-sdk/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	var networkCmd = &cobra.Command{
		Use:   "network",
		Short: "Virtual machine commands",
	}
	rootCmd.AddCommand(networkCmd)

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "Get a list of networks",
		Args:  cobra.NoArgs,
		RunE:  listNetwork,
	}
	networkCmd.AddCommand(cmdList)

	var cmdGet = &cobra.Command{
		Use:   "get",
		Short: "Get a network",
		Args:  cobra.ExactArgs(1),
		RunE:  getNetwork,
	}
	networkCmd.AddCommand(cmdGet)
}

func listNetwork(_ *cobra.Command, _ []string) error {
	var networks []client.VirtualNetwork
	for i := 0; ; i++ {
		page, network, err := previderClient.VirtualNetwork.Page(i)
		if err != nil {
			return fmt.Errorf("list virtual networks: %w", err)
		}
		networks = append(networks, *network...)

		if i == page.TotalPages {
			break
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "Type", "Group"})
	for _, network := range networks {
		table.Append([]string{
			network.Name,
			network.Id,
			network.Type,
			network.Group,
		})
	}
	table.Render()
	return nil
}

func getNetwork(_ *cobra.Command, args []string) error {
	content, err := previderClient.VirtualNetwork.Get(args[0])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content)
	return nil
}
