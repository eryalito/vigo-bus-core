package vitrasa

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/eryalito/vigo-bus-core/internal/sqlite"
	"github.com/eryalito/vigo-bus-core/pkg/api"

	"golang.org/x/net/html"
)

type VitrasaClient struct {
	// ScheduleEndpoint is the base ScheduleEndpoint of the Vitrasa API
	ScheduleEndpoint string
}

// NewVitrasaClient creates a new VitrasaClient with the given URL
func NewVitrasaClient() *VitrasaClient {
	return &VitrasaClient{
		ScheduleEndpoint: "http://infobus.vitrasa.es:8002/Default.aspx",
	}
}

func (c *VitrasaClient) GetSchedules(stopNumber int) ([]api.Schedule, error) {
	stopNumberStr := strconv.Itoa(stopNumber)
	endpoint := c.ScheduleEndpoint + "?parada=" + stopNumberStr

	// Perform the GET request
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Extract the schedules information from the HTML
	schedules, err := extractSchedule(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract schedule: %v", err)
	}

	return schedules, nil
}

// extractSchedule extracts the schedule information from the HTML document
func extractSchedule(n *html.Node) ([]api.Schedule, error) {
	targetNode := findNodeById(n, "GridView1")
	if targetNode == nil {
		fmt.Println("GridView1 node not found")
		return nil, errors.New("GridView1 node not found")
	}
	targetNode = findNode(targetNode, []string{"tbody"})
	if targetNode == nil {
		fmt.Println("Target node not found")
		return nil, errors.New("target node not found")
	}

	var schedules []api.Schedule

	count := 0
	for c := targetNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "tr" {
			if count == 0 {
				// Skip the header row
				count++
				continue
			}

			var schedule api.Schedule
			fieldCounter := 0

			for field := c.FirstChild; field != nil; field = field.NextSibling {
				if field.Type == html.ElementNode && field.Data == "td" {
					fieldCounter++
					childNode := findNode(field, []string{"font"}).FirstChild
					if childNode != nil && childNode.Type == html.TextNode {
						data := childNode.Data
						switch fieldCounter {
						case 1:
							line, err := retrieveLine(data)
							if err != nil {
								fmt.Printf("Error retrieving line %s: %v\n", data, err)
								return nil, errors.New("error retrieving line")
							}
							schedule.Line = line
						case 2:
							schedule.Route = data
						case 3:
							time, err := strconv.Atoi(data)
							if err != nil {
								fmt.Printf("Error converting time to integer: %v\n", err)
								return nil, err
							}
							schedule.Time = time
						}
					}
				}
			}

			schedules = append(schedules, schedule)
		}
	}

	return schedules, nil
}

func retrieveLine(name string) (api.Line, error) {
	// Perform the GET request

	sdb_conn, err := sqlite.NewBusConnector()
	if err != nil {
		return api.Line{}, fmt.Errorf("failed to get line by name: %v", err)
	}
	line, err := sdb_conn.GetLineByName(name)
	if err != nil {
		return api.Line{}, fmt.Errorf("failed to get line by name: %v", err)
	}
	return line, nil
}
func findNodeById(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findNodeById(c, id); result != nil {
			return result
		}
	}

	return nil
}

// findNode traverses the HTML tree to find the node at the specified path
func findNode(n *html.Node, path []string) *html.Node {
	if len(path) == 0 {
		return n
	}

	segment := path[0]
	var tagName string
	re := regexp.MustCompile(`^([a-zA-Z]+)\[(\d+)\]$`)
	matches := re.FindStringSubmatch(segment)
	if len(matches) == 3 {
		tagName = matches[1]
		index, err := strconv.Atoi(matches[2])
		if err != nil {
			fmt.Printf("Error converting index to integer: %v\n", err)
			return nil
		}

		// Handle indexed segment
		count := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == tagName {
				if count == index {
					return findNode(c, path[1:])
				}
				count++
			}
		}
	} else {
		// Handle non-indexed segment
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == segment {
				return findNode(c, path[1:])
			}
		}
	}

	return nil
}
