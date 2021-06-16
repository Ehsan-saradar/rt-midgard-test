package record

// In 2021-04 ThorNode had two bugs when withdrawing with impermanent loss.
// All the constants were generated by querying the real values from Thornode with:
// $ go run ./cmd/onetime/fetchunits [config.json]

// https://gitlab.com/thorchain/thornode/-/issues/912
// There was a bug in thornode, it withdraw units were more then the actually removed pool units,
// but the impermanent loss protection units were also added to them.
type withdrawUnitCorrection struct {
	TX          string
	ActualUnits int64
}

var withdrawUnitCorrectionsMainnet202104 = map[int64]withdrawUnitCorrection{
	47832:  {"4338F014E1FAC05C2248ECE0A36061D92CC76ADF13CCA773272AD70E00B56154", 9066450465},
	79082:  {"A1B155BD4F57DDF91200733EE2552C9E0E828E632F0D91EF69BCAF3D74D8D512", 169807962},
	81055:  {"7613CEC05CA9B3A4BEF864F22E51EA29EB377EF4EC00885F91377F6D74D1DA4D", 2267292958},
	81462:  {"5E02AE1FE7A777BC6CBE8F4FC2DAFC9F8A6464BAAC58697202EAE1A2271D91D2", 8002689544},
	84221:  {"8885C9AC8A26002DA29090D6173D6A1C340AC6BD96837146BDA4ED059EF0760F", 288123877},
	85406:  {"E6907237BFFDFD5F733E5B422D4BC3106A8BCF933A7547843E458580C625D5D5", 609672362},
	88797:  {"F552E27BC9774E546CA4024B8274C758FC6433F3A38B0DB16137196F55E58C73", 2208373135},
	89415:  {"2E48177404B36CE893240A5B0CFF3FA501CE914BBA1F7D3FFEFC75D44110ADCF", 767266632},
	90002:  {"4D41DA864AE89E8B4CC315360F145E33501B2C1534A5757C1104606C967AB54F", 19621520713},
	100196: {"C94BD47100E0C9983845735A3FA0C6C511713CB4486CBB3777F8DA386011A0C0", 8280457915},
	105465: {"E86DCD9FDD898A3F7781D049EE0442DCC69ACBC2FBB110125A501AF7CF3003D7", 911047010},
	109333: {"C1BD2175944D490D56755B37D1EB88385F9BF7A34EF609418A332526859C6EE2", 406716426},
	110069: {"8D5BBF31ABCB8297AB2804186D6AAA1B479E79B1CB0A0C1B2586F0F89225C28B", 13600885317},
	112985: {"DAC7FCA92A9B42B82BFBE9C03C756A1AFBEF178CF8D2F6F2E044407A6696D581", 117224625},
	128842: {"34B820F7158C3AB690C2DCF088356D1A70E6721551C2159C96729CE9FA97B698", 93675000000},
	128845: {"0754C907993E389BA7947CB775D456BB829E12B3D7EEB676413E749BB847068B", 146382616748},
	131366: {"8EEB3FBAA095F46E12207257C3CB0771BDB55C3EB2322F86FD75594ECC015AD1", 45078869167},
	138590: {"7058BA9B3FF1173D620773458F84C5EA247EBB38C74C505E1FB8069CDB8A6E27", 14950765467},
	147789: {"8CDA8459400D97CC436F1D19B6E42A4CEDDD21F2A231D1F9D4438B43A7750136", 4873515514},
	147798: {"EAF6064BD7CB29389917BF4FF0D499D8E99890D9B561D8FF63F610092FADA4A3", 814479987},
	151691: {"6A7A7C3A7A65F4704151DB1972EAFA6A237B03BA82D46721E761F3063753C42C", 345151887},
	153980: {"85A19DA310282D35A6C51F4C34F921D27F2DF090535790F0C533FE61EA980CD7", 1115323168},
	163137: {"0BA388B1BCF76C04B81D885ECB99E0E98A295778234FF9A88E9CA8ED69706DF4", 3086810573},
	166532: {"156CCFBC66F775C7FDF9D3E18F071C6CEC2ADFAB4F7F435094AA516ECD1C698A", 8288025767},
	257485: {"E6B6FBC73BFD62BC36F0E236BC065FCC18D328832908C240399E2DF2E2CB6565", 9702125229},
	260113: {"A6788765BCFBEC33F0F4585CB736105D4005AB81FEC30113231CF1D41F843AEA", 272714488439},
	260114: {"1296D15627331C78CA5BC7CEE014C98273C5B08D358FA451C8039B42EAD61054", 128877756350},
	260115: {"F6B4EDB5555CC4FAF513729D16F2D906DEC5C950DC95530F26ABFDC7ECD5DBCE", 75139724801},
	260116: {"86679B5EE155F2997251108713C96AE0AC91444BFD0883A99D0611A255F0F2D7", 41517402427},
	260119: {"04B6F0AEDFA9ABD9DD949541C2B7762DC2EA62026ACB39C8992482355318FB8C", 29065838793},
	265159: {"2BABF243911BA2CFF2551143131985515C4873C9D6C87E44027E0F7F14E29792", 18962634918},
	269611: {"CCE905915CEC65FD6FDC48E31E43E65FCD73ABDCF90A4419EFDFE7E43B63DDD0", 156734300},
	271635: {"02BA91CF8F6FF3E35A1C7F0F1991BB2A2E200B78B3CF7A77DAF77E66067B205F", 83041261241},
	271741: {"4B95FDF07545DF8BDD9B05982F013166E0BAB8B54F419548DEB0D3EE2E5F454E", 1539766365},
	277262: {"E0E67CF364BFDD9B312C1899C60582F720A44F1A8023333F7849E0AAD0B9E4DB", 9402258},
	292069: {"70558EA306ADA6C6705A4C15AA60BB06D9000F75F9C2FA85153027F0AC131357", 10046673124},
}

var withdrawUnitCorrections *(map[int64]withdrawUnitCorrection)

func correctWithdawsImpLoss(withdraw *Unstake, meta *Metadata) {
	if withdrawUnitCorrections == nil {
		return
	}
	correction, ok := (*withdrawUnitCorrections)[meta.BlockHeight]
	if ok {
		if correction.TX == string(withdraw.Tx) {
			withdraw.StakeUnits = correction.ActualUnits
		}
	}
}

func loadMainnetWithdrawImpLossUnitCorrections() {
	withdrawUnitCorrections = &withdrawUnitCorrectionsMainnet202104
	for k := range *withdrawUnitCorrections {
		WithdrawCorrections.Add(k, correctWithdawsImpLoss)
	}
}

// Sometimes when withdrawing the pool units of a member went up, not down:
// https://gitlab.com/thorchain/thornode/-/issues/896
type addInsteadWithdraw struct {
	Pool     string
	RuneAddr string
	Units    int64
}

var addInsteadWithdrawMapMainnet202104 = map[int64]addInsteadWithdraw{
	84876:  {"BTC.BTC", "thor1h7n7lakey4tah37226musffwjhhk558kaay6ur", 2029187601},
	170826: {"BNB.BNB", "thor1t5t5xg7muu3fl2lv6j9ck6hgy0970r08pvx0rz", 31262905},
}

func loadMainnetWithdrawImpLossAdds() {
	addAddEvent := func(d *Demux, meta *Metadata) {
		add, ok := addInsteadWithdrawMapMainnet202104[meta.BlockHeight]
		if ok {
			d.reuse.Stake = Stake{
				AddBase: AddBase{
					Pool:     []byte(add.Pool),
					RuneAddr: []byte(add.RuneAddr),
				},
				StakeUnits: add.Units,
			}
			Recorder.OnStake(&d.reuse.Stake, meta)
		}
	}
	for height := range addInsteadWithdrawMapMainnet202104 {
		AdditionalEvents.Add(height, addAddEvent)
	}
}

func loadMainnetCorrectionsWithdrawImpLoss() {
	loadMainnetWithdrawImpLossAdds()
	loadMainnetWithdrawImpLossUnitCorrections()
}
