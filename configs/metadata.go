package configs

import (
	"sync/atomic"

	"go.uber.org/zap"
)

type Metadata struct {
	Name      string
	Version   string
	Date      string
	GitCommit string
	GitState  string
	BuiltBy   string
}

func (d *Metadata) ToFields() (fields []zap.Field) {
	fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("name", d.Name))
	fields = append(fields, zap.String("version", d.Version))
	fields = append(fields, zap.String("date", d.Date))
	fields = append(fields, zap.String("git-commit", d.GitCommit))
	fields = append(fields, zap.String("git-state", d.GitState))
	fields = append(fields, zap.String("built-by", d.BuiltBy))

	return
}

var defaultMetadata atomic.Pointer[Metadata]

func GetMetadata() *Metadata {
	return defaultMetadata.Load()
}

func SetMetadata(d *Metadata) {
	defaultMetadata.Store(d)
}
