package test

import (
	"io/ioutil"
	"os"

	"github.com/trustedanalytics/tap-cli/api"
)

const testCredentialsPath = "/tmp/testTapCredentials.json"

//SwitchToTestCredentialsFile sets credentials_manager.CredentialsPath to test path
//	so that tests do not overrite regular tap-cli credentials.json file
func SwitchToTestCredentialsFile() {
	api.CredsPath = testCredentialsPath
}
func FillCredentialsTestFile(content string) {
	SwitchToTestCredentialsFile()
	ioutil.WriteFile(testCredentialsPath, []byte(content), api.PERMISSIONS)
}
func ReadCredentialsTestFile() ([]byte, error) {
	SwitchToTestCredentialsFile()
	return ioutil.ReadFile(testCredentialsPath)
}
func DeleteTestCredentialsFile() {
	SwitchToTestCredentialsFile()
	os.Remove(testCredentialsPath)
}
