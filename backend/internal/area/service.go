package area

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]AreaResponse, error) {
	areas, err := s.repo.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]AreaResponse, 0, len(areas))
	for _, item := range areas {
		responses = append(responses, item.ToResponse())
	}

	return responses, nil
}
