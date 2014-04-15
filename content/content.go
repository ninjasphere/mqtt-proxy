package content

import "github.com/GeertJohan/go.rice"

func ContentBox() *rice.Box {
	return rice.MustFindBox("pages")
}
