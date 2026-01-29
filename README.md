# Chain Analytics - 链上数据查询工具

一个简单的CLI工具，用于查询区块链上的数据（地址余额、交易详情、区块信息等）。

## 功能

- ✅ 查询地址余额
- ✅ 查询交易详情
- ✅ 查询区块信息
- ✅ 支持多链（ETH、BTC等）

## 安装

```bash
go install github.com/iBreez3/chian-analytics@latest
```

## 使用

### 查询ETH地址余额

```bash
chian-analytics balance eth 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
```

### 查询BTC地址余额

```bash
chian-analytics balance btc bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh
```

### 查询ETH交易详情

```bash
chian-analytics tx eth 0x5c504ed432cb51138bcf09aa5e8a410dd4a1e204ef84bfed1be16dfba1b22060
```

### 查询BTC交易详情

```bash
chian-analytics tx btc f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16
```

### 查询区块信息

```bash
chian-analytics block eth 12345678
chian-analytics block btc 700000
```

## 开发

```bash
# 克隆仓库
git clone https://github.com/iBreez3/chian-analytics.git
cd chian-analytics

# 运行
go run main.go

# 构建
go build -o chian-analytics
```

## License

MIT
