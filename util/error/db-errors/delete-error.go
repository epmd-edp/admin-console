package db_errors

type CodebaseIsUsedByCDPipeline struct {
	Status   StatusReason
	Message  string
	Codebase string
	Pipeline string
}

func (e CodebaseIsUsedByCDPipeline) Error() string {
	return string(e.Status)
}

type RemoveStageRestriction struct {
	Status  StatusReason
	Message string
}

func (e RemoveStageRestriction) Error() string {
	return string(e.Status)
}

type RemoveCDPipelineRestriction struct {
	Status  StatusReason
	Message string
}

func (e RemoveCDPipelineRestriction) Error() string {
	return string(e.Status)
}
