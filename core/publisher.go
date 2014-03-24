package core

type Publisher struct {
	Subscribers []string
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Subscribe() {

}
