package appraisal

const (
	CHARACTERS_COLLECTION_ID string = "characters"

	S_CHAR_APPRAISALS string = "shop_appraisals"
	B_CHAR_APPRAISALS string = "buyback_appraisals"
)

type CharacterAppraisalCodes struct {
	ShopAppraisals    []string `firestore:"shop_appraisals"`
	BuybackAppraisals []string `firestore:"buyback_appraisals"`
}
