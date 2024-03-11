package parserc

type Parser[T any] func(ParserInput) ParserResult[T]
