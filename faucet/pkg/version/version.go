package version

import "fmt"

var (
	gittag  = "unk"
	githash = "unk"
)

func Version() string {
	return fmt.Sprintf("%s-%s", gittag, githash)
}
