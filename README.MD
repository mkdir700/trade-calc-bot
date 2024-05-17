# TradeCalcBot

🚧WIP🚧

**TradeCalcBot** 是一个基于 Telegram 的机器人，旨在帮助用户进行仓位管理和风险管理。通过输入本金和亏损比例，TradeCalcBot 可以计算出合理的仓位大小和杠杆，帮助用户降低交易风险。

## 功能特点

- **仓位大小计算**：根据用户输入的本金和风险比例，计算出合理的仓位大小。
- **杠杆计算**：根据用户的风险承受能力，建议适当的杠杆倍数。
- **简便易用**：通过 Telegram 聊天界面进行交互，随时随地获取交易建议。

## 使用方法

1. **启动机器人**：在 Telegram 中找到你的机器人并点击 "Start" 按钮。
2. **输入本金和风险比例**：按照机器人的提示输入你的本金和你能承受的最大亏损比例。
3. **获取建议**：机器人将会计算并返回建议的仓位大小和杠杆倍数。

## 示例

![recording](https://github.com/mkdir700/trade-calc-bot/assets/56359329/1df4fe0a-f3b9-49ba-a820-c3738996ddf1)

## TODO

- [ ] 加入资金费率的计算

- [ ] 通过交易所 API，自动获取本金

- [ ] 通过交易所 API，获取资金费率

## 贡献

欢迎贡献代码和提出问题！请提交 Pull Request 或创建 Issue。

1. **克隆仓库**
    ```bash
    https://github.com/mkdir700/trade-calc-bot.git
    cd trade-calc-bot
    ```

2. **安装依赖**
    确保你已经安装了 Go 语言环境。然后在项目目录下运行：
    ```bash
    go mod tidy
    ```

3. **设置环境变量**
    修改 `.env.development` 文件，并添加你的 Telegram Bot Token：
    ```env
    BOT_TOKEN=your_telegram_bot_token
    ```

4. **运行机器人**
    ```bash
    make run
    ```


