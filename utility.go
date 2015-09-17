package webhook

import (
	"appengine"
)

func L(context appengine.Context, msg string) {
	context.Infof("%v", msg)
}
