package entity

type Investor struct {
	ID            string
	Name          string
	AssetPosition []*InvestorAssetPosition
}
