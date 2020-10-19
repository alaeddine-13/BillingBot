package main

import (
	"log"
	"fmt"
	"os"
	"github.com/aws/aws-lambda-go/lambda"
)



// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler() error {
	monthlythresholdString := os.Getenv("MONTHLY_THRESHOLD")
	montlythreshold := safeParseFloat32(monthlythresholdString)


	costResult, err := GetCostResult()
	if err != nil {
		fmt.Printf("Error getting this month's cost result: %v", err)
		return err
	}
	
	log.Printf("Monthly costs so far: $%f", costResult.Total)
	if costResult.Total > montlythreshold{
		reportMonthlyThresholdExceeded(costResult, montlythreshold)
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
