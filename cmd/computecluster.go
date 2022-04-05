package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	var computeClusterCmd = &cobra.Command{
		Use:   "computecluster",
		Short: "Compute cluster commands",
	}
	rootCmd.AddCommand(computeClusterCmd)

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "Get a list of compute clusters",
		Args:  cobra.NoArgs,
		RunE:  listComputeClusters,
	}
	computeClusterCmd.AddCommand(cmdList)
}

func listComputeClusters(_ *cobra.Command, _ []string) error {
	clusters, err := previderClient.VirtualMachine.ComputeClusterList()
	if err != nil {
		return fmt.Errorf("list compute clusters: %w", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description"})
	for _, cluster := range *clusters {
		table.Append([]string{
			cluster.Name,
			cluster.Description,
		})
	}
	table.Render()
	return nil
}
