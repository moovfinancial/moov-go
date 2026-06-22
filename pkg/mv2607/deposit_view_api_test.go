package mv2607_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2607"
)

func TestCreateDepositView(t *testing.T) {
	var (
		method       string
		path         string
		version      string
		sourceSystem string
		contentType  string
		body         []byte
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		sourceSystem = r.Header.Get("X-Source-System")
		contentType = r.Header.Get("Content-Type")
		body = b

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"moovAccountID": "account-123",
			"sourceSystem": "jh_silverlake",
			"sourceAccountID": "987654321",
			"ingestedAt": "2026-06-18T12:00:00Z"
		}`))
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	depositView := mv2607.NewDepositViewClient(client)

	document := []byte(`{"accountNumber":"12345"}`)
	resp, err := depositView.CreateDepositView(context.Background(), "account-123", mv2607.SourceSystemJHSilverlake, document)
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, method)
	require.Equal(t, "/underwriting/account-123/deposit-accounts", path)
	require.Equal(t, moov.Version2026_07.String(), version)
	require.Equal(t, string(mv2607.SourceSystemJHSilverlake), sourceSystem)
	require.Equal(t, "application/json", contentType)
	require.Equal(t, document, body)

	require.NotNil(t, resp)
	require.Equal(t, "account-123", resp.MoovAccountID)
	require.Equal(t, mv2607.SourceSystemJHSilverlake, resp.SourceSystem)
	require.Equal(t, "987654321", resp.SourceAccountID)
	require.Equal(t, "2026-06-18T12:00:00Z", resp.IngestedAt.UTC().Format(time.RFC3339))
}

func TestCreateDepositView_JHSilverlakeRecord(t *testing.T) {
	var received jhSilverlakeRecord

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, string(mv2607.SourceSystemJHSilverlake), r.Header.Get("X-Source-System"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&received))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"moovAccountID": "account-123",
			"sourceSystem": "jh_silverlake",
			"sourceAccountID": "987654321",
			"ingestedAt": "2026-06-18T12:00:00Z"
		}`))
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	depositView := mv2607.NewDepositViewClient(client)

	strPtr := func(s string) *string { return &s }
	numPtr := func(s string) *json.Number { n := json.Number(s); return &n }
	intPtr := func(i int) *int { return &i }

	record := jhSilverlakeRecord{
		SrcKey:    strPtr("src-1"),
		FornKey:   strPtr("cust-1"),
		AcctId:    numPtr("987654321"),
		AcctType:  strPtr("DDA"),
		CurBal:    numPtr("1234.56"),
		AvlBal:    numPtr("1000.00"),
		NumCrMTD:  intPtr(3),
		NumDrMTD:  intPtr(7),
		BrandCode: strPtr("MOOV"),
	}

	document, err := json.Marshal(record)
	require.NoError(t, err)

	resp, err := depositView.CreateDepositView(context.Background(), "account-123", mv2607.SourceSystemJHSilverlake, document)
	require.NoError(t, err)

	require.Equal(t, record, received)

	require.NotNil(t, resp)
	require.Equal(t, mv2607.SourceSystemJHSilverlake, resp.SourceSystem)
	require.Equal(t, "987654321", resp.SourceAccountID)
}

func TestParseSourceSystem(t *testing.T) {
	cases := []struct {
		input string
		want  mv2607.SourceSystem
	}{
		{"jh_silverlake", mv2607.SourceSystemJHSilverlake},
		{"jh_cif2020", mv2607.SourceSystemJHCIF2020},
		{"jh_coredirector", mv2607.SourceSystemJHCoreDirector},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, ok := mv2607.ParseSourceSystem(tc.input)
			require.True(t, ok)
			require.Equal(t, tc.want, got)
		})
	}

	_, ok := mv2607.ParseSourceSystem("unknown")
	require.False(t, ok)
}

// jhSilverlakeRecord models the deposit account payload Jack Henry SilverLake
// sends as the deposit view document. It is a test fixture only; CreateDepositView
// forwards the raw bytes, so callers may marshal whatever the source system emits.
//
// Field types deliberately mirror SilverLake's wire format rather than Moov's
// internal representations. Monetary amounts (CurBal, AvlBal, etc.) are json.Number
// because SilverLake emits them as decimal-string floats; this is NOT Moov's
// money model (integer cents) and must not be reused as a production type.
type jhSilverlakeRecord struct {
	Dlt          *string      `json:"Dlt"`          // Delete Record Flag
	SrcKey       *string      `json:"SrcKey"`       // Unique Record Key
	FornKey      *string      `json:"FornKey"`      // Customer Key
	AcctId       *json.Number `json:"AcctId"`       // Account number
	AcctType     *string      `json:"AcctType"`     // Account type
	AcctStat     *json.Number `json:"AcctStat"`     // Status
	AcctStatDesc *string      `json:"AcctStatDesc"` // Status Description
	CurBal       *json.Number `json:"CurBal"`       // Current balance
	CustId       *string      `json:"CustId"`       // Customer Id
	AvlBal       *json.Number `json:"AvlBal"`       // Available balance
	ColBal       *json.Number `json:"ColBal"`       // Collected balance
	HldAmt       *json.Number `json:"HldAmt"`       // Hold amount
	FltAmt       *json.Number `json:"FltAmt"`       // Total Accrual Float
	RegCCAmt     *json.Number `json:"RegCCAmt"`     // Tot Availibility Float
	ODProtAmt    *json.Number `json:"ODProtAmt"`    // Protection Amount
	AvlSweepAmt  *json.Number `json:"AvlSweepAmt"`  // Available Sweep Amount

	NetMemoPostAmt *json.Number `json:"NetMemoPostAmt"` // Net Memo Post Amount
	ODLmt          *json.Number `json:"ODLmt"`          // OD limit
	ODPrvlgAmt     *json.Number `json:"ODPrvlgAmt"`     // Privilege Amount
	LastStmtDt     *json.Number `json:"LastStmtDt"`     // Last Statement (CYYMMDD)
	LastDepDt      *json.Number `json:"LastDepDt"`      // Last Deposit (CYYMMDD)
	LastConDt      *json.Number `json:"LastConDt"`      // Last Contact (CYYMMDD)
	HldMailCode    *string      `json:"HldMailCode"`    // Hold mail
	HldMailDesc    *string      `json:"HldMailDesc"`    // Hold mail description
	LastActDt      *json.Number `json:"LastActDt"`      // Last Active (CYYMMDD)
	BrCode         *json.Number `json:"BrCode"`         // Branch number
	BrDesc         *string      `json:"BrDesc"`         // Bank/branch name

	ProdCode         *string `json:"ProdCode"`         // Service charge code
	ProdDesc         *string `json:"ProdDesc"`         // S/C description
	AcctClsfCode     *string `json:"AcctClsfCode"`     // Class
	AcctClsfDesc     *string `json:"AcctClsfDesc"`     // Class description
	IntBearAcct      *string `json:"IntBearAcct"`      // Accrue Int?
	TellerSICCode    *string `json:"TellerSICCode"`    // Special instruction code
	TellerSICDesc    *string `json:"TellerSICDesc"`    // Teller Special Inst Desc
	SigVerifyCode    *string `json:"SigVerifyCode"`    // Signature verify
	SigVerifyDesc    *string `json:"SigVerifyDesc"`    // Verify Sig Desc
	OffCode          *string `json:"OffCode"`          // Officer
	OffDesc          *string `json:"OffDesc"`          // Officer name
	SerChgWav        *string `json:"SerChgWav"`        // Service charge type
	SalesPerson      *string `json:"SalesPerson"`      // Sales Associate
	SerChgWavRsnCode *string `json:"SerChgWavRsnCode"` // Waive Reason Code
	SerChgWavRsnDesc *string `json:"SerChgWavRsnDesc"` // Waive Reason Description

	TINCode        *string      `json:"TINCode"`        // TIN status
	TINDesc        *string      `json:"TINDesc"`        // TIN Description
	TaxId          *json.Number `json:"TaxId"`          // Tax ID number
	DormantChgWav  *string      `json:"DormantChgWav"`  // Dormant Waive Code
	OpenDt         *json.Number `json:"OpenDt"`         // Date Opened (CYYMMDD)
	ClsDt          *json.Number `json:"ClsDt"`          // Date Closed (CYYMMDD)
	ElecStmtType   *string      `json:"ElecStmtType"`   // Only Statement Type
	LangType       *string      `json:"LangType"`       // Language Code
	MinBal         *json.Number `json:"MinBal"`         // Minimum balance
	HighStmtBal    *json.Number `json:"HighStmtBal"`    // Maximum balance
	AvgStmtColBal  *json.Number `json:"AvgStmtColBal"`  // Agr Col Bal This Qtr
	AvgStmtLdgrBal *json.Number `json:"AvgStmtLdgrBal"` // Agr Led Bal This Qtr
	AbbName        *string      `json:"AbbName"`        // Short name

	IntRate        *json.Number `json:"IntRate"`        // Interest Rate
	AccrInt        *json.Number `json:"AccrInt"`        // Accrued interest
	PrevYTDIntPaid *json.Number `json:"PrevYTDIntPaid"` // Previous YTD int
	YTDInt         *json.Number `json:"YTDInt"`         // YTD interest
	CallRptCode    *string      `json:"CallRptCode"`    // Call report code
	ClubPln        *json.Number `json:"ClubPln"`        // Club plan

	MTDODFee        *json.Number `json:"MTDODFee"`        // OD charge MTD
	MTDSerChgFee    *json.Number `json:"MTDSerChgFee"`    // M-T-D Fees
	MTDRetChkFee    *json.Number `json:"MTDRetChkFee"`    // Returned ck chg MTD
	MTDOthFee       *json.Number `json:"MTDOthFee"`       // Month to date fees charged
	MTDWavSerChgFee *json.Number `json:"MTDWavSerChgFee"` // M-T-D Service chg waived
	MTDWavODFee     *json.Number `json:"MTDWavODFee"`     // OD charge waived MTD
	MTDWavRetChkFee *json.Number `json:"MTDWavRetChkFee"` // Returned ck wve MTD
	MTDWavOthFee    *json.Number `json:"MTDWavOthFee"`    // Month to date Waived

	AnlysAcct     *string      `json:"AnlysAcct"`     // A/A code
	AccrMeth      *string      `json:"AccrMeth"`      // Acc Method
	AccrMethDesc  *string      `json:"AccrMethDesc"`  // Accrued Method Desc
	ClsOnZeroBal  *string      `json:"ClsOnZeroBal"`  // Close on zero balance
	FedWithCode   *json.Number `json:"FedWithCode"`   // Fed W/H code
	StateWithCode *json.Number `json:"StateWithCode"` // State W/H code
	ChgdOffAmt    *json.Number `json:"ChgdOffAmt"`    // Charged-off amount
	ClsRsnCode    *string      `json:"ClsRsnCode"`    // Closed reason Code
	ClsRsnDesc    *string      `json:"ClsRsnDesc"`    // Closed Description

	AggCurBalCycleAmt *json.Number `json:"AggCurBalCycleAmt"` // Aggregate ledg this stmt
	AggCycleDay       *json.Number `json:"AggCycleDay"`       // Aggregate days this stmt
	LastODDt          *json.Number `json:"LastODDt"`          // Last Overdrawn (CYYMMDD)
	ConsDaysNSF       *int         `json:"ConsDaysNSF"`       // Consec days NSF
	ConsDaysOD        *int         `json:"ConsDaysOD"`        // Consec days OD
	ODTimesMTD        *int         `json:"ODTimesMTD"`        // Times OD MTD
	ODTimesQTD        *int         `json:"ODTimesQTD"`        // Times OD this qtr
	ODTimesPrevQtr    *int         `json:"ODTimesPrevQtr"`    // Times OD 2nd qtr
	EscheatDt         *json.Number `json:"EscheatDt"`         // Date Last Escheat (CYYMMDD)
	GLCostCtr         *json.Number `json:"GLCostCtr"`         // GL Cost Center
	DormantDt         *json.Number `json:"DormantDt"`         // Date Last Dormant (CYYMMDD)

	AggCurBalIntCycleAmt *json.Number `json:"AggCurBalIntCycleAmt"` // Aggr ledger bal Int cyc
	AggIntCycleDay       *json.Number `json:"AggIntCycleDay"`       // Aggregate days Int cycle
	IntCycle             *json.Number `json:"IntCycle"`             // Interest Cr cycle
	LastDepAmt           *json.Number `json:"LastDepAmt"`           // Amt of last deposit
	LastWthdwlAmt        *json.Number `json:"LastWthdwlAmt"`        // Amt of last withdrawal
	LastSerChgDt         *json.Number `json:"LastSerChgDt"`         // Date Last Service Charge (CYYMMDD)
	LmtTrnAcctType       *string      `json:"LmtTrnAcctType"`       // Track Excess Drs?

	ACHCrCycleAmt *json.Number `json:"ACHCrCycleAmt"` // ACH credits s/c cycle amount
	ACHDrCycleAmt *json.Number `json:"ACHDrCycleAmt"` // ACH debits s/c cycle amount
	ACHCrCycleCnt *int         `json:"ACHCrCycleCnt"` // ACH credits s/c cycle count
	ACHDrCycleCnt *int         `json:"ACHDrCycleCnt"` // ACH debits s/c cycle count

	ATMOnUsCrMTDAmt *json.Number `json:"ATMOnUsCrMTDAmt"` // Amt of Onus ATM credits
	ATMOnUsDrMTDAmt *json.Number `json:"ATMOnUsDrMTDAmt"` // Amt of Onus ATM debits
	ATMOnUsMTDCnt   *int         `json:"ATMOnUsMTDCnt"`   // No. of Onus ATM credits
	ATMOnUsDrMTDCnt *int         `json:"ATMOnUsDrMTDCnt"` // No. of Onus ATM debits
	MTDIntPaid      *json.Number `json:"MTDIntPaid"`      // M-T-D Interest paid

	NSFItemsMTD      *int         `json:"NSFItemsMTD"`      // Items NSF This cyc
	NSFTimesMTD      *int         `json:"NSFTimesMTD"`      // Times NSF MTD
	AggODCycleAmt    *json.Number `json:"AggODCycleAmt"`    // Agg Led balance negative
	ODDaysMTD        *int         `json:"ODDaysMTD"`        // Days OD This cyc
	ODItemsMTD       *int         `json:"ODItemsMTD"`       // Items OD This cyc
	SerChgCrItemsCnt *int         `json:"SerChgCrItemsCnt"` // Number S/C credits
	SerChgDrItemsCnt *int         `json:"SerChgDrItemsCnt"` // Number S/C debits
	NSFAmt           *json.Number `json:"NSFAmt"`           // Amount NSF
	ODLmtCode        *json.Number `json:"ODLmtCode"`        // Overdraft limit code
	NSFItemsRet      *int         `json:"NSFItemsRet"`      // Items Returns This cyc
	CRASMSACode      *json.Number `json:"CRASMSACode"`      // SMSA Number
	CRASMSADesc      *string      `json:"CRASMSADesc"`      // SMSA Description
	AutoNSFFee       *string      `json:"AutoNSFFee"`       // Automatic NSF charge
	RateRevDt        *json.Number `json:"RateRevDt"`        // Rate review date (CYYMMDD)

	AmtCrMTD     *json.Number `json:"AmtCrMTD"`     // Amt of credits
	AmtDrMTD     *json.Number `json:"AmtDrMTD"`     // Amt of debits
	ChkGuar      *string      `json:"ChkGuar"`      // Check guaranty card
	ChgdOffDt    *json.Number `json:"ChgdOffDt"`    // Date Last Charged Off (CYYMMDD)
	NumCrMTD     *int         `json:"NumCrMTD"`     // Number of credits
	NumDrMTD     *int         `json:"NumDrMTD"`     // Number of debits
	NumRetChkMTD *int         `json:"NumRetChkMTD"` // No. returned check MTD
	StmtCycle    *json.Number `json:"StmtCycle"`    // Statement cycle
	StmtBal      *json.Number `json:"StmtBal"`      // Previous stmt bal

	ODTimes3Qtr     *int `json:"ODTimes3Qtr"`     // Times OD 3rd qtr
	ODTimes4Qtr     *int `json:"ODTimes4Qtr"`     // Times OD 4th qtr
	NSFTimesQTD     *int `json:"NSFTimesQTD"`     // Times NSF This Qtr
	NSFTimesPrevQtr *int `json:"NSFTimesPrevQtr"` // Times NSF 2nd Qtr
	NSFTimes3Qtr    *int `json:"NSFTimes3Qtr"`    // Times NSF 3rd Qtr
	NSFTimes4Qtr    *int `json:"NSFTimes4Qtr"`    // Times NSF 4th Qtr

	YestBal          *json.Number `json:"YestBal"`          // Yesterday balance
	YestColBal       *json.Number `json:"YestColBal"`       // Yesterday collected balance
	NSFFeesYTD       *json.Number `json:"NSFFeesYTD"`       // YTD NSF Fees
	ODFeesYTD        *json.Number `json:"ODFeesYTD"`        // YTD OD Fees
	SerChgFeesYTD    *json.Number `json:"SerChgFeesYTD"`    // YTD S/C
	SerChgFeesWavYTD *json.Number `json:"SerChgFeesWavYTD"` // Y-T-D Service chg waived

	OneDayFlt   *json.Number `json:"OneDayFlt"`   // Day 1 Float
	TwoDayFlt   *json.Number `json:"TwoDayFlt"`   // Day 2 Float
	ThreeDayFlt *json.Number `json:"ThreeDayFlt"` // Day 3 Float
	FourDayFlt  *json.Number `json:"FourDayFlt"`  // Day 4 Float
	FiveDayFlt  *json.Number `json:"FiveDayFlt"`  // Day 5 Float

	ChgBackItemsQTD *int         `json:"ChgBackItemsQTD"` // Chg Back Itm qtr
	LastPayDt       *json.Number `json:"LastPayDt"`       // Last Interest Paid (CYYMMDD)
	AmtLastIntPd    *json.Number `json:"AmtLastIntPd"`    // Amt of last int Pd
	OrigAcctId      *json.Number `json:"OrigAcctId"`      // Original Account Id
	ATMCard         *string      `json:"ATMCard"`         // ATM Card
	StdIndustCode   *json.Number `json:"StdIndustCode"`   // SIC Code
	StdIndustDesc   *string      `json:"StdIndustDesc"`   // SIC Description
	CurQtrMinBal    *json.Number `json:"CurQtrMinBal"`    // Quarter minimum balance
	CurQtrMaxBal    *json.Number `json:"CurQtrMaxBal"`    // Quarter maximum balance

	TrnAmtCurDayDr  *json.Number `json:"TrnAmtCurDayDr"`  // Trn Debit CurDay Amt
	TrnAmtPrevDayDr *json.Number `json:"TrnAmtPrevDayDr"` // Trn Debit PrevDay Amt
	TrnCntCurDayDr  *int         `json:"TrnCntCurDayDr"`  // Trn Debit CurDay Cnt
	TrnCntPrevDayDr *int         `json:"TrnCntPrevDayDr"` // Trn Debit PrevDay Cnt
	TrnAmtCurDayCr  *json.Number `json:"TrnAmtCurDayCr"`  // Trn Credit CurDay Amt
	TrnAmtPrevDayCr *json.Number `json:"TrnAmtPrevDayCr"` // Trn Credit PrevDay Amt
	TrnCntCurDayCr  *int         `json:"TrnCntCurDayCr"`  // Trn Credit CurDay Cnt
	TrnCntPrevDayCr *int         `json:"TrnCntPrevDayCr"` // Trn Credit PrevDay Cnt

	MonthAvgBalAmt *json.Number `json:"MonthAvgBalAmt"` // M-T-D Average Available Balance
	GroupCode      *json.Number `json:"GroupCode"`      // GL Group Code
	GLProdCode     *json.Number `json:"GLProdCode"`     // GL Product Code
	MTDAvgColBal   *json.Number `json:"MTDAvgColBal"`   // M-T-D Average Collected Balance
	RateVar        *json.Number `json:"RateVar"`        // Rate Variance
	RateVarCode    *string      `json:"RateVarCode"`    // Rate Variance Code
	PrimeRateIdx   *json.Number `json:"PrimeRateIdx"`   // Rate Number

	EstbPersonName  *string `json:"EstbPersonName"`  // Name of person opening account
	EstbPersonTitle *string `json:"EstbPersonTitle"` // Title of person opening account
	BrandCode       *string `json:"BrandCode"`       // Brand Code
	BrandDesc       *string `json:"BrandDesc"`       // Brand Description

	NSFItemsYTD *int `json:"NSFItemsYTD"` // NSF Year to Date Items
	NSFTimesYTD *int `json:"NSFTimesYTD"` // NSF Year to Date Times

	StopCondNotfType      *string `json:"StopCondNotfType"`      // Stop code
	HldCondNotfType       *string `json:"HldCondNotfType"`       // Hold Exists
	AlrtCondNotfType      *string `json:"AlrtCondNotfType"`      // Alert Exists
	SpecInstrCondNotfType *string `json:"SpecInstrCondNotfType"` // Special Instructions Exists

	CurBalDailyClsAmt *json.Number `json:"CurBalDailyClsAmt"` // Current Balance Daily Close Amount
	AvlBalDailyClsAmt *json.Number `json:"AvlBalDailyClsAmt"` // Available Balance Daily Close Amount
	ColBalDailyClsAmt *json.Number `json:"ColBalDailyClsAmt"` // Collected Balance Daily Close Amount
}
