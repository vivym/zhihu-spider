package model

import "github.com/Kamva/mgm/v2"

type HotTopicList struct {
	mgm.DefaultModel `bson:",inline"`
	Paging           Paging     `json:"paging"`
	Data             []HotTopic `json:"data" bson:"data"`
}

type HotTopic struct {
	DetailText string         `json:"detail_text" bson:"detail_text"`
	Target     HotTopicTarget `json:"target" bson:"target"`
}

type HotTopicTarget struct {
	BoundTopicIds []int  `json:"bound_topic_ids" bson:"bound_topic_ids"`
	Excerpt       string `json:"excerpt" bson:"excerpt"`
	AnswerCount   int    `json:"answer_count" bson:"answer_count"`
	ID            int    `json:"id" bson:"id"`
	Title         string `json:"title" bson:"title"`
	Created       int    `json:"created" bson:"created"`
	CommentCount  int    `json:"comment_count" bson:"comment_count"`
	FollowerCount int    `json:"follower_count" bson:"follower_count"`
}

func (h *HotTopicList) CollectionName() string {
	return "zhihu_hot_topic_list"
}
