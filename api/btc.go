package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Blockstream API客户端
type BlockstreamClient struct {
	baseURL string
	client  *http.Client
}

// BTC余额响应
type BTCBalanceResponse struct {
	Address string `json:"address"`
	ChainStats struct {
		Funded  int64 `json:"funded"`
		Spent   int64 `json:"spent"`
		TxCount int   `json:"tx_count"`
	} `json:"chain_stats"`
	MempoolStats struct {
		Funded  int64 `json:"funded"`
		Spent   int64 `json:"spent"`
		TxCount int   `json:"tx_count"`
	} `json:"mempool_stats"`
}

// BTC交易响应
type BTCTxResponse struct {
	Txid string `json:"txid"`
	Version int `json:"version"`
	Size int `json:"size"`
	Fee int64 `json:"fee"`
	Inputs []struct {
		Prevout struct {
			Scriptpubkey string `json:"scriptpubkey"`
			ScriptpubkeyType string `json:"scriptpubkey_type"`
			Value int64 `json:"value"`
		} `json:"prevout"`
		ScriptSig string `json:"scriptsig"`
		Witness []string `json:"witness"`
	} `json:"vinput"`
	Outputs []struct {
		Scriptpubkey string `json:"scriptpubkey"`
		ScriptpubkeyType string `json:"scriptpubkey_type"`
		Value int64 `json:"value"`
	} `json:"vout"`
	Status struct {
		Confirmed bool   `json:"confirmed"`
		BlockHeight int   `json:"block_height"`
		BlockHash string `json:"block_hash"`
		BlockTime int64 `json:"block_time"`
	} `json:"status"`
}

// BTC区块响应
type BTCBlockResponse struct {
	Id         string `json:"id"`
	Height     int    `json:"height"`
	Version    int    `json:"version"`
	Timestamp  int64  `json:"timestamp"`
	Nonce      int    `json:"nonce"`
	Bits       int    `json:"bits"`
	MerkleRoot string `json:"merkle_root"`
	TxCount    int    `json:"tx_count"`
	Size       int    `json:"size"`
	Weight     int    `json:"weight"`
}

func NewBlockstreamClient() *BlockstreamClient {
	return &BlockstreamClient{
		baseURL: "https://blockstream.info/api",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBalance 查询BTC地址余额
func (c *BlockstreamClient) GetBalance(address string) (string, error) {
	url := fmt.Sprintf("%s/address/%s", c.baseURL, address)

	resp, err := c.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result BTCBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析失败: %w", err)
	}

	// 计算余额（已确认的 + 内存池的）
	funded := result.ChainStats.Funded + result.MempoolStats.Funded
	spent := result.ChainStats.Spent + result.MempoolStats.Spent
	balance := funded - spent

	// 转换为BTC单位（1 BTC = 100,000,000 satoshis）
	btc := float64(balance) / 1e8
	return fmt.Sprintf("%.8f BTC", btc), nil
}

// GetTransaction 查询BTC交易详情
func (c *BlockstreamClient) GetTransaction(txid string) (*BTCTxResponse, error) {
	url := fmt.Sprintf("%s/tx/%s", c.baseURL, txid)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result BTCTxResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return &result, nil
}

// GetBlock 查询BTC区块信息
func (c *BlockstreamClient) GetBlock(hash string) (*BTCBlockResponse, error) {
	url := fmt.Sprintf("%s/block/%s", c.baseURL, hash)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result BTCBlockResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return &result, nil
}

// GetBlockByHeight 根据高度查询区块
func (c *BlockstreamClient) GetBlockByHeight(height int) (*BTCBlockResponse, error) {
	url := fmt.Sprintf("%s/block-height/%d", c.baseURL, height)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 返回的是区块hash
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	blockHash := string(body)

	// 用hash查询区块详情
	return c.GetBlock(blockHash)
}
