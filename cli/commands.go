package cli

import (
	"context"
	"fmt"
	"view_count/viewservice"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "practice"}
var viewService viewservice.Service

func Execute(svc viewservice.Service) error {
	viewService = svc
	return rootCmd.Execute()
}

// var inMemory = &cobra.Command{
// 	Use:   "inMemory",
// 	Short: "To run the server with inMemory",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		inmemory()
// 	},
// }

var getViewCmd = &cobra.Command{
	Use:   "get-view [id]",
	Short: "Get a specific view",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getView(args[0])
	},
}

var getAllViewsCmd = &cobra.Command{
	Use:   "get-all-views",
	Short: "Get all views",
	Run: func(cmd *cobra.Command, args []string) {
		getAllViews()
	},
}

var incrementViewCmd = &cobra.Command{
	Use:   "increment-view [id]",
	Short: "Increment a specific view",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		incrementView(args[0])
	},
}

var getTopTenCmd = &cobra.Command{
	Use:   "get-top-ten",
	Short: "Get Top 10 Viewed Video",
	Run: func(cmd *cobra.Command, args []string) {
		getTopViews()
	},
}

var getRecentCmd = &cobra.Command{
	Use:   "get-recent",
	Short: "Get recently Viewed Video",
	Run: func(cmd *cobra.Command, args []string) {
		getRecentViews()
	},
}

func init() {
	rootCmd.AddCommand(getViewCmd)
	rootCmd.AddCommand(getAllViewsCmd)
	rootCmd.AddCommand(incrementViewCmd)
	rootCmd.AddCommand(getTopTenCmd)
	rootCmd.AddCommand(getRecentCmd)
	// rootCmd.AddCommand(inMemory)
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

func getTopViews() {
	ctx := context.Background()
	videos, err := viewService.GetTopVideos(ctx, 10)
	if err != nil {
		fmt.Println("Error getting top videos.", err)
		return
	}
	fmt.Println(videos)
}

func getRecentViews() {
	ctx := context.Background()
	videos, err := viewService.GetRecentVideos(ctx, 10)
	if err != nil {
		fmt.Println("Error getting recent videos", err)
		return
	}
	fmt.Println(videos)
}

// func inmemory() {
// 	viewRepo := viewrepository.NewInmemoryRepo()
// 	vs := viewservice.NewService(viewRepo)
// 	h := NewHandler(vs)
// 	r := routeIntialiser(*h)
// 	log.Println("Server started on :8080")
// 	log.Fatal(http.ListenAndServe(":8080", r))

// }
