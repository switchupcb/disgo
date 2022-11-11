package wrapper

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	Logger = zerolog.New(os.Stdout)
)
