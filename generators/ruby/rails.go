package ruby

import (
	"github.com/ronyv89/temple/dockertools"
)

func NewRailsProject() {
	dockertools.PullImage("ruby")
}
