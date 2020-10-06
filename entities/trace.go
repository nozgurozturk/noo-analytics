package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Action int

const (
	view Action = 1 + iota
	click
	reveal
)

type Trace struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	IP        string             `json:"ip" bson:"ip"`
	Location  string             `json:"loc" bson:"loc"`
	TimeStamp int64              `json:"ts" bson:"ts"`
	Year      int                `json:"y" bson:"y"`
	Month     int                `json:"m" bson:"m"`
	Day       int                `json:"d" bson:"d"`
	Hour      int                `json:"h" bson:"h"`
	Agent     string             `json:"a" bson:"a"`
	Action    Action             `json:"act" bson:"act"`
	Tag       string             `json:"tag" bson:"tag"`
}

type TraceDTO struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	IP       string             `json:"ip"`
	Location string             `json:"loc"`
	Agent    string             `json:"a"`
	Action   Action             `json:"act"`
	Tag      string             `json:"tag"`
}

func ToTrace(traceDto *TraceDTO) *Trace {
	return &Trace{
		IP:       traceDto.IP,
		Location: traceDto.Location,
		Agent:    traceDto.Agent,
		Action:   traceDto.Action,
		Tag:      traceDto.Tag,
	}
}

func ToTraceDTO(trace *Trace) *TraceDTO {
	return &TraceDTO{
		ID:       trace.ID,
		IP:       trace.IP,
		Location: trace.Location,
		Agent:    trace.Agent,
		Action:   trace.Action,
		Tag:      trace.Tag,
	}
}

func ToTraceDTOs(traces []*Trace) []*TraceDTO {
	traceDTOs := make([]*TraceDTO, len(traces))

	for i, trace := range traces {
		traceDTOs[i] = ToTraceDTO(trace)
	}
	return traceDTOs
}
