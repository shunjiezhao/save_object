package objects

import (
	"ch2/apiServer/locate"
	"ch2/lib/objectstream"
	"fmt"
	"io"
)

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}
