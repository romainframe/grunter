package system

import (
	"errors"
)

func (s System) Build() (System, error) {
	if len(s.Systems) == 0 && len(s.Blocks) == 0 {
		return System{}, errors.New("no blocks defined")
	}

	for i, b := range s.Blocks {
		block, err := b.Build(s.Name)
		if err != nil {
			return System{}, err
		}
		s.Blocks[i] = block
	}

	for j, subSys := range s.Systems {
		system, err := subSys.Build()
		if err != nil {
			return System{}, err
		}
		s.Systems[j] = system
	}

	return s, nil
}
