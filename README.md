# retarus-go
The offical Golang SDK provided by Retarus to contact our messaging services.



# Quickstart
## Send an example fax job
```golang
package example

import (
	github.com/retarus/retarus-go/fax
)

func main() {
	client := fax.Client{
		Config: fax.Config{
			User:           "your-user@mail.de",
			Password:       "your_private_password",
			CustomerNumber: "your-customer-number",
			Endpoint:       fax.DE,
		},
		HTTPClient: http.Client{Timeout: 5 * time.Second},
	}

	job := fax.Job{
		Recipients: []fax.Recipient{
			{
				Number: "004989000000000", // number to send to
			},
		},
		Documents: []fax.Document{
			{
				Name:      "test.txt", // local document to send
				Charset:   fax.UTF_8,
				Reference: "testJob",
				Data:      "dGVzdGZheAo=", // testfax
			},
		},
	}

	jobID, err := client.Send(job)
	if err != nil {
		// log error message
	}
}
```

## Send an example sms job
```golang
package example

import (
	github.com/retarus/retarus-go/sms
)

func main() {
	client := sms.Client{
		Config: sms.Config{
			User:     "your-user@mail.de",
			Password: "your_private_password",
			Endpoint: sms.DE1,
		},
		HTTPClient: http.Client{Timeout: 5 * time.Second},
	}

	job := sms.Job{
		Messages: []sms.Message{
			{
				Text: "this is a test message",
				Recipients: []sms.Recipient{
					{
						Dst:         "0049176000000000", // number to send to
						CustomerRef: "retarus",
					},
				},
			},
		},
	}

	jobID, err := client.Send(job)
	if err != nil {
		// log error message
	}
}
```



# Examples
Further examples can be found in ```fax/example_test.go``` and ```sms/example_test.go```

# Help 
To get additional help visit our [knowledge center](https://developers.retarus.com/)
