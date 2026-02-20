package repository

type HealthRepository struct{}

func NewHealthRepository() *HealthRepository {
	return &HealthRepository{}
}

func (r *HealthRepository) GetStatus() string {
	return "ok"
}
