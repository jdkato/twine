package nlp

// A Model holds the structures and data used internally by prose.
type Model struct {
	Name string

	tagger *perceptronTagger
}

// DataSource provides training data to a Model.
type DataSource func(model *Model)

// ModelFromData creates a new Model from user-provided training data.
func ModelFromData(name string, sources ...DataSource) *Model {
	model := defaultModel(true)
	model.Name = name
	for _, source := range sources {
		source(model)
	}
	return model
}

func defaultModel(tagging bool) *Model {
	var tagger *perceptronTagger

	if tagging {
		tagger = newPerceptronTagger()
	}

	return &Model{
		Name: "en-v2.0.0",

		tagger: tagger,
	}
}
