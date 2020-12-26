aws-sdk-for-go-v2-sample
===

This repository is a collection of sample implementations using the AWS SDK for Go v2.

## Services

### CloudWatch Logs

- Create LogGroup if it does not exists
- Create LogStream if it does not exists
- Put a LogEvent

#### Execution sample

```bash
$ go run cwlogs/main.go
LogGroup:       sdk-for-go-v2-log-group
LogStream:      sdk-for-go-v2-log-stream-202012261538
EventMessage:   "This is a sample log event message."
```

### DynamoDB

- Create table if it does not exists
- Put items
- Scan table
- Scan table with specific attributes


#### Execution sample

```bash
$ go run dynamodb/main.go
Sample Item 1
Sample Item 2
Sample Item 3
Sample Item 4
Sample Item 5
TableName: sdk-for-go-sample-table
Scan with all attributes:
[
 {
  "name": "Sample Item 5",
  "created_at": 1608971155521931000,
  "message1": "This is a sample message 5-1",
  "message2": "This is a sample message 5-2",
  "message3": "This is a sample message 5-3"
 },
 {
  "name": "Sample Item 2",
  "created_at": 1608971155474955000,
  "message1": "This is a sample message 2-1",
  "message2": "This is a sample message 2-2",
  "message3": "This is a sample message 2-3"
 },
 {
  "name": "Sample Item 4",
  "created_at": 1608971155506741000,
  "message1": "This is a sample message 4-1",
  "message2": "This is a sample message 4-2",
  "message3": "This is a sample message 4-3"
 },
 {
  "name": "Sample Item 1",
  "created_at": 1608971155455575000,
  "message1": "This is a sample message 1-1",
  "message2": "This is a sample message 1-2",
  "message3": "This is a sample message 1-3"
 },
 {
  "name": "Sample Item 3",
  "created_at": 1608971155490865000,
  "message1": "This is a sample message 3-1",
  "message2": "This is a sample message 3-2",
  "message3": "This is a sample message 3-3"
 }
]

˜Scan with some attributes:
[
 {
  "name": "Sample Item 5",
  "created_at": 1608971155521931000,
  "message1": "This is a sample message 5-1"
 },
 {
  "name": "Sample Item 2",
  "created_at": 1608971155474955000,
  "message1": "This is a sample message 2-1"
 },
 {
  "name": "Sample Item 4",
  "created_at": 1608971155506741000,
  "message1": "This is a sample message 4-1"
 },
 {
  "name": "Sample Item 1",
  "created_at": 1608971155455575000,
  "message1": "This is a sample message 1-1"
 },
 {
  "name": "Sample Item 3",
  "created_at": 1608971155490865000,
  "message1": "This is a sample message 3-1"
 }
]
```