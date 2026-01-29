package main

import (
	"fmt"
	"os"

	"github.com/iBreez3/chian-analytics/api"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chian-analytics",
	Short: "链上数据查询工具",
	Long:  `一个简单的CLI工具，用于查询区块链上的数据（地址余额、交易详情、区块信息等）。支持ETH、BTC等多链。`,
}

var balanceCmd = &cobra.Command{
	Use:   "balance <chain> <address>",
	Short: "查询地址余额",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		chain := args[0]
		address := args[1]
		queryBalance(chain, address)
	},
}

var txCmd = &cobra.Command{
	Use:   "tx <chain> <txid>",
	Short: "查询交易详情",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		chain := args[0]
		txid := args[1]
		queryTransaction(chain, txid)
	},
}

var blockCmd = &cobra.Command{
	Use:   "block <chain> <height_or_hash>",
	Short: "查询区块信息",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		chain := args[0]
		heightOrHash := args[1]
		queryBlock(chain, heightOrHash)
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
	rootCmd.AddCommand(txCmd)
	rootCmd.AddCommand(blockCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func queryBalance(chain, address string) {
	fmt.Printf("查询 %s 地址余额: %s\n", chain, address)

	switch chain {
	case "eth", "ethereum":
		client := api.NewEtherscanClient()
		balance, err := client.GetBalance(address)
		if err != nil {
			fmt.Printf("❌ 查询失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 余额: %s\n", balance)

	case "btc", "bitcoin":
		client := api.NewBlockstreamClient()
		balance, err := client.GetBalance(address)
		if err != nil {
			fmt.Printf("❌ 查询失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 余额: %s\n", balance)

	default:
		fmt.Printf("❌ 暂不支持链: %s\n", chain)
		os.Exit(1)
	}
}

func queryTransaction(chain, txid string) {
	fmt.Printf("查询 %s 交易详情: %s\n", chain, txid)

	switch chain {
	case "eth", "ethereum":
		client := api.NewEtherscanClient()
		tx, err := client.GetTransaction(txid)
		if err != nil {
			fmt.Printf("❌ 查询失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 交易详情:\n")
		fmt.Printf("   Hash: %s\n", tx.Hash)
		fmt.Printf("   From: %s\n", tx.From)
		fmt.Printf("   To: %s\n", tx.To)
		fmt.Printf("   Value: %s wei\n", tx.Value)
		fmt.Printf("   Gas: %s\n", tx.Gas)
		fmt.Printf("   Gas Price: %s\n", tx.GasPrice)
		fmt.Printf("   Gas Used: %s\n", tx.GasUsed)
		fmt.Printf("   Confirmations: %s\n", tx.Confirmations)

	case "btc", "bitcoin":
		client := api.NewBlockstreamClient()
		tx, err := client.GetTransaction(txid)
		if err != nil {
			fmt.Printf("❌ 查询失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 交易详情:\n")
		fmt.Printf("   TxID: %s\n", tx.Txid)
		fmt.Printf("   Size: %d bytes\n", tx.Size)
		fmt.Printf("   Fee: %d satoshis\n", tx.Fee)
		fmt.Printf("   Inputs: %d\n", len(tx.Inputs))
		fmt.Printf("   Outputs: %d\n", len(tx.Outputs))
		fmt.Printf("   Confirmed: %v\n", tx.Status.Confirmed)
		if tx.Status.Confirmed {
			fmt.Printf("   Block Height: %d\n", tx.Status.BlockHeight)
		}

	default:
		fmt.Printf("❌ 暂不支持链: %s\n", chain)
		os.Exit(1)
	}
}

func queryBlock(chain, heightOrHash string) {
	fmt.Printf("查询 %s 区块信息: %s\n", chain, heightOrHash)

	switch chain {
	case "eth", "ethereum":
		client := api.NewEtherscanClient()
		block, err := client.GetBlock(heightOrHash)
		if err != nil {
			fmt.Printf("❌ 查询失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 区块详情:\n")
		fmt.Printf("   Block Number: %s\n", block.BlockNumber)
		fmt.Printf("   Miner: %s\n", block.BlockMiner)
		fmt.Printf("   Block Reward: %s\n", block.BlockReward)
		fmt.Printf("   Transaction Count: %s\n", block.TransactionCount)

	case "btc", "bitcoin":
		client := api.NewBlockstreamClient()
		block, err := client.GetBlock(heightOrHash)
		if err != nil {
			fmt.Printf("❌ 查询失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 区块详情:\n")
		fmt.Printf("   Hash: %s\n", block.Id)
		fmt.Printf("   Height: %d\n", block.Height)
		fmt.Printf("   Timestamp: %d\n", block.Timestamp)
		fmt.Printf("   Size: %d bytes\n", block.Size)
		fmt.Printf("   Weight: %d\n", block.Weight)
		fmt.Printf("   Tx Count: %d\n", block.TxCount)

	default:
		fmt.Printf("❌ 暂不支持链: %s\n", chain)
		os.Exit(1)
	}
}
