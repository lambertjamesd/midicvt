package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type channelMetadata struct {
	volumeChange float32
}

type midiChangeMetadata struct {
	volume   float32
	speed    float32
	channels [16]channelMetadata
}

func readMetadata(content string) midiChangeMetadata {
	var metadata = make(map[string]string)

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)

		var attrName = strings.TrimSpace(parts[0])

		if len(parts) == 2 {
			metadata[attrName] = strings.TrimSpace(parts[1])
		}
	}

	var result = midiChangeMetadata{
		1,
		1,
		[16]channelMetadata{
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
			channelMetadata{1},
		},
	}

	volumeString, has := metadata["volume"]

	if has {
		volume, err := strconv.ParseFloat(volumeString, 32)

		if err != nil {
			log.Fatal(err)
		}

		result.volume = float32(volume)
	}

	volumeString, has = metadata["speed"]

	if has {
		speed, err := strconv.ParseFloat(volumeString, 32)

		if err != nil {
			log.Fatal(err)
		}

		result.speed = float32(speed)
	}

	for i := 0; i < 16; i++ {
		volumeString, has = metadata[fmt.Sprintf("channelVolume%d", i)]

		if has {
			volume, err := strconv.ParseFloat(volumeString, 32)

			if err != nil {
				log.Fatal(err)
			}

			result.channels[i].volumeChange = float32(volume)
		}
	}

	return result
}

func applyMetadata(input *Midi, metdata *midiChangeMetadata) {
	input.TicksPerQuarter = uint16(float32(input.TicksPerQuarter) * metdata.speed)

	for _, track := range input.Tracks {
		for _, event := range track.Events {
			if event.EventType == MidiOn {
				var newVolume = float32(event.SecondParam) * metdata.volume * metdata.channels[event.Channel].volumeChange
				if newVolume > 127 {
					newVolume = 127
				}
				event.SecondParam = uint8(newVolume)
			}
		}
	}
}

func maxOutVolume(input *Midi) {
	var maxVolume = 0

	for _, track := range input.Tracks {
		for _, event := range track.Events {
			if event.EventType == MidiOn && int(event.SecondParam) > maxVolume {
				maxVolume = int(event.SecondParam)
			}
		}
	}

	for _, track := range input.Tracks {
		for _, event := range track.Events {
			if event.EventType == MidiOn {
				event.SecondParam = uint8(int(event.SecondParam) * 127 / maxVolume)
			}
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal(fmt.Sprintf("Usage: %s input.mid output.mid\n", os.Args[0]))
	}

	inputfile, err := os.Open(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	input, err := ReadMidi(inputfile)

	if err != nil {
		log.Fatal(err)
	}

	result := cleanupMidi(input)

	if len(os.Args) > 3 && os.Args[3] == "--max" {
		maxOutVolume(result)
	} else if len(os.Args) > 4 && os.Args[3] == "--metadata" {
		fileContent, err := ioutil.ReadFile(os.Args[4])

		if err != nil {
			log.Fatal(err)
		}

		var metadata = readMetadata(string(fileContent))
		applyMetadata(result, &metadata)
	}

	outputFile, err := os.OpenFile(os.Args[2], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)

	err = WriteMidi(outputFile, result)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Processed file %s and saved it to %s", os.Args[1], os.Args[2]))
}
