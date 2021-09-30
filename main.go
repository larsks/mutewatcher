package main

import (
	"github.com/larsks/go-decouple"
	"github.com/lawl/pulseaudio"
	"github.com/rs/zerolog/log"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/rtmididrv"
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

func SetMidiLed(wr *writer.Writer, control int, state bool) {
	if state {
		must(writer.CcOn(wr, uint8(control)))
	} else {
		must(writer.CcOff(wr, uint8(control)))
	}
}

func main() {
	_ = decouple.Load()
	deviceName, exists := decouple.GetString("MUTEWATCHER_MIDI_DEVICE", "")
	if !exists {
		panic("you must provide a device name in MUTEWATCHER_MIDI_DEVICE")
	}
	deviceChannel, _ := decouple.GetInt("MUTEWATCHER_MIDI_CHANNEL", 0)
	sinkLedControl, _ := decouple.GetInt("MUTEWATCHER_SINK_LED", 0x30)
	sourceLedControl, _ := decouple.GetInt("MUTEWATCHER_SOURCE_LED", 0x31)

	drv, err := driver.New()
	must(err)

	port, err := midi.OpenOut(drv, -1, deviceName)
	must(err)

	log := log.With().
		Str("device", port.String()).
		Int("channel", deviceChannel).Logger()

	wr := writer.New(port)
	wr.SetChannel(uint8(deviceChannel))

	client, err := pulseaudio.NewClient()
	must(err)

	updates, err := client.Updates()
	must(err)

	info, err := client.ServerInfo()
	must(err)

	sink_name := info.DefaultSink
	source_name := info.DefaultSource

	sink_muted := IsSinkMuted(client, sink_name)
	source_muted := IsSourceMuted(client, source_name)

	log.Info().Str("sink", sink_name).Bool("muted", sink_muted).Msgf("initial sink status")
	log.Info().Str("source", source_name).Bool("muted", source_muted).Msgf("initial source status")

	SetMidiLed(wr, sinkLedControl, sink_muted)
	SetMidiLed(wr, sourceLedControl, source_muted)

	for range updates {
		info, err := client.ServerInfo()
		must(err)

		if sink_muted != IsSinkMuted(client, info.DefaultSink) {
			sink_muted = !sink_muted
			SetMidiLed(wr, sinkLedControl, sink_muted)
			log.Info().Str("sink", sink_name).Bool("muted", sink_muted).Msgf("sink status changed")
		}

		if source_muted != IsSourceMuted(client, info.DefaultSource) {
			source_muted = !source_muted
			SetMidiLed(wr, sourceLedControl, source_muted)
			log.Info().Str("source", source_name).Bool("muted", source_muted).Msgf("source status changed")
		}
	}
}
