package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net"
	"time"
)

type RetryableError struct {
	Err error
}

func (e RetryableError) Error() string {
	return fmt.Sprintf("retryable error: %v", e.Err)
}

func (e RetryableError) Unwrap() error {
	return e.Err
}

func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return isRetriablePgCode(pgErr.Code)
	}

	return false
}

func isRetriablePgCode(code string) bool {
	switch code {
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure,
		pgerrcode.SQLClientUnableToEstablishSQLConnection,
		pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
		pgerrcode.ProtocolViolation:
		return true
	case pgerrcode.SerializationFailure,
		pgerrcode.DeadlockDetected,
		pgerrcode.TransactionRollback:
		return true
	case pgerrcode.CannotConnectNow,
		pgerrcode.AdminShutdown,
		pgerrcode.CrashShutdown:
		return true
	}
	return false
}

func WithRetry(ctx context.Context, intervals []time.Duration, fn func() error) error {
	var lastErr error

	if err := fn(); err == nil {
		return nil
	} else {
		lastErr = err
	}

	for _, delay := range intervals {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case <-time.After(delay):
		}

		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err

		if !IsRetryableError(err) {
			return fmt.Errorf("non-retryable error: %w", lastErr)
		}
	}
	return fmt.Errorf("last error: %w", lastErr)
}
