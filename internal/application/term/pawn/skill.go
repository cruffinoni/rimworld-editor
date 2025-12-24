package pawn

import (
	"github.com/iancoleman/strcase"

	"github.com/cruffinoni/rimworld-editor/internal/xml/algorithm"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/application/terminal/faction"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type Skills struct {
	sg     *generated.Savegame
	rp     PawnsRegisterer
	rf     faction.Registerer
	logger logging.Logger
}

func NewSkills(logger logging.Logger, sg *generated.Savegame, rp PawnsRegisterer, rf faction.Registerer) *Skills {
	return &Skills{
		sg:     sg,
		rp:     rp,
		rf:     rf,
		logger: logger,
	}
}

const (
	PassionNone  = "None"
	PassionMinor = "Minor"
	PassionMajor = "Major"
)

var skillCap = map[string]int64{
	"Shooting":     20,
	"Melee":        20,
	"Construction": 20,
	"Mining":       20,
	"Cooking":      20,
	"Plants":       20,
	"Animals":      20,
	"Crafting":     20,
	"Artistic":     20,
	"Medicine":     20,
	"Social":       20,
	"Intellectual": 20,
}

func (s *Skills) Edit(pawnID, skill, passion string, level int) {
	v, ok := s.rp[pawnID]
	if !ok {
		s.logger.WithField("pawn_id", pawnID).Error("Pawn not found")
		return
	}
	skill = strcase.ToCamel(skill)
	if skillCap[skill] <= 0 {
		s.logger.WithField("skill", skill).Error("Skill not found")
		return
	}
	if level < 0 || level > 20 {
		s.logger.WithField("level", level).Error("Level must be between 0 and 20")
		return
	}
	passion = strcase.ToCamel(passion)
	if passion != PassionMinor && passion != PassionMajor && passion != PassionNone {
		s.logger.WithField("passion", passion).Error("Passion must be either 'Minor', 'Major', or 'None'")
		return
	}
	for i, j := 0, v.Skills.Skills.Capacity(); i < j; i++ {
		val := v.Skills.Skills.At(i)
		if val.Def == skill {
			val.Level = int64(level)
			if passion != PassionNone {
				val.Passion = passion
				val.ValidateField("Passion")
			}
			val.XpSinceLastLevel = float64(s.calculateMaxXP(val.Level) * 1.0)
			v.Skills.Skills.Set(val, val.Attr, i)
			s.logger.WithFields(logging.Fields{
				"pawn":    getPawnFullName(v),
				"skill":   skill,
				"level":   level,
				"passion": passion,
			}).Info("Skill updated")
			return
		}
	}
}

/*
0-9 = 1000 + level*1000
10-19= 12000 + (level-10)*2000
20=30000
*/
func (s *Skills) calculateMaxXP(level int64) int64 {
	if level < 10 {
		return 1000 + level*1000
	} else if level < 20 {
		return 12000 + (level-10)*2000
	} else {
		return 30000
	}
}

func (s *Skills) ForceGraduate(pawnID string) {
	v, ok := s.rp[pawnID]
	if !ok {
		s.logger.WithField("pawn_id", pawnID).Error("Pawn not found")
		return
	}
	for i, j := 0, v.Skills.Skills.Capacity(); i < j; i++ {
		val := v.Skills.Skills.At(i)
		val.Level = int64(i + 1)
		if val.Passion == "" {
			val.Passion = PassionMinor
			val.ValidateField("Passion")
		}
		val.XpSinceLastLevel = float64(s.calculateMaxXP(val.Level) * 1.0)
		v.Skills.Skills.Set(val, val.Attr, i)
		s.logger.WithFields(logging.Fields{
			"pawn":    getPawnFullName(v),
			"skill":   val.Def,
			"level":   val.Level,
			"passion": val.Passion,
		}).Info("Skill updated")
	}
}

func (s *Skills) printPawnSkill(fullName string, p *generated.Thing) {
	s.logger.WithFields(logging.Fields{
		"pawn_id":  fullName,
		"fullName": getPawnFullName(p),
	}).Info("Pawn skills")
	algorithm.SliceForeach[*generated.SkillsSkillsLiPawnsMothballedWorldPawnsWorldGameSavegameInner](p.Skills.Skills, func(skill *generated.SkillsSkillsLiPawnsMothballedWorldPawnsWorldGameSavegameInner) {
		s.logger.WithFields(logging.Fields{
			"skill":   skill.Def,
			"level":   skill.Level,
			"passion": skill.Passion,
		}).Info("Skill entry")
	})
}

func (s *Skills) ListPawnsSkills(pawnID string) {
	if pawnID != "ALL" {
		v, ok := s.rp[pawnID]
		if !ok {
			s.logger.WithField("pawn_id", pawnID).Error("Pawn not found")
			return
		}
		s.printPawnSkill(pawnID, v)
	} else {
		for k, v := range s.rp {
			s.printPawnSkill(k, v)
		}
	}
}
