package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/whisper"
)

const (
	tcp = "tcp://"

	defaultNode  = "127.0.0.1:26657"
	defaultTopic = "12345678"
)

var (
	flagNode  string
	flagTopic string
)

func main() {
	cobra.EnableCommandSorting = false

	rootCMD := &cobra.Command{
		Use:   "tm-whisper",
		Short: "Tendermint Whisper Client",
	}

	rootCMD.AddCommand(
		pubEnvelopeCmd(),
	)

	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func pubEnvelopeCmd() *cobra.Command {
	pubCMD := &cobra.Command{
		Use: "pub-envelope",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("this command requires a string message")
			}
			return nil
		},
		Short: "Publish Envelope To the tendermint",
		RunE:  EnvelopePublishCMD,
	}

	pubCMD.Flags().StringVar(&flagNode, "node", defaultNode, "Connect to the tendermint node at this address. (default\"localhost:26657\")")
	pubCMD.Flags().StringVar(&flagTopic, "topic", defaultTopic, `tm-whisper pub-envelope --topic=topic`)

	return pubCMD
}

func EnvelopePublishCMD(cmd *cobra.Command, args []string) error {
	payload := args[0]
	node := tcp + flagNode
	client, err := http.New(node, "/websocket")
	if err != nil {
		return errors.New("connect to tendermint error")
	}

	topicString, err := hex.DecodeString(flagTopic)
	if err != nil {
		return errors.New("topic param is invalid")
	}

	topic := whisper.BytesToTopic(topicString)

	env := coretypes.Envelope{
		TTL:   whisper.DefaultTTL,
		Topic: topic,
		Data:  []byte(payload),
	}

	res, err := client.PublishEnvelope(context.Background(), env)
	if err != nil {
		return err
	}

	fmt.Printf("EnvelopePublish success, hash: %v\n", res.Hash)

	return nil
}
