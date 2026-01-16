package v1

import (
	"github.com/dubbikins/envy/v2/tag"
	"github.com/dubbikins/envy/v2/tag/_default"
	"github.com/dubbikins/envy/v2/tag/env"
	"github.com/dubbikins/envy/v2/tag/required"
)


var WalkFn = tag.Chained(_default.WalkFn, env.WalkFn,required.WalkFn)