package models

import (
    "github.com/dylanxhernandez/proto-padmd/internal/db"
)

type Document struct {
    ID int
    Title string
    Content string
}
