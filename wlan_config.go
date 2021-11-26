package toold

import (
	"fmt"
	"io/ioutil"
)

//WlanProfile WlanProfile
func WlanProfile(name string, path string) (int, string, error) {
	path = path + "/" + name + ".xml"
	if IsFile(path) {
		return 0, path, nil
	}
	xml := fmt.Sprintf(`<?xml version="1.0"?>
	<WLANProfile xmlns="http://www.microsoft.com/networking/WLAN/profile/v1">
		<name>%v</name>
		<SSIDConfig>
			<SSID>
				<hex>40424852475A4E2E434E</hex>
				<name>%v</name>
			</SSID>
		</SSIDConfig>
		<connectionType>ESS</connectionType>
		<connectionMode>auto</connectionMode>
		<MSM>
			<security>
				<authEncryption>
					<authentication>WPA2PSK</authentication>
					<encryption>AES</encryption>
					<useOneX>false</useOneX>
				</authEncryption>
				<sharedKey>
					<keyType>passPhrase</keyType>
					<protected>true</protected>
					<keyMaterial>01000000D08C9DDF0115D1118C7A00C04FC297EB01000000BB76AE41CF41B540BF7CC8E0E63C37EC00000000020000000000106600000001000020000000BB0D7D9C956E7233B1228831638DAA371BC181D15FE8040F741945DF2E7A554A000000000E80000000020000200000005C0B58297386C75AEEBDD47CFA13600E6B208E6A0D09F3AEF1E0E7F340F781A01000000002E89CBDEA375F0B328723A3ED7ED4BA40000000642442A76A2E6C2BF6AA8B98E702B9128263EA7C7D65F8A7D0B62940168A7036BF07C87449DB5ED79C0C92CE107766232733821DF1C2C778944A8B0667DE22FC</keyMaterial>
				</sharedKey>
			</security>
		</MSM>
		<MacRandomization xmlns="http://www.microsoft.com/networking/WLAN/profile/v3">
			<enableRandomization>false</enableRandomization>
		</MacRandomization>
	</WLANProfile>
	`, name, name)

	err := ioutil.WriteFile(path, []byte(xml), 0666)
	return 1, path, err
}

//IsWlanProfile IsWlanProfile
func IsWlanProfile(name, path string) bool {
	path = path + "" + name + ".xml"
	return IsFile(path)
}
