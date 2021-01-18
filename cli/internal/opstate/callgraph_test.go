package opstate

import (
	"testing"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/stretchr/testify/assert"
)

type noopOpFormatter struct{}

func (noopOpFormatter) FormatOpRef(opRef string) string {
	return opRef
}

func TestCallGraph(t *testing.T) {
	// arrange
	timestamp, err := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 4, 2014 at 6:05pm (PST)")
	if err != nil {
		t.Fatal(err)
	}
	objectUnderTest := CallGraph{}
	parentID := "id1"
	child1ID := "id2"
	containerRef := "containerRef"
	child1child1ID := "child1child1Id"
	child2If := false
	child3If := true

	// act
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID: parentID,
				Op: &model.OpCall{
					BaseCall: model.BaseCall{
						OpPath: "oppath",
					},
					OpID: "firstopid",
				},
			},
		},
		Timestamp: timestamp,
	})
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child1ID,
				ParentID: &parentID,
				Op: &model.OpCall{
					BaseCall: model.BaseCall{
						OpPath: "oppath2",
					},
					OpID: "secondopid",
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 1),
	})
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child1child1ID,
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "id1234567890",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 2),
	})
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child1Child2Id",
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "id0987654321",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 3),
	})
	objectUnderTest.HandleEvent(&model.Event{
		CallEnded: &model.CallEnded{
			Call: model.Call{
				ID:       "child1Child2Id",
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "id0987654321",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
			Outcome: model.OpOutcomeSucceeded,
		},
		Timestamp: timestamp.Add(time.Second * 4),
	})
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child2ID",
				ParentID: &parentID,
				If:       &child2If,
			},
		},
		Timestamp: timestamp.Add(time.Second * 5),
	})
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child3ID",
				ParentID: &parentID,
				If:       &child3If,
				Serial:   []*model.CallSpec{},
			},
		},
		Timestamp: timestamp.Add(time.Second * 6),
	})
	str := "\n" + objectUnderTest.String(
		noopOpFormatter{},
		StaticLoadingSpinner{},
		timestamp.Add(time.Second*6),
		false,
	)

	// assert
	expected := `
◎ oppath
├─◎ oppath2
│ ├─◉ ⋰ id123456 containerRef 4s
│ └─◉ ☑ id098765 containerRef 1s
├─◉ ⋰ if skipped 1s
└─◉ ⋰ if
    serial 0s`
	assert.Equal(t, expected, str)
}
