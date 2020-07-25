package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

func testSNS(sess *session.Session, iters int, payload, topic string) error {
	client := sns.New(sess)

	var minDur = time.Duration(math.MaxInt64)
	var maxDur time.Duration

	start := time.Now()
	for i := 0; i < iters; i++ {
		localStart := time.Now()
		input := &sns.PublishInput{
			Message:  &payload,
			TopicArn: &topic,
		}
		if _, err := client.Publish(input); err != nil {
			return err
		}
		reqDur := time.Since(localStart)
		if reqDur > maxDur {
			maxDur = reqDur
		}
		if reqDur < minDur {
			minDur = reqDur
		}
	}
	dur := time.Since(start)
	fmt.Printf("performed %d sns publishes in %v (%v per request, %v max and %v min)\n", iters, dur, dur/time.Duration(iters), maxDur, minDur)
	return nil
}

func testSQS(sess *session.Session, iters int, payload, queue string) error {
	client := sqs.New(sess)

	var minDur = time.Duration(math.MaxInt64)
	var maxDur time.Duration

	start := time.Now()
	for i := 0; i < iters; i++ {
		localStart := time.Now()
		input := &sqs.SendMessageInput{
			MessageBody: &payload,
			QueueUrl:    &queue,
		}
		if _, err := client.SendMessage(input); err != nil {
			return err
		}
		reqDur := time.Since(localStart)
		if reqDur > maxDur {
			maxDur = reqDur
		}
		if reqDur < minDur {
			minDur = reqDur
		}
	}
	dur := time.Since(start)
	fmt.Printf("performed %d sqs send in %v (%v per request, %v max and %v min)\n", iters, dur, dur/time.Duration(iters), maxDur, minDur)
	return nil
}

func HandleRequest() error {
	iters, err := strconv.Atoi(os.Getenv("ITERATIONS_TO_PERFORM"))
	if err != nil {
		return err
	}
	payloadLen, err := strconv.Atoi(os.Getenv("PAYLOAD_LENGTH"))
	if err != nil {
		return err
	}
	payload := make([]byte, payloadLen)
	if _, err := rand.Read(payload); err != nil {
		return err
	}
	for i := 0; i < payloadLen; i++ {
		char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		payload[i] = char[int(payload[i])%len(char)]
	}
	start := time.Now()
	rawHash := sha256.Sum256(payload)
	payloadHash := hex.EncodeToString(rawHash[:])
	fmt.Printf("payload's hash is %s (calculated in %s)\n", payloadHash, time.Since(start))

	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	snsTopic := os.Getenv("SNS_TOPIC")
	queueURL := os.Getenv("SQS_QUEUE")

	fmt.Printf("sns test begins (%d iters)\n", iters)
	if err := testSNS(sess, iters, string(payload), snsTopic); err != nil {
		return err
	}
	fmt.Printf("sqs test begins (%d iters)\n", iters)
	if err := testSQS(sess, iters, string(payload), queueURL); err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
