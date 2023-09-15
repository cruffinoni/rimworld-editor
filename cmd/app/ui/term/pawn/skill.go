package pawn

import (
	"github.com/cruffinoni/rimworld-editor/algorithm"
	"github.com/iancoleman/strcase"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/faction"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
	"github.com/cruffinoni/rimworld-editor/generated"
)

type Skills struct {
	sg *generated.Savegame
	rp PawnsRegisterer
	rf faction.Registerer
}

func NewSkills(sg *generated.Savegame, rp PawnsRegisterer, rf faction.Registerer) *Skills {
	return &Skills{
		sg: sg,
		rp: rp,
		rf: rf,
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
		printer.PrintErrorSf("Pawn '%s' not found", pawnID)
		return
	}
	skill = strcase.ToCamel(skill)
	if skillCap[skill] <= 0 {
		printer.PrintErrorSf("Skill '%s' not found", skill)
		return
	}
	if level < 0 || level > 20 {
		printer.PrintErrorSf("Level '%d' must be between 0 & 20", level)
		return
	}
	passion = strcase.ToCamel(passion)
	if passion != PassionMinor && passion != PassionMajor && passion != PassionNone {
		printer.PrintErrorS("Passion must be either 'Minor' or 'Major'")
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
			if passion != PassionNone {
				printer.Printf("Skill {-BOLD}%s{-RESET} of {-BOLD}%s{-RESET} set to {-BOLD}%d{-RESET} with passion '{-BOLD}%s{-RESET}'", skill, getPawnFullNameColorFormatted(v), level, passion)
			} else {
				printer.Printf("Skill {-BOLD}%s{-RESET} of {-BOLD}%s{-RESET} set to {-BOLD}%d{-RESET}", skill, getPawnFullNameColorFormatted(v), level)
			}
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
		printer.PrintErrorSf("Pawn '%s' not found", pawnID)
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
		printer.Printf("Skill {-BOLD}%s{-RESET} of {-BOLD}%s{-RESET} set to {-BOLD}%d{-RESET} with passion '{-BOLD}%s{-RESET}'", val.Def, getPawnFullNameColorFormatted(v), val.Level, val.Passion)
	}
}

func (s *Skills) printPawnSkill(fullName string, p *generated.Thing) {
	printer.Printf("Pawn {-BOLD}%s's{-RESET} (%s) skills", fullName, getPawnFullNameColorFormatted(p))
	algorithm.SliceForeach[*generated.SkillsSkillsLiPawnsAliveWorldPawnsWorldGameSavegameInner](p.Skills.Skills, func(skill *generated.SkillsSkillsLiPawnsAliveWorldPawnsWorldGameSavegameInner) {
		var color string
		if skill.Level <= 5 {
			color = "F_RED"
		} else if skill.Level > 5 && skill.Level < 12 {
			color = "F_YELLOW"
		} else {
			color = "F_GREEN"
		}
		printer.Printf("Skill: {-BOLD}%s", skill.Def)
		if skill.Passion != "" {
			printer.Printf("\tLevel: {-%s}%d{-RESET} | Passion: {-BOLD,F_CYAN}%s", color, skill.Level, skill.Passion)
		} else {
			printer.Printf("\tLevel: {-%s}%d", color, skill.Level)
		}
	})
}

func (s *Skills) ListPawnsSkills(pawnID string) {
	if pawnID != "ALL" {
		v, ok := s.rp[pawnID]
		if !ok {
			printer.PrintErrorSf("Pawn '%s' not found", pawnID)
			return
		}
		s.printPawnSkill(pawnID, v)
	} else {
		for k, v := range s.rp {
			s.printPawnSkill(k, v)
		}
	}
}
