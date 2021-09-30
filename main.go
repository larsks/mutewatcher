package main

import (
	"fmt"

	"github.com/lawl/pulseaudio"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func IsSourceMuted(client *pulseaudio.Client, name string) bool {
	sources, err := client.Sources()
	must(err)
	for _, source := range sources {
		if source.Name == name {
			return source.Muted
		}
	}
	return false
}

func IsSinkMuted(client *pulseaudio.Client, name string) bool {
	sinks, err := client.Sinks()
	must(err)
	for _, sink := range sinks {
		if sink.Name == name {
			return sink.Muted
		}
	}
	return false
}

func main() {
	client, err := pulseaudio.NewClient()
	must(err)

	updates, err := client.Updates()
	must(err)

	info, err := client.ServerInfo()
	must(err)

	sink_muted := IsSinkMuted(client, info.DefaultSink)
	source_muted := IsSourceMuted(client, info.DefaultSource)

	for range updates {
		info, err := client.ServerInfo()
		must(err)

		if sink_muted != IsSinkMuted(client, info.DefaultSink) {
			sink_muted = !sink_muted
			if sink_muted {
				fmt.Printf("sink %s is muted\n", info.DefaultSink)
			} else {
				fmt.Printf("sink %s is not muted\n", info.DefaultSink)
			}
		}

		if source_muted != IsSourceMuted(client, info.DefaultSource) {
			source_muted = !source_muted
			if source_muted {
				fmt.Printf("source %s is muted\n", info.DefaultSource)
			} else {
				fmt.Printf("source %s is not muted\n", info.DefaultSource)
			}
		}
	}
}
