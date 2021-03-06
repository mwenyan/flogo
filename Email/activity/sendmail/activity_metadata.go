package sendmail

import (
	"github.com/project-flogo/legacybridge"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var jsonMetadata = `{
	"title": "Send Mail",
	"name": "tibco-wi-sendmail",
	"author": "TIBCO Software Inc.",
	"version": "1.1.0",
	"type": "flogo:activity",
	"display": {
		"visible": true,
		"description": "Simple Send Mail Activity",
		"category": "General",
		"smallIcon": "icons/ic-tibco-wi-sendmail.svg",
		"largeIcon": "icons/ic-tibco-wi-sendmail@2x.png",
		"appPropertySupport": true
	},
	"ref": "git.tibco.com/git/product/ipaas/wi-contrib.git/contributions/General/activity/sendmail",
	"inputs": [
		{
			"name": "Server",
			"type": "string",
			"display": {
				"description": "Host name or IP address for the mail server",
				"name": "Server",
				"appPropertySupport": true
			},
			"required": true,
			"value": ""
		},
		{
			"name": "Port",
			"type": "integer",
			"required": true,
			"display": {
				"name": "Port",
				"description": "The port used to connect to the mail server",
				"appPropertySupport": true
			}
		},
		{
			"name": "Connection Security",
			"type": "string",
			"required": true,
			"allowed": [
				"TLS",
				"SSL",
				"NONE"
			],
			"value": "TLS",
			"display": {
				"description": "The type of connection to be used to communicate with the mail server.",
				"name": "Connection Security"
			}
		},
		{
			"name": "Username",
			"type": "string",
			"required": true,
			"value": "",
			"display": {
				"description": "The username to use when authenticating to the mail server",
				"name": "Username",
				"appPropertySupport": true
			}
		},
		{
			"name": "Password",
			"type": "string",
			"required": true,
			"value": "",
			"display": {
				"description": "The password to use when authenticating to the mail server",
				"name": "Password",
				"type": "password",
				"appPropertySupport": true
			}
		},
		{
			"name": "sender",
			"type": "string",
			"value": ""
		},
		{
			"name": "recipients",
			"type": "string",
			"required": true,
			"value": ""
		},
		{
			"name": "subject",
			"type": "string",
			"value": ""
		},
		{
			"name": "message",
			"type": "string",
			"value": ""
		}
	]
}`

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	legacybridge.RegisterLegacyActivity(NewActivity(md))
}
