package usecase

import "context"

type MinIOUseCase interface {
	GetFile(ctx context.Context, fileName string)
}
