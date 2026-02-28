package eventbus

import (
	"github.com/fundathyon/base/core/customctx"
	"github.com/fundathyon/base/core/result"
)

type Consumer interface {
	Consume(ctx *customctx.CustomContext, event Event) result.Result[DeliveryChannel]
}
