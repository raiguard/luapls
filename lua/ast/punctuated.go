package ast

import "github.com/raiguard/luapls/lua/token"

type Punctuated[T Node] struct {
	Pairs    []Pair[T]
	StartPos token.Pos // If there are zero pairs then we need an alternative way to get the position.
}

func SimplePunctuated[T Node](node T) Punctuated[T] {
	return Punctuated[T]{
		Pairs:    []Pair[T]{{node, nil}},
		StartPos: node.Pos(),
	}
}

func (p *Punctuated[T]) Pos() token.Pos {
	if len(p.Pairs) == 0 {
		return p.StartPos
	}
	return p.Pairs[0].Pos()
}

func (p *Punctuated[T]) End() token.Pos {
	if len(p.Pairs) == 0 {
		return p.StartPos
	}
	return p.Pairs[len(p.Pairs)-1].End()
}

func (p *Punctuated[T]) String() string {
	out := ""
	for _, pair := range p.Pairs {
		out += pair.String()
	}
	return out
}

type Pair[T Node] struct {
	Node      T
	Delimeter *Unit // Optional
}

func (p *Pair[T]) Pos() token.Pos {
	return p.Node.Pos()
}

func (p *Pair[T]) End() token.Pos {
	if p.Delimeter != nil {
		return p.Delimeter.End()
	}
	return p.Node.End()
}

func (p *Pair[T]) String() string {
	out := p.Node.String()
	if p.Delimeter != nil {
		out += p.Delimeter.String()
	}
	return out
}
