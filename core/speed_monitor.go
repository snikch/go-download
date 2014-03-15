package core

import "time"

type SpeedMonitor struct {
	downloaded ByteSize
	Speed      ByteSize // Bytes / Second
	updater    chan ByteSize
	running    chan bool
	stopper    chan bool
	times      []time.Time
	readings   []ByteSize
}

func NewSpeedMonitor(updater chan ByteSize) *SpeedMonitor {
	m := SpeedMonitor{
		updater:  updater,
		running:  make(chan bool),
		times:    make([]time.Time, 0, 10),
		readings: make([]ByteSize, 0, 10),
		stopper:  make(chan bool),
	}
	return &m
}
func (s *SpeedMonitor) Stop() {
	s.stopper <- true
}

func (s *SpeedMonitor) Start() {
	c := time.Tick(250 * time.Millisecond)
	for {
		select {
		// Add a new data point, and update speed if required
		case now := <-c:
			// Don't do anything until we've started downloading
			if s.downloaded == 0 {
				continue
			}

			// How big is our current array?
			index := len(s.readings)

			// Grab the last n-1 records if we're at capacity
			if index >= cap(s.readings) {
				s.readings = s.readings[1:index]
				s.times = s.times[1:index]
			}

			// Add our values
			s.readings = append(s.readings, s.downloaded)
			s.times = append(s.times, now)

			// Wait until we have a few readings
			if index < 3 {
				continue
			}

			// Update our speed
			s.updateSpeed()

		case bytes := <-s.updater:
			s.downloaded = bytes

		case stop := <-s.stopper:
			if stop == true {
				break
			}
		}
	}
}

func (s *SpeedMonitor) updateSpeed() {
	downloaded := s.readings[len(s.readings)-1] - s.readings[0]
	duration := s.times[len(s.times)-1].Sub(s.times[0])

	s.Speed = ByteSize(float64(downloaded) / duration.Seconds())
}
