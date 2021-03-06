package arn

import "github.com/aerogo/nano"

// ClientErrorReport saves JavaScript errors that happen in web clients like browsers.
type ClientErrorReport struct {
	ID           string `json:"id"`
	Message      string `json:"message"`
	Stack        string `json:"stack"`
	FileName     string `json:"fileName"`
	LineNumber   int    `json:"lineNumber"`
	ColumnNumber int    `json:"columnNumber"`

	HasCreator
}

// StreamClientErrorReports returns a stream of all characters.
func StreamClientErrorReports() chan *ClientErrorReport {
	channel := make(chan *ClientErrorReport, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("ClientErrorReport") {
			channel <- obj.(*ClientErrorReport)
		}

		close(channel)
	}()

	return channel
}

// AllClientErrorReports returns a slice of all characters.
func AllClientErrorReports() []*ClientErrorReport {
	var all []*ClientErrorReport

	stream := StreamClientErrorReports()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}
