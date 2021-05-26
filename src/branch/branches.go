package branch

// Branch type alias
type Branch int

type Branches []Branch

// All branches ordered
const (
	Concord Branch = iota
	Scarborough
	Hampden
	Bangor
	Patten
	IslandFalls
	Oakfield
	Houlton
	MarsHill
	PresqueIsle
	Ashland
	FortFairfield
	Caribou
	VanBuren
	EagleLake
	FortKent
)

// NoBranch is an invalid branch
const NoBranch Branch = -1

// Get returns the branch or NoBranch
func Get(branch int) Branch {
	if branch > 15 {
		return NoBranch
	}
	return Branch(branch)
}

// DistanceVector is an ordered list of distances to other branches following
// the defined branch order
type DistanceVector = []int

// DistanceVectors is a map of all of the distances from a branch to all other branches
var DistanceVectors = map[Branch]DistanceVector{
	Concord:       {0, 105, 235, 237, 326, 328, 338, 355, 380, 394, 376, 399, 407, 430, 406, 422},
	Scarborough:   {105, 0, 133, 135, 223, 226, 236, 253, 278, 292, 274, 297, 305, 328, 304, 320},
	Hampden:       {235, 133, 0, 6, 96, 99, 108, 125, 150, 164, 146, 170, 177, 201, 176, 192},
	Bangor:        {237, 135, 6, 0, 90, 92, 102, 118, 144, 158, 140, 163, 171, 194, 170, 186},
	Patten:        {326, 223, 96, 90, 0, 10, 20, 38, 63, 67, 50, 78, 78, 99, 80, 96},
	IslandFalls:   {328, 226, 99, 92, 10, 0, 11, 28, 53, 67, 56, 72, 80, 103, 86, 102},
	Oakfield:      {338, 236, 108, 102, 20, 11, 0, 17, 43, 57, 46, 62, 70, 93, 76, 92},
	Houlton:       {355, 253, 125, 118, 38, 28, 17, 0, 27, 42, 62, 46, 54, 77, 91, 92},
	MarsHill:      {380, 278, 150, 144, 63, 53, 43, 27, 0, 15, 37, 20, 27, 50, 68, 72},
	PresqueIsle:   {394, 292, 164, 158, 67, 67, 57, 42, 15, 0, 23, 11, 13, 36, 53, 58},
	Ashland:       {376, 274, 146, 140, 50, 56, 46, 62, 37, 23, 0, 33, 29, 50, 35, 58},
	FortFairfield: {399, 297, 170, 163, 78, 72, 62, 46, 20, 11, 33, 0, 12, 32, 66, 55},
	Caribou:       {407, 305, 177, 171, 78, 80, 70, 54, 27, 13, 29, 12, 0, 23, 55, 45},
	VanBuren:      {430, 328, 201, 194, 99, 103, 93, 77, 50, 36, 50, 32, 23, 0, 59, 42},
	EagleLake:     {406, 304, 176, 170, 80, 86, 76, 91, 68, 53, 35, 66, 55, 59, 0, 17},
	FortKent:      {422, 320, 192, 186, 96, 102, 92, 92, 72, 58, 58, 55, 45, 42, 17, 0},
}

// DistanceTo returns the distance from this branch to the given branch
func (b Branch) DistanceTo(branch Branch) int {
	return DistanceVectors[b][branch]
}

// DistanceTrip returns the total trip distance from this branch
// to all other branches given in order of visit
func (b Branch) DistanceTrip(branches ...Branch) int {
	dist := 0
	for i, branch := range branches {
		if i == 0 {
			dist += b.DistanceTo(branch)
			continue
		}
		dist += branches[i-1].DistanceTo(branches[i])
	}
	return dist
}

// RoundTrip calculates the total round trip distance from this
// branch to all given branches and then back to the origin
func (b Branch) RoundTrip(branches ...Branch) int {
	to := b.DistanceTrip(branches...)
	to += branches[len(branches)-1].DistanceTo(b)
	return to
}

func (b Branches) Len() int {
	return len(b)
}

func (b Branches) Less(i, j int) bool {
	return b[i] < b[j]
}

func (b Branches) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
