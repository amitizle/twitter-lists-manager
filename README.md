# Twitter Lists Manager

Manage Twitter lists from the CLI, using JSON files as a source of truth.

## Install

Either `go get -u github.com/amitizle/twitter-lists-manager`, or if you're using [Homebrew](https://brew.sh/) you can use my own `tap`:

```bash
$ brew tap amitizle/tap
$ brew update
$ brew install twitter-lists-manager
```

The _Homebrew_ binary is versioned (i.e. only updated on a new version of the tool, and not from `master`).

Also all binaries are released (using [goreleaser](https://github.com/goreleaser/goreleaser)) to the [releases page](https://github.com/amitizle/twitter-lists-manager/releases)
for both _Linux_ and _MacOS_.

## Usage

In order to use this tool you'd need to apply for a dev account on _Twitter_ ([here](https://developer.twitter.com/en/apply-for-access)) and add a new app.
Then you'll receive: `access_token`, `access_token_secret`, `consumer_key` and `consumer_key_secret`.

Use those with the flags (run `twitter-lists-manager --help` to see the exact flags), or use environment variables:

```
$ export TWITTER_ACCESS_TOKEN="ACCESS_TOKEN"
$ export TWITTER_ACCESS_TOKEN_SECRET="SECRET"
$ export TWITTER_CONSUMER_KEY="CONSUMER_KEY"
$ export TWITTER_CONSUMER_SECRET="CONSUMER_SECRET"
```

Then run the binary as you would any other binary. Run `--help` for more information.
