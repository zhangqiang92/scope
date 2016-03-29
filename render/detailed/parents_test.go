package detailed_test

import (
	"fmt"
	"testing"

	"github.com/weaveworks/scope/render"
	"github.com/weaveworks/scope/render/detailed"
	"github.com/weaveworks/scope/report"
	"github.com/weaveworks/scope/test"
	"github.com/weaveworks/scope/test/fixture"
	"github.com/weaveworks/scope/test/reflect"
)

func TestParents(t *testing.T) {
	for _, c := range []struct {
		name string
		node report.Node
		want []detailed.Parent
	}{
		{
			name: "Node accidentally tagged with itself",
			node: render.HostRenderer.Render(fixture.Report)[fixture.ClientHostNodeID].WithParents(
				report.EmptySets.Add(report.Host, report.MakeStringSet(fixture.ClientHostNodeID)),
			),
			want: nil,
		},
		{
			node: render.HostRenderer.Render(fixture.Report)[fixture.ClientHostNodeID],
			want: nil,
		},
		{
			node: render.ContainerImageRenderer.Render(fixture.Report)[fixture.ClientContainerImageNodeID],
			want: []detailed.Parent{
				{ID: fixture.ClientHostNodeID, Label: fixture.ClientHostName, TopologyID: "hosts"},
			},
		},
		{
			node: render.ContainerRenderer.Render(fixture.Report)[fixture.ClientContainerNodeID],
			want: []detailed.Parent{
				{ID: fixture.ClientContainerImageNodeID, Label: fixture.ClientContainerImageName, TopologyID: "containers-by-image"},
				{ID: fixture.ClientHostNodeID, Label: fixture.ClientHostName, TopologyID: "hosts"},
			},
		},
		{
			node: render.ProcessRenderer.Render(fixture.Report)[fixture.ClientProcess1NodeID],
			want: []detailed.Parent{
				{ID: fixture.ClientContainerNodeID, Label: fixture.ClientContainerName, TopologyID: "containers"},
				{ID: fixture.ClientContainerImageNodeID, Label: fixture.ClientContainerImageName, TopologyID: "containers-by-image"},
				{ID: fixture.ClientHostNodeID, Label: fixture.ClientHostName, TopologyID: "hosts"},
			},
		},
	} {
		name := c.name
		if name == "" {
			name = fmt.Sprintf("Node %q", c.node.ID)
		}
		if have := detailed.Parents(fixture.Report, c.node); !reflect.DeepEqual(c.want, have) {
			t.Errorf("%s: %s", name, test.Diff(c.want, have))
		}
	}
}
