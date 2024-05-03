package system

import (
	"github.com/romainframe/grunter/pkg/grunter/block"
)

type System struct {
	Name    string        `json:"name"`
	Systems []System      `json:"systems"`
	Blocks  []block.Block `json:"blocks"`
}
