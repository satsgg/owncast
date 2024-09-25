package lnd

import (
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc"
)

var Client lnrpc.LightningClient
var clientErr error

func Start() {
	// tlsPath := "/Users/chad/.polar/networks/2/volumes/lnd/alice/tls.cert"
	// macDir := "/Users/chad/.polar/networks/2/volumes/lnd/alice/data/chain/bitcoin/regtest/"
	tlsPath := "/Users/chad/.polar/networks/5/volumes/lnd/bob/tls.cert"
	macDir := "/Users/chad/.polar/networks/5/volumes/lnd/bob/data/chain/bitcoin/regtest/"


	// hex encoded lnd invoice macaroon
	// invoiceMacaroon := lndclient.MacaroonData("0201036c6e640258030a10e56658597e117eb7afb006dfbfd542d81201301a160a0761646472657373120472656164120577726974651a170a08696e766f69636573120472656164120577726974651a0f0a076f6e636861696e120472656164000006201f3e075f198f445a31db4d3bef265f00c238b81fad232df67a82b27b450dc0c5")

	// hex encoded lnd tls certificate
	// tlsCert := lndclient.TLSData("2d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d494943506a43434165536741774942416749524149444b667330785a535361783755777962594f6f6e5577436759494b6f5a497a6a3045417749774d5445660a4d4230474131554543684d576247356b494746316447396e5a57356c636d46305a575167593256796444454f4d4177474131554541784d4659577870593255770a4868634e4d6a51774e5445774d5455784e7a45315768634e4d6a55774e7a41314d5455784e7a4531576a41784d523877485159445651514b45785a73626d51670a595856306232646c626d56795958526c5a43426a5a584a304d51347744415944565151444577566862476c6a5a54425a4d424d4742797147534d3439416745470a43437147534d34394177454841304941424b6c4a41666c7759625941325a656c77654d304d466c6d434949375253784742626b5a6f42436b5661415078326c570a5a37733975624f6d4238596a6975502f38326a624d2f55703067354f7452514336637a393463476a6764777767646b7744675944565230504151482f424151440a41674b6b4d424d47413155644a51514d4d416f47434373474151554642774d424d41384741315564457745422f7751464d414d4241663877485159445652304f0a42425945464e52753946795830645177734d2b396e72423530376d30664857634d49474242674e5648524545656a42346767566862476c6a5a59494a6247396a0a5957786f62334e306767566862476c6a5a59494f6347397359584974626a4974595778705932574346476876633351755a47396a613256794c6d6c75644756790a626d467367675231626d6c3467677031626d6c346347466a613256306767646964575a6a623235756877522f41414142687841414141414141414141414141410a414141414141414268775373465141444d416f4743437147534d343942414d43413067414d4555434951436c673145632f4750316177665a61594444716739320a4f34754f42644f436c6b57495453703074674b784f5149675a54384e6b506c61592f5865516a694530643142563273567449497973332f2f3142554b6846366a0a3966383d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a")

	// lnd grpc url
	// lndHost := "127.0.0.1:10001"
	lndHost := "127.0.0.1:10002"

	macFilename := lndclient.MacFilename(
		"invoice.macaroon",
	)

	Client, clientErr = lndclient.NewBasicClient(
		lndHost, 
		// "",
		tlsPath,
		// "",
		macDir,
		// tlsPath,
		"regtest",
		// tlsCert,
		// invoiceMacaroon,
		macFilename,
	)

	if clientErr != nil {
		println("error creating basic client", clientErr)
	} else {
		println("connected to lnd!")
	}

	// invioces, err := client.ListInvoices(context.Background(), lndclient.Invoice{
	// invoice, err := client.AddInvoice(context.Background(), &lnrpc.Invoice{
	// 	Memo: "test",
	// 	Value: 100000,
	// })
	// if err != nil {
	// 	println("error creating invoice")
	// 	return nil
	// }
	// println("invoice", invoice)
	// println("invoice", invoice.String())
	// return client
}