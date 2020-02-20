package db_errors

type StatusReason string

const (
	StatusReasonNotFound                   StatusReason = "NotFound"
	StatusReasonUnknown                    StatusReason = "Unknown"
	StatusReasonCodebaseIsUsedByCDPipeline StatusReason = "Used"
	StatusRemoveStageRestriction           StatusReason = "RemoveStageRestriction"
	StatusCDStageIsNotTheLast              StatusReason = "StatusCDStageIsNotTheLast"
	StatusRemoveCDPipelineRestriction      StatusReason = "RemoveCDPipelineRestriction"
)

func IsNotFound(err error) bool {
	return reasonForError(err) == StatusReasonNotFound
}

func CodebaseIsUsed(err error) bool {
	return reasonForError(err) == StatusReasonCodebaseIsUsedByCDPipeline
}

func StageErrorOccurred(err error) bool {
	r := reasonForError(err)
	return r == StatusRemoveStageRestriction || r == StatusCDStageIsNotTheLast
}

func CDPipelineErrorOccurred(err error) bool {
	return reasonForError(err) == StatusRemoveCDPipelineRestriction
}

func reasonForError(err error) StatusReason {
	switch t := err.(type) {
	case StatusError:
		return t.Status
	case CodebaseIsUsedByCDPipeline:
		return t.Status
	case RemoveStageRestriction:
		return t.Status
	case RemoveCDPipelineRestriction:
		return t.Status
	}
	return StatusReasonUnknown
}
