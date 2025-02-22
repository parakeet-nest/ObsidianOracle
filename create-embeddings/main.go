package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/parakeet-nest/parakeet/content"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	//TODO: to be changed at the end
	contentPath := gear.GetEnvString("CONTENT_PATH", "./data")

	ollamaUrl := gear.GetEnvString("OLLAMA_BASE_URL", "http://localhost:11434")

	embeddingsModel := gear.GetEnvString("LLM_EMBEDDINGS", "mxbai-embed-large")

	elasticStore := embeddings.ElasticsearchStore{}
	errEsStore := elasticStore.Initialize(
		[]string{
			os.Getenv("ELASTICSEARCH_HOSTS"),
		},
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
		nil,
		os.Getenv("ELASTICSEARCH_INDEX"),
	)
	if errEsStore != nil {
		log.Fatalln("😡:", errEsStore)
	}

	// Parse the Obsidian vault
	contentFiles, errContentFiles := content.GetMapOfContentFiles(contentPath, ".md")
	if errContentFiles != nil {
		log.Fatalln("😡:", errContentFiles)
	}
	//fmt.Println("📚 Content files:", contentFiles)

	promptTemplate := `METADATA:
	Lineage: {{.Lineage}}
	File: {{.Metadata.FilePath}}
	## {{.Header}}
	{{.Content}}
	`

	// Iterate over the content files and create chunks then embeddings
	// key is the path of the markdown file
	// contentFile is the content of the markdown file
	counter := 0
	for key, contentFile := range contentFiles {
		// create chunks from the markdown file (contentFile)
		chunks := content.ParseMarkdownWithLineage(contentFile)
		
		for idx, chunk := range chunks {
			fmt.Println("📝 Creating embedding from document ", idx)
			fmt.Println("📝", key)
			fmt.Println("🖼️", chunk.Header)
			fmt.Println("🌲", chunk.Lineage)

			chunk.Metadata = make(map[string]interface{})

			chunk.Metadata["FilePath"] = key

			prompt, errInterpolation := content.InterpolateString(promptTemplate, chunk)
			if errInterpolation != nil {
				log.Println("😡:", errInterpolation)
				//continue
			} else {
				fmt.Println("================================================")
				fmt.Println("📝 Prompt:")
				fmt.Println(prompt)
				fmt.Println("================================================")
				// Create the embeddings
				embedding, errEmbedding := embeddings.CreateEmbedding(
					ollamaUrl,
					llm.Query4Embedding{
						Model:  embeddingsModel,
						Prompt: prompt,
					},
					strconv.Itoa(idx)+"-"+strconv.Itoa(counter),
				)
				if errEmbedding != nil {
					fmt.Println("😡 when generating embedding:", errEmbedding)
				} else {
					if _, errEsSave := elasticStore.Save(embedding); errEsSave != nil {
						log.Fatalln("😡 we have a problem with ES when saving embedding:", errEsSave)
					}
					fmt.Println("🎉 Document", embedding.Id, "indexed successfully")
					counter++
				}
			}

		}
	}
	fmt.Println("🎉 All documents indexed successfully")

}
