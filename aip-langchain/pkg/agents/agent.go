package agents

import "github.com/jbenet/goprocess"

type Agent interface {
	Run(proc goprocess.Process)
}
