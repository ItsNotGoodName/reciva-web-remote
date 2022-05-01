package bus

import (
	"fmt"

	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

func parseStateSubscribe(slug interface{}) (string, error) {
	uuid := fmt.Sprint(slug)
	if uuid == "" {
		return "", fmt.Errorf("invalid uuid")
	}

	return uuid, nil
}

func newErrorCommand(err error) dto.Command {
	return dto.Command{
		Type: dto.CommandTypeError,
		Slug: err.Error(),
	}
}

func newStateCommand(state *state.State) dto.Command {
	return dto.Command{
		Type: dto.CommandTypeState,
		Slug: state,
	}
}

func newStatePartialCommand(partial state.Partial) dto.Command {
	return dto.Command{
		Type: dto.CommandTypeStatePartial,
		Slug: partial,
	}
}
