package keys

const (
	IS_BUY_TRUE  string = "b"
	IS_BUY_FALSE string = "s"
)

func isBuyStr(isBuy bool) string {
	if isBuy {
		return IS_BUY_TRUE
	} else {
		return IS_BUY_FALSE
	}
}
