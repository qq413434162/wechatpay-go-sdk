package wechatpay

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	AppId                      = "x"
	MchID                      = "x"
	MchAPIv3Key                = "x"
	MchCertificateSerialNumber = "x"
	MchPrivateKeyPath          = "x/apiclient_key.pem"
)

func TestWechatPay_Bill(t *testing.T) {
	Convey("Test WechatPay Bill", t, func() {
		u, err := NewWechatPay(
			context.TODO(),
			AppId,
			MchID,
			MchAPIv3Key,
			MchCertificateSerialNumber,
			MchPrivateKeyPath,
		)
		So(
			err,
			ShouldBeNil,
		)

		ok, err := u.Bill(
			context.TODO(),
			"2023-08-14",
			"test.xlsx",
		)
		So(
			err,
			ShouldBeNil,
		)
		So(
			ok,
			ShouldBeTrue,
		)
	})
}
