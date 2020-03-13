package spiders

const (
	baseURL   = "https://www.zhihu.com/api"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36"
)

type Config struct {
	Delay     int
	SortBy    string
	MaxTopics int
	MaxPages  int
}
