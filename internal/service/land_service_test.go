package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/validator"
)

type mockLandRepo struct {
	saveCalled   bool
	updateCalled bool
	searchResult []domain.LandInfo
	searchErr    error
}

func (m *mockLandRepo) Save(ctx context.Context, req domain.LandInfoRequest) error {
	m.saveCalled = true
	return nil
}

func (m *mockLandRepo) Update(ctx context.Context, req domain.LandInfoRequest) error {
	m.updateCalled = true
	return nil
}

func (m *mockLandRepo) GetLandInfoData(ctx context.Context, criteria domain.LandSearchCriteria) ([]domain.LandInfo, error) {
	if m.searchErr != nil {
		return nil, m.searchErr
	}
	if m.searchResult != nil {
		return m.searchResult, nil
	}
	return nil, errors.New("not implemented")
}

type stubMDMS struct {
	data map[string]any
}

func (s stubMDMS) MDMSCall(context.Context, *domain.RequestInfo, string) (map[string]any, error) {
	return s.data, nil
}

func TestCreateRejectsStateLevelTenant(t *testing.T) {
	svc := NewLandService(
		validator.NewLandValidator(validator.NewLandMDMSValidator()),
		nil, nil, &mockLandRepo{}, stubMDMS{data: validMDMS()},
	)
	_, err := svc.Create(context.Background(), &domain.LandInfoRequest{
		RequestInfo: &domain.RequestInfo{UserInfo: &domain.ContractUser{UUID: "u1"}},
		LandInfo:    &domain.LandInfo{TenantID: "pb"},
	})
	if err == nil {
		t.Fatal("expected state-level tenant error")
	}
	if ce, ok := err.(*CustomException); !ok || ce.Code != InvalidTenant {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateRequiresID(t *testing.T) {
	svc := NewLandService(
		validator.NewLandValidator(validator.NewLandMDMSValidator()),
		nil, nil, &mockLandRepo{}, stubMDMS{data: validMDMS()},
	)
	_, err := svc.Update(context.Background(), &domain.LandInfoRequest{
		RequestInfo: &domain.RequestInfo{UserInfo: &domain.ContractUser{UUID: "u1"}},
		LandInfo:    &domain.LandInfo{TenantID: "pb.amritsar"},
	})
	if err == nil {
		t.Fatal("expected update id error")
	}
	if ce, ok := err.(*CustomException); !ok || ce.Code != UpdateError {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validMDMS() map[string]any {
	return map[string]any{
		"MdmsRes": map[string]any{
			"common-masters": map[string]any{
				"OwnerShipCategory": []any{map[string]any{"code": "INDIVIDUAL"}},
			},
		},
	}
}
