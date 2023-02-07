package articles

type TagSerializer struct {
	Tag
}

func (s *TagSerializer) Response() string {
	return s.Tag.Tag
}

type TagsSerializer struct {
	Tags []Tag
}

func (s *TagsSerializer) Response() []string {
	response := []string{}
	for _, tag := range s.Tags {
		serializer := TagSerializer{tag}
		response = append(response, serializer.Response())
	}
	return response
}
