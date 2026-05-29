package mdms

import (
	"encoding/json"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
)

type criteriaReq struct {
	RequestInfo  *domain.RequestInfo `json:"RequestInfo"`
	MdmsCriteria criteria          `json:"MdmsCriteria"`
}

type criteria struct {
	TenantID      string         `json:"tenantId"`
	ModuleDetails []moduleDetail `json:"moduleDetails"`
}

type moduleDetail struct {
	ModuleName    string         `json:"moduleName"`
	MasterDetails []masterDetail `json:"masterDetails"`
}

type masterDetail struct {
	Name   string `json:"name"`
	Filter string `json:"filter,omitempty"`
}

// BuildRequest mirrors LandUtil.getMDMSRequest.
func BuildRequest(requestInfo *domain.RequestInfo, tenantID string) criteriaReq {
	filter := landerrors.ActiveMasterFilter
	return criteriaReq{
		RequestInfo: requestInfo,
		MdmsCriteria: criteria{
			TenantID: tenantID,
			ModuleDetails: []moduleDetail{
				{
					ModuleName: landerrors.BPAModule,
					MasterDetails: []masterDetail{
						{Name: "ApplicationType", Filter: filter},
						{Name: "ServiceType", Filter: filter},
						{Name: "DocumentTypeMapping"},
						{Name: "RiskTypeComputation"},
						{Name: "OccupancyType", Filter: filter},
						{Name: "SubOccupancyType", Filter: filter},
						{Name: "Usages", Filter: filter},
						{Name: "CalculationType"},
						{Name: "CheckList"},
					},
				},
				{
					ModuleName: landerrors.CommonMastersModule,
					MasterDetails: []masterDetail{
						{Name: landerrors.OwnershipCategoryKey, Filter: filter},
						{Name: "OwnerType", Filter: filter},
						{Name: "DocumentType", Filter: filter},
					},
				},
			},
		},
	}
}

// ExtractMasterCodes reads common-masters codes from an MDMS search response.
func ExtractMasterCodes(mdmsData map[string]any) (map[string][]string, error) {
	raw, err := json.Marshal(mdmsData)
	if err != nil {
		return nil, err
	}
	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, err
	}

	mdmsRes, ok := root["MdmsRes"].(map[string]any)
	if !ok {
		return nil, landerrors.New(landerrors.InvalidTenantIDMDMSKey, landerrors.InvalidTenantIDMDMSMessage)
	}
	common, ok := mdmsRes[landerrors.CommonMastersModule].(map[string]any)
	if !ok {
		return nil, landerrors.New(landerrors.InvalidTenantIDMDMSKey, landerrors.InvalidTenantIDMDMSMessage)
	}

	out := make(map[string][]string)
	for master, val := range common {
		out[master] = extractCodes(val)
	}
	return out, nil
}

func extractCodes(value any) []string {
	switch v := value.(type) {
	case []any:
		codes := make([]string, 0, len(v))
		for _, item := range v {
			if code, ok := item.(string); ok {
				codes = append(codes, code)
				continue
			}
			if m, ok := item.(map[string]any); ok {
				if code, ok := m["code"].(string); ok {
					codes = append(codes, code)
				}
			}
		}
		return codes
	default:
		return nil
	}
}
