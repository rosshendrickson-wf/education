package examples

import (
	"log"
	"math"
)

// Simple Pearson Co-efficient
//
//def pearson(v1, v2):
//  # Simple sums
//      sum1 = sum(v1)
//	      sum2 = sum(v2)
//
//		    # Sums of the squares
//			    sum1Sq = sum([pow(v, 2) for v in v1])
//				    sum2Sq = sum([pow(v, 2) for v in v2])
//
//					  # Sum of the products
//					      pSum = sum([v1[i] * v2[i] for i in range(len(v1))])
//
//						    # Calculate r (Pearson score)
//							    num = pSum - sum1 * sum2 / len(v1)
//								    den = sqrt((sum1Sq - pow(sum1, 2) / len(v1)) * (sum2Sq - pow(sum2, 2)
//									               / len(v1)))
//												       if den == 0:
//													           return 0
//
//															       return 1.0 - num / den

func Pearson(v1 []float64, v2 []float64) float64 {

	// Simple Sums
	sum1 := float64(0.0)
	for _, value := range v1 {
		sum1 += value
	}
	sum2 := float64(0.0)
	for _, value := range v2 {
		sum2 += value
	}

	// Sum of the squares
	sum1Sq := float64(0.0)
	for _, value := range v1 {
		sum1Sq += math.Pow(value, 2)
	}

	sum2Sq := float64(0.0)
	for _, value := range v2 {
		sum2Sq += math.Pow(value, 2)
	}

	// Sum of the products
	pSum := float64(0.0)
	for i, value := range v1 {
		pSum += value * v2[i]
	}

	// Calculate r (Pearson score)
	viLen := float64(len(v1))
	num := pSum - sum1*sum2/viLen
	den := math.Sqrt((sum1Sq - math.Pow(sum1, 2)/viLen) * (sum2Sq - math.Pow(sum2, 2)/viLen))

	if den == 0 {
		log.Println(den)
		return 0
	}
	return 1.0 - num/den

}

func FasterPearson(v1 []float64, v2 []float64) float64 {

	sum1 := float64(0.0)
	sum2 := float64(0.0)
	sum1Sq := float64(0.0)
	sum2Sq := float64(0.0)
	pSum := float64(0.0)
	var value2 float64
	for i, value := range v1 {
		value2 = v2[i]
		// Simple Sums
		sum1 += value
		sum2 += value2

		// Sum of the squares
		sum1Sq += math.Pow(value, 2)
		sum2Sq += math.Pow(value2, 2)

		// Sum of the products
		pSum += value * value2
	}

	// Calculate r (Pearson score)
	viLen := float64(len(v1))
	num := pSum - sum1*sum2/viLen
	den := math.Sqrt((sum1Sq - math.Pow(sum1, 2)/viLen) * (sum2Sq - math.Pow(sum2, 2)/viLen))

	if den == 0 {
		log.Println(den)
		return 0
	}
	return 1.0 - num/den
}
