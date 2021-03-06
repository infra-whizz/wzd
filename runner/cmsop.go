package wzd_runner

import (
	nanocms_runners "github.com/infra-whizz/wzcmslib/nanorunners"
	nanocms_results "github.com/infra-whizz/wzcmslib/nanorunners/results"
	nanocms_state "github.com/infra-whizz/wzcmslib/nanostate"
	wzlib_utils "github.com/infra-whizz/wzlib/utils"
)

type WzCMS struct {
	//stateIndex *nanocms_state.NanoStateIndex
	compiler    *nanocms_state.StateCompiler
	pyexe       string
	modulesRoot string
}

// NewWzCMS creates a CMS runner
func NewWzCMS(path ...string) *WzCMS {
	cms := new(WzCMS)
	cms.compiler = nanocms_state.NewStateCompiler().Index(path...)

	return cms
}

// SetPyInterpreter shebang
func (cms *WzCMS) SetPyInterpreter(pyexe string) *WzCMS {
	cms.pyexe = pyexe
	return cms
}

//SetChrootedModules tells CMS to run all modules in own chroot rather then in "/"
func (cmd *WzCMS) SetChrootedModules(root string) *WzCMS {
	if root == "" {
		root = "/"
	}
	cmd.modulesRoot = root
	return cmd
}

// Call a loaded and compiled state
func (cmd *WzCMS) localCall(meta *nanocms_state.NanoStateMeta) (int, []*nanocms_results.ResultLogEntry, error) {
	retcode, err := cmd.compiler.Compile(meta.Path)
	if err != nil {
		return retcode, nil, err
	}

	localRunner := nanocms_runners.NewLocalRunner().SetPyInterpreter(cmd.pyexe).SetChrootPath(cmd.modulesRoot).
		AddStateRoots(cmd.compiler.GetStateIndex().GetStateRoots()...)
	localRunner.Run(cmd.compiler.GetState())
	results := nanocms_results.NewResultsToLog().LoadResults(localRunner.Response()).ToLog()

	return localRunner.Errcode(), results, nil

}

// OfflineCallById state by the Id from the completely downloaded state tree.
// If some reference files aren't there, this call supposed to fail.
func (cms *WzCMS) OfflineCallById(stateId string) (int, []*nanocms_results.ResultLogEntry, error) {
	meta, err := cms.compiler.GetStateIndex().GetStateById(stateId)
	if err != nil {
		return wzlib_utils.EX_UNAVAILABLE, nil, err
	}

	return cms.localCall(meta)
}
