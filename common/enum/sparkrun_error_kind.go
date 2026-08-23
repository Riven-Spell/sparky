package enum

import "github.com/Riven-Spell/enum/v2"

type eSparkrunErrorKind struct {
	enum.EnumImpl[SparkrunErrorKind, eSparkrunErrorKind]
}

var ESparkrunErrorKind eSparkrunErrorKind

type SparkrunErrorKind string

func (k SparkrunErrorKind) String() string {
	return string(k)
}

func (eSparkrunErrorKind) Parse() SparkrunErrorKind  { return "parse" }
func (eSparkrunErrorKind) Exec() SparkrunErrorKind   { return "exec" }
func (eSparkrunErrorKind) Exit() SparkrunErrorKind   { return "exit" }
func (eSparkrunErrorKind) Usage() SparkrunErrorKind  { return "usage" }
func (eSparkrunErrorKind) Target() SparkrunErrorKind { return "target" }
