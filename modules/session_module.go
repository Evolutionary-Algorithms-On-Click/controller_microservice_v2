package modules

import (
	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
)

// SessionModule encapsulates the business logic for sessions.
type SessionModule struct {
	Repo repository.SessionRepository
}

// NewSessionModule creates and returns a new SessionModule.
func NewSessionModule(repo repository.SessionRepository) *SessionModule {
	return &SessionModule{
		Repo: repo,
	}
}

// TODO: Implement business logic methods for SessionModule
// Example:
// func (m *SessionModule) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
// 	 return m.Repo.GetByID(ctx, id)
// }
