package webexteams

import (
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"github.com/peterhellberg/link"
	"time"
)

type Participant struct {
	ID                       string                    `json:"id"`
	OrgID                    string                    `json:"orgId"`
	Host                     bool                      `json:"host"`
	CoHost                   bool                      `json:"coHost"`
	SpaceModerator           bool                      `json:"spaceModerator"`
	Email                    string                    `json:"email"`
	DisplayName              string                    `json:"displayName"`
	Invitee                  bool                      `json:"invitee"`
	Muted                    bool                      `json:"muted"`
	MeetingStartTime         time.Time                 `json:"meetingStartTime"`
	Video                    string                    `json:"video"`
	State                    string                    `json:"state"`
	BreakoutSessionID        string                    `json:"breakoutSessionId"`
	JoinedTime               time.Time                 `json:"joinedTime"`
	LeftTime                 time.Time                 `json:"leftTime"`
	SiteURL                  string                    `json:"siteUrl"`
	MeetingID                string                    `json:"meetingId"`
	HostEmail                string                    `json:"hostEmail"`
	Devices                  []Device                  `json:"devices"`
	BreakoutSessionsAttended []BreakoutSessionAttended `json:"breakoutSessionsAttended"`
	SourceID                 string                    `json:"sourceId"`
}

// Meetings is the List of Meetings
type Participants struct {
	Items []Participant `json:"items,omitempty"`
}
type ParticipantDevice struct {
	CorrelationID  string    `json:"correlationId"`
	DeviceType     string    `json:"deviceType"`
	AudioType      string    `json:"audioType"`
	JoinedTime     time.Time `json:"joinedTime"`
	LeftTime       time.Time `json:"leftTime"`
	DurationSecond int       `json:"durationSecond"`
	CallType       string    `json:"callType"`
	PhoneNumber    string    `json:"phoneNumber"`
}

type BreakoutSessionAttended struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	JoinedTime time.Time `json:"joinedTime"`
	LeftTime   time.Time `json:"leftTime"`
}

type ParticipantsService service

func (p *Participants) AddParticipant(item Participant) []Participant {
	p.Items = append(p.Items, item)
	return p.Items
}

func (s *ParticipantsService) participantsPagination(linkHeader string, size, max int) *Participants {
	items := &Participants{}

	for _, l := range link.Parse(linkHeader) {
		if l.Rel == "next" {

			response, err := s.client.R().
				SetResult(&Participants{}).
				SetError(&Error{}).
				Get(l.URI)

			if err != nil {
				return nil
			}
			items = response.Result().(*Participants)
			if size != 0 {
				size = size + len(items.Items)
				if size < max {
					meetings := s.participantsPagination(response.Header().Get("Link"), size, max)
					for _, meeting := range meetings.Items {
						items.AddParticipant(meeting)
					}
				}
			} else {
				meetings := s.participantsPagination(response.Header().Get("Link"), size, max)
				for _, meeting := range meetings.Items {
					items.AddParticipant(meeting)
				}
			}

		}
	}

	return items
}

// ListParticipants lists participants in a meeting
type ListParticipantsQueryParams struct {
	MeetingID            string    `url:"meetingId,omitempty"`
	Max                  int       `url:"max,omitempty"`
	MeetingStartTimeFrom time.Time `url:"joinTimeFrom,omitempty"`
	//JoinTimeTo   string `url:"joinTimeTo,omitempty"`
	Paginate bool // Indicates if pagination is needed
}

func (w *ParticipantsService) ListParticipants(queryParams ListParticipantsQueryParams) (*Participants, *resty.Response, error) {
	path := "/meetingParticipants"
	queryParamsString, _ := query.Values(queryParams)
	response, err := w.client.R().
		SetQueryString(queryParamsString.Encode()).
		SetResult(&Participants{}).
		SetError(&Error{}).
		Get(path)

	if err != nil {
		return nil, response, err
	}

	result := response.Result().(*Participants)
	if queryParams.Paginate {
		items := w.participantsPagination(response.Header().Get("Link"), 0, 0)
		for _, meeting := range items.Items {
			result.AddParticipant(meeting)
		}
	} else {
		if len(result.Items) < queryParams.Max {
			items := w.participantsPagination(response.Header().Get("Link"), len(result.Items), queryParams.Max)
			for _, meeting := range items.Items {
				result.AddParticipant(meeting)
			}
		}
	}

	return result, response, err
}

//
//// QueryParticipants queries participants with specific parameters
//func (w *ParticipantsService) QueryParticipants(meetingID, queryParam string) ([]Participant, error) {
//	var response struct {
//		Items []Participant `json:"items"`
//	}
//
//	_, err := w.client.R().
//		SetQueryParams(map[string]string{
//			"meetingId":  meetingID,
//			"queryParam": queryParam, // Customize with additional query parameters as needed
//		}).
//		SetResult(&response).
//		Get("/meetingParticipants")
//
//	if err != nil {
//		return nil, err
//	}
//	return response.Items, nil
//}
//
//// GetParticipant retrieves a specific participant by ID
//func (w *ParticipantsService) GetParticipant(participantID string) (*Participant, error) {
//	var participant Participant
//
//	_, err := w.client.R().
//		SetPathParam("participantId", participantID).
//		SetResult(&participant).
//		Get("/meetingParticipants/{participantId}")
//
//	if err != nil {
//		return nil, err
//	}
//	return &participant, nil
//}
//
//// UpdateParticipant updates a participant's details
//func (w *ParticipantsService) UpdateParticipant(participantID string, participantData map[string]interface{}) (*Participant, error) {
//	var updatedParticipant Participant
//
//	_, err := w.client.R().
//		SetPathParam("participantId", participantID).
//		SetBody(participantData).
//		SetResult(&updatedParticipant).
//		Put("/meetingParticipants/{participantId}")
//
//	if err != nil {
//		return nil, err
//	}
//	return &updatedParticipant, nil
//}
//
//// AdmitParticipant admits a participant into a meeting
//func (w *ParticipantsService) AdmitParticipant(participantID string) error {
//	_, err := w.client.R().
//		SetPathParam("participantId", participantID).
//		Post("/meetingParticipants/{participantId}/admit")
//
//	return err
//}
