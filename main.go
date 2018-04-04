package main

import (
	"context"

	"github.com/rexray/gocsi"

	"github.com/princerachit/csi-openebs/provider"
	"github.com/princerachit/csi-openebs/service"
	"github.com/sirupsen/logrus"
)

// main is ignored when this package is built as a go plug-in.
func main() {
	logrus.Info("Running main")
	gocsi.Run(
		context.Background(),
		service.Name,
		"A description of the SP",
		"",
		provider.New())
}
