package parser

var parsers = []Parser{
	&TruyenDepParser{},
}

func GetParserForDomain(domain string) Parser {
	for _, p := range parsers {
		if p.DomainMatch(domain) {
			return p
		}
	}
	return nil
}
