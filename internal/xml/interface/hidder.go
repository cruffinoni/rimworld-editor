package _interface

type FieldValidator interface {
	ValidateField(field string)
	IsValidField(field string) bool
	CountValidatedField() int
}
