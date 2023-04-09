package io

import "io"

func TryToClose(closer io.Closer) {
	_ = closer.Close()
}
