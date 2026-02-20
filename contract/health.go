package contract

type HealthRepository interface {
	GetStatus() string
}

type HealthService interface {
	GetStatus() string
}
