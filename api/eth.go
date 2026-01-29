package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Etherscan API客户端
type EtherscanClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// ETH余额响应
type ETHBalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// ETH交易响应
type ETHTxResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  *ETHTx `json:"result"`
}

type ETHTx struct {
	BlockNumber      string `json:"blockNumber"`
	TimeStamp        string `json:"timestamp"`
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	GasUsed          string `json:"gasUsed"`
	Input            string `json:"input"`
	ContractAddress  string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Confirmations    string `json:"confirmations"`
}

// ETH区块响应
type ETHBlockResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Result  []ETHBlock `json:"result"`
}

type ETHBlock struct {
	BlockNumber      string `json:"blockNumber"`
	TimeStamp        string `json:"timeStamp"`
	BlockMiner       string `json:"blockMiner"`
	BlockReward      string `json:"blockReward"`
	TransactionCount string `json:"transactionCount"`
}

func NewEtherscanClient() *EtherscanClient {
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	if apiKey == "" {
		apiKey = "YourApiKeyToken" // 免费版使用默认key
	}

	return &EtherscanClient{
		apiKey:  apiKey,
		baseURL: "https://api.etherscan.io/api",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBalance 查询ETH地址余额
func (c *EtherscanClient) GetBalance(address string) (string, error) {
	url := fmt.Sprintf("%s?module=account&action=balance&address=%s&tag=latest&apikey=%s",
		c.baseURL, address, c.apiKey)

	resp, err := c.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result ETHBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析失败: %w", err)
	}

	if result.Status != "1" {
		return "", fmt.Errorf("API错误: %s", result.Message)
	}

	// 转换为ETH单位
	wei := result.Result
	if wei != "0" {
		weiInt, _ := strconv.ParseInt(wei, 10, 64)
		eth := float64(weiInt) / 1e18
		return fmt.Sprintf("%.6f ETH", eth), nil
	}

	return "0 ETH", nil
}

// GetTransaction 查询ETH交易详情
func (c *EtherscanClient) GetTransaction(txid string) (*ETHTx, error) {
	url := fmt.Sprintf("%s?module=proxy&action=eth_getTransactionByHash&txhash=%s&apikey=%s",
		c.baseURL, txid, c.apiKey)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result ETHTxResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	if result.Result == nil {
		return nil, fmt.Errorf("交易未找到")
	}

	return result.Result, nil
}

// GetBlock 查询ETH区块信息
func (c *EtherscanClient) GetBlock(height string) (*ETHBlock, error) {
	url := fmt.Sprintf("%s?module=block&action=getblockreward&blockno=%s&apikey=%s",
		c.baseURL, height, c.apiKey)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result ETHBlockResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	if result.Status != "1" || len(result.Result) == 0 {
		return nil, fmt.Errorf("区块未找到")
	}

	return &result.Result[0], nil
}
