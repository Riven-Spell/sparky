package vllm_enum

import "github.com/Riven-Spell/enum/v2"

type eInputAudioFormat struct {
	enum.EnumImpl[InputAudioFormat, eInputAudioFormat]
}

var EInputAudioFormat eInputAudioFormat

type InputAudioFormat string

func (f InputAudioFormat) String() string {
	return string(f)
}

func (eInputAudioFormat) Wav() InputAudioFormat { return "wav" }
func (eInputAudioFormat) Mp3() InputAudioFormat { return "mp3" }
