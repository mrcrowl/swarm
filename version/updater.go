package version

import (
	"errors"
	"io"

	update "github.com/inconshreveable/go-update"
)

// UpdaterLike is an interface that looks like a go-updater
type UpdaterLike interface {
	Apply(update io.Reader, opts update.Options) error
}

// real implementation

type realUpdater struct{}

func (h *realUpdater) Apply(updateReader io.Reader, opts update.Options) error {
	return update.Apply(updateReader, opts)
}

// mock implementation

type mockUpdater struct {
	simluateError bool
	successful    bool
}

func (mock *mockUpdater) SimulateError() {
	mock.simluateError = true
}

func (mock *mockUpdater) Apply(updateReader io.Reader, opts update.Options) error {
	if mock.simluateError {
		return errors.New("Pretend error in Apply")
	}
	mock.successful = true
	return nil
}
