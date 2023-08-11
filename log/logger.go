package log

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

func Error(ctx context.Context, projectID, logName string, err error) {
	var convertedErr interface{ StackTrace() errors.StackTrace }

	if errors.As(err, &convertedErr) {
		// スタックトレースがある場合
		logPrintln(ctx, projectID, logName, convertedErr)
	} else {
		logPrintln(ctx, projectID, logName, err)
	}
}

func logPrintln(ctx context.Context, projectID, logName string, v interface{}) {
	if len(projectID) == 0 {
		// 通常のログ出力
		log.Printf("%+v\n", v)
	} else {
		client, err2 := logging.NewClient(ctx, projectID)
		if err2 != nil {
			// ログクライアントのインスタンス作成時にエラーが発生した場合
			// インスタンス作成エラーを出力
			log.Println(err2)

			// 通常のログ出力
			log.Printf("%+v\n", v)
		}
		defer client.Close()

		logger := client.Logger(logName).StandardLogger(logging.Error)

		// GCPのフォーマットでログ出力
		logger.Println(fmt.Sprintf("%+v\n", v))
	}
}
