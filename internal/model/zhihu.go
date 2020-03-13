package model

import "github.com/Kamva/mgm/v2"

type ZhihuHotTopics struct {
	mgm.DefaultModel `bson:",inline"`
	Time             int32        `json:"time" bson:"time"`
	Keywords         []Keyword    `json:"keywords" bson:"keywords"`
	Topics           []ZhihuTopic `json:"topics" bson:"topics"`
}

type ZhihuTopic struct {
	Heat     int32     `json:"heat" bson:"heat"`
	QID      int32     `json:"qid" bson:"qid"`
	Title    string    `json:"title" bson:"title"`
	Excerpt  string    `json:"excerpt" bson:"excerpt"`
	Keywords []Keyword `json:"keywords" bson:"keywords"`
}

type Keyword struct {
	Name   string  `json:"name" bson:"name"`
	Weight float64 `json:"weight" bson:"weight"`
	POS    string  `json:"pos" bson:"pos"`
}

func (z *ZhihuHotTopics) CollectionName() string {
	return "zhihu_hot_topics"
}
