package model

type Union[T any, U any] struct {
	First  T
	Second U
}

func NewUnion[T any, U any](first T, second U) Union[T, U] {
	return Union[T, U]{
		First:  first,
		Second: second,
	}
}

type RequestWithArtifact struct {
	Request  Request
	Artifact *Artifact
}

func NewRequestWithArtifact(request Request, artifact *Artifact) RequestWithArtifact {
	return RequestWithArtifact{
		Request:  request,
		Artifact: artifact,
	}
}
