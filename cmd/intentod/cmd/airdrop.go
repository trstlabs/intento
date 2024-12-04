package cmd

// import (
// 	"bufio"
// 	"encoding/csv"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"os"
// 	"strconv"

// 	sdkmath "cosmossdk.io/math"
// 	"github.com/spf13/cobra"
// )

// // // SnapshotEntry represents an entry in the snapshot.
// // type SnapshotEntry struct {
// // 	Address string      `json:"address"`
// // 	Weight  sdkmath.Int `json:"weight"`
// // }

// // // Snapshot represents the overall snapshot as a slice of SnapshotEntry.
// // type Snapshot []SnapshotEntry

// // func ExportSnapshotCmd() *cobra.Command {
// // 	cmd := &cobra.Command{
// // 		Use:   "export-snapshot [user-list.csv] [nft-list1.csv] [nft-list2.csv] ... [output-file] --nft-weight <nft-weight> --user-weight <user-weight>",
// // 		Short: "Export a snapshot from user and NFT lists with specified weights",
// // 		Long: `Export a snapshot combining addresses from user and multiple NFT lists, assigning fixed weights.
// // - User list: A CSV file with user addresses.
// // - NFT lists: One or more CSV files with NFT addresses.
// // - Output file: The resulting snapshot JSON file.

// // Addresses from the user list will receive the weight specified by --user-weight.
// // Addresses from the NFT lists will receive the weight specified by --nft-weight.
// // If an address exists in both, it will receive the maximum of the two weights.

// // Example:
// // intentod export-snapshot user_list.csv nft_list1.csv nft_list2.csv snapshot_output.json --nft-weight 10 --user-weight 5`,
// // 		Args: cobra.MinimumNArgs(3), // At least user list, one NFT list, and the output file
// // 		RunE: func(cmd *cobra.Command, args []string) error {
// // 			// Parse input arguments
// // 			userListFile := args[0]
// // 			outputFile := args[len(args)-1]
// // 			nftListFiles := args[1 : len(args)-1]

// // 			// Parse weights from flags
// // 			nftWeight, err := cmd.Flags().GetFloat64("nft-weight")
// // 			if err != nil {
// // 				return fmt.Errorf("failed to get nft-weight flag: %w", err)
// // 			}

// // 			userWeight, err := cmd.Flags().GetFloat64("user-weight")
// // 			if err != nil {
// // 				return fmt.Errorf("failed to get user-weight flag: %w", err)
// // 			}

// // 			// Initialize a map to store the maximum weight per address
// // 			addressWeights := make(map[string]float64)

// // 			// Helper function to process a CSV file
// // 			processCSV := func(filePath string, weight float64) error {
// // 				file, err := os.Open(filePath)
// // 				if err != nil {
// // 					return fmt.Errorf("failed to open file %s: %w", filePath, err)
// // 				}
// // 				defer file.Close()

// // 				reader := csv.NewReader(file)
// // 				reader.FieldsPerRecord = -1 // Handle variable fields

// // 				// Read and process records
// // 				records, err := reader.ReadAll()
// // 				if err != nil {
// // 					return fmt.Errorf("failed to read CSV file %s: %w", filePath, err)
// // 				}

// // 				// Process each record
// // 				for _, record := range records {
// // 					if len(record) < 1 {
// // 						return fmt.Errorf("invalid record in %s: %v", filePath, record)
// // 					}
// // 					address := record[0]

// // 					// Assign weight and store maximum weight for the address
// // 					if existingWeight, exists := addressWeights[address]; !exists || weight > existingWeight {
// // 						addressWeights[address] = weight
// // 					}
// // 				}

// // 				return nil
// // 			}

// // 			// Process user list with userWeight
// // 			if err := processCSV(userListFile, userWeight); err != nil {
// // 				return fmt.Errorf("error processing user list: %w", err)
// // 			}

// // 			// Process each NFT list with nftWeight
// // 			for _, nftListFile := range nftListFiles {
// // 				if err := processCSV(nftListFile, nftWeight); err != nil {
// // 					return fmt.Errorf("error processing NFT list %s: %w", nftListFile, err)
// // 				}
// // 			}

// // 			// Create the snapshot
// // 			snapshot := Snapshot{
// // 				Accounts: make([]Account, 0, len(addressWeights)),
// // 			}
// // 			for address, weight := range addressWeights {
// // 				snapshot.Accounts = append(snapshot.Accounts, Account{
// // 					Address: address,
// // 					Weight:  weight,
// // 				})
// // 			}

// // 			// Write snapshot to output file
// // 			outputData, err := json.MarshalIndent(snapshot, "", "  ")
// // 			if err != nil {
// // 				return fmt.Errorf("failed to marshal snapshot: %w", err)
// // 			}

// // 			if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
// // 				return fmt.Errorf("failed to write snapshot to %s: %w", outputFile, err)
// // 			}

// // 			fmt.Printf("Snapshot exported successfully to %s\n", outputFile)
// // 			return nil
// // 		},
// // 	}

// // 	cmd.Flags().Float64("nft-weight", 1.0, "Weight assigned to NFT list addresses")
// // 	cmd.Flags().Float64("user-weight", 1.0, "Weight assigned to user list addresses")
// // 	return cmd
// // }

// // // Parse the user list CSV and assign weights
// // func parseUserListCSV(filePath string, weight int64) (map[string]int64, error) {
// // 	file, err := os.Open(filePath)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer file.Close()

// // 	users := make(map[string]int64)
// // 	scanner := bufio.NewScanner(file)
// // 	for scanner.Scan() {
// // 		address := scanner.Text()
// // 		convertedAddress, err := ConvertBech32(address)
// // 		if err != nil {
// // 			return nil, fmt.Errorf("invalid bech32 address %s: %w", convertedAddress, err)
// // 		}
// // 		users[convertedAddress] = weight
// // 	}
// // 	return users, scanner.Err()
// // }

// // // Parse the NFT list CSV and assign weights
// // func parseNFTListCSV(filePath string, weight int64) (map[string]int64, error) {
// // 	file, err := os.Open(filePath)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer file.Close()

// // 	nftOwners := make(map[string]int64)
// // 	reader := csv.NewReader(file)
// // 	for {
// // 		record, err := reader.Read()
// // 		if err == io.EOF {
// // 			break
// // 		}
// // 		if err != nil {
// // 			return nil, err
// // 		}

// // 		if len(record) < 2 {
// // 			continue
// // 		}

// // 		address := record[0]
// // 		convertedAddress, err := ConvertBech32(address)
// // 		if err != nil {
// // 			return nil, fmt.Errorf("invalid bech32 address %s: %w", convertedAddress, err)
// // 		}
// // 		amount, err := strconv.ParseInt(record[1], 10, 64)
// // 		if err != nil {
// // 			return nil, fmt.Errorf("invalid NFT amount for address %s: %w", address, err)
// // 		}
// // 		nftOwners[convertedAddress] = amount * weight
// // 	}
// // 	return nftOwners, nil
// // }

// // // Merge user weights and NFT weights into a Snapshot
// // func mergeWeightsToSnapshot(userWeights, nftWeights map[string]int64) Snapshot {
// // 	finalWeights := make(map[string]int64)
// // 	for addr, weight := range userWeights {
// // 		finalWeights[addr] = weight
// // 	}
// // 	for addr, weight := range nftWeights {
// // 		if existingWeight, ok := finalWeights[addr]; ok {
// // 			finalWeights[addr] = sdkmath.Max(existingWeight, weight)
// // 		} else {
// // 			finalWeights[addr] = weight
// // 		}
// // 	}

// // 	var snapshot Snapshot
// // 	for addr, weight := range finalWeights {
// // 		snapshot = append(snapshot, SnapshotEntry{
// // 			Address: addr,
// // 			Weight:  sdkmath.NewInt(weight),
// // 		})
// // 	}
// // 	return snapshot
// // }
