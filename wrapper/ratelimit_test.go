package wrapper

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// TestRateGlobalLimitTypeHandle is an experimental test to determine whether the Global Rate Limit
// uses a bucket-rate limit.
func TestRateGlobalLimitTypeHandle(t *testing.T) {
	bot := &Client{
		ApplicationID:  os.Getenv("APPID"),
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
	}
	request := new(GetCurrentBotApplicationInformation)
	responses := make(chan *Application, 50)
	errs := make(chan error)
	ct := 0

	ticker := time.NewTicker(time.Second)
	for {
		ct++
		fmt.Println("Spawned request routine", ct)
		go func(count int) {
			app, err := request.Send(bot)
			if err != nil {
				errs <- err
			}

			fmt.Println(count, "APP", app)
			responses <- app
		}(ct)

		select {
		case <-ticker.C:
			goto TICK
		case e := <-errs:
			t.Fatalf("%v", e)
		default:
			// when 50 requests have been sent,
			if ct == 50 {

				// wait until all 50 requests have received a response.
				for {
					// when all 50 requests have received a response,
					if len(responses) == cap(responses) {
						// wait 1/50th of a second and see whether a request can be sent.
						<-time.NewTimer(time.Second / 50).C
						ct++
						fmt.Println("Sending the 51st request.")
						app, err := request.Send(bot)
						if err != nil {
							fmt.Println("Discord uses a bucket rate limit for global rate limits.")
							t.Fatalf("%v", err)
						}

						fmt.Println(ct, "APP", app)
						fmt.Println("Discord uses a constant rate limit for global rate limits.")
						break
					}

					select {
					case <-ticker.C:
						fmt.Println(len(responses), cap(responses))
						t.Fatalf("Failed to send the 51st request in 1 second.")
					default:
						continue
					}
				}

				// wait for the remainder of the second.
				for {
					select {
					case <-ticker.C:
						goto TICK
					}
				}
			}

			continue
		}
	}

TICK:
	fmt.Println("Sent", ct, "requests in 1 second.")
}
