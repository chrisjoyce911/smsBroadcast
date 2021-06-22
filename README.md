# smsBroadcast

## Introduction

The SMS Broadcast Advanced HTTP API can be used to integrate SMS messaging into your own systems. This API allows for tracking of SMS messages and includes functions for receiving inbound message replies.

SMS Broadcast also offers a simple HTTP API. Details of this can be found in the Documentation section of the SMS Broadcast Client Portal.

### Sending SMS Messages

The send messages function can be used to send a single message or multiple messages in a single API request.

### Message Status Updates

Status updates allow you to keep track of the delivery of the messages you send. A confirmation will be sent to you to show if the message was successfully delivered or failed.
This option must be activated by SMS Broadcast. If you require this function, please contact us and provide us with the URL on your server to send the data.
The SMS Broadcast API will send a GET request to your URL (as provided when setup) with the below parameters

### Inbound SMS Messages

Our API will send a request to your URL with the inbound SMS data for each message.
This option must be activated by SMS Broadcast. If you require this function, please contact us and provide us with the URL on your server to send the data.
SMS Broadcast will send an HTTP GET request to the URL you provided when this service was activated.

### Account Balance

The account balance function can be used to lookup the number of credits left in your account.

#### Input Parameters

| Parameter | Description |
| action    | This must be set to “balance” for this function. |
| username  | Your SMS Broadcast username. This is the same username that you would use to login to the SMS Broadcast website. |
| password  | Your SMS Broadcast password. This is the same password that you would use to login to the SMS Broadcast website. |

#### Output Parameters

| Parameter | Description |
| Status    | Will show the status of your request. The possible statuses are: OK: This request was accepted. ERROR: There is a problem with the request. |
| Balance or Error | Will show your account balance or a reason for the error. | 
