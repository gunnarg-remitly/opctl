package opstate

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/sdks/go/model"
)

// CallGraph maintains a record of the current state of an op
type CallGraph struct {
	rootNode *callGraphNode
	errors   []error
}

type callGraphNode struct {
	call      *model.Call
	startTime *time.Time
	endTime   *time.Time
	state     string
	children  []*callGraphNode
}

func newCallGraphNode(call *model.Call, timestamp time.Time) *callGraphNode {
	return &callGraphNode{
		call:      call,
		startTime: &timestamp,
		children:  []*callGraphNode{},
	}
}

var errNotFoundInGraph = errors.New("not found in graph")

func (n *callGraphNode) insert(call *model.Call, startTime time.Time) error {
	if n.call.ID == *call.ParentID {
		n.children = append(n.children, newCallGraphNode(call, startTime))
		return nil
	}
	for _, child := range n.children {
		err := child.insert(call, startTime)
		if err == nil {
			return nil
		}
	}
	return errNotFoundInGraph
}

func (n *callGraphNode) find(call *model.Call) (*callGraphNode, error) {
	if call.ID == n.call.ID {
		return n, nil
	}
	for _, child := range n.children {
		c, _ := child.find(call)
		if c != nil {
			return c, nil
		}
	}
	return nil, errNotFoundInGraph
}

func (n *callGraphNode) isLeaf() bool {
	return len(n.children) == 0
}

func (n *callGraphNode) countChildren() int {
	count := 0
	for _, child := range n.children {
		if child.isLeaf() {
			count++
		} else {
			count += child.countChildren()
		}
	}
	return count
}

func (n callGraphNode) String(opFormatter clioutput.OpFormatter, loader LoadingSpinner, collapseCompleted bool) string {
	var str strings.Builder

	// Graph node indicator
	if n.isLeaf() {
		str.WriteString("◉ ")
	} else {
		str.WriteString("◎ ")
	}

	// Leading "status"
	switch n.state {
	case model.OpOutcomeSucceeded:
		str.WriteString(success.Sprint("☑ "))
	case model.OpOutcomeFailed:
		str.WriteString(failed.Sprint("⚠ "))
	case model.OpOutcomeKilled:
		str.WriteString("️☒ ")
	case model.OpOutcomeSkipped:
		str.WriteString("☐ ")
	case "":
		// only display loading spinner on leaf nodes
		if n.isLeaf() {
			str.WriteString(loader.String() + " ")
		}
	default:
		str.WriteString(n.state + " ")
	}

	call := *n.call

	// "Named" ops
	if call.Name != nil {
		str.WriteString(highlighted.Sprint(*call.Name) + " ")
	}

	// Main node description
	var desc string
	if call.Container != nil {
		desc = muted.Sprint(call.Container.ContainerID[:8]) + " "
		if call.Container.Name != nil {
			desc += highlighted.Sprint(*call.Container.Name)
		} else {
			desc += *call.Container.Image.Ref
		}
	} else if call.Op != nil {
		desc = highlighted.Sprint(opFormatter.FormatOpRef(call.Op.OpPath))
	} else if call.Parallel != nil {
		desc = "parallel"
	} else if call.ParallelLoop != nil {
		desc = "parallel loop"
	} else if call.Serial != nil {
		desc = "serial"
	} else if call.SerialLoop != nil {
		desc = "serial loop"
	}

	collapsed := n.state == model.OpOutcomeSucceeded && !n.isLeaf() && collapseCompleted

	if call.If != nil {
		str.WriteString("if")
		// this means it was skipped
		if desc == "" {
			str.WriteString(" " + muted.Sprint("skipped"))
		} else {
			str.WriteString("\n")
			if n.isLeaf() || collapsed {
				str.WriteString("  ")
			} else {
				str.WriteString("│ ")
			}
		}
	}

	str.WriteString(desc)

	// Time elapsed (if done)
	if n.startTime != nil && n.endTime != nil {
		str.WriteString(" " + n.endTime.Sub(*n.startTime).String())
	}

	// Add the command invoked by a container if it's not named
	if call.Container != nil && call.Container.Name == nil && len(call.Container.Cmd) > 0 {
		str.WriteString(" " + muted.Sprint(strings.ReplaceAll(strings.Join(call.Container.Cmd, " "), "\n", "\\n")))
	}

	// Collapsed nodes
	if collapsed {
		str.WriteString(" ")
		childCount := n.countChildren()
		if childCount == 1 {
			str.WriteString(muted.Sprint("(1 child)"))
		} else {
			str.WriteString(muted.Sprintf("(%d children)", childCount))
		}
		return str.String()
	}

	// Children
	childLen := len(n.children)
	for i, child := range n.children {
		childLines := strings.Split(child.String(opFormatter, loader, collapseCompleted), "\n")
		for j, part := range childLines {
			if j == 0 {
				if i < childLen-1 {
					str.WriteString(fmt.Sprintf("\n├─%s", part))
				} else {
					str.WriteString(fmt.Sprintf("\n└─%s", part))
				}
			} else if i < childLen-1 {
				str.WriteString(fmt.Sprintf("\n│ %s", part))
			} else {
				str.WriteString(fmt.Sprintf("\n  %s", part))
			}
		}
	}

	return str.String()
}

// String returns a visual representation of the current state of the call graph
func (g CallGraph) String(opFormatter clioutput.OpFormatter, loader LoadingSpinner, collapseCompleted bool) string {
	var str strings.Builder
	str.WriteString(g.rootNode.String(opFormatter, loader, collapseCompleted))
	for _, err := range g.errors {
		str.WriteString("\n" + warning.Sprint("⚠️  ") + err.Error())
	}
	return str.String()
}

// HandleEvent accepts an opctl event and updates the call graph appropriately
func (g *CallGraph) HandleEvent(event *model.Event) error {
	if event.CallStarted != nil {
		if event.CallStarted.Call.ParentID == nil {
			if g.rootNode == nil {
				g.rootNode = newCallGraphNode(&event.CallStarted.Call, event.Timestamp)
				return nil
			}
			return errors.New("parent node already set")
		}
		return g.rootNode.insert(&event.CallStarted.Call, event.Timestamp)
	} else if event.CallEnded != nil {
		node, err := g.rootNode.find(&event.CallEnded.Call)
		if err != nil {
			err = fmt.Errorf("bad ended event %s, %v: %v", event.CallEnded.Call.ID, event.CallEnded.Ref, err)
			g.errors = append(g.errors, err)
			return err
		}
		node.endTime = &event.Timestamp
		node.state = event.CallEnded.Outcome
	}
	return nil
}
