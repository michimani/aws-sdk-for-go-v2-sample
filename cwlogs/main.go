package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func main() {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	client := cloudwatchlogs.NewFromConfig(config)

	groupName := "sdk-for-go-v2-log-group"
	if err := createLogGroup(client, groupName); err != nil {
		log.Fatalln(err)
		return
	}

	streamName := "sdk-for-go-v2-log-stream-" + time.Now().Format("200601021504")
	if err := createLogStream(client, groupName, streamName); err != nil {
		log.Fatalln(err)
		return
	}

	message := "This is a sample log event message."
	if err := putLogEvent(client, groupName, streamName, message); err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("LogGroup\t: %s\nLogStream\t: %s\nEventMessage\t: \"%s\"", groupName, streamName, message)
}

// createLogGroup creates new log group if it does not exists.
func createLogGroup(client *cloudwatchlogs.Client, name string) error {
	describeIn := cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(name),
	}
	out, err := client.DescribeLogGroups(context.TODO(), &describeIn)
	if err != nil {
		return err
	}

	for _, group := range out.LogGroups {
		if *group.LogGroupName == name {
			return nil
		}
	}

	in := cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(name),
	}

	_, cerr := client.CreateLogGroup(context.TODO(), &in)
	if cerr != nil {
		return cerr
	}

	return nil
}

// createLogStream creates new log stream if it does not exists.
func createLogStream(client *cloudwatchlogs.Client, groupName, streamName string) error {
	describeIn := cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(groupName),
		LogStreamNamePrefix: aws.String(streamName),
	}
	out, err := client.DescribeLogStreams(context.TODO(), &describeIn)
	if err != nil {
		return err
	}

	for _, ls := range out.LogStreams {
		if *ls.LogStreamName == streamName {
			return nil
		}
	}

	createIn := cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(groupName),
		LogStreamName: aws.String(streamName),
	}
	_, cerr := client.CreateLogStream(context.TODO(), &createIn)
	if cerr != nil {
		return cerr
	}

	return nil
}

// putLogEvent puts a log event with message.
func putLogEvent(client *cloudwatchlogs.Client, groupName, streamName, message string) error {
	eventTime := time.Now().UnixNano() / int64(time.Millisecond)
	in := cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(groupName),
		LogStreamName: aws.String(streamName),
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(message),
				Timestamp: aws.Int64(eventTime),
			},
		},
	}

	nextSeqToken, err := getNextSequeneToken(client, groupName, streamName)
	if err != nil {
		return err
	}

	if nextSeqToken != "" {
		in.SequenceToken = aws.String(nextSeqToken)
	}

	_, perr := client.PutLogEvents(context.TODO(), &in)
	if perr != nil {
		return perr
	}

	return nil
}

// getNextSequeneToken return UploadSequenceToken of log stream if it exists.
func getNextSequeneToken(client *cloudwatchlogs.Client, groupName, streamName string) (string, error) {
	describeIn := cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(groupName),
		LogStreamNamePrefix: aws.String(streamName),
	}
	out, err := client.DescribeLogStreams(context.TODO(), &describeIn)
	if err != nil {
		return "", err
	}

	if len(out.LogStreams) == 0 {
		return "", nil
	}

	for _, ls := range out.LogStreams {
		if *ls.LogStreamName == streamName {
			if ls.UploadSequenceToken != nil {
				return *ls.UploadSequenceToken, nil
			} else {
				return "", nil
			}
		}
	}

	return "", nil
}
