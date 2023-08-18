# 微信支付 Go SDK


微信支付 Go SDK 是一个便捷易用的软件包，基于[微信官方sdk](https://github.com/wechatpay-apiv3/wechatpay-go)深度封装，用于处理支付和通知。该SDK简化了与微信支付API的工作流程，使开发人员能够快速将支付功能集成到他们的应用程序中。

## 特性

- 与微信支付API的无缝集成，用于支付和通知。
- 提供创建订单、处理支付、查询订单、关闭订单和下载账单等功能的易于使用的函数。
- 内置支持解析和验证通知回调。
- 使用`http`软件包高效处理HTTP请求和响应。
- 支持可自定义选项，以更好地控制集成。

## 安装

要安装微信支付 Go SDK，只需在您的项目目录中运行以下命令：

```sh
go get github.com/qq413434162/wechatpay-go-sdk
```
## 使用
### 并在您的Go代码中导入SDK软件包：
```go
import "github.com/qq413434162/wechatpay-go-sdk"
```
### 创建一个新的微信支付客户端实例：
```go
client, err := wechatpay.NewClient(
    context.Background(),
    "your-app-id",
    "your-mch-id",
    "your-mch-api-v3-key",
    "your-mch-certificate-serial-number",
    "path-to-mch-private-key.pem",
)
if err != nil {
    // 处理错误
}
```

## 更多
有关更多示例和详细文档，请参阅微信的[API文档](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/index.shtml)和[SDK的GitHub](https://github.com/wechatpay-apiv3/wechatpay-go)。

## 常见问题
常见问题请见 [FAQ.md](https://github.com/wechatpay-apiv3/wechatpay-go/blob/main/FAQ.md)。

## 贡献
欢迎贡献、提交问题和功能请求！请随时在[GitHub](https://github.com/qq413434162/wechatpay-go-sdk)上打开拉取请求或提交问题。

## 许可证
此SDK采用Apache-2.0 license许可证发布。有关详情，请参阅微信官方SDK的LICENSE。

**免责声明：** 此SDK未经微信支付官方维护。请自行承担风险。