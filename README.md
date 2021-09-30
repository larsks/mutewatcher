# mutewatcher

Mutewatcher sets MIDI controls values to reflect the state of the
Pulseaudio default sink and default source mute state.

## Configuration

You can configure mutewatcher using the following environment
variables:

- `MUTEWATCHER_MIDI_DEVICE` - Device name (or unique substring)
- `MUTEWATCHER_MIDI_CHANNEL` - Channel on which to send messages
- `MUTEWATCHER_SINK_LED` - Control number to reflect sink (output) mute status
- `MUTEWATCHER_SOURCE_LED` - Control number to reflect source (input) mute status
