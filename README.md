# check-pls
Testing GH App for the Checks API

## How to run

### Get a webhook URL

Use [smee](https://smee.io) to get a webhook URL that you can use for testing. Follow the directions on the site to install the smee CLI locally if you haven't already.

### Create a GitHub App
Follow [these](https://docs.github.com/en/developers/apps/building-github-apps/creating-a-github-app) steps to create a GitHub App. Use the URL generated from smee as the webhook URL. Save the changes and make sure to create a private key. Save this file somewhere you'll remember!

### Set up local environment
Clone this repo locally. Copy the `.envconfig` file and name it `.env`. For `APP_GITHUB_APP_INTEGRATION_ID`, you'll use the App ID listed in the about information of your app on GitHub. For `APP_GITHUB_APP_PRIVATE_KEY`, you'll use the contents of that pem file that you created earlier (replace the newlines with `\n` and wrap the key with quotes).

### Run the app locally
Run `make build && make run`.

Run `smee -u <SMEE_URL> -t http://127.0.0.1:8080/api/github/hook`.

Make sure the app is installed to a repo of your choosing. Now you are set up!

## Adding an event handler

In [`cmd/server/handlers.go`](cmd/server/handlers.go), add a struct that implements the [`githubapp.EventHandler` interface](https://github.com/palantir/go-githubapp#usage) with your desired event.
