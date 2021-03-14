package branch

// FedMileage2021 is the federal mileage reimbursement per mile in dollars
const FedMileage2021 float64 = 0.56

// AverageSpeedLimit is an estimation of the average speed limit for a route
const AverageSpeedLimit int = 65

// Calculation return type
type Calculation struct {
	Distance     int      // distance travelled (mi)
	Cost         float64  // mileage reimbursement (dollars)
	Time         float64  // time travelled (hours)
	Origin       Branch   // origin branch
	Destinations []Branch // destination branches
}
