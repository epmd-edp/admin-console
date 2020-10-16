package query

type PerfServer struct {
	Id   int    `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
}
