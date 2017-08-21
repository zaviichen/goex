package common

//import (
//	"github.com/sirupsen/logrus"
//	"os"
//	"log"
//)
//
////var log = logrus.New()
//
//func init() {
//	log.Formatter = new(logrus.JSONFormatter)
//	log.Formatter = new(logrus.TextFormatter) // default
//
//	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
//	if err == nil {
//		log.Out = file
//	} else {
//		log.Info("Failed to log to file, using default stderr")
//	}
//
//	log.Level = logrus.DebugLevel
//}
