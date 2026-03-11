package repository

import (
	"context"
	"database/sql"
	"fmt"
	models "github.com/krtech-it/metricagent/internal/model"
	"time"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) Storage {
	return &DBStorage{
		db: db,
	}
}

func (m *DBStorage) Create(ctx context.Context, metric *models.Metrics) error {
	retryIntervals := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

	return WithRetry(ctx, retryIntervals, func() error {
		if _, err := m.Get(ctx, metric.ID); err == nil {
			return fmt.Errorf("metric %v already exists", metric.ID)
		}
		_, err := m.db.ExecContext(ctx, "INSERT INTO metrics (id, m_type, delta, value) values ($1, $2, $3, $4)",
			metric.ID, metric.MType, metric.Delta, metric.Value)
		if err != nil {
			return err
		}
		return nil
	})
}

func (m *DBStorage) Update(ctx context.Context, metric *models.Metrics) error {
	retryIntervals := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

	return WithRetry(ctx, retryIntervals, func() error {
		if _, err := m.Get(ctx, metric.ID); err != nil {
			return fmt.Errorf("metric %v does not exist", metric.ID)
		}
		_, err := m.db.ExecContext(ctx, "update metrics set m_type = $2, delta = $3, value = $4 WHERE id = $1",
			metric.ID, metric.MType, metric.Delta, metric.Value)
		if err != nil {
			return err
		}
		return nil
	})

}

func (m *DBStorage) Upsert(ctx context.Context, metrics []*models.Metrics) error {
	retryIntervals := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

	return WithRetry(ctx, retryIntervals, func() error {
		tx, err := m.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		stmt, err := tx.PrepareContext(ctx, `
INSERT INTO metrics
(id, m_type, delta, value) values ($1, $2, $3, $4)
ON CONFLICT(id) DO UPDATE
set 
    m_type = EXCLUDED.m_type,
        delta = CASE
            WHEN EXCLUDED.m_type = 'counter' THEN
                COALESCE(metrics.delta, 0) + COALESCE(EXCLUDED.delta, 0)
            ELSE
                NULL
        END,
        value = CASE
            WHEN EXCLUDED.m_type = 'gauge' THEN
                EXCLUDED.value
            ELSE
                NULL
        END`)
		if err != nil {
			return fmt.Errorf("failed to prepare statement %w", err)
		}
		defer stmt.Close()
		for _, metric := range metrics {
			if _, err := stmt.ExecContext(ctx, metric.ID, metric.MType,
				metric.Delta, metric.Value); err != nil {
				return err
			}
		}
		return tx.Commit()
	})

}

func (m *DBStorage) Get(ctx context.Context, ID string) (*models.Metrics, error) {
	metric := &models.Metrics{}
	var (
		deltaNull sql.NullInt64
		valueNull sql.NullFloat64
	)
	row := m.db.QueryRowContext(ctx, "select id, m_type, delta, value from metrics where id = $1", ID)
	err := row.Scan(&metric.ID, &metric.MType, &deltaNull, &valueNull)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("metric %v does not exist", ID)
		}
		return nil, err
	}
	if deltaNull.Valid {
		metric.Delta = &deltaNull.Int64
	}
	if valueNull.Valid {
		metric.Value = &valueNull.Float64
	}
	return metric, nil
}

func (m *DBStorage) GetAll(ctx context.Context) ([]*models.Metrics, error) {
	var metrics []*models.Metrics
	retryIntervals := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

	err := WithRetry(ctx, retryIntervals, func() error {

		rows, err := m.db.QueryContext(ctx, "select id, m_type, delta, value from metrics")
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var (
				metric    models.Metrics
				deltaNull sql.NullInt64
				valueNull sql.NullFloat64
			)
			err = rows.Scan(&metric.ID, &metric.MType, &deltaNull, &valueNull)
			if err != nil {
				return err
			}
			if deltaNull.Valid {
				metric.Delta = &deltaNull.Int64
			}
			if valueNull.Valid {
				metric.Value = &valueNull.Float64
			}
			metrics = append(metrics, &metric)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	})
	return metrics, err
}

func (m *DBStorage) Ping(ctx context.Context) error {
	if m.db == nil {
		return fmt.Errorf("db is not configured")
	}
	ctxPing, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := m.db.PingContext(ctxPing); err != nil {
		return err
	}
	return nil
}
