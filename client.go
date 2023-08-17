package wechatpay

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/wechatpay-apiv3/wechatpay-go/core/consts"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type Client struct {
	client      *core.Client
	AppId       string `description:"应用ID"`
	MchID       string `description:"商户号"`
	MchAPIv3Key string `description:"apiv3密钥"`
}

type billDownInfo struct {
	DownloadUrl string `json:"download_url" description:"下载地址"`
}

func NewClient(
	ctx context.Context,
	appId,
	mchID,
	mchAPIv3Key string,
	mchCertificateSerialNumber string,
	mchPrivateKeyPath string,
) (
	client *Client,
	err error,
) {
	client = &Client{
		AppId:       appId,
		MchID:       mchID,
		MchAPIv3Key: mchAPIv3Key,
	}
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(mchPrivateKeyPath)
	if err != nil {
		return
	}

	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(
			client.MchID,
			mchCertificateSerialNumber,
			mchPrivateKey,
			client.MchAPIv3Key,
		),
	}
	client.client, err = core.NewClient(ctx, opts...)
	return
}

// ParseNotifyRequest 解释响应请求
func (s *Client) ParseNotifyRequest(
	ctx context.Context,
	request *http.Request,
) (
	notifyReq *notify.Request,
	resource *payments.Transaction,
	err error,
) {
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithClient(
		ctx,
		s.client,
		s.MchID,
		s.MchAPIv3Key,
	)
	if err != nil {
		return
	}

	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(s.MchID)

	// 3. 使用证书访问器初始化 `notify.Handler`
	handler, err := notify.NewRSANotifyHandler(s.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		return
	}

	notifyReq, err = handler.ParseNotifyRequest(ctx, request, &resource)
	return
}

// Prepay 创建订单获取订单预备号
func (s *Client) Prepay(
	ctx context.Context,
	req app.PrepayRequest,
) (
	prepayId string,
	err error,
) {
	svc := app.AppApiService{Client: s.client}
	// 自动填入无需外部再一次填入
	req.Appid = core.String(s.AppId)
	req.Mchid = core.String(s.MchID)
	resp, _, err := svc.Prepay(ctx,
		req,
	)
	if err != nil {
		return
	}
	if resp == nil || resp.PrepayId == nil {
		return
	}
	prepayId = *resp.PrepayId
	return
}

// Sign 加密
func (s *Client) Sign(
	ctx context.Context,
	timestamp string,
	nonceStr string,
	prepayId string,
) (
	signature string,
	err error,
) {
	var signStr = s.AppId + "\n" +
		timestamp + "\n" +
		nonceStr + "\n" +
		prepayId + "\n" +
		""
	sign, err := s.client.Sign(ctx, signStr)
	if err != nil {
		return
	}
	signature = sign.Signature
	return
}

// Bill 加密
func (s *Client) Bill(
	ctx context.Context,
	date string,
	filePath string,
) (
	ok bool,
	err error,
) {
	// 获取账单文件
	result, err := s.client.Get(ctx, consts.WechatPayAPIServer+"/v3/bill/tradebill?bill_date="+date)
	if err != nil {
		return
	}

	bytes, err := io.ReadAll(result.Response.Body)
	if err != nil {
		return
	}

	v := &billDownInfo{}
	err = json.Unmarshal(bytes, &v)
	bytes = nil
	if err != nil {
		return
	}
	if v.DownloadUrl == "" {
		return
	}

	// 下载账单文件
	result, err = s.client.Get(ctx, v.DownloadUrl)
	if result.Response == nil || result.Response.Body == nil {
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	written, err := io.Copy(file, result.Response.Body)
	if err != nil {
		return
	}
	if written > 0 {
		ok = true
	}
	return
}

// CloseOrder 关闭订单
func (s *Client) CloseOrder(
	ctx context.Context,
	outTradeNo string,
) (
	ok bool,
	err error,
) {
	wechatSvc := app.AppApiService{Client: s.client}
	resp, err := wechatSvc.CloseOrder(ctx, app.CloseOrderRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(s.MchID),
	})
	if err != nil {
		return
	}
	if resp.Response.StatusCode == 204 {
		ok = true
	}
	return
}

// QueryOrderByOutTradeNo 查询订单
func (s *Client) QueryOrderByOutTradeNo(
	ctx context.Context,
	outTradeNo string,
) (
	result *payments.Transaction,
	err error,
) {
	wechatSvc := app.AppApiService{Client: s.client}
	result, _, err = wechatSvc.QueryOrderByOutTradeNo(ctx, app.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(s.MchID),
	})
	return
}
