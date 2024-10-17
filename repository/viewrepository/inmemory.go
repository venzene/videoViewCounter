package viewrepository

import (
	"container/heap"
	"context"
	"sync"
	"time"
	"view_count/model"
)

// type inmemoryRepo struct {
// TODO find a Bug
// 	data map[string]int
// }

type inmemoryRepo struct {
	mu   sync.RWMutex
	data map[string]*videoData
	// TODO: create 2 heap. one for count, one for time : DONE
	viewHeap VideoViewHeap
	timeHeap VideoTimeHeap
}

type videoData struct {
	Id          string
	Views       int
	LastUpdated time.Time // TODO time.Time : done
}

func NewInmemoryRepo() *inmemoryRepo {
	return &inmemoryRepo{
		data:     make(map[string]*videoData),
		viewHeap: make(VideoViewHeap, 0),
		timeHeap: make(VideoTimeHeap, 0),
	}
}

type VideoViewHeap []*videoData
type VideoTimeHeap []*videoData

func (h VideoViewHeap) Len() int            { return len(h) }
func (h VideoViewHeap) Less(i, j int) bool  { return h[i].Views > h[j].Views }
func (h VideoViewHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *VideoViewHeap) Push(x interface{}) { *h = append(*h, x.(*videoData)) }
func (h *VideoViewHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h VideoTimeHeap) Len() int            { return len(h) }
func (h VideoTimeHeap) Less(i, j int) bool  { return h[i].LastUpdated.After(h[j].LastUpdated) }
func (h VideoTimeHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *VideoTimeHeap) Push(x interface{}) { *h = append(*h, x.(*videoData)) }
func (h *VideoTimeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// TODO: Write unit test cases
func (repo *inmemoryRepo) GetView(ctx context.Context, videoId string) (view int, err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	video, ok := repo.data[videoId]
	if !ok {
		video = &videoData{Id: videoId, Views: 0}
	}
	return video.Views, nil
}

func (repo *inmemoryRepo) GetAllViews(ctx context.Context) (info []model.VideoInfo, err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	info = make([]model.VideoInfo, len(repo.data))

	c := 0
	for k, v := range repo.data {
		info[c] = model.VideoInfo{
			Id:    k,
			Views: v.Views,
		}
		c++
	}
	return info, nil
}

func (repo *inmemoryRepo) Increment(ctx context.Context, videoId string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	video, exists := repo.data[videoId]
	if !exists {
		video = &videoData{Id: videoId, Views: 0}
	}

	video.Views++
	video.LastUpdated = time.Now()
	repo.data[videoId] = video

	// TODO use fix : done
	if !exists {
		heap.Push(&repo.viewHeap, video)
		heap.Push(&repo.timeHeap, video)
	} else {
		for i, v := range repo.viewHeap {
			if v.Id == videoId {
				repo.viewHeap[i].Views = video.Views
				heap.Fix(&repo.viewHeap, i)
			}
		}

		for i, v := range repo.timeHeap {
			if v.Id == videoId {
				repo.timeHeap[i].Views = video.Views
				heap.Fix(&repo.viewHeap, i)
			}
		}
	}
	return nil
}

func (repo *inmemoryRepo) GetTopVideos(ctx context.Context, n int) (info []model.VideoInfo, err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if n > repo.viewHeap.Len()  {
		n = repo.viewHeap.Len()
	}

	tempHeap := make(VideoViewHeap, 0, repo.viewHeap.Len())
	tempHeap = append(tempHeap, repo.viewHeap...)

	topVideos := make([]model.VideoInfo, n)

	for i := 0; i < n; i++ {
		video := heap.Pop(&tempHeap).(*videoData)
		topVideos[i] = model.VideoInfo{
			Id:    video.Id,
			Views: video.Views,
		}
	}

	return topVideos, nil

}

func (repo *inmemoryRepo) GetRecentVideos(ctx context.Context, n int) ([]model.VideoInfo, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if n > repo.timeHeap.Len(){
		n = repo.timeHeap.Len()
	}

	tempHeap := make(VideoTimeHeap, 0, repo.timeHeap.Len())
	tempHeap = append(tempHeap, repo.timeHeap...)

	recentVideos := make([]model.VideoInfo, n)
	for i := 0; i < n; i++ {
		video := heap.Pop(&tempHeap).(*videoData)
		recentVideos[i] = model.VideoInfo{
			Id:    video.Id,
			Views: video.Views,
		}
	}
	return recentVideos, nil
}
