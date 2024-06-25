package dto

import "template/internal/entity"

type Test struct {
	Tests []entity.Test `json:"tests"`
}

//dto convert
