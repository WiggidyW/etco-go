package implrdb

type PreviousContracts struct {
	Buyback []string `firestore:"buyback_codes"`
	Shop    []string `firestore:"shop_codes"`
	Haul    []string `firestore:"haul_codes"`
}
