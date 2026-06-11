package module

//go:generate go run ./gen

import (
	"fmt"

	"github.com/dinhtp/lee-goo/modules"
	"github.com/dinhtp/lee-goo/system/config"
	"github.com/dinhtp/lee-goo/system/logger"
	"github.com/dinhtp/lee-goo/system/migrator"
)

// buildRunner looks up the named module, loads config and logger, and returns
// a Runner scoped to that single source. Returns an error for unknown names
// before any DB connection is opened.
func buildRunner(name string) (*migrator.Runner, error) {
	s, ok := modules.Registry()[name]
	if !ok {
		return nil, fmt.Errorf("unknown module %q", name)
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("%s: load config: %w", name, err)
	}
	log, err := logger.NewLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: init logger: %w", name, err)
	}

	return migrator.NewRunner(cfg.Database.DSN(), log, []migrator.Source{s}), nil
}
