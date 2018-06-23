package vote

import (
	"encoding/csv"
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/divan/num2words"

	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/domain/slack"
)

var _ plugin.Plugin = (*Plugin)(nil)
var _ plugin.Generator = (*Generator)(nil)

const defaultSubcommand = "vote"
const chartMax = 40

var words2num = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

// Vote is interface for one of vote
type Vote interface {
	GetTitle() string
	GetOwner() string
	GetOwnerName() string
	GetTimestamp() string
	GetChannel() string
	GetHash() string
	GetAllowDup() bool
	GetChoiseTexts() []string
	GetUserChoises() [][]string
}

// Repository for vote plugin
type Repository interface {
	Create(hash, title, owner, ownerName, timestamp, channel string, allowDup bool, choises []string)
	AddUserChoise(user string, number int, timestamp string, channel string)
	DeleteUserChoise(user string, number int, timestamp string, channel string)
	Find(timestamp string, channel string) Vote
	FindByHash(hash string) Vote
}

// Plugin for vote
type Plugin struct {
	config     plugin.Config
	client     slack.Client
	repository Repository
}

// Generator for vote plugin
type Generator struct {
	Repository Repository
}

func (p *Plugin) createAttachment(hash string, title string, allowDup bool, choises []string, userChoises [][]string) slack.Attachment {

	userVoteCounts := make(map[string]int)
	if userChoises != nil {
		for _, us := range userChoises {
			for _, u := range us {
				userVoteCounts[u]++
			}
		}
	}

	fields := make([]slack.AttachmentField, len(choises))
	for i, choise := range choises {
		total := ""
		additionalText := ""
		dupCount := 0
		if userChoises != nil && len(userChoises) > i && len(userChoises[i]) != 0 {
			additionalText = "\n"
			for _, u := range userChoises[i] {
				if !allowDup && userVoteCounts[u] > 1 {
					additionalText = fmt.Sprintf("%s ~<@%s>~", additionalText, u)
					dupCount++
				} else {
					additionalText = fmt.Sprintf("%s <@%s>", additionalText, u)
				}
			}
			count := len(userChoises[i]) - dupCount
			if count > 0 {
				total = fmt.Sprintf("`%s`", strconv.Itoa(count))
			}
		}
		fields[i] = slack.AttachmentField{
			Value: fmt.Sprintf(":%s: %s %s%s", num2words.Convert(i+1), choise, total, additionalText),
		}
	}

	return slack.Attachment{
		Pretext: fmt.Sprintf("[%s] %s", hash, title),
		Color:   "#3e6cf7",
		Fields:  fields,
	}
}

// ReceiveReactionAdded run received reaction_added
func (p *Plugin) ReceiveReactionAdded(reactionAdded slack.Reaction) {
	if p.client.GetBotUserID() == reactionAdded.User {
		return
	}

	vote := p.repository.Find(reactionAdded.Item.Timestamp, reactionAdded.Item.Channel)
	if vote == nil {
		return
	}
	num, ok := words2num[reactionAdded.Reaction]
	if !ok || len(vote.GetChoiseTexts()) < num {
		return
	}
	p.repository.AddUserChoise(reactionAdded.User, num, vote.GetTimestamp(), vote.GetChannel())

	hash := vote.GetHash()
	title := vote.GetTitle()
	choises := vote.GetChoiseTexts()
	userChoises := vote.GetUserChoises()
	userChoises[num-1] = append(userChoises[num-1], reactionAdded.User)
	timestamp := vote.GetTimestamp()
	attachment := p.createAttachment(hash, title, vote.GetAllowDup(), choises, userChoises)
	p.client.UpdateAttachments(
		vote.GetOwnerName(),
		[]slack.Attachment{attachment},
		reactionAdded.Item.Channel,
		timestamp,
	)
}

// ReceiveReactionRemoved run received reaction_added
func (p *Plugin) ReceiveReactionRemoved(reactionRemoved slack.Reaction) {
	if p.client.GetBotUserID() == reactionRemoved.User {
		return
	}

	vote := p.repository.Find(reactionRemoved.Item.Timestamp, reactionRemoved.Item.Channel)
	if vote == nil {
		return
	}
	num, ok := words2num[reactionRemoved.Reaction]
	if !ok || len(vote.GetChoiseTexts()) < num {
		return
	}

	p.repository.DeleteUserChoise(reactionRemoved.User, num, vote.GetTimestamp(), vote.GetChannel())

	vote = p.repository.Find(reactionRemoved.Item.Timestamp, reactionRemoved.Item.Channel)
	if vote == nil {
		return
	}

	hash := vote.GetHash()
	title := vote.GetTitle()
	choises := vote.GetChoiseTexts()
	userChoises := vote.GetUserChoises()
	timestamp := vote.GetTimestamp()
	attachment := p.createAttachment(hash, title, vote.GetAllowDup(), choises, userChoises)
	p.client.UpdateAttachments(
		vote.GetOwnerName(),
		[]slack.Attachment{attachment},
		reactionRemoved.Item.Channel,
		timestamp,
	)
}

// ReceiveMessage run received message
func (p *Plugin) ReceiveMessage(msg slack.Message) bool {
	if !p.config.CheckEnabledMessage(msg) {
		return false
	}
	subcommand := p.config.GetSubcommand(defaultSubcommand)
	if strings.HasPrefix(msg.Text, fmt.Sprintf("@%s %s.add_nd", p.config.MentionName, subcommand)) {
		p.addVote(msg, false)
		return true
	} else if strings.HasPrefix(msg.Text, fmt.Sprintf("@%s %s.add", p.config.MentionName, subcommand)) {
		p.addVote(msg, true)
		return true
	} else if strings.HasPrefix(msg.Text, fmt.Sprintf("@%s %s.chart", p.config.MentionName, subcommand)) {
		p.showChart(msg)
		return true
	}

	return false
}

func (p *Plugin) showChart(msg slack.Message) {
	splits := strings.SplitN(msg.Text, " ", 3)
	if len(splits) < 3 {
		message := p.config.ResponseTemplates.GetText("not_found_vote", nil)
		p.client.SendMessage(message, msg.Channel)
		return
	}
	hash := splits[2]
	vote := p.repository.FindByHash(hash)
	if vote == nil {
		return
	}
	userChoises := vote.GetUserChoises()
	dupChoiseCounts := make([]int, len(userChoises))
	if userChoises != nil && !vote.GetAllowDup() {
		userVoteCounts := make(map[string][]int)
		for i, us := range userChoises {
			for _, u := range us {
				if userVoteCounts[u] == nil {
					userVoteCounts[u] = make([]int, 0)
				}
				userVoteCounts[u] = append(userVoteCounts[u], i)
			}
		}
		for _, us := range userVoteCounts {
			if len(us) > 1 {
				for _, c := range us {
					dupChoiseCounts[c]++
				}
			}
		}
	}

	total := 0
	choiseCounts := make([]int, len(userChoises))
	for i := range userChoises {
		minus := 0
		if !vote.GetAllowDup() {
			minus = dupChoiseCounts[i]
		}
		l := len(userChoises[i]) - minus
		total += l
		choiseCounts[i] = l
	}
	fields := make([]slack.AttachmentField, len(userChoises))
	for i := range userChoises {
		rate := float64(choiseCounts[i]) / float64(total)
		shift := math.Pow(10, float64(1))
		percentage := math.Floor(rate*100*shift+.5) / shift
		value := fmt.Sprintf(":%s: `.` 0.0%%", num2words.Convert(i+1))
		if rate > 0 {
			value = fmt.Sprintf(
				":%s: `%s` %0.1f%%",
				num2words.Convert(i+1),
				strings.Repeat("━", int(rate*chartMax)),
				percentage,
			)
		}
		fields[i] = slack.AttachmentField{
			Value: value,
		}
	}

	attachments := []slack.Attachment{
		slack.Attachment{
			Pretext: fmt.Sprintf("[%s] %s", vote.GetHash(), vote.GetTitle()),
			Color:   "#3e6cf7",
			Fields:  fields,
		},
	}
	p.client.SendAttachments(msg.User, attachments, msg.Channel)
}

func (p *Plugin) addVote(msg slack.Message, allowDup bool) {
	// split by space with considered quote
	r := csv.NewReader(strings.NewReader(msg.Text))
	r.Comma = ' '
	splits, err := r.Read()
	if err != nil {
		log.Println(err)
		return
	}
	if len(splits) < 4 {
		message := p.config.ResponseTemplates.GetText("not_found_vote_choises", nil)
		p.client.SendMessage(message, msg.Channel)
	} else if len(splits) > 12 {
		message := p.config.ResponseTemplates.GetText("over_ten_vote_choises", nil)
		p.client.SendMessage(message, msg.Channel)
	}

	title := splits[2]
	choises := splits[3:len(splits)]

	h := fnv.New32a()
	h.Write([]byte(msg.TimeStamp))
	hash := strconv.FormatUint(uint64(h.Sum32()), 10)

	attachment := p.createAttachment(hash, title, allowDup, choises, nil)
	result := p.client.SendAttachments(
		msg.UserName,
		[]slack.Attachment{attachment},
		msg.Channel,
	)
	for i := range choises {
		r := num2words.Convert(i + 1)
		p.client.AddReaction(result.Channel, result.Timestamp, r)
	}

	p.repository.Create(hash, title, msg.User, msg.UserName, result.Timestamp, msg.Channel, allowDup, choises)
}

// Generate to vote plugin
func (g *Generator) Generate(config plugin.Config, client slack.Client) plugin.Plugin {
	return &Plugin{
		config:     config,
		client:     client,
		repository: g.Repository,
	}
}
