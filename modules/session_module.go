package modules

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/Thanus-Kumaar/controller_microservice_v2/db"
// 	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
// 	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
// 	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v4/pgxpool"
// 	"github.com/rs/zerolog"
// )

// // SessionModule holds dependencies for business logic operations.
// type SessionModule struct {
// 	JupyterClient *jupyterclient.Client
// 	DBPool        *pgxpool.Pool
// 	Logger        *zerolog.Logger
// }

// // NewSessionModule creates and returns a new SessionModule instance.
// func NewSessionModule(client *jupyterclient.Client) *SessionModule {
// 	return &SessionModule{
// 		JupyterClient: client,
// 		DBPool:        db.Pool,
// 		Logger:        pkg.Logger,
// 	}
// }

// // CreateNewSession orchestrates the creation of resources and persistence.
// func (m *SessionModule) CreateNewSession(ctx context.Context, language string, notebookID uuid.UUID, userID uuid.UUID) (*models.SessionConnectionResponse, error) {
// 	// 1. Generate the unique, persistent session ID
// 	sessionID := uuid.New()

// 	// Use a shorter context for the external kernel call
// 	kernelCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
// 	defer cancel()

// 	// 2. Start the kernel via the Gateway
// 	kernelInfo, err := m.JupyterClient.StartKernel(kernelCtx, language)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to start kernel: %w", err)
// 	}

// 	// Clean-up defer is essential: if DB transaction fails, kill the kernel.
// 	defer func() {
// 		if r := recover(); r != nil || err != nil {
// 			m.Logger.Error().Str("kernel_id", kernelInfo.ID).Msg("DB transaction failed, cleaning up kernel")
// 			// NOTE: DeleteKernel should run asynchronously to avoid blocking recovery, but here we run it inline for guaranteed cleanup.
// 			m.JupyterClient.DeleteKernel(context.Background(), kernelInfo.ID)
// 		}
// 	}()

// 	// 3. Database Transaction (Atomic Write)
// 	// TODO: Replace this with actual database persistence
// 	/*
// 		conn, err := m.DBPool.Acquire(ctx)
// 		if err != nil { return nil, fmt.Errorf("could not acquire db connection: %w", err) }
// 		defer conn.Release()

// 		// Start Transaction
// 		tx, err := conn.Begin(ctx)
// 		if err != nil { return nil, fmt.Errorf("failed to start db transaction: %w", err) }

// 		// 3a. Store new Session record
// 		_, err = tx.Exec(ctx, `INSERT INTO sessions (id, notebook_id, current_kernel_id, status, last_active_at) VALUES ($1, $2, $3, $4, $5)`,
// 			sessionID, notebookID, kernelInfo.ID, "active", time.Now())
// 		if err != nil { tx.Rollback(ctx); return nil, fmt.Errorf("failed to insert session: %w", err) }

// 		// 3b. Update Notebook's modified time
// 		_, err = tx.Exec(ctx, `UPDATE notebooks SET last_modified_at = $1 WHERE id = $2`, time.Now(), notebookID)
// 		if err != nil { tx.Rollback(ctx); return nil, fmt.Errorf("failed to update notebook: %w", err) }

// 		// 3c. Commit Transaction
// 		err = tx.Commit(ctx)
// 		if err != nil { return nil, fmt.Errorf("failed to commit transaction: %w", err) }
// 	*/

// 	// 4. Construct response for frontend (using placeholder logic for now)
// 	webSocketPath := fmt.Sprintf("/ws/session/%s", sessionID.String())

// 	m.Logger.Info().Str("session_id", sessionID.String()).Str("kernel_id", kernelInfo.ID).Msg("New session and kernel persistence complete")

// 	return &models.SessionConnectionResponse{
// 		SessionID:     sessionID.String(),
// 		KernelID:      kernelInfo.ID,
// 		WebSocketPath: webSocketPath,
// 	}, nil
// }
