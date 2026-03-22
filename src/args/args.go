package args

import (
	"errors"
	"fmt"
)

type Mode string

const (
	ModeDetails Mode = "details"
	ModeLite    Mode = "lite"
)

// Parsed holds the result of parsing os.Args[1:].
type Parsed struct {
	Help       bool
	Mode       Mode
	Alias      string
	ConfigPath string
}

// Parse parses argv (os.Args[1:]) into a Parsed value.
func Parse(argv []string) (Parsed, error) {
	var p Parsed
	for i := 0; i < len(argv); i++ {
		arg := argv[i]
		switch arg {
		case "help", "-h", "--help":
			p.Help = true
		case "-m", "--mode":
			if i+1 >= len(argv) {
				return p, errors.New("--mode requires a value (details|lite)")
			}
			i++
			switch argv[i] {
			case "details":
				p.Mode = ModeDetails
			case "lite":
				p.Mode = ModeLite
			default:
				return p, fmt.Errorf("unknown mode %q; valid values: details, lite", argv[i])
			}
		case "-c", "--config":
			if i+1 >= len(argv) {
				return p, errors.New("--config requires a value")
			}
			i++
			p.ConfigPath = argv[i]
		case "-a", "--alias":
			if i+1 >= len(argv) {
				return p, errors.New("--alias requires a value")
			}
			i++
			p.Alias = argv[i]
		default:
			return p, fmt.Errorf("unknown argument %q", arg)
		}
	}
	return p, nil
}
