## webHooks-Svix
Webhooks and Svix Implementation in Python

## Cases to test here

# Case 1
- change the secret key of sender and receiver here and then test
- in this case it should throw error and not retry
- in case of invalid payload i.e. bad request do not retry

# Case 2
- Simulate internal server errors (e.g., 500 status code) in the receiver.
- In this case, the sender should retry the webhook with exponential backoff.

# Current status of the project
- Manual webhook with retries fully implemented
- Integrated svix to handle retries and encryption of the payload
- Calling the svix client, svix dashboard is receiving request
- But app is running locally so svix not able to access the API
- Need to host this app, so that svix app is able to call the app
- livspace email has the svix app