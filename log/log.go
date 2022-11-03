package log

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

const (
	mode_dev     = "dev"
	mode_release = "release"
)

func init() {
	fmt.Println("please set mode(dev/release) print log info ")
	fmt.Println("please set output log file.default print to stdio")
}
func Log() *log.Logger {
	return log.StandardLogger()
}
func SetOutput(out io.Writer) {
	log.SetOutput(out)
}
func SetMode(mode string) {
	if mode == mode_dev {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}
