package audit

import (
	"encoding/json"
	"os"
	"sync"
)

// FileObserver дописывает событие аудита в конец файла построчным JSON.
type FileObserver struct {
	file *os.File
	mu   sync.Mutex
}

func NewFileObserver(fileName string) (*FileObserver, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &FileObserver{file: file}, nil
}

func (o *FileObserver) Notify(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	o.mu.Lock()
	defer o.mu.Unlock()
	_, err = o.file.Write(data)
	return err
}
