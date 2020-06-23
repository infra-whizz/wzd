package wzd_runner

import (
	"fmt"
	"path"

	nanocms_compiler "github.com/infra-whizz/wzcmslib/nanostate/compiler"
)

type WzCMS struct {
	root string
}

func NewWzCMS() *WzCMS {
	cmsOp := new(WzCMS)
	return cmsOp
}

func (cms *WzCMS) SetStateRoot(root string) {
	cms.root = root
}

func (cms *WzCMS) Call(state string) {
	cmd := nanocms_compiler.NewNstCompiler()
	if err := cmd.LoadFile(path.Join(cms.root, state+".st")); err != nil {
		fmt.Println("Loaded state failed:", err.Error())
	}
}
