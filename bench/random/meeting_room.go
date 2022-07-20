package random

import (
	"math/rand"
	"time"

	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
)

type timeRange struct {
	StartAt int64
	EndAt   int64
}

const (
	stepForMeetingRoom        = int(5 * time.Minute / time.Second)
	minDurationForMeetingRoom = int(5 * time.Minute / time.Second)
	maxDurationForMeetingRoom = int(1 * time.Hour / time.Second)
)

var (
	meetingRoomList = []model.MeetingRoom{
		model.AlphaMeetingRoom,
		model.BravoMeetingRoom,
		model.CharlieMeetingRoom,
		model.DeltaMeetingRoom,
		model.EchoMeetingRoom,
		model.FoxtrotMeetingRoom,
		model.GolfMeetingRoom,
		model.HotelMeetingRoom,
		model.IndiaMeetingRoom,
		model.JulietMeetingRoom,
		model.KiloMeetingRoom,
		model.LimaMeetingRoom,
		model.MikenovemberMeetingRoom,
		model.NovemberMeetingRoom,
		model.OscarMeetingRoom,
		model.PapaMeetingRoom,
		model.QuebecMeetingRoom,
		model.RomeoMeetingRoom,
		model.SierraMeetingRoom,
		model.TangoMeetingRoom,
		model.UniformMeetingRoom,
		model.VictorMeetingRoom,
		model.WhiskeyMeetingRoom,
		model.XrayMeetingRoom,
		model.YankeeMeetingRoom,
		model.ZuluMeetingRoom,
	}
	meetingTimeRangeMap map[model.MeetingRoom][]timeRange
)

func init() {
	SetStep(stepForMeetingRoom)
	SetScheduleDuration(minDurationForMeetingRoom, maxDurationForMeetingRoom)

	meetingTimeRangeMap = make(map[model.MeetingRoom][]timeRange, len(meetingRoomList))

	for _, meetingRoom := range meetingRoomList {
		timeRanges := []timeRange{}

		startAt := int64(minDate)
		for {
			endAt := startAt + duration()
			if endAt > int64(maxDate) {
				endAt = int64(maxDate)
			}
			timeRanges = append(timeRanges, timeRange{StartAt: startAt, EndAt: endAt})
			if endAt == int64(maxDate) {
				break
			}
			startAt = endAt
		}

		meetingTimeRangeMap[meetingRoom] = timeRanges
	}

	SetStep(step)
	SetScheduleDuration(minDuration, maxDuration)
}

func MeetingRoom() model.MeetingRoom {
	return meetingRoomList[rand.Intn(len(meetingRoomList))]
}

func TimeRange(room model.MeetingRoom) timeRange {
	timeRanges := meetingTimeRangeMap[room]
	return timeRanges[rand.Intn(len(timeRanges))]
}
