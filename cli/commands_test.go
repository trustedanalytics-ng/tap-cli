package cli

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/trustedanalytics/tap-api-service/client"
	"github.com/trustedanalytics/tap-cli/api"
	"os"
	"testing"
)

func TestApiAndLoginServiceSetters(t *testing.T) {
	Convey("Test OAuth2 login2", t, func() {
		Convey("Should fail when no credentials.json file", func() {
			os.Remove(api.CredsPath)

			So(func() {
				NewOAuth2Service()
			}, ShouldPanicWith, "Please login first!")
		})
		Convey("Should fail when wrong format in credentials.json file", func() {
			wrongContent := "@"
			fillCredentialsFile(wrongContent)

			So(func() {
				NewOAuth2Service()
			}, ShouldPanicWith, "invalid character '"+wrongContent+"' looking for beginning of value")
		})
	})
}

func TestNewBasicAuthService(t *testing.T) {
	Convey("Should add https address if address not provided", t, func() {
		basicAuth := NewBasicAuthService("myaddress.com", "user", "password")
		basicCreds := basicAuth.ApiServiceLogin.(*client.TapApiServiceApiBasicAuthConnector)

		So(basicCreds.Address, ShouldEqual, "https://myaddress.com")
	})
	Convey("Should not add https", t, func(){
		Convey("when there is http:// ", func(){

			basicAuth := NewBasicAuthService("http://myaddress.com", "user", "password")
			basicCreds := basicAuth.ApiServiceLogin.(*client.TapApiServiceApiBasicAuthConnector)

			So(basicCreds.Address, ShouldEqual, "http://myaddress.com")
		})
		Convey("when there is ftp:// ", func(){

			basicAuth := NewBasicAuthService("ftp://myaddress.com", "user", "password")
			basicCreds := basicAuth.ApiServiceLogin.(*client.TapApiServiceApiBasicAuthConnector)

			So(basicCreds.Address, ShouldEqual, "ftp://myaddress.com")
		})
	})
}
