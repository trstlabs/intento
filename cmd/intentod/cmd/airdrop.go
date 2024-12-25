package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cobra"
)

func ExportSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-snapshot [user-list.csv] [nft-list1.csv] [nft-list2.csv] ... [output-file] --nft-weight <nft-weight> --user-weight <user-weight>",
		Short: "Export a snapshot from user and NFT lists with specified weights",
		Long: `Export a snapshot combining addresses from user and multiple NFT lists, assigning fixed weights.
- User list: A CSV file with user addresses.
- NFT lists: One or more CSV files with NFT addresses.
- Output file: The resulting snapshot JSON file.

Addresses from the user list will receive the weight specified by --user-weight.
Addresses from the NFT lists will receive the weight specified by --nft-weight.
If an address exists in both, it will receive the maximum of the two weights.

Example:
intentod export-snapshot user_list.csv nft_list1.csv nft_list2.csv snapshot_output.json --nft-weight 10 --user-weight 5`,
		Args: cobra.MinimumNArgs(3), // At least user list, one NFT list, and the output file
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse input arguments
			userListFile := args[0]
			outputFile := args[len(args)-1]
			nftListFiles := args[1 : len(args)-1]

			// Parse weights from flags
			nftWeight, err := cmd.Flags().GetInt("nft-weight")
			if err != nil {
				return fmt.Errorf("failed to get nft-weight flag: %w", err)
			}

			// Parse weights from flags
			nftWeight1, err := cmd.Flags().GetInt("nft-weight-1")
			if err != nil {
				return fmt.Errorf("failed to get nft-weight 2 flag: %w", err)
			}

			userWeight, err := cmd.Flags().GetInt("user-weight")
			if err != nil {
				return fmt.Errorf("failed to get user-weight flag: %w", err)
			}

			// Initialize a map to store the maximum weight per address
			addressWeights := make(map[string]sdkmath.Int)

			// Helper function to process a CSV file
			// Helper function to process a CSV file
			processCSV := func(filePath string, weight sdkmath.Int) error {
				file, err := os.Open(filePath)
				if err != nil {
					return fmt.Errorf("failed to open file %s: %w", filePath, err)
				}
				defer file.Close()

				reader := csv.NewReader(file)
				reader.FieldsPerRecord = -1 // Handle variable fields

				// Skip the first line (header)
				_, err = reader.Read() // Read and discard the first record
				if err != nil {
					return fmt.Errorf("failed to read header in CSV file %s: %w", filePath, err)
				}

				// Read and process remaining records
				records, err := reader.ReadAll()
				if err != nil {
					return fmt.Errorf("failed to read CSV file %s: %w", filePath, err)
				}

				// Process each record
				for _, record := range records {
					if len(record) < 1 {
						return fmt.Errorf("invalid record in %s: %v", filePath, record)
					}
					address := record[0]

					if existingWeight, exists := addressWeights[address]; exists {
						// Use the maximum weight
						if weight.GT(existingWeight) {
							addressWeights[address] = weight
						}
					} else {
						// New address, assign the initial weight
						addressWeights[address] = weight
					}
				}

				return nil
			}
			// Convert weights to sdkmath.Int
			nftWeightInt := sdkmath.NewInt(int64(nftWeight))
			nftWeight1Int := sdkmath.NewInt(int64(nftWeight1))
			userWeightInt := sdkmath.NewInt(int64(userWeight))

			// Process user list with userWeight
			if err := processCSV(userListFile, userWeightInt); err != nil {
				return fmt.Errorf("error processing user list: %w", err)
			}

			// Process each NFT list with nftWeight
			for i, nftListFile := range nftListFiles {
				if i == 0 {
					if err := processCSV(nftListFile, nftWeight1Int); err != nil {
						return fmt.Errorf("error processing NFT list 2 %s: %w", nftListFile, err)
					}
				} else if err := processCSV(nftListFile, nftWeightInt); err != nil {
					return fmt.Errorf("error processing NFT list %s: %w", nftListFile, err)
				}
			}

			// Create the snapshot
			snapshot := Snapshot{}
			for address, weight := range addressWeights {
				snapshot = append(snapshot, SnapshotEntry{
					Address: address,
					Weight:  weight,
				})
			}

			// Write snapshot to output file
			outputData, err := json.MarshalIndent(snapshot, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
				return fmt.Errorf("failed to write snapshot to %s: %w", outputFile, err)
			}

			fmt.Printf("Snapshot exported successfully to %s\n", outputFile)
			return nil
		},
	}

	cmd.Flags().Int("nft-weight-1", 1, "Weight assigned to NFT list addresses")
	cmd.Flags().Int("user-weight", 1, "Weight assigned to user list addresses")
	cmd.Flags().Int("nft-weight", 1, "Weight assigned to NFT list addresses")
	return cmd
}
