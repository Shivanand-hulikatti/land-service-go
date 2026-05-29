package service

import (
	"context"
	"testing"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/validator"
)

func TestSearchReturnsEmptySliceWhenNoRows(t *testing.T) {
	repo := &mockLandRepo{searchResult: []domain.LandInfo{}}
	svc := NewLandService(
		validator.NewLandValidator(validator.NewLandMDMSValidator()),
		nil, nil, repo, stubMDMS{data: validMDMS()},
	)

	lands, err := svc.Search(context.Background(), domain.LandSearchCriteria{
		TenantID: "pb.amritsar",
	}, &domain.RequestInfo{
		UserInfo: &domain.ContractUser{Type: "EMPLOYEE"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lands) != 0 {
		t.Fatalf("expected empty slice, got %d", len(lands))
	}
}

func TestSearchEmployeeEmptyCriteriaRejected(t *testing.T) {
	svc := NewLandService(
		validator.NewLandValidator(validator.NewLandMDMSValidator()),
		nil, nil, &mockLandRepo{}, stubMDMS{data: validMDMS()},
	)

	_, err := svc.Search(context.Background(), domain.LandSearchCriteria{}, &domain.RequestInfo{
		UserInfo: &domain.ContractUser{Type: "EMPLOYEE"},
	})
	if err == nil {
		t.Fatal("expected invalid search")
	}
	if ce, ok := err.(*landerrors.CustomException); !ok || ce.Code != landerrors.InvalidSearch {
		t.Fatalf("got %v", err)
	}
}

func TestSetOwnerStatusFromActive(t *testing.T) {
	owners := []domain.OwnerInfo{
		{Active: boolPtr(true)},
		{Active: boolPtr(false)},
		{},
	}
	setOwnerStatusFromActive(owners)
	if owners[0].Status == nil || !*owners[0].Status {
		t.Fatal("expected active owner status true")
	}
	if owners[1].Status == nil || *owners[1].Status {
		t.Fatal("expected inactive owner status false")
	}
	if owners[2].Status == nil || *owners[2].Status {
		t.Fatal("nil active should map to status false")
	}
}
