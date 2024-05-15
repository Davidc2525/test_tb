package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	types "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func main() {
	start := time.Now()
	tbAddress := os.Getenv("TB_ADDRESS")
	if len(tbAddress) == 0 {
		tbAddress = "3000"
	}
	client, err := tb.NewClient(types.ToUint128(0), []string{tbAddress}, 256)
	if err != nil {
		fmt.Printf("Error creating client: %s", err)
		return
	}
	defer client.Close()

	if len(os.Args) <= 1 {
		show_help()
		return
	}
	switch os.Args[1] {
	case "create_bank":
		id_account, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("create_bank", id_account)
		create_bank(client, id_account)
	case "new":
		id_account, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("new", id_account)
		new_account(client, id_account)

	case "get":
		id_account, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("get")
		get_account(client, id_account)
	case "transfer":
		debit, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		credit, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		amount, err := strconv.ParseUint(os.Args[4], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("transfer", debit, credit, amount)
		new_transfer(client, debit, credit, amount)
	case "loop":
		debit, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		credit, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		amount, err := strconv.ParseUint(os.Args[4], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		i_, err := strconv.ParseInt(os.Args[5], 10, 32)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("loop", i_, debit, credit, amount)
		new_transfer_loop(client, int(i_), debit, credit, amount)
	case "help":
		show_help()
	default:
		show_help()

	}

	elapsed := time.Since(start)
	fmt.Printf("took %s\n", elapsed)

}
func show_help() {
	fmt.Println(`
uso
./main help #show this help msg
./main new <account_number> #create new account
./main get <account_number>
./main create_bank <account_number>
./main transfer <debit_account_number> <credit_account_number> <amount> # use the bank account to transfer
./loop <debit_account_number>  <credit_account_number> <ammount> <iterate>  #iterate x 8000 

If it is the first time you must create the file from the tiger directory

	cd tiger

and create the file

	./tigerbeetle format --cluster=0 --replica=0 --replica-count=1 0_0.tigerbeetle

then in that same directory you execute TB

	./tigerbeetle start --addresses=3000 0_0.tigerbeetle

		`)
}
func get_account(client tb.Client, id uint64) (uint64, error) {
	accounts, err := client.LookupAccounts([]types.Uint128{types.ToUint128(id)})
	if err != nil {
		fmt.Printf("Could not fetch accounts: %s", err)
		return 0, err
	}
	account := accounts[0]
	big_credit := account.CreditsPosted.BigInt()
	big_debit := account.DebitsPosted.BigInt()

	balance := big_credit.Sub(&big_credit, &big_debit)
	fmt.Println("balance", balance)

	return balance.Uint64(), nil
}
func create_bank(client tb.Client, id uint64) {
	accountsRes, err := client.CreateAccounts([]types.Account{
		{
			ID:             types.ToUint128(id),
			DebitsPending:  types.ToUint128(0),
			DebitsPosted:   types.ToUint128(0),
			CreditsPending: types.ToUint128(0),
			CreditsPosted:  types.ToUint128(0),
			UserData128:    types.ToUint128(0),
			UserData64:     0,
			UserData32:     0,
			Reserved:       0,
			Ledger:         1,
			Code:           718,
			Flags:          0,
			Timestamp:      0,
		},
	})
	if err != nil {
		fmt.Printf("Error creating accounts: %s", err)

	}

	for _, err_r := range accountsRes {
		fmt.Printf("Error creating account %d: %s", err_r.Index, err_r.Result)

	}

}
func new_transfer(client tb.Client, debit uint64, credit uint64, amount uint64) error {
	transfers := []types.Transfer{{
		ID:              types.ID(),
		DebitAccountID:  types.ToUint128(debit),
		CreditAccountID: types.ToUint128(credit),
		Amount:          types.ToUint128(amount),
		PendingID:       types.ToUint128(0),
		UserData128:     types.ToUint128(2),
		UserData64:      0,
		UserData32:      0,
		Timeout:         0,
		Ledger:          1,
		Code:            1,
		Flags:           0,
		Timestamp:       0,
	}}

	transfersRes, err := client.CreateTransfers(transfers)
	if err != nil {
		fmt.Printf("Error creating transfer batch: %s", err)
		return err
	}
	fmt.Println(transfersRes)
	return nil
}
func new_transfer_loop(client tb.Client, i_ int, debit uint64, credit uint64, amount uint64) error {
	var transfers = make([]types.Transfer, 8000)
	for i := 0; i < i_; i++ {
		start := time.Now()

		for x := 0; x < 8000; x++ {
			transfers[x] = types.Transfer{
				ID:              types.ID(),
				DebitAccountID:  types.ToUint128(debit),
				CreditAccountID: types.ToUint128(credit),
				Amount:          types.ToUint128(amount),
				PendingID:       types.ToUint128(0),
				UserData128:     types.ToUint128(2),
				UserData64:      0,
				UserData32:      0,
				Timeout:         0,
				Ledger:          1,
				Code:            1,
				Flags:           0,
				Timestamp:       0,
			}

		}
		r, err := client.CreateTransfers(transfers)
		elapsed := time.Since(start)
		fmt.Printf("[%v]batch by 8000 took %s %v\n", i, elapsed, r)
		if err != nil {
			fmt.Printf("Error creating transfer batch: %s\n", err)
		}

	}

	return nil
}
func new_account(client tb.Client, id uint64) error {
	accountsRes, err := client.CreateAccounts([]types.Account{
		{
			ID:             types.ToUint128(id),
			DebitsPending:  types.ToUint128(0),
			DebitsPosted:   types.ToUint128(0),
			CreditsPending: types.ToUint128(0),
			CreditsPosted:  types.ToUint128(0),
			UserData128:    types.ToUint128(0),
			UserData64:     0,
			UserData32:     0,
			Reserved:       0,
			Ledger:         1,
			Code:           718,
			Flags:          types.AccountFlags{DebitsMustNotExceedCredits: true}.ToUint16(),
			Timestamp:      0,
		},
	})
	if err != nil {
		fmt.Printf("Error creating accounts: %s", err)

	}

	for _, err_r := range accountsRes {
		fmt.Printf("Error creating account %d: %s", err_r.Index, err_r.Result)

	}

	return err
}
