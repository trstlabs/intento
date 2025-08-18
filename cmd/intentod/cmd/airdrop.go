package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cobra"
)

func ExportSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-snapshot [nft-list1.csv] [nft-list2.csv] ... [output-file] --nft-weights <weight1>,<weight2>,...",
		Short: "Export a snapshot from NFT lists with specified weights",
		Long: `Export a snapshot combining addresses from multiple NFT lists, assigning specified weights.
- NFT lists: One or more CSV files with NFT addresses.
- Output file: The resulting snapshot JSON file.

Each NFT list will be assigned a weight from the --nft-weights flag in order.
If an address appears in multiple lists, it will receive the maximum weight from all lists.

Example:
intentod export-snapshot nft_list1.csv nft_list2.csv nft_list3.csv snapshot_output.json --nft-weights 10,20,30`,
		Args: cobra.MinimumNArgs(2), // At least one NFT list and the output file
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse input arguments
			outputFile := args[len(args)-1]
			nftListFiles := args[:len(args)-1]

			// Parse NFT weights flag
			nftWeightsStr, err := cmd.Flags().GetString("nft-weights")
			if err != nil {
				return fmt.Errorf("failed to get nft-weights flag: %w", err)
			}

			// Parse weights from comma-separated string
			nftWeights, err := parseWeights(nftWeightsStr)
			if err != nil {
				return fmt.Errorf("failed to parse nft-weights: %w", err)
			}

			// Ensure we have enough weights for all NFT lists
			if len(nftWeights) < len(nftListFiles) {
				return fmt.Errorf("not enough weights provided: got %d weights for %d NFT lists",
					len(nftWeights), len(nftListFiles))
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
				for i, record := range records {
					if len(record) < 1 {
						return fmt.Errorf("invalid record in %s: %v", filePath, record)
					}
					//address := record[0]
					address, err := ConvertBech32(record[0])
					if err != nil {
						fmt.Printf("Invalid address in snapshot: %s %d\n", record[0], i)
						continue
					}
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
			// Process each NFT list with its corresponding weight
			for i, nftListFile := range nftListFiles {
				weight := sdkmath.NewInt(int64(nftWeights[i]))
				filename := filepath.Base(nftListFile)
				fmt.Printf("Processing NFT list: %s with weight %s\n", filename, weight.String())
				
				// Process the file with the actual weight
				if err := processCSV(nftListFile, weight); err != nil {
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

			// Print final statistics
			fmt.Println("\n=== Airdrop Distribution Summary ===")
			for i, nftListFile := range nftListFiles {
				weight := sdkmath.NewInt(int64(nftWeights[i]))
				fmt.Printf("List %d: %-30s Weight: %s\n", 
					i+1, 
					filepath.Base(nftListFile), 
					weight.String())
			}

			fmt.Printf("Snapshot exported successfully to %s\n", outputFile)
			return nil
		},
	}

	cmd.Flags().String("nft-weights", "1,2,3,4,5,6", "Comma-separated list of weights for NFT lists (must have at least as many weights as NFT lists)")
	return cmd
}

// parseWeights parses a comma-separated string of integers into a slice of ints
func parseWeights(weightsStr string) ([]int, error) {
	if weightsStr == "" {
		return nil, fmt.Errorf("empty weights string")
	}

	parts := strings.Split(weightsStr, ",")
	weights := make([]int, 0, len(parts))

	for i, part := range parts {
		weight, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid weight at position %d: %w", i, err)
		}
		if weight <= 0 {
			return nil, fmt.Errorf("weight at position %d must be positive, got %d", i, weight)
		}
		weights = append(weights, weight)
	}

	return weights, nil
}
