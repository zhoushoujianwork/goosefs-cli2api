package dingtalk

import (
	"context"
	"goosefs-cli2api/config"

	"github.com/xops-infra/noop/log"

	dt "github.com/xops-infra/go-dingtalk-sdk-wrapper"
)

var robotClient *dt.RobotClient

func init() {
	robotClient = dt.NewRobotClient()
}

func SendAlert(msg string) {
	if config.Config.DingtalkAlert == nil || config.Config.DingtalkAlert.Token == "" {
		log.Errorf("dingtalk token is empty, skip send dingtalk msg")
		return
	}
	err := robotClient.SendMessage(context.Background(), &dt.SendMessageRequest{
		AccessToken: config.Config.DingtalkAlert.Token,
		MessageContent: dt.MessageContent{
			MsgType: "text",
			Text: dt.TextBody{
				Content: msg,
			},
		},
	})
	if err != nil {
		log.Errorf("ERROR: send dingtalk msg failed: %v", err)
	}
	log.Infof("INFO: send dingtalk msg: %s", msg)
}
