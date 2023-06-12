package exp

import (
	"fmt"

	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/mat"

	"github.com/greenboxal/aip/aip-langchain/pkg/providers/openai"
)

type EmbeddingTestCommand struct {
	oai      *openai.Client
	embedder *openai.Embedder
}

func NewEmbeddingTestCommand(
	oai *openai.Client,
) *EmbeddingTestCommand {
	return &EmbeddingTestCommand{
		oai: oai,

		embedder: &openai.Embedder{
			Client: oai,
			Model:  openai.AdaEmbeddingV2,
		},
	}
}

func (etc *EmbeddingTestCommand) Run(cmd *cobra.Command, args []string) error {
	var letters []string

	for i := 'a'; i <= 'c'; i++ {
		letters = append(letters, string(i))
	}

	letters = []string{"dead", "alive", "zombie"}

	embeddings, err := etc.embedder.GetEmbeddings(cmd.Context(), letters)

	if err != nil {
		return err
	}

	banana := embeddings[0].Vector()
	apple := embeddings[1].Vector()
	dog := embeddings[2].Vector()

	bananaApple := mat.NewVecDense(1536, nil)
	bananaApple.SubVec(banana, apple)
	bananaAppleLength := bananaApple.Norm(2)

	bananaDog := mat.NewVecDense(1536, nil)
	bananaDog.SubVec(banana, dog)
	bananaDogLength := bananaApple.Norm(2)

	fmt.Printf("banana - apple = %f\n", bananaAppleLength)
	fmt.Printf("banana - dog = %f\n", bananaDogLength)

	for i, e := range embeddings {
		letter := letters[i]

		fmt.Printf("%s = %f\n", letter, e.Embeddings)
	}

	return nil
}
