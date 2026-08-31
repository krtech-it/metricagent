package audit

import "go.uber.org/zap"

//go:generate mockgen -source=internal/audit/audit.go -destination=internal/audit/mock_audit.go -package=audit

// Event описывает событие аудита запроса на обновление метрик.
type Event struct {
	Timestamp int64    `json:"ts"`
	Metrics   []string `json:"metrics"`
	IPAddress string   `json:"ip_address"`
}

// Observer — подписчик на события аудита (файл, удалённый сервер и т. д.).
type Observer interface {
	Notify(event Event) error
}

// Publisher — издатель, хранит список подписчиков и рассылает им события.
type Publisher struct {
	observers []Observer
	logger    *zap.Logger
}

func NewPublisher(logger *zap.Logger, observers ...Observer) *Publisher {
	return &Publisher{observers: observers, logger: logger}
}

func (p *Publisher) Register(o Observer) {
	p.observers = append(p.observers, o)
}

// Notify рассылает событие всем зарегистрированным наблюдателям.
// Ошибка одного наблюдателя не прерывает рассылку остальным.
func (p *Publisher) Notify(event Event) {
	if p == nil {
		return
	}
	for _, o := range p.observers {
		if err := o.Notify(event); err != nil {
			if p.logger != nil {
				p.logger.Error("audit observer failed", zap.Error(err))
			}
		}
	}
}
