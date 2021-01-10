package opstate

import (
	"testing"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
)

type noopOpFormatter struct{}

func (noopOpFormatter) FormatOpRef(opRef string) string {
	return opRef
}

func TestCallGraph(t *testing.T) {
	timestamp, err := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 4, 2014 at 6:05pm (PST)")
	if err != nil {
		t.Fatal(err)
	}
	graph := CallGraph{}
	parentID := "id1"
	graph.HandleEvent(&model.Event{
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
	child1ID := "id2"
	graph.HandleEvent(&model.Event{
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
	containerRef := "containerRef"
	child1child1ID := "child1child1Id"
	graph.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child1child1ID,
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "container1id",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 2),
	})
	graph.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child1Child2Id",
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "container1id",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 3),
	})
	graph.HandleEvent(&model.Event{
		CallEnded: &model.CallEnded{
			Call: model.Call{
				ID:       "child1Child2Id",
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "container1id",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
			Outcome: model.OpOutcomeSucceeded,
		},
		Timestamp: timestamp.Add(time.Second * 4),
	})
	child2If := false
	graph.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child2ID",
				ParentID: &parentID,
				If:       &child2If,
			},
		},
		Timestamp: timestamp.Add(time.Second * 5),
	})
	child3If := true
	graph.HandleEvent(&model.Event{
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

	str := "\n" + graph.String(
		noopOpFormatter{},
		StaticLoadingSpinner{},
		timestamp.Add(time.Second*6),
		false,
	)
	expected := `
◎ oppath
├─◎ oppath2
│ ├─◉ ⋰ containe containerRef 4s
│ └─◉ ☑ containe containerRef 1s
├─◉ ⋰ if skipped 1s
└─◉ ⋰ if
    serial 0s`
	if str != expected {
		t.Errorf("call graph string not correct: expected\n```\n%s\n```\nactual\n```\n%s\n```", expected, str)
	}
}
