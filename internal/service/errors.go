package service

import "github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"

type CustomException = landerrors.CustomException

var (
	NewCustomException    = landerrors.New
	NewCustomExceptionMap = landerrors.NewMap
)
