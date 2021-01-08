package t

import (
	"MediaImage/test"
	"os"
	"test/logger"
	"testing"
)

func TestTempDir(t *testing.T) {

	logger.Debug(os.TempDir())

	s1 := t.TempDir()
	logger.Debug(s1)

	logger.Debug(t.TempDir())

	logger.Debug(test.Path("a"))
}
