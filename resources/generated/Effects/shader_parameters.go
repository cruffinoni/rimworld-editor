// Code generated by rimworld-editor. DO NOT EDIT.

package effects

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type ShaderParameters struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DistortionTex_             string  `xml:"_DistortionTex"`
	DistortionScrollSpeed_     float64 `xml:"_distortionScrollSpeed"`
	BounceTex_                 string  `xml:"_BounceTex"`
	BounceSpeed_               float64 `xml:"_bounceSpeed"`
	BounceAmplitude_           float64 `xml:"_bounceAmplitude"`
	RotationTex_               string  `xml:"_RotationTex"`
	RotationFromOffset_        int64   `xml:"_RotationFromOffset"`
	RotationSpeed_             float64 `xml:"_rotationSpeed"`
	RotationAmplitude_         float64 `xml:"_rotationAmplitude"`
	DistortionIntensity_       float64 `xml:"_distortionIntensity"`
	TexScale_                  float64 `xml:"_texScale"`
	TexBscale_                 float64 `xml:"_texBScale"`
	TexBscrollSpeed_           float64 `xml:"_texBScrollSpeed"`
	Intensity_                 float64 `xml:"_Intensity"`
	BrightnessMultiplier_      float64 `xml:"_brightnessMultiplier"`
	Clip_                      int64   `xml:"_Clip"`
	MultiplyTex_               string  `xml:"_MultiplyTex"`
	ExtraTex_                  string  `xml:"_ExtraTex"`
	TexAscrollSpeed_           float64 `xml:"_texAScrollSpeed"`
	MultiplyTexB_              string  `xml:"_MultiplyTexB"`
	TexAscale_                 float64 `xml:"_texAScale"`
	SmokeTex_                  string  `xml:"_SmokeTex"`
	SmokeScrollSpeed_          float64 `xml:"_smokeScrollSpeed"`
	ChangeSpeed_               float64 `xml:"_ChangeSpeed"`
	DetailScrollSpeed_         string  `xml:"_detailScrollSpeed"`
	ScanTex_                   string  `xml:"_ScanTex"`
	SmokeTex2_                 string  `xml:"_SmokeTex2"`
	NumFrames_                 int64   `xml:"_NumFrames"`
	DetailIntensity_           float64 `xml:"_detailIntensity"`
	DetailScale_               float64 `xml:"_DetailScale"`
	VerticalScale_             float64 `xml:"_VerticalScale"`
	SmokeAmount_               float64 `xml:"_smokeAmount"`
	DistortionScale_           float64 `xml:"_distortionScale"`
	MainTex_                   string  `xml:"_MainTex"`
	WordSpaceDistortionToggle_ int64   `xml:"_wordSpaceDistortionToggle"`
	PulseSpeed_                float64 `xml:"_pulseSpeed"`
	FramesPerSec_              int64   `xml:"_FramesPerSec"`
	GradientScale_             float64 `xml:"_GradientScale"`
	MaskTex_                   string  `xml:"_MaskTex"`
	GradientTex_               string  `xml:"_GradientTex"`
	Thickness_                 int64   `xml:"_Thickness"`
	SmokeTex1_                 string  `xml:"_SmokeTex1"`
	ScrollSpeed_               float64 `xml:"_ScrollSpeed"`
	Distortion_                float64 `xml:"_Distortion"`
	Interval_                  float64 `xml:"_Interval"`
	DetailTex_                 string  `xml:"_DetailTex"`
	Detail_                    int64   `xml:"_Detail"`
	MultiplyTexA_              string  `xml:"_MultiplyTexA"`
	InnerRingSize_             float64 `xml:"_innerRingSize"`
	OutTime_                   float64 `xml:"_outTime"`
	AgeOffset_                 float64 `xml:"_AgeOffset"`
	Noise_                     string  `xml:"_Noise"`
	InTime_                    float64 `xml:"_inTime"`
	SolidTime_                 float64 `xml:"_solidTime"`
	OuterRingSize_             float64 `xml:"_outerRingSize"`
	NoiseTex_                  string  `xml:"_NoiseTex"`
}

func (s *ShaderParameters) Assign(*xml.Element) error {
	return nil
}

func (s *ShaderParameters) CountValidatedField() int {
	if s.FieldValidated == nil {
		return 0
	}
	return len(s.FieldValidated)
}

func (s *ShaderParameters) Equal(*ShaderParameters) bool {
	return false
}

func (s *ShaderParameters) GetAttributes() attributes.Attributes {
	return s.Attr
}

func (s *ShaderParameters) GetPath() string {
	return ""
}

func (s *ShaderParameters) Greater(*ShaderParameters) bool {
	return false
}

func (s *ShaderParameters) IsValidField(field string) bool {
	return s.FieldValidated[field]
}

func (s *ShaderParameters) Less(*ShaderParameters) bool {
	return false
}

func (s *ShaderParameters) SetAttributes(attr attributes.Attributes) {
	s.Attr = attr
	return
}

func (s *ShaderParameters) Val() *ShaderParameters {
	return nil
}

func (s *ShaderParameters) ValidateField(field string) {
	if s.FieldValidated == nil {
		s.FieldValidated = make(map[string]bool)
	}
	s.FieldValidated[field] = true
	return
}
