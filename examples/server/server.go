package main

import (
	"github.com/matrixxsoftware/go-mdd/cmdc"
	"github.com/matrixxsoftware/go-mdd/dictionary"
	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
	"github.com/matrixxsoftware/go-mdd/transport/tcp"

	log "github.com/sirupsen/logrus"
)

func main() {
	formatter := &log.TextFormatter{}
	formatter.DisableQuote = true
	log.SetFormatter(formatter)
	log.SetLevel(log.TraceLevel)

	dict := dictionary.New()
	dict.Add(&dictionary.ContainerDefinition{
		Key:           101,
		Name:          "Request",
		SchemaVersion: 5222,
		ExtVersion:    2,
		Fields: []dictionary.FieldDefinition{
			{Name: "Field1", Type: field.Int32},
			{Name: "Field2", Type: field.String},
			{Name: "Field3", Type: field.Decimal},
			{Name: "Field4", Type: field.UInt32},
			{Name: "Field5", Type: field.UInt16},
			{Name: "Field6", Type: field.UInt64},
		},
	})
	dict.Add(&dictionary.ContainerDefinition{
		Key:           88,
		Name:          "Response",
		SchemaVersion: 5222,
		ExtVersion:    2,
		Fields: []dictionary.FieldDefinition{
			{Name: "ResultCode", Type: field.UInt32},
			{Name: "ResultMessage", Type: field.String},
		},
	})

	addr := "0.0.0.0:14060"
	log.Infof("Server listening on %s", addr)

	transport, err := tcp.NewServerTransport(addr)
	if err != nil {
		panic(err)
	}
	defer transport.Close()

	codec := cmdc.NewCodecWithDict(dict)
	server := &mdd.Server{
		Codec:     codec,
		Transport: transport,
	}

	server.MessageHandler(func(request *mdd.Containers) (*mdd.Containers, error) {
		log.Infof("Server received request:\n%s", request.Dump())

		return &mdd.Containers{
			Containers: []mdd.Container{
				{
					Header: mdd.Header{
						Version:       1,
						TotalField:    3,
						Depth:         0,
						Key:           88,
						SchemaVersion: 5222,
						ExtVersion:    2,
					},
					Fields: []mdd.Field{
						{Data: []byte("0")},
						{Data: []byte("(2:OK)")},
					},
				},
			},
		}, nil
	})

	err = transport.Listen()
	if err != nil {
		panic(err)
	}
}
