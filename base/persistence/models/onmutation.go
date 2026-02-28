package models

import (
	"github.com/fundathyon/base/core/customctx"
)

type OnMutationFunc func(cc *customctx.CustomContext, table string, operation string, model map[string]interface{}) error
