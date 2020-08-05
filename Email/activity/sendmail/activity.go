package sendmail

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLog is the default logger for the SendMail Activity
var activityLog = logger.GetLogger("general-activity-sendmail")

const (
	ivServer     = "Server"
	ivConnType   = "Connection Security"
	ivPort       = "Port"
	ivUserName   = "Username"
	ivPassword   = "Password"
	ivSender     = "sender"
	ivRecipients = "recipients"
	ivMessage    = "message"
	ivSubject    = "subject"
)

type SendMailActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &SendMailActivity{metadata: metadata}
}

//func init() {
//	md := activity.NewMetadata(jsonMetadata)
//	activity.Register(&SendMailActivity{metadata: md})
//}

// Metadata returns the activity's metadata
func (a *SendMailActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *SendMailActivity) Eval(context activity.Context) (done bool, err error) {
	activityName := context.TaskName()
	server := context.GetInput(ivServer)
	if server == nil || server.(string) == "" {
		return false, activity.NewError(fmt.Sprintf("SMTP Server host must be configured for SendMail Activity [%s] in Flow[%s].", activityName, context.FlowDetails().Name()), "", nil)
	}

	portValue := context.GetInput(ivPort)
	if portValue == nil {
		return false, activity.NewError(fmt.Sprintf("SMTP Server port must be configured for SendMail Activity[%s] in Flow[%s].", activityName, context.FlowDetails().Name()), "", nil)
	}

	port, err := data.CoerceToValue(portValue, data.TypeString)
	if err != nil {
		return false, activity.NewError(fmt.Sprintf("Invalid SMTP Server port [%s] configured for SendMail Activity[%s] in Flow[%s].", portValue, activityName, context.FlowDetails().Name()), "", nil)
	}

	sender := context.GetInput(ivSender)
	if sender == nil || sender.(string) == "" {
		sender = context.GetInput(ivUserName)
		if sender == nil || sender.(string) == "" {
			return false, activity.NewError(fmt.Sprintf("Sender or UserName must be configured for SendMail Activity[%s] in Flow[%s].", activityName, context.FlowDetails().Name()), "", nil)
		}

	}

	recipients := context.GetInput(ivRecipients)
	if recipients == nil || recipients.(string) == "" {
		return false, activity.NewError(fmt.Sprintf("One or more recipients must be configured for SendMail Activity[%s] in Flow[%s].", activityName, context.FlowDetails().Name()), "", nil)
	}

	for _, recp := range strings.Split(recipients.(string), ",") {
		_, err := mail.ParseAddress(recp)
		if err != nil {
			return false, activity.NewError(fmt.Sprintf("Invalid Email Address[%s]. A valid Email address must be configured for SendMail Activity[%s] in Flow[%s].", recp, activityName, context.FlowDetails().Name()), "", nil)
		}
	}

	message := context.GetInput(ivMessage)
	msg := ""
	if message != nil {
		msg = message.(string)
	}

	subject := context.GetInput(ivSubject)
	sbj := ""
	if subject != nil {
		sbj = subject.(string)
	}

	userName := context.GetInput(ivUserName).(string)
	password := context.GetInput(ivPassword).(string)

	connectionType := context.GetInput(ivConnType)
	if connectionType == nil {
		connectionType = "TLS"
	}

	if strings.HasSuffix(password, "=") {
		data, err := base64.StdEncoding.DecodeString(password)
		if err == nil {
			password = string(data)
		}
	}

	activityLog.Debugf("Client Connection Type - [%s]", connectionType)

	if connectionType == "TLS" {
		//Use authentication

		recipinets1 := strings.TrimRight(recipients.(string), ",")
		msg1 := []byte("To: " + getRecipentsStr(recipinets1) + "\r\n" +
			"Subject: " + sbj + "\r\n" +
			"\r\n" +
			msg + "\r\n")
		err = smtp.SendMail(
			server.(string)+":"+port.(string),
			getAuth(userName, password, server.(string)),
			sender.(string),
			getRecipents(recipinets1),
			msg1,
		)
		if err != nil {
			if err.Error() == "EOF" {
				activityLog.Errorf("Failed to send email due to error - [%s]. Try SSL connection type.", err.Error())
				return false, err
			} else {
				activityLog.Errorf("Failed to send email due to error - [%s]", err.Error())
				return false, err
			}
		}
	} else if connectionType == "SSL" {
		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         server.(string),
		}

		conn, err := tls.Dial("tcp", server.(string)+":"+port.(string), tlsconfig)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "tls") {
				activityLog.Errorf("Failed to connect to the mail server due to error - [%s]. Try TLS connection type.", err.Error())
				return false, err
			} else {
				activityLog.Errorf("Failed to connect to the mail server due to error - [%s]", err.Error())
				return false, err
			}
		}

		c, err := smtp.NewClient(conn, server.(string))
		if err != nil {
			activityLog.Errorf("Failed to create a mail client due error - [%s]", err.Error())
			return false, err
		}

		defer c.Quit()
		// Auth

		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(getAuth(userName, password, server.(string))); err != nil {
				activityLog.Errorf("Failed to authenticate with the mail server due to error - [%s]", err.Error())
				return false, err
			}
		}

		// To && From
		if err = c.Mail(sender.(string)); err != nil {
			activityLog.Errorf("Failed to send email due to error - [%s]", err.Error())
			return false, err
		}

		for _, recp := range strings.Split(recipients.(string), ",") {
			if recp == "" {
				continue
			}
			if err = c.Rcpt(recp); err != nil {
				activityLog.Errorf("Failed to send email due to error - [%s]", err.Error())
				return false, err
			}
		}

		msg1 := []byte("From: " + sender.(string) + "\r\n" +
			"To: " + recipients.(string) + "\r\n" +
			"Subject: " + sbj + "\r\n" +
			"\r\n" +
			msg + "\r\n")

		// Data
		w, err := c.Data()
		if err != nil {
			activityLog.Errorf("Failed to send email due to error - [%s]", err.Error())
			return false, err
		}

		_, err = w.Write([]byte(msg1))
		if err != nil {
			activityLog.Errorf("Failed to send email due to error - [%s]", err.Error())
			return false, err
		}

		err = w.Close()
		if err != nil {
			activityLog.Errorf("Failed to send email due to error - [%s]", err.Error())
			return false, err
		}

	} else {
		client, err := smtp.Dial(server.(string) + ":" + port.(string))
		if err != nil {
			activityLog.Errorf("Failed to connect to the mail server due to error - [%s].", err.Error())
			return false, err
		}
		defer client.Close()
		// Set the sender and recipient.

		if err = client.Mail(sender.(string)); err != nil {
			activityLog.Errorf("Failed to send email due to error - [%s]", err.Error())
			return false, err
		}
		for _, recp := range strings.Split(getRecipentsStr(recipients.(string)), ",") {
			if recp == "" {
				continue
			}
			if err = client.Rcpt(recp); err != nil {
				activityLog.Errorf("Failed to send email due to error - [%s]", err.Error())
				return false, err
			}
		}
		// Send the email body.
		wc, err := client.Data()
		if err != nil {
			return false, err
		}
		defer wc.Close()
		buf := bytes.NewBufferString("To: " + strings.Replace(getRecipentsStr(recipients.(string)), ",", ";", -1) + "\r\n" +
			"Subject: " + sbj + "\r\n" +
			"\r\n" +
			msg + "\r\n")
		if _, err = buf.WriteTo(wc); err != nil {
			activityLog.Errorf("Failed to send email due to error - [%s]", err.Error())
			return false, err
		}

	}
	activityLog.Info("Mail successfully sent")

	return true, nil
}

func getRecipents(recipinets string) []string {
	newRecipinets := strings.Split(recipinets, ",")
	var finalRecips []string
	for _, rec := range newRecipinets {
		if rec != "" {
			finalRecips = append(finalRecips, strings.TrimSpace(rec))
		}
	}
	return finalRecips
}

func getRecipentsStr(recipinets string) string {
	newRecipinets := getRecipents(recipinets)
	return strings.Join(newRecipinets, ",")
}

func getAuth(userName, password, server string) smtp.Auth {
	if strings.Contains(server, "smtp.office365.com") {
		return office365Auth(userName, password)
	}
	return smtp.PlainAuth(
		"",
		userName,
		password,
		server,
	)
}

type loginAuth struct {
	username, password string
}

func office365Auth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}
