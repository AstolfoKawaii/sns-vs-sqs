# SNS vs SQS
It is a simple lambda function to benchmark Amazon's (AWS) SNS publish and SQS send performance.

To perform benchmarks you need:
* Create SNS topic
* Create SQS queue
* Create lambda function
* Grant access for function to SNS topic and SQS queue
* Set setup lambda's properties
  * Upload zipped executable (and set the handler name)
  * Configure available memory
  * Set environment variables:
    * ITERATIONS_TO_PERFORM - how many messages to send/publish
    * PAYLOAD_LENGTH - size of the payload in bytes
    * SNS_TOPIC - sns topic arn
    * SQS_QUEUE - sqs queue url
* Run the function and view logs in CloudWatch
