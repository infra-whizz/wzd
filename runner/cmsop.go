package wzd_runner

import (
	nanocms_runners "github.com/infra-whizz/wzcmslib/nanorunners"
	nanocms_state "github.com/infra-whizz/wzcmslib/nanostate"
	wzlib_utils "github.com/infra-whizz/wzlib/utils"
)

type WzCMS struct {
	//stateIndex *nanocms_state.NanoStateIndex
	compiler *nanocms_state.StateCompiler
}

// NewWzCMS creates a CMS runner
func NewWzCMS(path ...string) *WzCMS {
	cms := new(WzCMS)
	cms.compiler = nanocms_state.NewStateCompiler().Index(path...)

	return cms
}

// Call a loaded and compiled state
func (cms *WzCMS) localCall(meta *nanocms_state.NanoStateMeta) (int, string) {
	retcode, err := cms.compiler.Compile(meta.Path)
	if err != nil {
		return retcode, err.Error()
	}

	localRunner := nanocms_runners.NewLocalRunner()
	localRunner.AddStateRoots(cms.compiler.GetStateIndex().GetStateRoots()...)
	localRunner.Run(cms.compiler.GetState())
	return localRunner.Errcode(), localRunner.Response().PrettyJSON()

}

// OfflineCallById state by the Id from the completely downloaded state tree.
// If some reference files aren't there, this call supposed to fail.
func (cms *WzCMS) OfflineCallById(stateId string) (int, string) {
	meta, err := cms.compiler.GetStateIndex().GetStateById(stateId)
	if err != nil {
		return wzlib_utils.EX_UNAVAILABLE, err.Error()
	}

	return cms.localCall(meta)
}
