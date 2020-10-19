package main

import (
	"strconv"
	"time"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"os"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
)

func safeParseFloat32(s string) float32 {
	amount, _ := strconv.ParseFloat(s, 32)
	return float32(amount)
}

func getFirstOfMonth() time.Time{
	now := time.Now()
    currentYear, currentMonth, _ := now.Date()
    currentLocation := now.Location()

    firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

    return firstOfMonth
}

func getDateInterval(day time.Time) costexplorer.DateInterval{
	firstDayDate := day.Format("2006-01-02")
	dayAfterDate := day.AddDate(0, 0, 1).Format("2006-01-02")
	dateInterval := costexplorer.DateInterval{
		Start: &firstDayDate,
		End:   &dayAfterDate,
	}
	return dateInterval
}

func getWebhookURL() string{
	parameter_name := os.Getenv("WEBHOOK_URL_PARAMETER_NAME")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	}))
	ssm_client := ssm.New(sess)
	decrypt := true
	params := &ssm.GetParameterInput{
		Name:				&parameter_name,
		WithDecryption:		&decrypt,
	}

	res, _ := ssm_client.GetParameter(params)

	return *res.Parameter.Value
}