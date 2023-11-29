package main

import (
    b64 "encoding/base64"
    "context"
    "encoding/json"
    "fmt"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ses"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "os"
)

type ContactForm struct {
    From    string  `json:"from"`
    Email   string  `json:"email"`
    Message string  `json:"message"`
}

const (
    Subject = "Email from Website Form"
    CharSet = "UTF-8"
)

func handler(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

    body := []byte(event.Body)
    bodyDecoded, _ := b64.StdEncoding.DecodeString(string(body))
    var form ContactForm
    json.Unmarshal(bodyDecoded, &form)

    sess, err := session.NewSession(&aws.Config{
        Region:aws.String("eu-west-2")},
    )

    svc := ses.New(sess)

    EmailTo := os.Getenv("EMAIL_TO")
    EmailFrom := os.Getenv("EMAIL_FROM")
    PlainBody := "Email From: "+form.From+" with email "+form.Email+"\r\n\r\n"+form.Message
    HtmlBody := "<b>"+form.From+"</b> with email "+form.Email+" says <br /><br />"+form.Message

    input := &ses.SendEmailInput{
        Destination: &ses.Destination{
            CcAddresses: []*string{
            },
            ToAddresses: []*string{
                aws.String(EmailTo),
            },
        },
        Message: &ses.Message{
            Body: &ses.Body {
                Html: &ses.Content {
                    Charset:    aws.String(CharSet),
                    Data:       aws.String(HtmlBody),
                },
                Text: &ses.Content {
                    Charset:    aws.String(CharSet),
                    Data:       aws.String(PlainBody),
                },
            },
            Subject: &ses.Content{
                Charset:    aws.String(CharSet),
                Data:       aws.String(Subject),
            },
        },
        Source: aws.String(EmailFrom),
    }

    result, err := svc.SendEmail(input)

 // Display error messages if they occur.
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case ses.ErrCodeMessageRejected:
                fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
            case ses.ErrCodeMailFromDomainNotVerifiedException:
                fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
            case ses.ErrCodeConfigurationSetDoesNotExistException:
                fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
        return events.LambdaFunctionURLResponse{
            StatusCode: 400,
            Body:       "There was an issue with the request",
        }, err
    }

    fmt.Println(result)

    response := events.LambdaFunctionURLResponse{
        StatusCode: 200,
        Body:       "Sent",
    }
    return response, nil
}


func main() {
    lambda.Start(handler)
}
