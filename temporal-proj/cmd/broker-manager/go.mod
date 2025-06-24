module temporal-proj/cmd/broker-manager

go 1.23.0

require (
	github.com/gin-gonic/gin v1.10.0
	go.temporal.io/sdk v1.34.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
	temporal-proj v0.0.0
	nomaproj/pkg/models v0.0.0
	nomaproj/pkg/utils v0.0.0
)

replace temporal-proj => ../..
replace nomaproj/pkg/models => ../../../pkg/models
replace nomaproj/pkg/utils => ../../../pkg/utils 