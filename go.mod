module transParserTest

go 1.14

require (
	github.com/fatih/astrewrite v0.0.0-20191207154002-9094e544fcef
	github.com/pingcap/check v0.0.0-20190102082844-67f458068fc8
	github.com/pingcap/errors v0.11.4
	github.com/pingcap/log v0.0.0-20200828042413-fce0951f1463 // indirect
	github.com/pingcap/parser v0.0.0-20200902143951-126c14c456eb
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20200904185747-39188db58858 // indirect
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
)

replace github.com/pingcap/parser => github.com/xiongjiwei/parser v0.0.0-20200908003518-25ce2f61fa9c
