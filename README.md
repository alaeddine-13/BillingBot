# Billing Bot
## Overview

Golang based billing bot.
This bot retrieves AWS costs consumed by a project using the Cost Explorer API and when threshold is exceeded, sends notifications to a slack channel.
This app is built using the Serverless framework, deployed using AWS Lambda, built using Golang in Go modules mode.

## Installation
* Install Golang:
```bash
wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.14.4.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```
* Install serverless framework: ```npm install -g serverless```
* Clone the repository
* Go under `src/billing` and build the project: ```make```
* You can run the serverless project locally: ```sudo serverless invoke local -f billing```
* You can manually deploy the project: ``` sudo serverless deploy```

## Usage:
The project consists in a lambda function writting with golang which uses the Cost Explorer API to retrieve the costs of resources located in the `eu-central-1` region and belonging to services `ecs` and `ecr`.
The cost is reported to the console. If threshold is exceeded, a report of cost is sent to the slack channel.
When deployed, the project runs on a daily basis as specified in the cron schedule, configured in the serverless.yml file. The log outputs can be viewed in the cloudwatch console under log group `/aws/lambda/billing-bot-prod-billing`.



