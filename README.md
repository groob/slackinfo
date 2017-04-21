slackinfo is a utility to print information about a slack team.
Requires an API token: https://api.slack.com/custom-integrations/legacy-tokens

Open an issue or pull request if you'd like additional functionality.

# Usage

```
./slackinfo -api.token=xoxp-abcd-secret-token

Name               Creator     CreatedDate                    NumMembers  Purpose
```

For `csv` output, use the `-csv` flag.
```
./slackinfo -api.token=xoxp-abcd-secret-token -csv > output.csv
```



