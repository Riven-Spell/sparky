package sparkrun

// RecipeNameOrFile is satisfied by [RecipeName] and [RecipeFile]. It
// is accepted by methods that operate on a recipe by name or by
// recipe-YAML path (Run, RecipeShow, RecipeVram, RecipeValidate).
type RecipeNameOrFile interface {
	recipeRef() string
}

// RecipeNameOrJobID is satisfied by [RecipeName] and [JobID]. It is
// accepted by methods that target either a recipe (by name) or a
// specific running job (by cluster ID such as "sparkrun_<hex>") --
// notably Stop, Logs, ClusterCheckJob, and the export verbs.
type RecipeNameOrJobID interface {
	workloadRef() string
}

type recipeName string

func (r recipeName) recipeRef() string   { return string(r) }
func (r recipeName) workloadRef() string { return string(r) }

// RecipeName references a recipe by its registered name. It satisfies
// both RecipeNameOrFile and RecipeNameOrJobID, so it can be passed to
// any verb that takes a recipe-or-target.
func RecipeName(name string) recipeName { return recipeName(name) }

type recipeFile string

func (r recipeFile) recipeRef() string { return string(r) }

// RecipeFile references a recipe by filesystem path to its YAML file.
func RecipeFile(path string) recipeFile { return recipeFile(path) }

type jobID string

func (j jobID) workloadRef() string { return string(j) }

// JobID references a running workload by its cluster ID
// (e.g. "sparkrun_3ba7d381d8c8").
func JobID(id string) jobID { return jobID(id) }
