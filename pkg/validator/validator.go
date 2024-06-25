package validator

import (
	"strings"
	"template/internal/entity"
	"unicode/utf8"
)

type Validator struct {
	UserErrors entity.UserFormError `json:"user_error,omitempty" bson:"user_error,omitempty"`
	BookErrors entity.BookFormError `json:"book_error,omitempty" bson:"book_error,omitempty"`
}

func (v *Validator) ValidUser() bool {
	return !NotBlank(v.UserErrors.Username) && !NotBlank(v.UserErrors.Password) && !NotBlank(v.UserErrors.Secret)
}

func (v *Validator) ValidBook() bool {
	return !NotBlank(v.BookErrors.ISBN) && !NotBlank(v.BookErrors.Title) && !NotBlank(v.BookErrors.Author) && !NotBlank(v.BookErrors.Publisher)

}

func (v *Validator) CheckField(ok bool, key *string, message string) {
	if !ok {
		*key = message
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func CheckAdmin(value string, key string) bool {
	return value == key
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(strings.TrimSpace(value)) >= n
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func CheckISBN(value string, n int) bool {

	str := strings.ReplaceAll(value, "-", "")
	return utf8.RuneCountInString(str) == n
}

func CheckArr(arr []string) bool {
	if len(arr) == 0 {
		return false
	}
	for i := range arr {
		if !NotBlank(arr[i]) {
			return false
		}
	}
	return true
}
