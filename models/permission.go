package models

type Access int

const (
	NO_ACCESS 	Access = 0
	READ 		Access = 2
	WRITE 		Access = 5
	ACCESS 		Access = 7
)

type Permission struct {
	View string
	Access Access
}