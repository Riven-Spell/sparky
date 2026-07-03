package vllm_models

import "github.com/Riven-Spell/sparky/common/vllm_enum"

type ChatCompletionContentPartTextParam struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type ChatCompletionContentPartImageParam struct {
	ImageURL ImageURL `json:"image_url"`
	Type     string   `json:"type"`
}

type ChatCompletionContentPartRefusalParam struct {
	Refusal string `json:"refusal"`
	Type    string `json:"type"`
}

type ChatCompletionContentPartAudioParam struct {
	AudioURL AudioURL `json:"audio_url"`
	Type     string   `json:"type"`
}

type ChatCompletionContentPartInputAudioParam struct {
	InputAudio InputAudio `json:"input_audio"`
	Type       string     `json:"type"`
}

type ChatCompletionContentPartVideoParam struct {
	VideoURL VideoURL `json:"video_url"`
	Type     string   `json:"type"`
}

type ImageURL struct {
	URL    string                  `json:"url"`
	Detail *vllm_enum.ImageDetail  `json:"detail,omitempty"`
}

type InputAudio struct {
	Data   string                    `json:"data"`
	Format vllm_enum.InputAudioFormat `json:"format"`
}

type AudioURL struct {
	URL string `json:"url"`
}

type VideoURL struct {
	URL string `json:"url"`
}

type File struct {
	File FileFile `json:"file"`
	Type string   `json:"type"`
}

type FileFile struct {
	FileData *string `json:"file_data,omitempty"`
	FileID   *string `json:"file_id,omitempty"`
	Filename *string `json:"filename,omitempty"`
}
