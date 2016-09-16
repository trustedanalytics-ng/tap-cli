package cli

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"os"
	"github.com/trustedanalytics/tap-cli/api"
)

func TestApiAndLoginServiceSetters(t *testing.T) {
	Convey("Test OAuth2 login2", t, func() {
		Convey("Should fail when no credentials.json file", func() {
			os.Remove(api.CredsPath)

			So(func(){
				NewOAuth2Service()}, ShouldPanicWith, "Please login first!")
		})
		Convey("Should fail when wrong format in credentials.json file", func() {
			wrongContent := "@"
			fillCredentialsFile(wrongContent)

			So(func(){
				NewOAuth2Service()}, ShouldPanicWith, "invalid character '" + wrongContent + "' looking for beginning of value")
		})
	})
}