package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"view_count/viewservice"
)

var viewService viewservice.Service

func Execute(svc viewservice.Service, input string) error {
	viewService = svc

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return fmt.Errorf("no command provided")
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "getView":
		if len(args) != 1 {
			return fmt.Errorf("get-view requires exactly one argument")
		}
		getView(args[0])

	case "getAll":
		getAllViews()

	case "incre":
		if len(args) != 1 {
			return fmt.Errorf("increment-view requires exactly one argument")
		}
		incrementView(args[0])

	case "top":
		n, _ := strconv.Atoi(args[0])
		getTopViews(n)

	case "recent":
		n, _ := strconv.Atoi(args[0])
		getRecentViews(n)

	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return nil
}

func getView(id string) {
	ctx := context.Background()

	views, err := viewService.GetView(ctx, id)
	if err != nil {
		fmt.Printf("Error getting view for ID: %s, error: %v\n", id, err)
		return
	}
	fmt.Printf("View count for ID %s: %d\n", id, views)
}

func getAllViews() {
	ctx := context.Background()
	videos, err := viewService.GetAllViews(ctx)
	if err != nil {
		fmt.Println("Error getting all videos.", err)
		return
	}
	fmt.Println(videos)
}

func incrementView(id string) {
	ctx := context.Background()
	err := viewService.Increment(ctx, id)
	if err != nil {
		fmt.Println("Error incrementing the views of this video.", err)
		return
	}
}

func getTopViews(n int) {
	ctx := context.Background()
	videos, err := viewService.GetTopVideos(ctx, n)
	if err != nil {
		fmt.Println("Error getting top videos.", err)
		return
	}
	fmt.Println(videos)
}

func getRecentViews(n int) {
	ctx := context.Background()
	videos, err := viewService.GetRecentVideos(ctx, n)
	if err != nil {
		fmt.Println("Error getting recent videos", err)
		return
	}
	fmt.Println(videos)
}
