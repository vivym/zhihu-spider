package spiders

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/Kamva/mgm/v2"
	"github.com/go-resty/resty/v2"
	"github.com/k3a/html2text"

	"github.com/vivym/zhihu-spider/internal/model"
	"github.com/vivym/zhihu-spider/internal/nlp"
)

type Spider struct {
	config *Config
	http   *resty.Client
	nlp    *nlp.NLPToolkit
}

func New(config Config, nlpToolkit *nlp.NLPToolkit) *Spider {
	http := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(5*time.Second).
		SetRetryMaxWaitTime(20*time.Second).
		SetHostURL(baseURL).
		SetHeader("User-Agent", userAgent)

	return &Spider{
		config: &config,
		http:   http,
		nlp:    nlpToolkit,
	}
}

func (s *Spider) Go() error {
	hotTopicList, err := s.fetchHotTopicList()
	if err != nil {
		return err
	}
	fmt.Println("hotTopicList done.")

	var zhihuHotTopics model.ZhihuHotTopics
	var sentences []string
	for i, hotTopic := range hotTopicList.Data {
		fmt.Println("topic:", hotTopic.Target.Title)
		if i >= s.config.MaxTopics {
			break
		}
		answers, err := s.fetchAnswersAll(hotTopic.Target.ID)
		if err != nil {
			if len(answers) == 0 {
				return err
			}
			log.Printf("warning: %v\n", err)
		}

		zhihuTopic := model.ZhihuTopic{
			Heat:    parseHeat(hotTopic.DetailText),
			QID:     int32(hotTopic.Target.ID),
			Title:   hotTopic.Target.Title,
			Excerpt: hotTopic.Target.Excerpt,
		}
		var texts []string
		for _, answer := range answers {
			text := []rune(html2text.HTML2Text("<div>" + answer.Content + "</div>"))
			limit := 1000
			if len(text) < 1000 {
				limit = len(text)
			}
			texts = append(texts, string(text[0:limit]))
		}
		sentence := strings.Join(texts, "\n")
		limit := 100
		if len(texts) < 100 {
			limit = len(texts)
		}
		sentences = append(sentences, strings.Join(texts[0:limit], "\n"))

		keywords, err := s.nlp.ExtractKeywords(sentence, 100)
		if err != nil {
			log.Printf("warning: %v\n", err)
		}
		zhihuTopic.Keywords = keywords

		zhihuHotTopics.Topics = append(zhihuHotTopics.Topics, zhihuTopic)

		delay := time.Duration(s.config.Delay + rand.Intn(300))
		time.Sleep(delay * time.Millisecond)
	}

	sentence := strings.Join(sentences, "\n")
	keywords, err := s.nlp.ExtractKeywords(sentence, 100)
	if err != nil {
		log.Printf("warning: %v\n", err)
	}
	zhihuHotTopics.Keywords = keywords

	if err := mgm.Coll(&zhihuHotTopics).Create(&zhihuHotTopics); err != nil {
		return err
	}

	if err := mgm.Coll(&hotTopicList).Create(&hotTopicList); err != nil {
		return err
	}

	return nil
}

func (s *Spider) fetchAnswersAll(questionID int) ([]model.Answer, error) {
	const limit, offset = 20, 0

	url := fmt.Sprintf(
		"/v4/questions/%d/answers?include=%s&limit=%d&offset=%d&platform=desktop&sort_by=%s",
		questionID,
		"data[*].is_normal,admin_closed_comment,reward_info,is_collapsed,annotation_action,annotation_detail,collapse_reason,is_sticky,collapsed_by,suggest_edit,comment_count,can_comment,content,editable_content,voteup_count,reshipment_settings,comment_permission,created_time,updated_time,review_info,relevant_info,question,excerpt,relationship.is_authorized,is_author,voting,is_thanked,is_nothelp,is_labeled,is_recognized,paid_info,paid_info_content;data[*].mark_infos[*].url;data[*].author.follower_count,badge[*].topics",
		limit,
		offset,
		s.config.SortBy, // "default, updated"
	)

	answers := []model.Answer{}
	for i := 0; i < s.config.MaxPages; i++ {
		fmt.Println("page", i+1)
		answerList, err := s.fetchAnswers(url)
		for _, answer := range answerList.Data {
			answers = append(answers, answer)
		}

		if err != nil || answerList.Paging.Next == "" || answerList.Paging.IsEnd {
			return answers, err
		}
		url = answerList.Paging.Next

		delay := time.Duration(s.config.Delay + rand.Intn(300))
		time.Sleep(delay * time.Millisecond)
	}

	return answers, nil
}

func (s *Spider) fetchAnswers(url string) (model.AnswerList, error) {
	rsp, err := s.http.R().
		SetResult(&model.AnswerList{}).
		Get(url)
	if err != nil {
		return model.AnswerList{}, err
	}
	return *rsp.Result().(*model.AnswerList), nil
}

func (s *Spider) fetchHotTopicList() (model.HotTopicList, error) {
	rsp, err := s.http.R().
		SetResult(&model.HotTopicList{}).
		Get("/v3/feed/topstory/hot-lists/total?limit=50&desktop=true")
	if err != nil {
		return model.HotTopicList{}, err
	}
	return *rsp.Result().(*model.HotTopicList), nil
}
