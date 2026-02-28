package bun

import (
	"goserve/internal/action"
	"goserve/internal/pagination"
)

func getFilter(actor *action.Actor) (string, []any) {
	var filter string
	var args []any
	switch actor.Scope {
	case "owner":
		filter = "created_by = ?"
		args = append(args, actor.UID)
	}
	return filter, args
}

func getSortOrder(order pagination.SortOrder) string {
	if order == pagination.SortOrderAscending {
		return "ASC"
	}
	return "DESC"
}
