# A BOT TO MAKE YOU GO ALL CAPS

## SLACK TOKEN

USED [THIS GUIDE](https://github.com/threatsimple/smug-broker/blob/master/doc/slack.md)
TO OBTAIN A SLACK TOKEN AND GET THE "APP" INSTALLED.

## GOLANG SLACK LIBRARY

HAD TO RUN THE FOLLOWING TO MAKE THE SLACK LIBRARY PLAY NICE
(CURRENT VERSION PICKED UP BY GO CONTAINS [THIS BUG](https://github.com/nlopes/slack/pull/618)):

```bash
go get github.com/nlopes/slack@d06c2a2b3249b44a9c5dee8485f5a87497beb9ea
```

## RUN

```bash
SLACKTOKEN=YOURTOKENHERE go run main.go
```