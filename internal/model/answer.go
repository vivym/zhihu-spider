package model

type AnswerList struct {
	Paging Paging   `json:"paging"`
	Data   []Answer `json:"data"`
}

type Answer struct {
	ID              int      `json:"id" bson:"id"`
	AnswerType      string   `json:"answer_type" bson:"answer_type"`
	Question        Question `json:"question" bson:"question"`
	URL             string   `json:"url" bson:"url"`
	CreatedTime     int      `json:"created_time" bson:"created_time"`
	UpdatedTime     int      `json:"updated_time" bson:"updated_time"`
	Content         string   `json:"content" bson:"content"`
	EditableContent string   `json:"editable_content" bson:"editable_content"`
	Excerpt         string   `json:"excerpt" bson:"excerpt"`
}

type Question struct {
	ID int `json:"id" bson:"id"`
}
