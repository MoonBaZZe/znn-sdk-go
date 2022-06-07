# Zenon Go SDK implementation

[![Go Report Card](https://goreportcard.com/badge/github.com/MoonBaZZe/znn-sdk-go)](https://goreportcard.com/report/github.com/MoonBaZZe/znn-sdk-go)
[![GoDoc](https://godoc.org/github.com/MoonBaZZe/znn-sdk-go?status.svg)](https://godoc.org/github.com/MoonBaZZe/znn-sdk-go)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](#contributing)
[![GitHub license](https://img.shields.io/github/license/MoonBaZZe/znn-sdk-go)](LICENSE)

It follows the [official Dart SDK](https://github.com/zenon-network/znn_sdk_dart) code structure. Tested with `Go v1.18`. It is `100%` compatible with [go-zenon](https://github.com/zenon-network/go-zenon/commit/c0f931d3cd9844a487ae08ed15f5c52896a9bb44).

The Go SDK features a client that will connect to a Zenon full node via `websockets`. We recommend [setting up](https://github.com/zenon-network/go-zenon/blob/master/README.md) and running a local Zenon full node.

## Usage

Install `Go v1.18` and start the [latest znnd node version](https://github.com/zenon-network/go-zenon/releases/latest) on your local host. By default the RPC calls are enabled, otherwise please check the `config.json`.

Some RPC calls require the `keyFile` to modify the state of the ledger for example moving funds, while others only read the on-chain data such as getting the latest momentum or a particular account-block.

### Commands with keyFile

- Note that there must be a valid `keyFile` with the specified name in `DefaultWalletDir`. The name of the `keyFile` can be specified by the user, otherwise by default it is the `baseAddress`.

```go
// Load your keyFile (wallet file)
z, err := zenon.NewZenon("keyfile-sdk")
if err != nil {
    zenon.CommonLogger.Error("", err)
}

// Connect to local node and decrypt your wallet from defaultKeyFilePath
if err := z.Start("123456", "ws://127.0.0.1:35998", 0); err != nil {
    zenon.CommonLogger.Error("", err)
}

// Issue RPC calls to the node

// Stop the client
if err := z.Stop(); err != nil {
    zenon.CommonLogger.Error("", err)
}
```

### Commands without keyFile

- Calls that can be issued without a `keyFile`. Note that calls that require a `keyFile` cannot be used.

```go
// Initialize zenon client (without keyFile)
z, err := zenon.NewZenon("")
if err != nil {
    zenon.CommonLogger.Error("", err)
}

// Connect to local node
if err := z.Start("", "ws://127.0.0.1:35998", 0); err != nil {
    zenon.CommonLogger.Error("", err)
}

// Issue RPC calls to the node

// Stop the client
if err := z.Stop(); err != nil {
    zenon.CommonLogger.Error("", err)
}
```

## Wallet commands

### Read existing keyFile

```go
keyFile, err := wallet.ReadKeyFile("keyfile-sdk", "keyfilePassword")
if err != nil {
    fmt.Printf("err: %s\n", err)
} else {
    fmt.Printf("baseAddress: %s\n", keyFile.BaseAddress)
    _, kp, err := keyFile.DeriveForIndexPath(0)
    if err != nil {
        fmt.Printf("err: %s\n", err)
    } else {
        fmt.Printf("kp address: %s\n", kp.Address.String())
    }
}
```

### Create a new keyFile

```go
kf, err := wallet.NewKeyFile()
if err != nil {
    zenon.WalletLogger.Error("", err)
} else {
    if err := wallet.WriteKeyFile(kf, "keyfile-sdk", "123456"); err != nil {
        zenon.WalletLogger.Error("", "wallet", err)
    }
}
```

### Send and receive commands

- Note that you must fill in the `ZnnTokenStandard` accordingly for `ZNN`, `QSR` or `ZTS`

```go
if err := z.Send(z.Client.LedgerApi.SendTemplate(toAddress, types.ZnnTokenStandard, amount, []byte{})); err != nil {
      fmt.Println(err)
}
if err := z.Send(z.Client.LedgerApi.ReceiveTemplate(types.HexToHashPanic("HASH"))); err != nil {
    fmt.Println(err)
}
```

### Pillar commands

#### Register Pillar

```go
if err := z.Send(z.Client.PillarApi.Register("keyfile-sdk", z.Address(), z.Address(), 0, 50)); err != nil {
    fmt.Println(err)
}
```

#### Update Pillar

```go
if err := z.Send(z.Client.PillarApi.UpdatePillar("keyfile-sdk", types.ParseAddressPanic("z1qqmqp78duzxhpvg7dwxph7724mqu2t3mru297p"), z.Address(), 10, 10)); err != nil {
    fmt.Println(err)
}
```

### Delegate commands

#### Delegate to Pillar

```go
z.Send(z.Client.PillarApi.Delegate("Pillar"))
```

#### Undelegate

```go
z.Send(z.Client.PillarApi.Undelegate())
```

### Sentinel

#### Deposit QSR for Sentinel slot setup

```go
amount := big.NewInt(50000 * constants.Decimals)
if err := z.Send(z.Client.SentinelApi.DepositQsr(amount)); err != nil {
    fmt.Println(err)
}
```

#### Withdraw QSR from Sentinel slot setup

```go
if err := z.Send(z.Client.SentinelApi.WithdrawQsr()); err != nil {
    fmt.Println(err)
}
```

#### Register Sentinel

```go
if err := z.Send(z.Client.SentinelApi.Register()); err != nil {
    fmt.Println(err)
}
```

#### Revoke Sentinel

```go
if err := z.Send(z.Client.SentinelApi.Revoke()); err != nil {
    fmt.Println(err)
}
```

### Plasma commands

#### Fuse Plasma for address

```go
if err := z.Send(z.Client.PlasmaApi.Fuse(toAddress, amount)); err != nil {
    fmt.Println(err)
}
```

### Staking commands

#### Stake for 6 months

```go
amount := big.NewInt(15000 * constants.Decimals)
if err := z.Send(z.Client.StakeApi.Stake(6*30*24*60*60, amount)); err != nil {
    fmt.Println(err)
}
```

### Token commands

#### Issue token

```go
totalSupply := big.NewInt(15000 * constants.Decimals)
maxSupply := big.NewInt(20000 * constants.Decimals)
if err := z.Send(z.Client.TokenApi.IssueToken("keyfile-sdk", "SDK test", "sdk-test.com", totalSupply, maxSupply, 8, true, true, true)); err != nil {
    fmt.Println(err)
}
```

#### Mint token

```go
tokenZts, _ := types.ParseZTS("ZTS")
if err := z.Send(z.Client.TokenApi.Mint(tokenZts, amount, z.Address())); err != nil {
    fmt.Println(err)
}
```

#### Burn token

```go
amount := big.NewInt(1250 * constants.Decimals)
tokenZts, _ := types.ParseZTS("ZTS")
if err := z.Send(z.Client.TokenApi.Burn(tokenZts, amount)); err != nil {
    fmt.Println(err)
}
```

#### Update token

```go
tokenZts, _ := types.ParseZTS("ZTS")
if err := z.Send(z.Client.TokenApi.UpdateToken(tokenZts, z.Address(), false, false)); err != nil {
    fmt.Println(err)
}
```

### Accelerator commands

#### Create project

```go
amountZnn := big.NewInt(5000 * constants.Decimals)
amountQsr := big.NewInt(50000 * constants.Decimals)
if err := z.Send(z.Client.AcceleratorApi.CreateProject("sdk-project-test", "sdk test description", "github.com/sdk/test", amountZnn, amountQsr)); err != nil {
    fmt.Println(err)
}
```
