package domain

// ResponseInfo mirrors org.egov.common.contract.response.ResponseInfo.
type ResponseInfo struct {
	APIID     string `json:"apiId,omitempty"`
	Ver       string `json:"ver,omitempty"`
	Ts        *int64 `json:"ts,omitempty"`
	ResMsgID  string `json:"resMsgId,omitempty"`
	MsgID     string `json:"msgId,omitempty"`
	Status    string `json:"status,omitempty"`
}

// NewResponseInfoFromRequest ports org.egov.land.util.ResponseInfoFactory.
func NewResponseInfoFromRequest(requestInfo *RequestInfo, success bool) ResponseInfo {
	apiID := ""
	ver := ""
	var ts *int64
	msgID := ""

	if requestInfo != nil {
		apiID = requestInfo.APIID
		ver = requestInfo.Ver
		ts = requestInfo.Ts
		msgID = requestInfo.MsgID
	}

	status := "failed"
	if success {
		status = "successful"
	}

	// Matches Java hard-coded value in ResponseInfoFactory.
	const resMsgID = "uief87324"

	return ResponseInfo{
		APIID:    apiID,
		Ver:      ver,
		Ts:       ts,
		ResMsgID: resMsgID,
		MsgID:    msgID,
		Status:   status,
	}
}
