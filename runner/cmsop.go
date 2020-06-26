package wzd_runner

import (
	"fmt"

	wzlib_utils "github.com/infra-whizz/wzlib/utils"

	nanocms_runners "github.com/infra-whizz/wzcmslib/nanorunners"
	nanocms_state "github.com/infra-whizz/wzcmslib/nanostate"
	nanocms_compiler "github.com/infra-whizz/wzcmslib/nanostate/compiler"
)

type WzCMS struct {
	stateIndex *nanocms_state.NanoStateIndex
}

// NewWzCMS creates a CMS runner
func NewWzCMS(path ...string) *WzCMS {
	cms := new(WzCMS)
	cms.stateIndex = nanocms_state.NewNanoStateIndex().AddStateRoots(path...).Index()

	return cms
}

// Call a loaded and compiled state
func (cms *WzCMS) localCall(meta *nanocms_state.NanoStateMeta) (int, string) {
	var err error = nil
	cmp := nanocms_compiler.NewNstCompiler()
	fmt.Println(meta.Path)
	if err = cmp.LoadFile(meta.Path); err != nil {
		return wzlib_utils.EX_GENERIC, err.Error()
	}

	// Load the entire chain of the local caller
	for {
		cMeta, _ := cms.stateIndex.GetStateById(cmp.Cycle())
		if cMeta != nil {
			if err = cmp.LoadFile(cMeta.Path); err != nil {
				return wzlib_utils.EX_GENERIC, err.Error()
			}
		} else {
			break
		}
	}

	state := nanocms_state.NewNanostate()
	if err := state.Load(cmp.Tree()); err != nil {
		return wzlib_utils.EX_GENERIC, err.Error()
	}

	localRunner := nanocms_runners.NewLocalRunner()
	localRunner.AddStateRoots(cms.stateIndex.GetStateRoots()...)
	localRunner.Run(state)
	return localRunner.Errcode(), localRunner.Response().PrettyJSON()
}

// OfflineCallById state by the Id from the completely downloaded state tree.
// If some reference files aren't there, this call supposed to fail.
func (cms *WzCMS) OfflineCallById(stateId string) (int, string) {
	meta, err := cms.stateIndex.GetStateById(stateId)
	if err != nil {
		return wzlib_utils.EX_UNAVAILABLE, err.Error()
	}

	return cms.localCall(meta)
}
