package bot

import "errors"

var ErrHostIsIncorrect = errors.New("Repo host is incorrrect")
var ErrCantParseRepo = errors.New("Can't parse repo")
