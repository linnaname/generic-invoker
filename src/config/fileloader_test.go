package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	loader := &FileLoader{"config.json", "."}
	configs, err := loader.Load()
	assert.NoError(t, err)
	assert.NotNil(t, configs)
	assert.Len(t, configs, 2)
}

func BenchmarkFile_Load(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loader := &FileLoader{"config.json", "."}
		loader.Load()
	}
}
