package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)
type DayResult struct {
	Groups []*costexplorer.Group
	TimePeriod costexplorer.DateInterval
}
func (DayResult DayResult) String() string {
    s := "{\nGroups:\t"
    s += fmt.Sprintf("%+v", DayResult.Groups)
    s += "\nTimePeriod:\t"
    s += fmt.Sprintf("%+v", DayResult.TimePeriod)
    return s + "}"
}

type CostResult struct {
	Total      float32 `jsong:"total"`
	Days     []*DayResult
}

// GetDayCostResult returns the estimated cost result for the given day.
func GetCostResult() (*CostResult, error) {
	sess := session.Must(session.NewSession())
	explorer := costexplorer.New(sess)

	// Only today's report

	// daily granularity
	granularity := costexplorer.GranularityDaily

	// We are grouping based on dimensions not tags
	dimensionType := costexplorer.GroupDefinitionTypeDimension

	serviceDimension := costexplorer.DimensionService
	serviceGroup := &costexplorer.GroupDefinition{
		Key:  &serviceDimension,
		Type: &dimensionType,
	}

	regionDimension := costexplorer.DimensionRegion
	regionGroup := &costexplorer.GroupDefinition{
		Key:  &regionDimension,
		Type: &dimensionType,
	}

	groupBy := []*costexplorer.GroupDefinition{serviceGroup, regionGroup}

	// Get both BlendedCost and UnblendedCost metrics
	blendedCost := costexplorer.MetricBlendedCost
	unblendedCost := costexplorer.MetricUnblendedCost
	metrics := []*string{&blendedCost, &unblendedCost}

	
	ecs_service := "Amazon EC2 Container Service"
	ecr_service := "Amazon EC2 Container Registry (ECR)"
	elb_service := "Elastic Load Balancing (ELB)"
	service_code_values := []*string{&ecs_service, &ecr_service, &elb_service}
	service_code_dimension := costexplorer.DimensionService
	service_code_dimension_values := &costexplorer.DimensionValues{
		Key:  &service_code_dimension,
		Values: service_code_values,
	}
	service_code_expression := &costexplorer.Expression{
		Dimensions:	service_code_dimension_values,
	}

	frankfurt_region_code := "eu-central-1"
	region_values := []*string{&frankfurt_region_code}
	region_dimension := costexplorer.DimensionRegion
	region_dimension_values := &costexplorer.DimensionValues{
		Key:	&region_dimension,
		Values:	region_values,
	}
	region_expression := &costexplorer.Expression{
		Dimensions:	region_dimension_values,
	}


	dimension_filters := []*costexplorer.Expression{service_code_expression, region_expression}

	and_filter := &costexplorer.Expression{
		And:	dimension_filters,
	}

	var total float32
	var days []*DayResult


	firstDay := getFirstOfMonth()
	lastDay := time.Now()
	for curr := firstDay; curr.Before(lastDay); curr = curr.AddDate(0, 0, 1) {


		var groups []*costexplorer.Group
		dateInterval := getDateInterval(curr)
		params := &costexplorer.GetCostAndUsageInput{
			Granularity: &granularity,
			TimePeriod:  &dateInterval,
			Metrics:     metrics,
			GroupBy:     groupBy,
			Filter: 	 and_filter,
		}

		for {
			result, err := explorer.GetCostAndUsage(params)
			if err != nil {
				return nil, fmt.Errorf("Error getting cost and usage report: %v", err)
			}

			if len(result.ResultsByTime) == 0 {
				return nil, fmt.Errorf("Empty result when getting cost and usage report")
			}

			for _, g := range result.ResultsByTime[0].Groups {
				// For some reasons costexplorer.MetricBlendedCost is equal to `BLENDED_COST`,
				// but the keys used in the result is `BlendedCost`.
				amount := safeParseFloat32(*g.Metrics["BlendedCost"].Amount)
				total += amount
				groups = append(groups, g)
			}

			// There is no next page, break from the loop.
			if result.NextPageToken == nil {
				break
			}

			// There are more pages, get them.
			params.NextPageToken = result.NextPageToken
		}
		day := DayResult{
			Groups:			groups,
			TimePeriod:		dateInterval,
		}
		days = append(days, &day)


	}

	

	return &CostResult{
		Total:      total,
		Days:		days,
	}, nil
}
