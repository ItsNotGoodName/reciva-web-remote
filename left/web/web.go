package web

import (
	"mime"
)

func init() {
	mime.AddExtensionType(".js", "application/javascript")
}
