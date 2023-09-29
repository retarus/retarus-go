## Retarus-Go: 
The Official Golang Library for Retarus Messaging Services
Retarus-Go allows you to easily integrate Retarus's suite of messaging services into your Golang applications. Whether you need to send SMS or Fax messages, this library has got you covered.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
    - [Initialize the Client](#initialize-the-client)
    - [Send an SMS](#send-an-sms)
    - [Send a Fax](#send-a-fax)
- [Examples](#examples)
- [Supported Services](#supported-services)
- [Regions](#regions)
- [Help and Support](#help-and-support)

## Installation

To install the Retarus-Go library, run the following command:

```bash
go get github.com/retarus/retarus-go
```

## Usage

> [!WARNING]
> At this time, we do not offer test accounts. Therefore, each test run will incur charges that will be billed to your company account.

### Initialize the Client
Before making any requests, you'll need to initialize the client. Here's how you can do it, but it is the same for all services:
```go
import "github.com/retarus/retarus-go"

// reads the credentials from system env variables.
config := retarus.fax.NewConfigFromEnv(Europe)
client := retarus.fax.NewClient(config)
```

As you'll observe, we employ the NewConfigFromEnv function. This approach utilizes the credentials specified in the operating system's environment variables, necessitating that these values be exported accordingly.
```bash
export retarus_username=value
export retarus_password=value
export retarus_cuno=yourCuno
```

### Send a Fax
Here's a basic example to send a Fax:
```go
job := fax.Job{
		Recipients: []fax.Recipient{
			{
				Number: "00498923422342", // number to send to
			},
		},
		Documents: []fax.Document{
			{
				Name:      "test.txt", // local document to send
				Data:      "dGVzdGZheAo=", // test fax
			},
		},
	}
jobID, err := client.Send(job)
if err != nil {
fmt.Println("Error: ", err)
}
fmt.Println("JobId: ", jobID)
```
## Examples
For more comprehensive examples, please refer to the [`examples`](/examples) directory in the repository.

## Supported Services

Retarus-Go currently supports the following services:

- Fax
- SMS

## Regions

The library allows you to specify a datacenter region when configuring your service client. Supported regions are:

- Europe
- America
- Switzerland
- Singapore

## Help and Support

For additional information or to get support, visit our [Knowledge Center](https://developers.retarus.com/).