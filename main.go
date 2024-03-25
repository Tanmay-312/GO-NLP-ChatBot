package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
)

//var wolframClient *wolfram.Client

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Event:")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}

	go printCommandEvents(bot.CommandEvents())

	bot.Command("query - <message>", &slacker.CommandDefinition{
		Description: "send any question to wolfram",
		Examples:    []string{"who is the president of India", "who is the prime minister of India"},
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			query := r.Param("message")
			//fmt.Println(query)
			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			//fmt.Println(msg)

			data, _ := json.MarshalIndent(msg, "", "    ")
			//fmt.Println(data)
			rough := string(data[:])
			//fmt.Println(rough)
			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			answer := value.String()
			//fmt.Println(answer)
			res, err := wolframClient.GetSpokentAnswerQuery(answer, wolfram.Metric, 1000)

			if err != nil {
				fmt.Println("Error in wolfram search query")
			}

			w.Reply(res)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}

}
