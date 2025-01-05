package sheets

import "go.uber.org/zap"

var logger = zap.Must(zap.NewDevelopment()).Sugar()
