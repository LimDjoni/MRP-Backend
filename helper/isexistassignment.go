package helper

import (
	"ajebackend/model/cafassignmentenduser"
	"ajebackend/model/electricassignmentenduser"
)

func IsExistElectricAssignment(id uint, list []electricassignmentenduser.ElectricAssignmentEndUser) bool {
	var result = false

	for _, value := range list {
		if value.ID == id {
			result = true
			return result
		}
	}

	return result
}

func IsExistCafAssignment(id uint, list []cafassignmentenduser.CafAssignmentEndUser) bool {
	var result = false

	for _, value := range list {
		if value.ID == id {
			result = true
			return result
		}
	}

	return result
}
