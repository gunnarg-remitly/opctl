package opstate

import (
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
)

type noopOpFormatter struct{}

func (noopOpFormatter) FormatOpRef(opRef string) string {
	return opRef
}

type staticLoadingSpinner struct{}

func (staticLoadingSpinner) String() string {
	return "⋰"
}

func TestCallGraph(t *testing.T) {
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
	})

	str := "\n" + graph.String(noopOpFormatter{}, staticLoadingSpinner{}, false)
	expected := `
◎ oppath
├─◎ oppath2
│ ├─◉ ⋰ containe containerRef
│ └─◉ ☑ containe containerRef 0s
├─◉ ⋰ if skipped
└─◉ ⋰ if
    serial`
	if str != expected {
		t.Errorf("call graph string not correct: expected\n```\n%s\n```\nactual\n```\n%s\n```", expected, str)
	}
}
