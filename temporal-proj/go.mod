module temporal-proj

go 1.23

require (
	gorm.io/gorm v1.25.12
	github.com/PuerkitoBio/goquery v1.10.1
	github.com/gin-gonic/gin v1.10.0
	go.temporal.io/sdk v1.34.0
	nomaproj/pkg/models v0.0.0
)

replace nomaproj/pkg/models => ../pkg/models 