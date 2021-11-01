package main

import "sort"

type ByEventTime []*MidiEvent

func (eventList ByEventTime) Len() int {
	return len(eventList)
}

func (eventList ByEventTime) Swap(i, j int) {
	eventList[i], eventList[j] = eventList[j], eventList[i]
}

func (eventList ByEventTime) Less(i, j int) bool {
	return eventList[i].AbsoluteTime < eventList[j].AbsoluteTime
}

func shouldIncludeEvent(event *MidiEvent) bool {
	return event.EventType != Metadata ||
		event.FirstParam == MetaTempo
}

func cleanupMidi(midi *Midi) *Midi {
	var allEvents ByEventTime = nil

	for _, track := range midi.Tracks {
		for _, event := range track.Events {
			if shouldIncludeEvent(event) {
				allEvents = append(allEvents, event)
			}
		}

	}

	sort.Stable(allEvents)

	var endTime uint32 = 0

	if len(allEvents) != 0 {
		endTime = allEvents[len(allEvents)-1].AbsoluteTime
	}

	allEvents = append(allEvents, &MidiEvent{
		endTime,
		Metadata,
		0xF,
		MetaEnd,
		0,
		nil,
	})

	return &Midi{
		SingleTrack,
		midi.TicksPerQuarter,
		[]*Track{&Track{allEvents}},
	}
}
