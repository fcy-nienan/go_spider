module go_spider

go 1.24.1

require (
	github.com/PuerkitoBio/goquery v1.10.3
	github.com/fcy-nienan/go_mq/mq_client v1.1.1
	github.com/fcy-nienan/go_mq/mq_server v1.1.1
	github.com/go-resty/resty/v2 v2.16.5
	github.com/lib/pq v1.10.9
)

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/fcy-nienan/go_mq/mq_common v1.1.0 // indirect
	golang.org/x/net v0.41.0 // indirect
)
replace github.com/fcy-nienan/go_mq/mq_client => ../go_mq/mq_client
replace github.com/fcy-nienan/go_mq/mq_server => ../go_mq/mq_server