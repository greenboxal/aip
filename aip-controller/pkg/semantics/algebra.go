package semantics

type SemanticBasis struct {
	// Max tokens in a single SemanticUnit
	MaxTokens int
	// Embedding dimensions
	EmbeddingDimensions int
	// Embedding model
	EmbeddingModel string
}

// SemanticUnit is a semantic set of tokens.
// This set is defined by a set of tokens and a set of weights.
// The weights are used to calculate the semantic vector of the set.
// The semantic vector is used to calculate the semantic distance between two semantic sets.
// The semantic distance is used to calculate the semantic similarity between two semantic sets.
// A SemanticUnit can be either a flat semantic set or a hierarchical semantic set.
// A flat semantic set is a semantic set that contains only tokens.
// A hierarchical semantic set is a semantic set that contains other semantic sets.
// A hierarchical semantic set can be factored into a flat semantic set and a set of hierarchical semantic sets by calling Factor().
type SemanticUnit interface {
	Basis() SemanticBasis

	// Add concatenates two semantic units
	// This is equivalent to concatenating two potentially related semantic sets
	// into a single semantic set.
	// The total token count of the resulting semantic set is approximately the sum of the token counts of the two semantic sets.
	Add(unit SemanticUnit) SemanticUnit

	// Sub removes one semantic unit from another
	// This is equivalent to `a.Add(b.Negate())`
	Sub(unit SemanticUnit) SemanticUnit

	// Scale scales the semantic unit by the given factor.
	// A positive factor increases the token count proportionally to |factor|.
	// A negative factor decreases the token count proportionally to |factor|.
	Scale(factor float64) SemanticUnit

	// DivInteger splits the semantic set into a number of smaller equal parts.
	// Definition:
	// 	Let S be a semantic set.
	// 	Let n be a positive integer.
	// 	Let S' be the result of splitting S into n equal parts.
	// 	Then S' is defined as:
	// 		S' = {s_i | s_i = s_j for all i, j in [1, n]}
	// 	Where:
	// 		s_i is the i-th semantic unit in S'
	// 		s_j is the j-th semantic unit in S
	// 	Note:
	// 		For any semantic unit s in S, there exists a semantic unit s' in S'
	// 		such that s' = s / n
	// 	Note:
	// 		For any semantic unit s in S, there exists a semantic unit s' in S'
	// 		such that s' = s * (1 / n)
	// 	Note:
	// 		For any semantic unit s in S, there exists a semantic unit s' in S'
	// 		such that s' = s * (1 / n)
	// 	Note:
	DivInteger(count int) SemanticUnit

	// Factor returns the factors of the current semantic unit.
	// This is equivalent to returning every distinguishable semantic unit contained in the semantic set.
	Factor() []SemanticUnit

	// Mul multiplies two semantic units.
	// This is equivalent to merging two potentially related semantic sets into a single semantic set.
	// The total token count of the resulting semantic set is proportional to (1 - |a.Dot(b)|) * |a| * |b|
	Mul(other SemanticUnit) SemanticUnit

	// Dot returns the dot product of two semantic units.
	// This is equivalent to measuring the similarity between two potentially related semantic sets.
	Dot(other SemanticUnit) float64

	// Wedge returns the wedge product of two semantic units.
	// This is equivalent to measuring the dissimilarity between two potentially related semantic sets.
	Wedge(other SemanticUnit) SemanticUnit

	// OuterProduct returns the outer product of two semantic units.
	// This is equivalent to merging two potentially unrelated semantic sets into a single semantic set.
	// The total token count of the resulting semantic set is approximately the product of the token counts of the two semantic sets.
	OuterProduct(other SemanticUnit) SemanticUnit

	// Invert inverts the current semantic unit.
	// This is equivalent to inverting every distinguishable semantic unit contained in the semantic set.
	// The total token count of the resulting semantic set is approximately the same.
	Invert() SemanticUnit

	// Negate returns the negation of the current semantic unit.
	// This is equivalent to negating every distinguishable semantic unit contained in the semantic set.
	// The total token count of the resulting semantic set is approximately the same.
	Negate() SemanticUnit

	// Square returns the square of the current semantic unit.
	// This is equivalent to x.Mul(x)
	Square() SemanticUnit

	Root() SemanticUnit
}
