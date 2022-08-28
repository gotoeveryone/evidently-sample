package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/evidently"
	"github.com/aws/aws-sdk-go-v2/service/evidently/types"
)

var (
	entityID    = flag.String("user", "dummy", "Specify user")
	subCommands = map[string]func() error{
		"hoge": func() error {
			fmt.Println("Hello hoge!!")
			return nil
		},
	}
)

func enableCommand(c *evidently.Client, entityID string, command string) (bool, error) {
	resp, err := c.EvaluateFeature(context.TODO(), &evidently.EvaluateFeatureInput{
		EntityId: aws.String(entityID),
		Project:  aws.String("evidently-sample"),
		Feature:  aws.String(fmt.Sprintf("enable-%s-command", command)),
	})

	if err != nil {
		return false, err
	}

	return resp.Value.(*types.VariableValueMemberBoolValue).Value, nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := evidently.NewFromConfig(cfg)
	if len(os.Args) > 1 {
		flag.CommandLine.Parse(os.Args[2:])

		sc := os.Args[1]
		fn, ok := subCommands[sc]
		if !ok {
			log.Fatalf("invalid command: %s", sc)
		}
		if flag, err := enableCommand(svc, *entityID, sc); err != nil {
			log.Fatal(err)
		} else if !flag {
			log.Fatalf("invalid command: %s", sc)
		}
		if err := fn(); err != nil {
			log.Fatal(err)
		}
	}
}
