package backuper

import (
	"encoding/json"
	"github.com/krtech-it/metricagent/internal/delivery/http/dto"
	"go.uber.org/zap"
	"os"
)

type Backuper struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
	logger  *zap.Logger
}

func NewBackuper(fileName string, logger *zap.Logger) (*Backuper, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Backuper{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
		logger:  logger,
	}, nil
}

func (p *Backuper) WriteEvent(metrics []*dto.ResponseGetMetric) error {
	p.file.Truncate(0)
	p.file.Seek(0, 0)
	return p.encoder.Encode(metrics)
}

func (p *Backuper) ReadEvent() ([]*dto.ResponseGetMetric, error) {
	p.file.Seek(0, 0)
	var events []*dto.ResponseGetMetric
	if err := p.decoder.Decode(&events); err != nil {
		return nil, err
	}
	return events, nil
}
